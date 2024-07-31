package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"os"
	"os/user"
	"text/template"

	"github.com/jvzantvoort/tmux-project/projecttype"
)

type Project struct {
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

func NewProject(projectname string, projtype projecttype.ProjectTypeConfig) *Project {

	retv := &Project{}

	retv.ProjectName = projectname

	retv.Init(projtype)

	return retv
}

func (proj *Project) Init(projtype projecttype.ProjectTypeConfig) {
	proj.HomeDir, _ = os.UserHomeDir()

	proj.ProjectDir = projtype.Workdir
	proj.ProjectType = projtype.ProjectType

	buildContext := build.Default

	proj.GOARCH = buildContext.GOARCH
	proj.GOOS = buildContext.GOOS
	proj.GOPATH = buildContext.GOPATH

	if currentUser, err := user.Current(); err == nil {
		proj.USER = currentUser.Username
	}

}

// buildConfig construct the text from the template definition and arguments.
func (t Project) ParseTemplateString(templatestring string) (string, error) {
	var retv string

	tmpl, err := template.New("prompt").Parse(templatestring)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, t)
	if err != nil {
		return "", err
	}
	retv = buf.String()
	return retv, nil
}

func (t Project) ParseTemplateFile(target string) (string, error) {
	var retv string
	var err error

	content, err := os.ReadFile(target)
	if err != nil {
		return retv, err
	}

	return t.ParseTemplateString(string(content))
}

// Write json output to an [io.Writer] compatible handle. It returns nil or the
// error of [json.MarshalIndent]
func (proj Project) Write(writer io.Writer) error {
	content, err := json.MarshalIndent(proj, "", "  ")
	if err == nil {
		fmt.Fprint(writer, string(content))
		fmt.Fprintf(writer, "\n")
	}
	return err
}

// Read session content from a [io.Reader] object.
func (proj *Project) Read(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &proj)
	if err != nil {
		return err
	}
	return nil
}

// Write session configuration to a projectfile
func (proj Project) WriteToFile(projectfile string) error {
	filehandle, err := os.OpenFile(projectfile, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer filehandle.Close()
	return proj.Write(filehandle)
}

// vim: noexpandtab filetype=go
