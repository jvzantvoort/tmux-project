package project

import (
	"bytes"
	"os"
	"text/template"

	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
)

// Template actions
// --------------------------------------------------------------------------

// buildConfig construct the text from the template definition and arguments.
func (proj Project) ParseTemplateString(templatestring string) (string, error) {
	functionname := utils.FunctionName()
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)
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
	functionname := utils.FunctionName()
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)
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
	functionname := utils.FunctionName()
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)
	tmpl, err := template.New("prompt").Parse(templatestring)
	utils.ErrorExit(err)
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, proj)
	utils.ErrorExit(err)
	return buf.String()
}

func (proj Project) LoadFile(target string) (string, error) {
	functionname := utils.FunctionName()
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)
	var retv string
	content, err := os.ReadFile(target)
	if err != nil {
		return "", err
	}
	retv = proj.Parse(string(content))
	return retv, nil
}
