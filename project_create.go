package tmuxproject

import (
	"fmt"
	"os"
	"strconv"
)

// CreateProject create a new project
func CreateProject(projecttype, projectname string) error {
	configuration := NewProjectTypeConfig(projecttype, projectname)
	configuration.Describe()

	tmplvars := NewProjTmplVars(projectname, configuration)
	tmplvars.ProjectDescription = Ask("Description")

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
	return nil
}
