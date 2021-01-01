package tmuxproject

import (
	"fmt"
	"os"
	"os/exec"
)

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
