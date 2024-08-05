/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// EditCmd represents the edit command
var EditCmd = &cobra.Command{
	Use:   messages.GetUse("edit"),
	Short: messages.GetShort("edit"),
	Long:  messages.GetLong("edit"),
	Run:   handleEditCmd,
}

func handleEditCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) != 1 {
		log.Error("No project provided")
		cmd.Help()
		os.Exit(1)
	}
	ProjectName := args[0]

	proj := project.NewProject(ProjectName)
	opts := []string{"-O"}
	opts = append(opts, filepath.Join(config.SessionDir(), proj.ProjectName+".rc"))
	opts = append(opts, filepath.Join(config.SessionDir(), proj.ProjectName+".env"))

	utils.Edit(opts...)

}

func init() {
	rootCmd.AddCommand(EditCmd)
}
