package git

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/jvzantvoort/tmux-project/utils"
)

// GitCmd object for git
type GitCmd struct {
	path    []string
	cwd     string
	command string
}

func (g GitCmd) which(command string) (string, error) {
	for _, dirname := range g.path {
		fpath := path.Join(dirname, command)
		if _, err := os.Stat(fpath); os.IsNotExist(err) {
			continue
		} else {
			return fpath, nil
		}
	}
	if runtime.GOOS == "windows" {
		return "git.exe", nil
	}
	return command, errors.New("unable to find command " + command)
}

func (g GitCmd) exec(args ...string) ([]string, error) {
	retv := []string{}
	cmnd := []string{}

	// cmnd = append(cmnd, g.command)
	cmnd = append(cmnd, args...)
	log.Debugf("command: %s %s", g.command, strings.Join(cmnd, " "))

	cmd := exec.Command(g.command, cmnd...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Dir = g.cwd
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		retv = append(retv, m)
	}
	eerror := cmd.Wait()
	if eerror != nil {
		log.Errorf(strings.Join(cmnd, " "))
		log.Errorf("  command failed, %v", eerror)
		log.Errorf("  cwd: %s", g.cwd)
	}
	return retv, eerror
}

// URL function returning the git url
func (g GitCmd) URL() (string, error) {
	retv, err := g.exec("config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}
	if len(retv) == 0 {
		return "", err
	}
	return string(retv[0]), err
}

func (g GitCmd) GetStatus() map[string]int {
	retv := make(map[string]int)

	lines, err := g.exec("status", "-s")

	if err != nil {
		log.Errorf("Error in GetStatus: %v", err)
		return retv
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		fcol := fields[0]
		var fvar int
		if val, ok := retv[fcol]; ok {
			fvar = val
		} else {
			fvar = 0
		}
		fvar++
		retv[fcol] = fvar
		log.Debugf("%s\n", fields[0])
	}

	return retv
}

// Branch function returning the current git branch
func (g GitCmd) Branch() (string, error) {
	retv, err := g.exec("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	if len(retv) == 0 {
		return "", err
	}
	return string(retv[0]), err
}

// Root function returning the git root
func (g GitCmd) Root() (string, error) {
	retv, err := g.exec("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	if len(retv) == 0 {
		return "", err
	}
	return string(retv[0]), err
}

func (g GitCmd) IsGit() bool {

	gitdir := path.Join(g.cwd, ".git")
	target_stat, err := os.Stat(g.cwd)
	// fmt.Printf("%v %v\n", target_stat, err)
	if err == nil {
		if !target_stat.IsDir() {
			return false
		}
		if err != nil {
			return false
		}
	} else {
		log.Errorf("target: %s, error: %v", g.cwd, err)
	}

	gitdir_stat, nerr := os.Stat(gitdir)
	if nerr != nil {
		return false
	}

	if !gitdir_stat.IsDir() {
		return false
	}

	return true
}

// NewGitCmd create a new git object
func NewGitCmd(dir string) *GitCmd {
	retv := &GitCmd{}
	path := strings.Split(os.Getenv("PATH"), ":")
	homedir.Dir()
	for _, dirn := range path {
		dpath, err := homedir.Expand(dirn)
		if err == nil {
			retv.path = append(retv.path, dpath)
		}
	}

	retv.cwd = dir
	var err error

	retv.command, _ = retv.which("git")
	utils.ErrorExit(err)

	return retv
}
