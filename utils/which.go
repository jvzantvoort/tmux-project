package utils

import (
	"os"
	"path"
	"strings"
)

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
