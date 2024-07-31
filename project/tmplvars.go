package project

import (
	"bytes"
	"os"
	"text/template"

	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/utils"
)

type ProjTmplVars struct {
	HomeDir            string `json:"homedir"`
	ProjectDescription string `json:"description"`
	ProjectDir         string `json:"directory"`
	ProjectName        string `json:"name"`
	ProjectType        string `json:"type"`
	GOARCH             string `json:"-"` // target architecture
	GOOS               string `json:"-"` // target operating system
	GOPATH             string `json:"-"` // Go paths
	USER               string `json:"-"`
}

func NewProjTmplVars(projectname string, conf projecttype.ProjectTypeConfig) *ProjTmplVars {

	v := &ProjTmplVars{}
	v.HomeDir, _ = os.UserHomeDir()
	v.ProjectDir = conf.Workdir
	v.ProjectName = projectname
	v.ProjectType = conf.ProjectType

	return v
}

// buildConfig construct the text from the template definition and arguments.
func (t ProjTmplVars) Parse(templatestring string) string {
	tmpl, err := template.New("prompt").Parse(templatestring)
	utils.ErrorExit(err)
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, t)
	utils.ErrorExit(err)
	return buf.String()
}

func (t ProjTmplVars) LoadFile(target string) (string, error) {
	var retv string
	content, err := os.ReadFile(target)
	if err != nil {
		return "", err
	}
	retv = t.Parse(string(content))
	return retv, nil
}

// vim: noexpandtab filetype=go
