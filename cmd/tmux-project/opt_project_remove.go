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
	swAssumeYes, _ := cmd.Flags().GetBool("yes")

	err := ProjectRemove(ProjectName, swArchive, swAssumeYes)
	cobra.CheckErr(err)

}

func init() {
	ProjectCmd.AddCommand(ProjectRemoveCmd)
	ProjectRemoveCmd.Flags().BoolP("archive", "a", false, "Archive the project before removal")
	ProjectRemoveCmd.Flags().BoolP("yes", "y", false, "Do not ask for confirmation")

	// --noarchive never had an effect; not archiving is the default now
	ProjectRemoveCmd.Flags().BoolP("noarchive", "x", true, "Do not archive before delete")
	cobra.CheckErr(ProjectRemoveCmd.Flags().MarkDeprecated("noarchive", "not archiving is the default, use --archive to archive"))
}
