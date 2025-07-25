package project

import (
	"go/build"
	"os"
	"os/user"
	"regexp"
	"strings"

	errno "github.com/jvzantvoort/tmux-project/errors"
	"github.com/jvzantvoort/tmux-project/git"
	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/utils"
)

// NewProject initialize a Project object
func NewProject(projectname string) *Project {

	utils.LogStart()
	defer utils.LogEnd()
	utils.LogArgument("projectname", projectname)

	retv := &Project{}

	retv.Name = projectname

	return retv
}

func (proj Project) NameIsValid() bool {

	utils.LogStart()
	defer utils.LogEnd()

	pattern := regexp.MustCompile(proj.Pattern)

	if pattern.MatchString(proj.Name) {
		utils.Debugf("project name matches pattern")
		return true
	} else {
		utils.Warningf("project name %s does not matches pattern %s", proj.Name, proj.Pattern)
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
	proj.Directory = ptobj.Directory
	utils.LogVariable("proj.Directory", proj.Directory)

	proj.ProjectType = ptobj.ProjectType
	utils.LogVariable("proj.ProjectType", proj.ProjectType)

	proj.Pattern = ptobj.Pattern
	utils.LogVariable("proj.Pattern", proj.Pattern)

	proj.ProjectTypeDir = ptobj.ProjectTypeDir
	utils.LogVariable("proj.ProjectTypeDir", proj.ProjectTypeDir)

	proj.Root = ptobj.Root
	utils.LogVariable("proj.Root", proj.Root)

	proj.SetupActions = ptobj.SetupActions
	proj.Repos = []Repo{}
	for _, inob := range ptobj.Repos {
		obj := &Repo{}
		obj.Url = inob.Url
		obj.Destination = inob.Destination
		obj.Branch = inob.Branch
		proj.Repos = append(proj.Repos, *obj)

	}

	for _, element := range ptobj.Targets {
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

// RefreshStruct refreshes the project structure by loading the project configuration file.
// It checks if the project exists and if not, it injects the project type configuration.
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

	if proj.Root == "" {
		proj.Root = "directory" // Default to "directory" if root is not specified
		utils.Debugf("no root defined, using default: %s", proj.Root)
	}

	proj.InjectExternal()

	// Translate some stuff
	proj.Directory = proj.Parse(proj.Directory)
	utils.LogVariable("proj.Directory", proj.Directory)
	return nil
}

func (proj *Project) SetDescription(instr ...string) {
	proj.Description = strings.Join(instr, " ")
	proj.Description = strings.TrimSpace(proj.Description)
	if len(proj.Description) == 0 {
		proj.Description = utils.Ask("Description")
	}
}

func (proj *Project) InitializeProject(projtype string, safe bool) error {

	utils.LogStart()
	defer utils.LogEnd()
	utils.LogArgument("projtype", projtype)
	utils.LogArgument("safe", safe)

	// SetupSessionDir setup the ~/.tmux.d directory
	err := utils.SetupSessionDir(false)
	if err != nil {
		utils.Errorf("Error: %s", err)

	}

	err = proj.RefreshStruct(projtype)
	if err != nil {
		utils.Errorf("Error: %s", err)
	}

	utils.LogVariable("proj.Root", proj.Root)

	if safe && !proj.Exists {
		if !proj.NameIsValid() {
			utils.Abort("Name %s is invalid", proj.Name)
		}
	}

	if proj.Root != "directory" {
		// if root is not "directory" the project is a git repository
		utils.Debugf("project root is not 'directory', cloning from %s", proj.Root)

		// get the parent directory of the project
		ngit := git.NewGitCmd(proj.HomeDir)
		if nerr := ngit.Clone(proj.Root, proj.Directory); nerr != nil {
			return nerr
		}
	} else {

		// Fail if directory already exists
		if _, err := os.Stat(proj.Directory); !os.IsNotExist(err) {
			utils.Abort("%s already exists", proj.Directory)
		}

		if err := utils.MkdirAll(proj.Directory); err != nil {
			utils.Abort("directory cannot be created: %s", proj.Directory)
		}
	}

	// Write the proj files
	for _, target := range proj.Targets {
		err = proj.ProcessTarget(&target)
		if err != nil {
			utils.Errorf("Error in target: %s", err)

		}
	}

	gqueue := git.NewQueue()
	for _, repo := range proj.Repos {
		gqueue.Add(repo.Url, proj.Directory, repo.Destination, repo.Branch)
	}
	gqueue.Run()

	queue := utils.NewQueue()
	for _, step := range proj.SetupActions {
		step = proj.Parse(step)
		queue.Add(proj.Directory, step)
	}

	queue.Run()

	return proj.Save()
}
