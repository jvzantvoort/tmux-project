package project

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/jvzantvoort/tmux-project/utils"
)

//go:embed scripts/*

var Content embed.FS

func GetScriptContent(name string) string {
	filename := fmt.Sprintf("scripts/%s", name)

	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Errorf("%s", err)
		msgstr = []byte("undefined")
	}
	return strings.TrimSuffix(string(msgstr), "\n")

}

// buildConfig construct the text from the template definition and arguments.
func (proj Project) Parse(templatestring string) string {

	buf := new(bytes.Buffer)

	utils.LogStart()
	defer utils.LogEnd()

	tmpl, err := template.New("project").Parse(templatestring)
	utils.ErrorExit(err)

	err = tmpl.Execute(buf, proj)
	utils.ErrorExit(err)

	return buf.String()
}

func (proj Project) LoadFile(target string) (string, error) {

	var retv string

	utils.LogStart()
	defer utils.LogEnd()

	content, err := os.ReadFile(target)
	if err != nil {
		return "", err
	}

	retv = proj.Parse(string(content))

	return retv, nil
}
