package utils

import (
	"os"
	"path"
	"strings"
)

// Which searches for a command in PATH and returns its full path.
// If the command is not found, returns the command name itself.
func Which(command string) string {
	Path := strings.Split(os.Getenv("PATH"), ":")

	for _, dirname := range Path {
		fullpath := path.Join(dirname, command)
		if _, err := os.Stat(fullpath); err == nil {
			if !os.IsNotExist(err) {
				return fullpath
			}
		}
	}
	return command
}
