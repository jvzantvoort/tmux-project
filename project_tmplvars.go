package tmuxproject

import (
	"bytes"
	"os"
	"text/template"

	"github.com/jvzantvoort/tmux-project/projecttype"
)

type ProjTmplVars struct {
	HomeDir            string
	ProjectDescription string
	ProjectDir         string
	ProjectName        string
}

func NewProjTmplVars(projectname string, conf projecttype.ProjectTypeConfig) *ProjTmplVars {

	v := &ProjTmplVars{}
	v.HomeDir = mainconfig.HomeDir
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
