package main

import (
	"fmt"

	"github.com/jvzantvoort/tmux-project/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Display the version number, commit hash, and build time of tmux-project",
	Run: func(cmd *cobra.Command, args []string) {
		info := version.GetVersion()
		fmt.Println(info.String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
