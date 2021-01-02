package tmuxproject

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

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
		stdout_list, stderr_list, eerror := Exec(configuration.Workdir, action)
		for _, stdout_line := range stdout_list {
			log.Infof("<stdout> %s", stdout_line)
		}
		for _, stderr_line := range stderr_list {
			log.Infof("<stderr> %s", stderr_line)
		}
		if eerror != nil {
			log.Errorf("Failed: %s", eerror)
		}
	}
	log.Debug("CreateProject: end")
	return nil
}
