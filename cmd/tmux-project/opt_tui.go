package main

import (
	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   messages.GetUse("tui"),
	Short: messages.GetShort("tui"),
	Long:  messages.GetLong("tui"),
	Run:   runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func runTUI(cmd *cobra.Command, args []string) {
	RunTUIApp()
}
