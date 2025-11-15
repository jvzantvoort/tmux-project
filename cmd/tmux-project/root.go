// Package main provides the command-line interface for tmux-project.
// It uses the Cobra library to implement subcommands for managing tmux-based projects.
package main

import (
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	verbose       bool
	OutputDirOpt  string
	OutputDir     string
	OutputFileOpt string
	OutputFile    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     messages.GetUse("root"),
	Short:   messages.GetShort("root"),
	Long:    messages.GetLong("root"),
	Version: version.GetVersion().Short(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)

}
