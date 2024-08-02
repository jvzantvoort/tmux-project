package project

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
)

func CheckPattern(projectname, patstr string) bool {
	pattern := regexp.MustCompile(patstr)
	if pattern.MatchString(projectname) {
		log.Debugf("project name matches pattern")
		return true
	} else {
		log.Warningf("project name %s does not matches pattern %s", projectname, patstr)
		return false
	}
}

// NewProjectConfig derives from ProjectTypeConfig and returns an updated
// object with translated values.
func NewProjectConfig(ptname, projectname string) projecttype.ProjectTypeConfig {

	retv := projecttype.NewProjectTypeConfig(ptname)

	projtypeconfigdir := path.Join(config.ConfigDir(), ptname)
	projtmplvars := NewProjTmplVars(projectname, retv)

	retv.Workdir = projtmplvars.Parse(retv.Workdir)

	// Fail if directory already exists
	if _, err := os.Stat(retv.Workdir); !os.IsNotExist(err) {
		utils.ErrorExit(fmt.Errorf("%s already exists", retv.Workdir))
	}

	CheckPattern(projectname, retv.Pattern)

	var err error
	for indx, cfgfile := range retv.Files {
		// Translate source names
		name := cfgfile.Name
		retv.Files[indx].Name, err = filepath.Abs(path.Join(projtypeconfigdir, name))
		if err != nil {
			log.Fatal(err)
		}

		// Translate destination names
		dest := cfgfile.Destination
		dest = projtmplvars.Parse(dest)
		retv.Files[indx].Destination, err = filepath.Abs(path.Join(config.SessionDir(), dest))
		if err != nil {
			log.Fatal(err)
		}

	}

	for indx, action := range retv.SetupActions {
		retv.SetupActions[indx] = projtmplvars.Parse(action)
	}

	return retv
}

// CreateProject create a new project
func CreateProject(projecttype, projectname string) error {
	log.Debug("CreateProject: start")
	projconf := NewProjectConfig(projecttype, projectname)
	projconf.Describe()

	tmplvars := NewProjTmplVars(projectname, projconf)
	tmplvars.ProjectDescription = utils.Ask("Description")
	log.Debugf("CreateProject: description \"%s\"", tmplvars.ProjectDescription)

	// Write the projconf files
	for _, target := range projconf.Files {
		srccontent, _ := tmplvars.LoadFile(target.Name)
		file, _ := os.Create(target.Destination)
		var err error
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

	if err := utils.MkdirAll(projconf.Workdir); err != nil {
		return fmt.Errorf("directory cannot be created: %s", projconf.Workdir)
	}

	queue := utils.NewQueue()
	for _, step := range projconf.SetupActions {
		queue.Add(projconf.Workdir, step)
	}
	queue.Run()

	log.Debug("CreateProject: end")
	return nil
}
