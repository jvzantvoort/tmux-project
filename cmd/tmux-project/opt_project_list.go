package main

import (
	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ProjectListCmd represents the list command
var ProjectListCmd = &cobra.Command{
	Use:   messages.GetUse("project/list"),
	Short: messages.GetShort("project/list"),
	Long:  messages.GetLong("project/list"),
	Run:   handleProjectListCmd,
}

func handleProjectListCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	PrintFull, _ := cmd.Flags().GetBool("full")

	if PrintFull {
		project.PrintFullList()
	} else {
		project.PrintShortList()
	}

}

func init() {
	// add shortcut alias
	rootCmd.AddCommand(ProjectListCmd)

	ProjectCmd.AddCommand(ProjectListCmd)
	ProjectListCmd.Flags().BoolP("full", "f", false, "Print full")
}
