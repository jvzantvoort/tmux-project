package utils

import (
	"os/exec"
)

// Which searches for a command in PATH and returns its full path.
// If the command is not found, returns the command name itself.
func Which(command string) string {
	if fullpath, err := exec.LookPath(command); err == nil {
		return fullpath
	}
	return command
}
