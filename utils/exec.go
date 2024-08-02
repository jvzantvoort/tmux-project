package utils

import (
	"bufio"
	"os/exec"
	"strings"

	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	wg sync.WaitGroup
)

func cleanup() {
	if r := recover(); r != nil {
		log.Errorf("Paniced %s", r)
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
		go item.Run()
	}
	wg.Wait()
}

func NewQueue() *Queue {
	return &Queue{}
}

// Run runs an item
func (e QueueElement) Run() {
	defer wg.Done() // lower counter
	defer cleanup() // handle panics

	stdout_list, stderr_list, eerror := Exec(e.Cwd, e.Command)
	for _, stdout_line := range stdout_list {
		log.Infof("<stdout> %s", stdout_line)
	}
	for _, stderr_line := range stderr_list {
		log.Errorf("<stderr> %s", stderr_line)
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
	log.Debugf("command: %s %s", cmnd, strings.Join(commandlist, " "))
	log.Debugf("         cwd: %s", cwd)
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
		log.Debugln(msg)
		stdout_list = append(stdout_list, msg)
	}

	stderr_scan := bufio.NewScanner(stderr)
	stderr_scan.Split(bufio.ScanLines)
	for stderr_scan.Scan() {
		msg := stderr_scan.Text()
		log.Errorln(msg)
		stderr_list = append(stderr_list, msg)
	}

	eerror := command.Wait()
	if eerror != nil {
		log.Errorf("command failed, %v", err)
	}

	return stdout_list, stderr_list, eerror
}
