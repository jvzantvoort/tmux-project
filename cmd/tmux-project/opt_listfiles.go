/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ListfileCmd represents the list command
var ListfileCmd = &cobra.Command{
	Use:   messages.GetUse("listfiles"),
	Short: messages.GetShort("listfiles"),
	Long:  messages.GetLong("listfiles"),
	Run:   handleListfileCmd,
}

func handleListfileCmd(cmd *cobra.Command, args []string) {
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
	for _, ink := range session.TargetPaths() {
		fmt.Printf("%s\n", ink)
	}

}

func init() {
	rootCmd.AddCommand(ListfileCmd)
}
