/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// InitProjCmdCmd represents the create command
var InitProjCmdCmd = &cobra.Command{
	Use:   messages.GetUse("init"),
	Short: messages.GetShort("init"),
	Long:  messages.GetLong("init"),
	Run:   handleInitProjCmdCmd,
}

func handleInitProjCmdCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) != 1 {
		log.Error("No project type provided")
		cmd.Help()
		os.Exit(1)
	}
	ProjectType := args[0]

	Force, _ := cmd.Flags().GetBool("force")

	if ProjectType == "default" {
		if !Force {
			log.Fatalf("Cannot overwrite default")
		}
	}
	ptc := projecttype.NewProjectTypeConfig(ProjectType)
	err := ptc.SetupProjectTypeConfig()
	if err != nil {
		utils.Abort("Error: %s", err)
	}

}

func init() {
	rootCmd.AddCommand(InitProjCmdCmd)
	InitProjCmdCmd.Flags().BoolP("force", "f", false, "Force")
}
