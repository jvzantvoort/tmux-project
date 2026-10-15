package main

import (
	"errors"
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// TypeInitCmd represents the type command
var TypeInitCmd = &cobra.Command{
	Use:   messages.GetUse("type/init"),
	Short: messages.GetShort("type/init"),
	Long:  messages.GetLong("type/init"),
	Run:   handleTypeInitCmd,
}

func handleTypeInitCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)

	if len(args) != 1 {
		log.Error("No project type provided")
		cobra.CheckErr(cmd.Help())
		log.Debugf("%s: end", cmd.Use)
		os.Exit(1)
	}
	defer log.Debugf("%s: end", cmd.Use)
	ProjectType := args[0]

	Force, _ := cmd.Flags().GetBool("force")

	ptc, err := projecttype.New(ProjectType)
	if err != nil && !errors.Is(err, projecttype.ErrProjectNotExists) {
		utils.Abort("Error: %s", err)
	}
	// init writes config.yml and the templates, replacing any changes
	if err == nil && !Force {
		log.Fatalf("project type %s already exists, use --force to overwrite its configuration", ProjectType)
	}
	err = ptc.Setup()
	if err != nil {
		utils.Abort("Error: %s", err)
	}

}

func init() {
	TypeCmd.AddCommand(TypeInitCmd)
	TypeInitCmd.Flags().BoolP("force", "f", false, "Overwrite an existing project type")
}
