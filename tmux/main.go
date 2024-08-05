package tmux

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jvzantvoort/tmux-project/utils"
)

func New(name, configfile string) {
	command := []string{"-f", configfile, "new", "-s", name}

	cmd := exec.Command(utils.Which("tmux"), command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	utils.ErrorExit(cmd.Run())
}

func Attach(name, configfile string) {
	command := []string{"-f", configfile, "attach", "-d", "-t", name}

	cmd := exec.Command(utils.Which("tmux"), command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	utils.ErrorExit(cmd.Run())
}

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

func ListActive() ([]string, error) {
	command := fmt.Sprintf("%s ls -F #{session_name}", utils.Which("tmux"))
	cwd, _ := os.UserHomeDir()
	stdout_lines, _, err := utils.Exec(cwd, command)

	if err != nil {
		return stdout_lines, err
	}
	return stdout_lines, nil
}
