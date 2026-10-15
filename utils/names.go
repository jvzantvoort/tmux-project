package utils

import (
	"fmt"
	"strings"
	"unicode"
)

// CheckPathName returns an error when name cannot safely be used as a
// single path element, e.g. to build ~/.tmux.d/<name>.json.
func CheckPathName(kind, name string) error {
	if len(name) == 0 {
		return fmt.Errorf("%s name is empty", kind)
	}
	if name == "." || name == ".." {
		return fmt.Errorf("%s name %q is not allowed", kind, name)
	}
	if strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("%s name %q must not contain a path separator", kind, name)
	}
	for _, char := range name {
		if unicode.IsControl(char) {
			return fmt.Errorf("%s name %q must not contain control characters", kind, name)
		}
	}
	return nil
}

// CheckSessionName returns an error when name cannot be used both as a
// path element and as a tmux session name. tmux silently replaces "." and
// ":" in session names with "_", and the shell integration looks up the
// project files by session name, so those characters are rejected.
func CheckSessionName(kind, name string) error {
	if err := CheckPathName(kind, name); err != nil {
		return err
	}
	if strings.HasPrefix(name, "-") {
		return fmt.Errorf("%s name %q must not start with \"-\"", kind, name)
	}
	if strings.ContainsAny(name, ".:") {
		return fmt.Errorf("%s name %q must not contain \".\" or \":\", tmux replaces them in session names", kind, name)
	}
	for _, char := range name {
		if unicode.IsSpace(char) {
			return fmt.Errorf("%s name %q must not contain whitespace", kind, name)
		}
	}
	return nil
}
