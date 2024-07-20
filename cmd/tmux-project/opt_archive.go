/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ArchiveCmd represents the archive command
var ArchiveCmd = &cobra.Command{
	Use:   "archive <projectname>",
	Short: "Archive a project",
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

	session := sessions.NewTmuxSession(ProjectName)

	ArchiveName := GetString(*cmd, "archivename")

	if ArchiveName == "" {
		ArchiveName = session.Workdir + ".tar.gz"
	}

	err := session.Archive(ArchiveName)
	if err != nil {
		log.Fatalf("Encountered error: %q", err)
	}
	// name: cmd.Use
}

func init() {
	rootCmd.AddCommand(ArchiveCmd)
	ArchiveCmd.Flags().StringP("archivename", "a", "", "Archive file")
	// ArchiveCmd.Flags().StringP("projectname", "n", "", "Name of project")
}
