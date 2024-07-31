package utils

import (
	"bufio"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Exec(cwd, args string) ([]string, []string, error) {
	commandlist := []string{}
	stdout_list := []string{}
	stderr_list := []string{}
	cmndargs := strings.Split(args, " ")
	cmnd := cmndargs[0]
	cmndargs = cmndargs[1:]
	commandlist = append(commandlist, cmndargs...)
	command := exec.Command(cmnd, commandlist...)
	log.Debugf("command: %s %s", cmnd, strings.Join(commandlist, " "))
	log.Debugf("         cwd: %s", cwd)
	command.Dir = cwd

	stdout, err := command.StdoutPipe()
	ErrorExit(err)

	stderr, err := command.StderrPipe()
	ErrorExit(err)

	err = command.Start()
	ErrorExit(err)

	stdout_scan := bufio.NewScanner(stdout)
	stdout_scan.Split(bufio.ScanLines)
	for stdout_scan.Scan() {
		msg := stdout_scan.Text()
		log.Debugln(msg)
		stdout_list = append(stdout_list, msg)
	}

	stderr_scan := bufio.NewScanner(stderr)
	stderr_scan.Split(bufio.ScanLines)
	for stderr_scan.Scan() {
		msg := stderr_scan.Text()
		log.Errorln(msg)
		stderr_list = append(stderr_list, msg)
	}

	eerror := command.Wait()
	if eerror != nil {
		log.Errorf("command failed, %v", err)
	}

	return stdout_list, stderr_list, eerror
}
