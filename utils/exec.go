package utils

import (
	"bufio"
	"os/exec"
	"strings"

	"fmt"
	"sync"
)

var (
	wg sync.WaitGroup
)

func cleanup() {
	if r := recover(); r != nil {
		Errorf("Paniced %s", r)
	}
}

type QueueElement struct {
	Cwd     string
	Command string
}

type Queue struct {
	Queue []QueueElement
}

// Add adds an item to the queue
func (q *Queue) Add(cwd, command string) {
	q.Queue = append(q.Queue, QueueElement{Cwd: cwd, Command: command})
}

// Run runs all items in the queue
func (q *Queue) Run() {
	for _, item := range q.Queue {
		wg.Add(1)
		go func(el QueueElement) {
			defer wg.Done()
			el.Run()
		}(item)
	}
	wg.Wait()
}

func NewQueue() *Queue {
	return &Queue{}
}

// Run runs an item
func (e QueueElement) Run() {
	defer cleanup() // handle panics

	stdout_list, stderr_list, eerror := Exec(e.Cwd, e.Command)
	for _, stdout_line := range stdout_list {
		Infof("<stdout> %s", stdout_line)
	}
	for _, stderr_line := range stderr_list {
		Errorf("<stderr> %s", stderr_line)
	}
	if eerror != nil {
		panic(fmt.Sprintf("Action \"%s\" failed", e.Command))
	}
}

// Exec executes a command
func Exec(cwd, args string) ([]string, []string, error) {
	commandlist := []string{}
	stdout_list := []string{}
	stderr_list := []string{}
	cmndargs := strings.Split(args, " ")
	cmnd := cmndargs[0]
	cmndargs = cmndargs[1:]
	commandlist = append(commandlist, cmndargs...)
	command := exec.Command(cmnd, commandlist...)
	Debugf("command: %s %s", cmnd, strings.Join(commandlist, " "))
	Debugf("         cwd: %s", cwd)
	command.Dir = cwd

	stdout, err := command.StdoutPipe()
	ErrorExit(err)

	stderr, err := command.StderrPipe()
	ErrorExit(err)

	err = command.Start()
	ErrorExit(err)

	stdout_scan := bufio.NewScanner(stdout)
	stdout_scan.Split(bufio.ScanLines)
	for stdout_scan.Scan() {
		msg := stdout_scan.Text()
		Debugf(msg)
		stdout_list = append(stdout_list, msg)
	}

	stderr_scan := bufio.NewScanner(stderr)
	stderr_scan.Split(bufio.ScanLines)
	for stderr_scan.Scan() {
		msg := stderr_scan.Text()
		Errorf(msg)
		stderr_list = append(stderr_list, msg)
	}

	eerror := command.Wait()
	if eerror != nil {
		Errorf("command failed, %v", err)
	}

	return stdout_list, stderr_list, eerror
}

// ExecSilent execute a command but don't care too much about the result.
func ExecSilent(cwd, args string) ([]string, []string, error) {
	commandlist := []string{}
	stdout_list := []string{}
	stderr_list := []string{}

	cmndargs := strings.Split(args, " ")

	cmnd := cmndargs[0]
	cmndargs = cmndargs[1:]
	commandlist = append(commandlist, cmndargs...)
	command := exec.Command(cmnd, commandlist...)
	command.Dir = cwd

	stdout, _ := command.StdoutPipe()
	stderr, _ := command.StderrPipe()

	err := command.Start()
	if err != nil {
		return stdout_list, stderr_list, err
	}

	stdout_scan := bufio.NewScanner(stdout)
	stdout_scan.Split(bufio.ScanLines)
	for stdout_scan.Scan() {
		msg := stdout_scan.Text()
		Debugf("stdout: %s", msg)
		stdout_list = append(stdout_list, msg)
	}

	stderr_scan := bufio.NewScanner(stderr)
	stderr_scan.Split(bufio.ScanLines)
	for stderr_scan.Scan() {
		msg := stderr_scan.Text()
		Debugf("stderr: %s", msg)
		stderr_list = append(stderr_list, msg)
	}

	eerror := command.Wait()
	if eerror != nil {
		Debugf("Error: command failed, %v", err)
	}

	return stdout_list, stderr_list, eerror
}
