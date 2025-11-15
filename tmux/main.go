// Package tmux provides utilities for managing tmux sessions and integrating
// with the tmux terminal multiplexer. It handles session creation, attachment,
// and listing active sessions.
package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jvzantvoort/tmux-project/utils"
)

// New creates a new tmux session with the specified name and configuration file
func New(name, configfile string) {
	command := []string{"-f", configfile, "new", "-s", name}

	cmd := exec.Command(utils.Which("tmux"), command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	utils.ErrorExit(cmd.Run())
}

// Attach attaches to an existing tmux session with the specified name and configuration
func Attach(name, configfile string) {
	command := []string{"-f", configfile, "attach", "-d", "-t", name}

	cmd := exec.Command(utils.Which("tmux"), command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	utils.ErrorExit(cmd.Run())
}

// Resume attaches to an existing tmux session if it exists, otherwise creates a new one
func Resume(name, configfile string) {
	active, err := ListActive()
	utils.ErrorExit(err)

	for _, active_name := range active {
		if active_name == name {
			Attach(name, configfile)
			return
		}
	}
	New(name, configfile)
}

// ListActive returns a list of currently active tmux session names
func ListActive() ([]string, error) {
	command := fmt.Sprintf("%s ls -F \"#{session_name}\"", utils.Which("tmux"))
	cwd, _ := os.UserHomeDir()
	retv := []string{}
	stdout_lines, _, _ := utils.ExecSilent(cwd, command)
	for _, line := range stdout_lines {
		retv = append(retv, strings.Trim(line, "\""))
	}
	return retv, nil
}
