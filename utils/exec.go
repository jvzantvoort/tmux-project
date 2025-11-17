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

// QueueElement represents a single command to be executed with its working directory
type QueueElement struct {
	Cwd     string
	Command string
}

// Queue manages a collection of commands to be executed concurrently
type Queue struct {
	Queue []QueueElement
}

// Add adds a command to the execution queue
func (q *Queue) Add(cwd, command string) {
	q.Queue = append(q.Queue, QueueElement{Cwd: cwd, Command: command})
}

// Run executes all queued commands concurrently using goroutines
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

// NewQueue creates and returns a new empty Queue
func NewQueue() *Queue {
	return &Queue{}
}

// Run executes a single queue element command
func (e QueueElement) Run() {
	defer cleanup()

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

// Exec executes a command in the specified working directory and returns stdout, stderr, and any error
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
		Debugf("%s", msg)
		stdout_list = append(stdout_list, msg)
	}

	stderr_scan := bufio.NewScanner(stderr)
	stderr_scan.Split(bufio.ScanLines)
	for stderr_scan.Scan() {
		msg := stderr_scan.Text()
		Errorf("%s", msg)
		stderr_list = append(stderr_list, msg)
	}

	eerror := command.Wait()
	if eerror != nil {
		Errorf("command failed, %v", err)
	}

	return stdout_list, stderr_list, eerror
}

// ExecSilent executes a command silently without extensive error handling,
// useful for commands where failure is acceptable
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
