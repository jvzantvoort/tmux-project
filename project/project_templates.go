package project

import (
	"bytes"
	"os"
	"text/template"

	"github.com/jvzantvoort/tmux-project/utils"
)

// Template actions
// --------------------------------------------------------------------------

// buildConfig construct the text from the template definition and arguments.
func (proj Project) ParseTemplateString(templatestring string) (string, error) {
	utils.LogStart()
	defer utils.LogEnd()
	var retv string

	tmpl, err := template.New("prompt").Parse(templatestring)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, proj)
	if err != nil {
		return "", err
	}
	retv = buf.String()
	return retv, nil
}

func (proj Project) ParseTemplateFile(target string) (string, error) {
	utils.LogStart()
	defer utils.LogEnd()
	var retv string
	var err error

	content, err := os.ReadFile(target)
	if err != nil {
		return retv, err
	}

	return proj.ParseTemplateString(string(content))
}

// buildConfig construct the text from the template definition and arguments.
func (proj Project) Parse(templatestring string) string {
	utils.LogStart()
	defer utils.LogEnd()
	tmpl, err := template.New("prompt").Parse(templatestring)
	utils.ErrorExit(err)
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, proj)
	utils.ErrorExit(err)
	return buf.String()
}

func (proj Project) LoadFile(target string) (string, error) {
	utils.LogStart()
	defer utils.LogEnd()
	var retv string
	content, err := os.ReadFile(target)
	if err != nil {
		return "", err
	}
	retv = proj.Parse(string(content))
	return retv, nil
}
