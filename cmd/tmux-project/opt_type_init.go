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
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) != 1 {
		log.Error("No project type provided")
		cobra.CheckErr(cmd.Help())
		os.Exit(1)
	}
	ProjectType := args[0]

	Force, _ := cmd.Flags().GetBool("force")

	if ProjectType == "default" {
		if !Force {
			log.Fatalf("Cannot overwrite default")
		}
	}
	ptc, err := projecttype.New(ProjectType)
	if err != nil && !errors.Is(err, projecttype.ErrProjectNotExists) {
		utils.Abort("Error: %s", err)
	}
	err = ptc.Setup()
	if err != nil {
		utils.Abort("Error: %s", err)
	}

}

func init() {
	TypeCmd.AddCommand(TypeInitCmd)
	TypeInitCmd.Flags().BoolP("force", "f", false, "Force")
}
