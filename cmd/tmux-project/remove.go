package main

import (
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/utils"
	"github.com/spf13/cobra"
)

func ProjectRemove(projectname string, interactive bool) error {
	utils.LogStart()
	defer utils.LogEnd()

	proj := project.NewProject(projectname)
	cobra.CheckErr(proj.RefreshStruct())
	var err error = nil
	// err := proj.RemoveProject()
	if err != nil {
		utils.Fatalf("Encountered error: %q", err)
	} else {
		utils.Infof("RemoveProject completed")

	}
	return nil
}
