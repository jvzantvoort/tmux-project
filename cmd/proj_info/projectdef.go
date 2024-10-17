package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jvzantvoort/tmux-project/git"
	"github.com/jvzantvoort/tmux-project/utils"
)

type ProjectDef struct {
	Path       string
	ProjectDir string
	Name       string
	AbsPath    string
	Branch     string
	Expected   bool
	Status     map[string]int
	SubPath    string
	Chapter    string
	Info       os.FileInfo
}

func (p ProjectDef) GetFields() []string {
	retv := []string{}
	retv = append(retv, p.Name)

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
		retv = append(retv, fmt.Sprintf("[%s]", stat_str))
	} else {
		retv = append(retv, " ")
	}

	retv = append(retv, br_col.Sprint(p.Branch))

	return retv
}

func (p *ProjectDef) GetGitInfo() {
	gitcmnd := git.NewGitCmd(p.Path)
	if gitcmnd.IsGit() {
		p.Branch, _ = gitcmnd.Branch()
		p.Status = gitcmnd.GetStatus()
	}
}

func NewProjectDef(projdir, dirname string) *ProjectDef {
	retv := &ProjectDef{}

	retv.Path = dirname
	retv.ProjectDir = projdir

	return retv
}

func (pd *ProjectDef) Init() {
	utils.LogStart()
	defer utils.LogEnd()

	pd.GetGitInfo()

	pd.Info, _ = os.Lstat(pd.Path)
	if pd.ProjectDir == pd.Path {
		pd.Name = "."
		pd.Chapter = "root"
	} else {
		pd.Name = pd.Path[len(pd.ProjectDir)+1:]
		pd.Chapter = "rest"
	}

}
