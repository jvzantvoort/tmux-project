package project

import (
	"go/build"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
)

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
	proj.ProjectDescription = strings.Join(instr, " ")
	proj.ProjectDescription = strings.TrimSpace(proj.ProjectDescription)
	if len(proj.ProjectDescription) == 0 {
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

	return proj.Save()
}
