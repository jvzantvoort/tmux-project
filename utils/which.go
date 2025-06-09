package utils

import (
	"os"
	"path"
	"strings"
)

// Which returns the full path of the command
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
