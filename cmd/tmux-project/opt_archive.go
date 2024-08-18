package main

import (
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ProjectArchiveCmd represents the archive command
var ProjectArchiveCmd = &cobra.Command{
	Use:   messages.GetUse("project/archive"),
	Short: messages.GetShort("project/archive"),
	Long:  messages.GetLong("project/archive"),
	Run:   handleProjectArchiveCmd,
}

func handleProjectArchiveCmd(cmd *cobra.Command, args []string) {
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

	ArchiveName := GetString(*cmd, "archivename")

	err := CreateArchive(ProjectName, ArchiveName)
	cobra.CheckErr(err)

}

func init() {
	ProjectCmd.AddCommand(ProjectArchiveCmd)
	ProjectArchiveCmd.Flags().StringP("archivename", "a", "", "Archive file")
}
