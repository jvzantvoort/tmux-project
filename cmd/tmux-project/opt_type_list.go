package main

import (
	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// TypeListCmd represents the type command
var TypeListCmd = &cobra.Command{
	Use:   messages.GetUse("type/list"),
	Short: messages.GetShort("type/list"),
	Long:  messages.GetLong("type/list"),
	Run:   handleTypeListCmd,
}

func handleTypeListCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)
	err := projecttype.ListProjectTypeConfigs()
	utils.LogIfError(err)
}

func init() {
	TypeCmd.AddCommand(TypeListCmd)
}
