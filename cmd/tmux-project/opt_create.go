/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"os"

	tp "github.com/jvzantvoort/tmux-project"
	"github.com/jvzantvoort/tmux-project/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create <projectname>",
	Short: "Create a project",
	Long:  messages.GetLong("create"),
	Run:   handleCreateCmd,
}

func handleCreateCmd(cmd *cobra.Command, args []string) {
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
	ProjectType := GetString(*cmd, "type")

	err := tp.CreateProject(ProjectType, ProjectName)
	if err != nil {
		log.Fatalf("Encountered error: %q", err)
	}

}

func init() {
	rootCmd.AddCommand(CreateCmd)
	CreateCmd.Flags().StringP("type", "t", "default", "Type of project")
}
