package main

import (
	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/spf13/cobra"
)

// ProjectCmd represents the type command
var ProjectCmd = &cobra.Command{
	Use:   messages.GetUse("project/root"),
	Short: messages.GetShort("project/root"),
	Long:  messages.GetLong("project/root"),
}

func init() {
	rootCmd.AddCommand(ProjectCmd)
}
