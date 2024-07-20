// Package tmux provides interface
//
//	import (
//	    "fmt"
//	    "github.com/jvzantvoort/tmux-project/tmux"
//	)
//	tmux := tmux.NewTmux()
package tmux

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jvzantvoort/tmux-project/sessions"
	"github.com/jvzantvoort/tmux-project/utils"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

type Tmux struct {
	CommandPath string
	CommandCwd  string
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
		utils.ErrorExit(err)
	}

	// Setup stderr pipe
	stderr, err := cmnd.StderrPipe()
	if err != nil {
		log.Errorf("stderr pipe failed, %v", err)
		log.Fatal(err)
		utils.ErrorExit(err)
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

func (t Tmux) CreateSession(sess sessions.TmuxSession) {
	//  SESSION=$1;
	//  CONFIGFILE="${CONST_CONFDIR}/${SESSION}.rc"
	//  [[ -f "${CONFIGFILE}" ]] || CONFIGFILE="${CONST_CONFDIR}/default.rc"
	//  TERM="${CONST_TERM}" tmux -f $CONFIGFILE new -s $SESSION
	command := []string{}
	command = append(command, "-f")
	command = append(command, sess.Configfile)
	command = append(command, "new")
	command = append(command, "-s")
	command = append(command, sess.Name)
	log.Debugf("%v\n", command)
	cmd := exec.Command(t.CommandPath, command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func (t Tmux) ResumeSession(sess sessions.TmuxSession) {
	//  SESSION=$1;
	//  CONFIGFILE="${CONST_CONFDIR}/${SESSION}.rc"
	//  [[ -f "${CONFIGFILE}" ]] || CONFIGFILE="${CONST_CONFDIR}/default.rc"
	//  TERM="${CONST_TERM}" tmux -f $CONFIGFILE new -s $SESSION
	command := []string{}
	command = append(command, "-f")
	command = append(command, sess.Configfile)
	command = append(command, "attach")
	command = append(command, "-d")
	command = append(command, "-t")
	command = append(command, sess.Name)

	log.Debugf("%v\n", command)

	cmd := exec.Command(t.CommandPath, command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func NewTmux() *Tmux {

	t := &Tmux{}
	t.CommandPath = utils.Which("tmux")

	t.CommandCwd, _ = homedir.Dir()

	return t
}
