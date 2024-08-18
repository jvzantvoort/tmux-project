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
	utils.LogArgument("projectname", projectname)

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
	utils.LogStart()
	defer utils.LogEnd()

	// load homedir object info
	proj.HomeDir, _ = os.UserHomeDir()
	utils.LogVariable("proj.HomeDir", proj.HomeDir)

	// load build object info
	buildContext := build.Default
	proj.GOARCH = buildContext.GOARCH
	utils.LogVariable("proj.GOARCH", proj.GOARCH)

	proj.GOOS = buildContext.GOOS
	utils.LogVariable("proj.GOOS", proj.GOOS)

	proj.GOPATH = buildContext.GOPATH
	utils.LogVariable("proj.GOPATH", proj.GOPATH)

	// load user info
	if currentUser, err := user.Current(); err == nil {
		proj.USER = currentUser.Username
		utils.LogVariable("proj.USER", proj.USER)
	}

}

func (proj *Project) InjectProjectType(projtype string) error {
	utils.LogStart()
	defer utils.LogEnd()

	utils.LogArgument("projtype", projtype)

	// load project type object info
	ptobj, err := projecttype.New(projtype)
	if err != nil {
		return err
	}
	proj.ProjectDir = ptobj.Workdir
	utils.LogVariable("proj.ProjectDir", proj.ProjectDir)

	proj.ProjectType = ptobj.ProjectType
	utils.LogVariable("proj.ProjectType", proj.ProjectType)

	proj.Pattern = ptobj.Pattern
	utils.LogVariable("proj.Pattern", proj.Pattern)

	proj.ProjectTypeDir = ptobj.ProjectTypeDir
	utils.LogVariable("proj.ProjectTypeDir", proj.ProjectTypeDir)

	proj.SetupActions = ptobj.SetupActions

	for _, element := range ptobj.Files {
		content, _ := ptobj.Content(element.Name)
		obj := Target{
			Name:        element.Name,
			Destination: element.Destination,
			Mode:        element.Mode,
			Content:     content,
		}
		proj.Targets = append(proj.Targets, obj)
	}
	return nil
}

func (proj *Project) RefreshStruct(args ...string) error {
	utils.LogStart()
	defer utils.LogEnd()

	var project_type string
	if len(args) == 1 {
		project_type = args[0]
	}

	proj.Exists = true

	// try to load the configfile
	err := proj.Open()
	if err == nil {
		utils.Debugf("succesfully read configfile")
	} else {

		// cannot find or open the project
		proj.Exists = false

		// The error was *not* that the file does not exist
		if errno.IsProjectNotExist(err) {
			if len(project_type) == 0 {
				return errno.ErrProjectTypeNotDefined
			}
			err := proj.InjectProjectType(project_type)
			if err != nil {
				return err
			}
		} else {
			utils.Errorf("failed to open project file: %s", err)
			return err
		}
	}

	proj.InjectExternal()

	// Translate some stuff
	proj.ProjectDir = proj.Parse(proj.ProjectDir)
	utils.LogVariable("proj.ProjectDir", proj.ProjectDir)
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
	utils.LogArgument("projtype", projtype)
	utils.LogArgument("safe", safe)

	err := utils.SetupSessionDir(false)
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
		err = proj.ProcessTarget(&target)
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
