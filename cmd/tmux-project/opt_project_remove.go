package main

import (
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
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
	log.Debugf("%s: start", cmd.Use)

	if len(args) != 1 {
		log.Error("No project provided")
		cobra.CheckErr(cmd.Help())
		log.Debugf("%s: end", cmd.Use)
		os.Exit(1)
	}
	defer log.Debugf("%s: end", cmd.Use)
	ProjectName := args[0]

	swArchive, _ := cmd.Flags().GetBool("archive")
	swInteractive, _ := cmd.Flags().GetBool("yes")

	if swArchive {
		err := CreateArchive(ProjectName, "")
		cobra.CheckErr(err)
	}
	err := ProjectRemove(ProjectName, swInteractive)
	cobra.CheckErr(err)

}

func init() {
	ProjectCmd.AddCommand(ProjectRemoveCmd)
	ProjectRemoveCmd.Flags().BoolP("noarchive", "x", true, "Do not archive before delete")
	ProjectRemoveCmd.Flags().BoolP("yes", "y", false, "Assume yes")
}
