package main

import (
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/utils"
	"github.com/spf13/cobra"
)

func ProjectRemove(projectname string, _ bool) error {
	utils.LogStart()
	defer utils.LogEnd()

	proj := project.NewProject(projectname)
	cobra.CheckErr(proj.RefreshStruct())
	// err := proj.RemoveProject()
	utils.Infof("RemoveProject completed")
	return nil
}
