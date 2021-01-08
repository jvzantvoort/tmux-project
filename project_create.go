package tmuxproject

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
	pt "github.com/jvzantvoort/tmux-project/projecttype"
)

var (
	wg sync.WaitGroup
)

func cleanup() {
	if r := recover(); r != nil {
		log.Errorf("Paniced %s", r)
	}
}

func RunSetupAction(workdir, action string) {
	defer wg.Done() // lower counter
	defer cleanup() // handle panics

	stdout_list, stderr_list, eerror := Exec(workdir, action)
	for _, stdout_line := range stdout_list {
		log.Infof("<stdout> %s", stdout_line)
	}
	for _, stderr_line := range stderr_list {
		log.Errorf("<stderr> %s", stderr_line)
	}
	if eerror != nil {
		panic(fmt.Sprintf("Action \"%s\" failed", action))
	}
}

// NewProjectConfig derives from ProjectTypeConfig and returns an updated
// object with translated values.
func NewProjectConfig(projecttype, projectname string) pt.ProjectTypeConfig {

	ptc := pt.NewProjectTypeConfig(projecttype)

	projtypeconfigdir := path.Join(mainconfig.ProjTypeConfigDir, projecttype)
	projtmplvars := NewProjTmplVars(projectname, ptc)

	ptc.Workdir = projtmplvars.Parse(ptc.Workdir)
	pattern := regexp.MustCompile(ptc.Pattern)
	if pattern.MatchString(projectname) {
		log.Debugf("project name matches pattern")
	} else {
		log.Warningf("project name %s does not matches pattern %s", projectname, ptc.Pattern)
	}

	var err error
	for indx, cfgfile := range ptc.Files {
		// Translate source names
		name := cfgfile.Name
		ptc.Files[indx].Name, err = filepath.Abs(path.Join(projtypeconfigdir, name))
		if err != nil {
			log.Fatal(err)
		}

		// Translate destination names
		dest := cfgfile.Destination
		dest = projtmplvars.Parse(dest)
		ptc.Files[indx].Destination, err = filepath.Abs(path.Join(mainconfig.TmuxDir, dest))
		if err != nil {
			log.Fatal(err)
		}

	}

	for indx, action := range ptc.SetupActions {
		ptc.SetupActions[indx] = projtmplvars.Parse(action)
	}

	return ptc
}

// CreateProject create a new project
func CreateProject(projecttype, projectname string) error {
	log.Debug("CreateProject: start")
	configuration := NewProjectConfig(projecttype, projectname)
	configuration.Describe()

	tmplvars := NewProjTmplVars(projectname, configuration)
	tmplvars.ProjectDescription = Ask("Description")
	log.Debugf("CreateProject: description \"%s\"", tmplvars.ProjectDescription)

	// Write the configuration files
	for _, target := range configuration.Files {
		srccontent, _ := tmplvars.LoadFile(target.Name)
		file, err := os.Create(target.Destination)
		_, err = file.WriteString(srccontent)
		if err != nil {
			return err
		}
		defer file.Close()

		num, err := strconv.Atoi(target.Mode)
		if err != nil {
			return err
		}

		mode, _ := strconv.ParseUint(fmt.Sprintf("%04d", num), 8, 32)
		if err := os.Chmod(target.Destination, os.FileMode(mode)); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(configuration.Workdir, os.FileMode(int(0755))); err != nil {
		return fmt.Errorf("Directory cannot be created: %s", configuration.Workdir)
	}

	for _, action := range configuration.SetupActions {
		wg.Add(1)
		go RunSetupAction(configuration.Workdir, action)
	}
	wg.Wait()
	log.Debug("CreateProject: end")
	return nil
}
