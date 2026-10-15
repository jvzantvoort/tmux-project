package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// QueueElement represents a single command to be executed with its working directory
type QueueElement struct {
	Cwd     string
	Command string
}

// Queue manages an ordered list of commands to be executed
type Queue struct {
	Queue []QueueElement
}

// Add adds a command to the execution queue
func (q *Queue) Add(cwd, command string) {
	q.Queue = append(q.Queue, QueueElement{Cwd: cwd, Command: command})
}

// Run executes the queued commands one after another, in the order they
// were added. It stops at the first command that fails and returns its
// error, as later commands may depend on earlier ones.
func (q *Queue) Run() error {
	for _, item := range q.Queue {
		if err := item.Run(); err != nil {
			return err
		}
	}
	return nil
}

// NewQueue creates and returns a new empty Queue
func NewQueue() *Queue {
	return &Queue{}
}

// Run executes a single queue element command through the shell
func (e QueueElement) Run() error {
	stdout_list, stderr_list, eerror := ExecShell(e.Cwd, e.Command)
	for _, stdout_line := range stdout_list {
		Infof("<stdout> %s", stdout_line)
	}
	for _, stderr_line := range stderr_list {
		Errorf("<stderr> %s", stderr_line)
	}
	if eerror != nil {
		return fmt.Errorf("action %q failed: %w", e.Command, eerror)
	}
	return nil
}

// splitLines splits command output into lines
func splitLines(data []byte) []string {
	retv := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		retv = append(retv, scanner.Text())
	}
	return retv
}

// run executes command in cwd and returns its stdout and stderr lines.
// Both streams are collected at the same time, so a command writing a lot
// to one stream cannot block on the other.
func run(cwd string, command *exec.Cmd) ([]string, []string, error) {
	var stdout, stderr bytes.Buffer

	command.Dir = cwd
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	return splitLines(stdout.Bytes()), splitLines(stderr.Bytes()), err
}

// Exec executes a command in the specified working directory and returns stdout, stderr, and any error.
// The command is split on whitespace; no shell quoting is applied.
func Exec(cwd, args string) ([]string, []string, error) {
	cmndargs := strings.Fields(args)
	if len(cmndargs) == 0 {
		return []string{}, []string{}, fmt.Errorf("empty command")
	}
	Debugf("command: %s", strings.Join(cmndargs, " "))
	Debugf("         cwd: %s", cwd)

	stdout_list, stderr_list, eerror := run(cwd, exec.Command(cmndargs[0], cmndargs[1:]...)) // #nosec G204 - commands come from the user's own configuration
	for _, msg := range stdout_list {
		Debugf("%s", msg)
	}
	for _, msg := range stderr_list {
		Errorf("%s", msg)
	}
	if eerror != nil {
		Errorf("command failed, %v", eerror)
	}
	return stdout_list, stderr_list, eerror
}

// ExecShell executes command with "sh -c" in the specified working
// directory, so quoting, pipes and && work as they do in a shell.
func ExecShell(cwd, command string) ([]string, []string, error) {
	Debugf("command: sh -c %q", command)
	Debugf("         cwd: %s", cwd)
	return run(cwd, exec.Command("sh", "-c", command)) // #nosec G204 - commands come from the user's own configuration
}

// ExecSilent executes a command silently without extensive error handling,
// useful for commands where failure is acceptable
func ExecSilent(cwd, args string) ([]string, []string, error) {
	cmndargs := strings.Fields(args)
	if len(cmndargs) == 0 {
		return []string{}, []string{}, fmt.Errorf("empty command")
	}

	stdout_list, stderr_list, eerror := run(cwd, exec.Command(cmndargs[0], cmndargs[1:]...)) // #nosec G204 - commands come from the user's own configuration
	for _, msg := range stdout_list {
		Debugf("stdout: %s", msg)
	}
	for _, msg := range stderr_list {
		Debugf("stderr: %s", msg)
	}
	if eerror != nil {
		Debugf("Error: command failed, %v", eerror)
	}
	return stdout_list, stderr_list, eerror
}
