package project

import (
	"go/build"
	"os"
	"os/user"
	"regexp"
	"strings"

	errno "github.com/jvzantvoort/tmux-project/errors"
	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/utils"
)

// NewProject initialize a Project object
func NewProject(projectname string) *Project {

	utils.LogStart()
	defer utils.LogEnd()

	retv := &Project{}

	retv.ProjectName = projectname

	return retv
}

func (proj Project) NameIsValid() bool {

	utils.LogStart()
	defer utils.LogEnd()

	pattern := regexp.MustCompile(proj.Pattern)

	if pattern.MatchString(proj.ProjectName) {
		utils.Debugf("project name matches pattern")
		return true
	} else {
		utils.Warningf("project name %s does not matches pattern %s", proj.ProjectName, proj.Pattern)
		return false
	}

}

func (proj *Project) InjectExternal() {
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

func (proj *Project) RefreshStruct(args ...string) error {
	utils.LogStart()
	defer utils.LogEnd()

	proj.Confess()

	proj.Exists = true

	// try to load the configfile
	err := proj.Open()
	if err == nil {
		utils.Debugf("read configfile")
	} else {

		// cannot find or open the project
		proj.Exists = false

		// The error was *not* that the file does not exist
		if errno.IsProjectNotExist(err) {
			if len(args) != 1 {
				return errno.ErrProjectTypeNotDefined
			}
			utils.Errorf("project not found")
			proj.InjectProjectType(args[0])
		} else {
			utils.Errorf("failed to open project file: %s", err)
			return err
		}
	}

	proj.InjectExternal()

	proj.ProjectDir = proj.Parse(proj.ProjectDir)
	return nil
}

func (proj *Project) SetDescription(instr ...string) {
	proj.ProjectDescription = strings.Join(instr, " ")
	proj.ProjectDescription = strings.TrimSpace(proj.ProjectDescription)
	if len(proj.ProjectDescription) == 0 {
		proj.ProjectDescription = utils.Ask("Description")
	}
}

func (proj *Project) InitializeProject(projtype string, safe bool) error {

	utils.LogStart()
	defer utils.LogEnd()

	err := utils.SetupSessionDir()
	if err != nil {
		utils.Errorf("Error: %s", err)

	}

	err = proj.RefreshStruct(projtype)
	if err != nil {
		utils.Errorf("Error: %s", err)

	}

	if safe && !proj.Exists {
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
		err = proj.ProcessProjectTarget(&target)
		if err != nil {
			utils.Errorf("Error in target: %s", err)

		}
	}

	queue := utils.NewQueue()
	for _, step := range proj.SetupActions {
		step = proj.Parse(step)
		queue.Add(proj.ProjectDir, step)
	}

	queue.Run()

	return proj.Save()
}
