/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/sessions"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// EditCmd represents the edit command
var EditCmd = &cobra.Command{
	Use:   messages.GetUse("edit"),
	Short: "Edit a project",
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

	session := sessions.NewTmuxSession(ProjectName)

	utils.Edit("-O", session.Configfile, session.Environment)

}

func init() {
	rootCmd.AddCommand(EditCmd)
}
