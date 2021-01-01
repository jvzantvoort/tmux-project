package tmuxproject

import (
	"fmt"
	"os"
	"strconv"
)

// CreateProject create a new project
func CreateProject(projecttype, projectname string) error {
	configuration := GetProjectTypeConfig(projecttype, projectname)
	DescribeProjectType(configuration)

	tmplvars := NewProjTmplVars(projectname, configuration)
	tmplvars.ProjectDescription = Ask("Description")
	for _, target := range configuration.Files {
		srccontent, _ := LoadFile(target.Name, *tmplvars)
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
