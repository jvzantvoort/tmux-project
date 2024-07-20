package tmuxproject

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jvzantvoort/tmux-project/sessions"

	log "github.com/sirupsen/logrus"
)

type ListTable struct {
	Name        string
	Description string
	Workdir     string
	Sane        bool
}

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
	if err != nil {
		log.Errorf("stdout pipe failed, %v", err)
		log.Fatal(err)
		panic(err)
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		log.Errorf("stderr pipe failed, %v", err)
		log.Fatal(err)
		panic(err)
	}

	command.Start()

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

// targetExists return true if target exists
func targetExists(targetpath string) bool {
	_, err := os.Stat(targetpath)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// which return path
func which(command string) string {
	Path := strings.Split(os.Getenv("PATH"), ":")
	var retv string
	for _, dirname := range Path {
		fullpath := path.Join(dirname, command)
		if targetExists(fullpath) {
			retv = fullpath
			break
		}
	}
	return retv
}

// ExitOnError check error and exit if not nil
func ExitOnError(err error) {
	if err != nil {
		log.Errorf("error %v\n", err)
		os.Exit(1)
	}
}

func Ask(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", question)
	text, _ := reader.ReadString('\n')
	return strings.TrimSuffix(text, "\n")
}

func ListTmuxConfigs() []ListTable {
	var retv []ListTable
	targets, err := ioutil.ReadDir(mainconfig.TmuxDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, target := range targets {
		target_name := target.Name()

		// we only want the session names
		if !strings.HasSuffix(target_name, ".rc") {
			continue
		}

		// "common" is shared by all others
		if target_name == "common.rc" {
			continue
		}

		target_name = strings.TrimSuffix(target_name, ".rc")

		session := sessions.NewTmuxSession(target_name)

		t := ListTable{}
		t.Name = session.Name
		t.Description = session.Description
		t.Workdir = session.Workdir
		t.Sane = session.IsSane()
		retv = append(retv, t)
	}
	return retv
}

// ExpandHome expand the tilde in a given path.
func ExpandHome(pathstr string) (string, error) {
	if len(pathstr) == 0 {
		return pathstr, nil
	}

	if pathstr[0] != '~' {
		return pathstr, nil
	}

	return filepath.Join(mainconfig.HomeDir, pathstr[1:]), nil

}

func LoadStringLines(target string) ([]string, error) {
	var retv []string
	content, err := ioutil.ReadFile(target)
	if err != nil {
		return retv, err
	}

	for _, line := range strings.Split(string(content), "\n") {
		retv = append(retv, line)
	}
	return retv, nil
}

// GetMatches parses each line and sticks the findings in a map
func GetMatches(regEx string, lines []string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)

	paramsMap = make(map[string]string)

	for _, line := range lines {
		match := compRegEx.FindStringSubmatch(line)
		for i, name := range compRegEx.SubexpNames() {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}
		}
	}
	return
}

func Edit(args ...string) {
	editor := which(Editor)
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
