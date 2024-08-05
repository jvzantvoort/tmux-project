package project

import (
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
	fmt_str := "  %-32s %s"
	args := make(map[string]string)

	args["ProjectDescription"] = proj.ProjectDescription
	args["ProjectDir"] = proj.ProjectDir
	args["ProjectName"] = proj.ProjectName
	args["ProjectType"] = proj.ProjectType
	args["HomeDir"] = proj.HomeDir
	args["ProjectTypeDir"] = proj.ProjectTypeDir
	args["Pattern"] = proj.Pattern
	args["GOARCH"] = proj.GOARCH
	args["GOOS"] = proj.GOOS
	args["GOPATH"] = proj.GOPATH
	args["USER"] = proj.USER

	for keyn, keyv := range args {
		if len(keyv) == 0 {
			continue
		}
		utils.Debugf(fmt_str, keyn, keyv)
	}
	for keyn, keyv := range args {
		if len(keyv) == 0 {
			utils.Debugf(fmt_str, keyn, "EMPTY!!!")
		}
	}
	if len(proj.SetupActions) == 0 {
		utils.Debugf(fmt_str, "SetupActions", "EMPTY!!!")
	}
	if len(proj.Targets) == 0 {
		utils.Debugf(fmt_str, "Targets", "EMPTY!!!")
	} else {
		for _, tgt := range proj.Targets {
			tgt.Confess()
		}

	}

}

func (projt ProjectTarget) Confess() {
	utils.Debugf("  %-32s %s", "Name", projt.Name)
	utils.Debugf("      %-28s %s", "Name", projt.Destination)
	utils.Debugf("      %-28s %s", "Name", projt.Mode)
}
