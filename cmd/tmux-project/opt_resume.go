/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/sessions"
	"github.com/jvzantvoort/tmux-project/tmux"
	"github.com/jvzantvoort/tmux-project/utils"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ResumeCmd represents the resume command
var ResumeCmd = &cobra.Command{
	Use:   messages.GetUse("resume"),
	Short: messages.GetShort("resume"),
	Long:  messages.GetLong("resume"),
	Run:   handleResumeCmd,
}

func handleResumeCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) != 1 {
		log.Error("No project provided")
		cmd.Help()
		os.Exit(1)
	}
	ProjectName := args[0]

	if len(ProjectName) == 0 {
		prompt := promptui.Select{
			Label: "Select project",
			Size:  20,
			Items: ListSessions(),
		}
		_, result, err := prompt.Run()

		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		result = strings.Split(result, " ")[0]
		ProjectName = result
	}
	sess := sessions.NewTmuxSessions()
	xsess, err := sess.Find(ProjectName)
	utils.ErrorExit(err)
	tmux.Resume(ProjectName, xsess.Configfile)
}

func ListSessions() []string {
	retv := []string{}
	active, _ := tmux.ListActive()

	sess := sessions.NewTmuxSessions()
	for _, sesi := range sess.Sessions {
		state := " "
		if utils.StringInSlice(sesi.Name, active) {
			state = "active"
		}
		message := fmt.Sprintf("%-32s %-6s %s", sesi.Name, state, sesi.Description)

		retv = append(retv, message)
	}

	return retv
}

func init() {
	rootCmd.AddCommand(ResumeCmd)
}
