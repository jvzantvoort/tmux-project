package main

import (
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ProjectCreateCmd represents the create command
var ProjectCreateCmd = &cobra.Command{
	Use:   messages.GetUse("project/create"),
	Short: messages.GetShort("project/create"),
	Long:  messages.GetLong("project/create"),
	Run:   handleProjectCreateCmd,
}

// handleProjectCreateCmd handles the project create command
func handleProjectCreateCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) != 1 {
		log.Error("No project provided")
		cobra.CheckErr(cmd.Help())
		os.Exit(1)
	}
	ProjectName := args[0]
	ProjectType := GetString(*cmd, "type")
	project_description := GetString(*cmd, "description")

	proj := project.NewProject(ProjectName)
	proj.SetDescription(project_description)
	err := proj.InitializeProject(ProjectType, true)
	if err != nil {
		utils.Fatalf("Encountered error: %q", err)
	} else {
		utils.Infof("InitializeProject completed")

	}
}

func init() {
	ProjectCmd.AddCommand(ProjectCreateCmd)
	ProjectCreateCmd.Flags().StringP("type", "t", "default", "Type of project")
	ProjectCreateCmd.Flags().StringP("description", "d", "", "Description of the project")
}
