package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// Edit opens the specified files in the user's preferred editor
func Edit(args ...string) {
	editor := Which(Editor)
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
