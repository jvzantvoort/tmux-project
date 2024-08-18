package main

import (
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   messages.GetUse("list"),
	Short: messages.GetShort("list"),
	Long:  messages.GetLong("list"),
	Run:   handleListCmd,
}

func handleListCmd(cmd *cobra.Command, args []string) {
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
	rootCmd.AddCommand(ListCmd)
	ListCmd.Flags().BoolP("full", "f", false, "Print full")
}
