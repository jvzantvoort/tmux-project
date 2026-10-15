package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EditorCommand returns the editor command line: $VISUAL, then $EDITOR,
// then the default Editor. The variables may contain arguments, e.g.
// "code --wait".
func EditorCommand() []string {
	for _, name := range []string{"VISUAL", "EDITOR"} {
		if fields := strings.Fields(os.Getenv(name)); len(fields) != 0 {
			return fields
		}
	}
	return []string{Editor}
}

// isVimLike reports whether editor understands vim's -O option
func isVimLike(editor string) bool {
	switch filepath.Base(editor) {
	case "vim", "nvim", "gvim", "mvim", "view":
		return true
	}
	return false
}

// EditArgs returns the command and arguments to edit files. vim-like
// editors open the files side by side.
func EditArgs(files ...string) (string, []string) {
	command := EditorCommand()
	args := append([]string{}, command[1:]...)
	if isVimLike(command[0]) && len(files) > 1 {
		args = append(args, "-O")
	}
	return command[0], append(args, files...)
}

// Edit opens the specified files in the user's preferred editor
func Edit(files ...string) error {
	editor, args := EditArgs(files...)
	cmd := exec.Command(Which(editor), args...) // #nosec G204 - the user's own editor
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor %s: %w", editor, err)
	}
	return nil
}
