package main

import (
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ArchiveCmd represents the archive command
var ArchiveCmd = &cobra.Command{
	Use:   messages.GetUse("archive"),
	Short: messages.GetShort("archive"),
	Long:  messages.GetLong("archive"),
	Run:   handleArchiveCmd,
}

func handleArchiveCmd(cmd *cobra.Command, args []string) {
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
	rootCmd.AddCommand(ArchiveCmd)
	ArchiveCmd.Flags().StringP("archivename", "a", "", "Archive file")
}
