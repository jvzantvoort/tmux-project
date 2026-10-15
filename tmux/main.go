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

// run executes tmux with args, connected to the terminal
func run(args ...string) error {
	cmd := exec.Command(utils.Which("tmux"), args...) // #nosec G204 - arguments are built by this package
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tmux %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

// sessionNameReplacer maps the characters tmux does not allow in session
// names the way tmux itself does
var sessionNameReplacer = strings.NewReplacer(".", "_", ":", "_")

// SessionName returns the name tmux uses for a session created as name.
func SessionName(name string) string {
	return sessionNameReplacer.Replace(name)
}

// target returns the tmux target that matches session name exactly,
// instead of tmux's default prefix and pattern matching.
func target(name string) string {
	return "=" + name
}

// newArgs returns the arguments to create session name. directory is the
// start directory of the session; it is left out when empty.
func newArgs(name, configfile, directory string, detached bool) []string {
	retv := []string{"-f", configfile, "new-session"}
	if detached {
		retv = append(retv, "-d")
	}
	retv = append(retv, "-s", name)
	if len(directory) != 0 {
		retv = append(retv, "-c", directory)
	}
	return retv
}

// resumeCommands returns the tmux invocations that bring the user to
// session name. Inside tmux the client is switched, as attaching from
// within a session would nest tmux.
func resumeCommands(name, configfile, directory string, exists, insideTmux bool) [][]string {
	if !insideTmux {
		if exists {
			return [][]string{{"-f", configfile, "attach-session", "-d", "-t", target(name)}}
		}
		return [][]string{newArgs(name, configfile, directory, false)}
	}

	retv := [][]string{}
	if !exists {
		retv = append(retv, newArgs(name, configfile, directory, true))
	}
	return append(retv, []string{"switch-client", "-t", target(name)})
}

// New creates a new tmux session with the specified name and configuration
// file, starting in directory when it is not empty.
func New(name, configfile, directory string) error {
	return run(newArgs(name, configfile, directory, false)...)
}

// Attach attaches to an existing tmux session with the specified name and configuration
func Attach(name, configfile string) error {
	return run(resumeCommands(name, configfile, "", true, false)[0]...)
}

// Resume brings the user to session name: it attaches to (or, inside tmux,
// switches to) the session when it exists, otherwise it creates it first.
func Resume(name, configfile, directory string) error {
	name = SessionName(name)

	active, err := ListActive()
	if err != nil {
		return err
	}

	exists := utils.StringInSlice(name, active)
	insideTmux := len(os.Getenv("TMUX")) != 0

	for _, args := range resumeCommands(name, configfile, directory, exists, insideTmux) {
		if err := run(args...); err != nil {
			return err
		}
	}
	return nil
}

// CurrentSession returns the name of the tmux session this process runs
// in, or an empty string when it does not run inside tmux.
func CurrentSession() string {
	if len(os.Getenv("TMUX")) == 0 {
		return ""
	}
	cwd, _ := os.UserHomeDir()
	stdout_lines, _, err := utils.ExecSilent(cwd, utils.Which("tmux")+" display-message -p #S")
	if err != nil || len(stdout_lines) == 0 {
		return ""
	}
	return strings.TrimSpace(stdout_lines[0])
}

// ListActive returns a list of currently active tmux session names
func ListActive() ([]string, error) {
	command := fmt.Sprintf("%s ls -F #{session_name}", utils.Which("tmux"))
	cwd, _ := os.UserHomeDir()
	retv := []string{}
	stdout_lines, _, _ := utils.ExecSilent(cwd, command)
	for _, line := range stdout_lines {
		retv = append(retv, strings.Trim(line, "\""))
	}
	return retv, nil
}
