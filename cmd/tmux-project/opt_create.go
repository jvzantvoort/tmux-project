/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   messages.GetUse("create"),
	Short: messages.GetShort("create"),
	Long:  messages.GetLong("create"),
	Run:   handleCreateCmd,
}

func handleCreateCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	utils.LogStart()
	defer utils.LogEnd()

	if len(args) != 1 {
		log.Error("No project provided")
		cmd.Help()
		os.Exit(1)
	}
	ProjectName := args[0]
	ProjectType := GetString(*cmd, "type")
	ProjectDescription := GetString(*cmd, "description")

	proj := project.NewProject(ProjectName)
	proj.SetDescription(ProjectDescription)
	err := proj.InitializeProject(ProjectType, true)
	if err != nil {
		utils.Fatalf("Encountered error: %q", err)
	} else {
		utils.Infof("InitializeProject completed")

	}
}

func init() {
	rootCmd.AddCommand(CreateCmd)
	CreateCmd.Flags().StringP("type", "t", "default", "Type of project")
	CreateCmd.Flags().StringP("description", "d", "", "Description of the project")
}
