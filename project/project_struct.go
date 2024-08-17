package project

import (
	"fmt"

	"github.com/jvzantvoort/tmux-project/utils"
)

type ProjectTarget struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	Content     string `json:"-"`
}

type Project struct {
	// Stored variables
	ProjectDescription string          `json:"description"`
	ProjectDir         string          `json:"directory"` // Workdir for the project
	ProjectName        string          `json:"name"`
	ProjectType        string          `json:"type"`
	SetupActions       []string        `json:"setupactions"`
	Targets            []ProjectTarget `json:"targets"`

	// Derived variables
	HomeDir        string `json:"-"`
	ProjectTypeDir string `json:"-"` // Directory where project type files are located
	Pattern        string `json:"-"` // pattern obtained from ProjectType
	GOARCH         string `json:"-"` // target architecture
	GOOS           string `json:"-"` // target operating system
	GOPATH         string `json:"-"` // Go paths
	USER           string `json:"-"` // Username
	Exists         bool   `json:"-"` // project exists
}

func (proj Project) Confess() {
	content := proj.Parse(GetScriptContent("confess"))
	fmt.Print(content)
}

func (projt ProjectTarget) Confess() {
	utils.Debugf("  %-32s %s", "Name", projt.Name)
	utils.Debugf("      %-28s %s", "Name", projt.Destination)
	utils.Debugf("      %-28s %s", "Name", projt.Mode)
}
