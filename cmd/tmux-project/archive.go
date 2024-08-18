package main

import (
	"fmt"

	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
)

func CreateArchive(projectname, archivename string) error {
	utils.LogStart()
	defer utils.LogEnd()

	obj := project.NewProject(projectname)
	obj.RefreshStruct()

	if len(archivename) == 0 {
		if len(obj.ProjectDir) == 0 {
			return fmt.Errorf("projectdir is empty")
		}
		archivename = obj.ProjectDir + ".tar.gz"
	}

	log.Debugf("Outputfile: %s", archivename)

	err := obj.Archive(archivename)
	if err != nil {
		return err
	}

	fmt.Printf("Created %s\n", archivename)

	return nil
}
