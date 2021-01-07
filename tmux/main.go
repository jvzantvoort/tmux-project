package tmux

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Tmux struct {
	CommandPath string
	CommandCwd  string
}

func Which(command string) string {
	Path := strings.Split(os.Getenv("PATH"), ":")
	var retv string
	for _, dirname := range Path {
		fullpath := path.Join(dirname, command)
		log.Debugf("test path: %s", fullpath)

		_, err := os.Stat(fullpath)
		if err == nil {
			if os.IsNotExist(err) {
				continue
			}
			retv = fullpath
		}
	}
	return retv

}

func (t Tmux) Exec(args ...string) ([]string, []string, error) {
	cmndargs := []string{}
	stdout_list := []string{}
	stderr_list := []string{}

	cmndargs = append(cmndargs, args...)
	message := fmt.Sprintf("%s %s", t.CommandPath, strings.Join(cmndargs, " "))
	log.Debugf("Exec[%s], start", message)

	cmnd := exec.Command(t.CommandPath, cmndargs...)
	cmnd.Dir = t.CommandCwd

	// Setup stdout pipe
	stdout, err := cmnd.StdoutPipe()
	if err != nil {
		log.Errorf("stdout pipe failed, %v", err)
		log.Fatal(err)
		panic(err)
	}

	// Setup stderr pipe
	stderr, err := cmnd.StderrPipe()
	if err != nil {
		log.Errorf("stderr pipe failed, %v", err)
		log.Fatal(err)
		panic(err)
	}

	// Start the command
	cmnd.Start()

	// readout stdout lines
	stdout_scan := bufio.NewScanner(stdout)
	stdout_scan.Split(bufio.ScanLines)
	for stdout_scan.Scan() {
		msg := stdout_scan.Text()
		log.Debugln(msg)
		stdout_list = append(stdout_list, msg)
	}

	// readout stderr lines
	stderr_scan := bufio.NewScanner(stderr)
	stderr_scan.Split(bufio.ScanLines)
	for stderr_scan.Scan() {
		msg := stderr_scan.Text()
		log.Errorln(msg)
		stderr_list = append(stderr_list, msg)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		log.Printf(">>    %s", m)
	}

	eerror := cmnd.Wait()
	if eerror != nil {
		log.Errorf("command failed, %v", err)
	}

	log.Debugf("Exec[%s], end", message)
	return stdout_list, stderr_list, eerror

}

func (t Tmux) ListActive() ([]string, error) {
	var stdout_lines []string
	//var stderr_lines []string
	var err error
	stdout_lines, _, err = t.Exec("ls", "-F", "#{session_name}")
	if err != nil {
		log.Errorf("command failed, %v", err)
		return stdout_lines, err
	}
	return stdout_lines, nil
}

func (t Tmux) SessionExists(sessionname string) bool {
	sessions, err := t.ListActive()
	if err != nil {
		return false
	}
	for _, sess := range sessions {
		if sess == sessionname {
			return true
		}
	}
	return false
}

func NewTmux() *Tmux {

	t := &Tmux{}
	t.CommandPath = Which("tmux")

	// Get user info
	usrobj, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	t.CommandCwd = usrobj.HomeDir

	return t
}
