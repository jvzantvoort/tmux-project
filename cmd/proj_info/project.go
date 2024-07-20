package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/jvzantvoort/tmux-project/git"
	"github.com/jvzantvoort/tmux-project/utils"
)

type Project struct {
	Name     string
	AbsPath  string
	Path     string
	Branch   string
	Expected bool
	Status   map[string]int
	SubPath  string
	Chapter  string
	Info     os.FileInfo
}

func (p Project) PrintLine() {
	br_col := color.New(BranchChangedColor)
	if utils.StringInSlice(p.Branch, []string{"master", "main", "develop"}) {
		br_col = color.New(BranchDefaultColor)
	}
	var stat_str string
	for status, amount := range p.Status {
		stat_str = fmt.Sprintf("%s %s:%d", stat_str, status, amount)
		stat_str = strings.TrimSpace(stat_str)
	}
	if len(stat_str) != 0 {
		stat_str = fmt.Sprintf(" [%s]", stat_str)
	}

	br_str := br_col.Sprint(p.Branch)

	fmt.Printf("   %-32s %s%s\n", p.Name, br_str, stat_str)
	// fmt.Printf("%q\n", p)
}

func NewProject(projdir, dirname string) *Project {
	retv := &Project{}

	retv.Path = dirname
	retv.Info, _ = os.Lstat(dirname)
	if projdir == dirname {
		retv.Name = "."
		retv.Chapter = "root"
	} else {
		retv.Name = dirname[len(projdir)+1:]
		retv.Chapter = "rest"
	}

	gitcmnd := git.NewGitCmd(retv.Path)
	if gitcmnd.IsGit() {
		log.Debugf("%s is a gitdir", retv.Path)
		retv.Branch, _ = gitcmnd.Branch()
		retv.Status = gitcmnd.GetStatus()
	}
	return retv
}

func findAllProjects(projdir string) []Project {
	var retv []Project

	filepath.Walk(projdir, func(file string, fi os.FileInfo, inerr error) error {
		err := inerr
		if err != nil {
			log.Errorf("this passed an error: %q", err)
		}
		if fi.IsDir() && fi.Name() == ".git" {
			dirname := filepath.Dir(file)
			repos := NewProject(projdir, dirname)
			retv = append(retv, *repos)
		}
		return nil
	})
	return retv

}
