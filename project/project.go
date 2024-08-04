package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
)

type ProjectTarget struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	Content     string `json:"-"`
}

type Project struct {
	HomeDir            string          `json:"homedir"`
	ProjectDescription string          `json:"description"`
	ProjectDir         string          `json:"directory"` // Workdir for the project
	ProjectName        string          `json:"name"`
	ProjectType        string          `json:"type"`
	SetupActions       []string        `json:"setupactions"`
	Targets            []ProjectTarget `json:"targets"`
	ProjectTypeDir     string          `json:"-"`
	Pattern            string          `json:"-"` // pattern obtained from ProjectType
	GOARCH             string          `json:"-"` // target architecture
	GOOS               string          `json:"-"` // target operating system
	GOPATH             string          `json:"-"` // Go paths
	USER               string          `json:"-"`
}

// NewProject initialize a Project object
func NewProject(projectname string) *Project {
	functionname := utils.FunctionName()
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)

	retv := &Project{}

	retv.ProjectName = projectname

	return retv
}

func (proj Project) NameIsValid() bool {

	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)

	pattern := regexp.MustCompile(proj.Pattern)

	if pattern.MatchString(proj.ProjectName) {
		log.Debugf("project name matches pattern")
		return true
	} else {
		log.Warningf("project name %s does not matches pattern %s", proj.ProjectName, proj.Pattern)
		return false
	}

}

func (proj *Project) InjectProjectType(projtype string) {
	// load project type object info
	ptobj := projecttype.NewProjectTypeConfig(projtype)
	proj.ProjectDir = ptobj.Workdir
	proj.ProjectType = ptobj.ProjectType
	proj.Pattern = ptobj.Pattern
	proj.ProjectTypeDir = ptobj.ProjectTypeDir
	proj.SetupActions = ptobj.SetupActions

	for _, element := range ptobj.Files {
		content, _ := ptobj.Content(element.Name)
		obj := ProjectTarget{
			Name:        element.Name,
			Destination: element.Destination,
			Mode:        element.Mode,
			Content:     content,
		}
		proj.Targets = append(proj.Targets, obj)
	}
}

func (proj *Project) RefreshStruct(projtype string) {
	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)

	proj.InjectProjectType(projtype)

	// load homedir object info
	proj.HomeDir, _ = os.UserHomeDir()

	// load build object info
	buildContext := build.Default
	proj.GOARCH = buildContext.GOARCH
	proj.GOOS = buildContext.GOOS
	proj.GOPATH = buildContext.GOPATH

	// load user info
	if currentUser, err := user.Current(); err == nil {
		proj.USER = currentUser.Username
	}

	proj.ProjectDir = proj.Parse(proj.ProjectDir)
}

func (proj *Project) SetDescription(instr ...string) {
	if len(instr) != 0 {
		proj.ProjectDescription = strings.Join(instr, " ")
	} else {
		proj.ProjectDescription = utils.Ask("Description")
	}
}

func (proj *Project) InitializeProject(projtype string, safe bool) error {
	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)

	proj.RefreshStruct(projtype)

	if safe {
		if !proj.NameIsValid() {
			utils.Abort("Name %s is invalid", proj.ProjectName)
		}
	}

	// Fail if directory already exists
	if _, err := os.Stat(proj.ProjectDir); !os.IsNotExist(err) {
		utils.Abort("%s already exists", proj.ProjectDir)
	}

	if err := utils.MkdirAll(proj.ProjectDir); err != nil {
		utils.Abort("directory cannot be created: %s", proj.ProjectDir)
	}

	// Write the proj files
	for _, target := range proj.Targets {
		proj.ProcessProjectTarget(&target)
	}

	queue := utils.NewQueue()
	for _, step := range proj.SetupActions {
		step = proj.Parse(step)
		queue.Add(proj.ProjectDir, step)
	}

	queue.Run()

	return nil
}

func (proj Project) ProcessProjectTarget(element *ProjectTarget) error {
	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)

	element.Name = proj.Parse(element.Name)
	element.Destination = proj.Parse(element.Destination)

	mode, err := utils.GetMode(element.Mode)
	if err != nil {
		return err
	}

	// source_file := filepath.Join(proj.ProjectTypeDir, element.Name)
	dest_file := filepath.Join(config.SessionDir(), element.Destination)
	content := proj.Parse(element.Content)

	filehandle, _ := os.Create(dest_file)
	_, err = filehandle.WriteString(content)
	if err != nil {
		return err
	}
	defer filehandle.Close()
	if err := os.Chmod(dest_file, mode); err != nil {
		return err
	}

	return nil
}

// Template actions
// --------------------------------------------------------------------------

// buildConfig construct the text from the template definition and arguments.
func (proj Project) ParseTemplateString(templatestring string) (string, error) {
	functionname := utils.FunctionName(2)
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
	functionname := utils.FunctionName(2)
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
	functionname := utils.FunctionName(2)
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
	functionname := utils.FunctionName(2)
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

// Write json output to an [io.Writer] compatible handle. It returns nil or the
// error of [json.MarshalIndent]
func (proj Project) Write(writer io.Writer) error {
	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)
	content, err := json.MarshalIndent(proj, "", "  ")
	if err == nil {
		fmt.Fprint(writer, string(content))
		fmt.Fprintf(writer, "\n")
	}
	return err
}

// Read session content from a [io.Reader] object.
func (proj *Project) Read(reader io.Reader) error {
	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)
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
	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)
	filehandle, err := os.OpenFile(projectfile, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer filehandle.Close()
	return proj.Write(filehandle)
}

// vim: noexpandtab filetype=go
