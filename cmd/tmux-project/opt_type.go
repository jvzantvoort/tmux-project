package main

import (
	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/spf13/cobra"
)

// TypeCmd represents the type command
var TypeCmd = &cobra.Command{
	Use:   messages.GetUse("type/root"),
	Short: messages.GetShort("type/root"),
	Long:  messages.GetLong("type/root"),
}

func init() {
	rootCmd.AddCommand(TypeCmd)
}
