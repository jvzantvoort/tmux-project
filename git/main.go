// Package git provides utilities for interacting with git repositories,
// including cloning, checking out branches, and managing repository state.
package git

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jvzantvoort/tmux-project/utils"
)

// GitCmd object for git
type GitCmd struct {
	path    []string
	cwd     string
	command string
}

// which finds the path to the provided command
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

// exec runs a git command, and returns the output
func (g GitCmd) exec(args ...string) ([]string, error) {
	retv := []string{}
	cmnd := []string{}

	// cmnd = append(cmnd, g.command)
	cmnd = append(cmnd, args...)
	utils.Debugf("command: %s %s", g.command, strings.Join(cmnd, " "))

	cmd := exec.Command(g.command, cmnd...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Dir = g.cwd
	err := cmd.Start()
	if err != nil {
		utils.Errorf(strings.Join(cmnd, " "))
		utils.Errorf("  command failed to start, %v", err)
		utils.Errorf("  cwd: %s", g.cwd)
		return retv, err
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		retv = append(retv, m)
	}
	eerror := cmd.Wait()
	if eerror != nil {
		utils.Errorf(strings.Join(cmnd, " "))
		utils.Errorf("  command failed, %v", eerror)
		utils.Errorf("  cwd: %s", g.cwd)
	}
	return retv, eerror
}

// GetStatus function returns the git status as a map of files and their status
func (g GitCmd) GetStatus() map[string]int {
	retv := make(map[string]int)

	lines, err := g.exec("status", "-s")

	if err != nil {
		utils.Errorf("Error in GetStatus: %v", err)
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
		utils.Debugf("%s\n", fields[0])
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

func (g *GitCmd) Clone(url, destination string) error {

	// make the path absolute
	destination, err := filepath.Abs(destination)
	if err != nil {
		return err
	}

	// when using nested directory make sure they exist
	directory := filepath.Dir(destination)
	err = utils.MkdirAll(directory)
	if err != nil {
		return err
	}

	_, err = g.exec("clone", url, destination)
	if err == nil {
		g.cwd = destination
	}
	return err

}

func (g GitCmd) Checkout(branch string) error {
	curBranch, err := g.Branch()
	if err != nil {
		return err
	}
	if curBranch == branch {
		return nil
	}

	_, err = g.exec("checkout", branch)
	return err
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
		utils.Errorf("target: %s, error: %v", g.cwd, err)
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
	for _, dirn := range path {
		dpath, err := utils.Expand(dirn)
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
