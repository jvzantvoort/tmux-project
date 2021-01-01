package tmuxproject

import (
	"bytes"
	"text/template"

	"github.com/jvzantvoort/tmux-project/config"
)

type ProjTmplVars struct {
	HomeDir            string
	ProjectDescription string
	ProjectDir         string
	ProjectName        string
}

func NewProjTmplVars(projectname string, conf config.ProjectTypeConfig) *ProjTmplVars {

	v := &ProjTmplVars{}
	v.HomeDir = GetHomeDir()
	v.ProjectDir = conf.Workdir
	v.ProjectName = projectname

	return v
}

// buildConfig construct the text from the template definition and arguments.
func (t ProjTmplVars) Parse(templatestring string) string {
	tmpl, err := template.New("prompt").Parse(templatestring)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, t)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// vim: noexpandtab filetype=go
