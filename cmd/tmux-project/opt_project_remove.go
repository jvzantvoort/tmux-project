package main

import (
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ProjectRemoveCmd represents the remove command
var ProjectRemoveCmd = &cobra.Command{
	Use:   messages.GetUse("project/remove"),
	Short: messages.GetShort("project/remove"),
	Long:  messages.GetLong("project/remove"),
	Run:   handleProjectRemoveCmd,
}

func handleProjectRemoveCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	utils.LogStart()
	defer utils.LogEnd()

	if len(args) != 1 {
		log.Error("No project provided")
		cobra.CheckErr(cmd.Help())
		os.Exit(1)
	}
	ProjectName := args[0]

	swArchive, _ := cmd.Flags().GetBool("archive")
	swInteractive, _ := cmd.Flags().GetBool("yes")

	if swArchive {
		err := CreateArchive(ProjectName, "")
		cobra.CheckErr(err)
	}
	err := RemoveArchive(ProjectName, swInteractive)
	cobra.CheckErr(err)

}

func init() {
	ProjectCmd.AddCommand(ProjectRemoveCmd)
	ProjectRemoveCmd.Flags().BoolP("noarchive", "x", true, "Do not archive before delete")
	ProjectRemoveCmd.Flags().BoolP("yes", "y", false, "Assume yes")
}
