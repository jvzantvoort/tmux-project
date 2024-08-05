/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/project"
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

	project := project.NewProject(ProjectName)
	project.RefreshStruct()
	project.Confess()

	if ArchiveName == "" {
		if len(project.ProjectDir) != 0 {
			ArchiveName = project.ProjectDir + ".tar.gz"
		}
	}
	if ArchiveName == "" {
		cobra.CheckErr(fmt.Errorf("no archive name provided"))
	}

	log.Debugf("Outputfile: %s", ArchiveName)

	err := project.Archive(ArchiveName)
	if err == nil {
		fmt.Printf("Created %s\n", ArchiveName)
	} else {
		log.Fatalf("Encountered error: %q", err)
	}
	// name: cmd.Use
}

func init() {
	rootCmd.AddCommand(ArchiveCmd)
	ArchiveCmd.Flags().StringP("archivename", "a", "", "Archive file")
}
