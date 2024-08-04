/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"fmt"

	"github.com/jvzantvoort/tmux-project/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ShellCmd represents the shell command
var ShellCmd = &cobra.Command{
	Use:   messages.GetUse("shell"),
	Short: "Shell output",
	Long:  messages.GetLong("shell"),
	Run:   handleShellCmd,
}

func handleShellCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	ShellName := "bash"

	if len(args) == 1 {
		ShellName = args[0]
	}
	fmt.Print(messages.GetShell(ShellName))
}

func init() {
	rootCmd.AddCommand(ShellCmd)
}
