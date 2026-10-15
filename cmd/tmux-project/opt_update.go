package main

import (
	"github.com/jvzantvoort/tmux-project/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	forceUpdate bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update tmux-project from GitHub",
	Long:  "Download and install the latest version of tmux-project from GitHub",
	RunE:  runUpdate,
	// errors are about the update, not about how the command was used
	SilenceUsage: true,
}

func init() {
	updateCmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "Force update even if already up to date")
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	return update.Execute(forceUpdate)
}
