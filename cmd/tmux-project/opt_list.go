/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	tp "github.com/jvzantvoort/tmux-project"
	"github.com/jvzantvoort/tmux-project/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List a project",
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
		tp.PrintFullList()
	} else {
		tp.PrintShortList()
	}

}

func init() {
	rootCmd.AddCommand(ListCmd)
	ListCmd.Flags().BoolP("full", "f", false, "Print full")
}
