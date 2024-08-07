package project

import (
	"bytes"
	"os"
	"text/template"

	"github.com/jvzantvoort/tmux-project/utils"
)

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
