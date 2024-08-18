package main

import (
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/tmux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ProjectResumeCmd represents the resume command
var ProjectResumeCmd = &cobra.Command{
	Use:   messages.GetUse("project/resume"),
	Short: messages.GetShort("project/resume"),
	Long:  messages.GetLong("project/resume"),
	Run:   handleProjectResumeCmd,
}

func handleProjectResumeCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) != 1 {
		log.Error("No project provided")
		cobra.CheckErr(cmd.Help())
		os.Exit(1)
	}
	sessionname := args[0]

	if sessionname == "ls" {
		project.PrintFullList()
		return
	}

	proj := project.NewProject(sessionname)
	configfile := filepath.Join(config.SessionDir(), proj.ProjectName+".rc")

	tmux.Resume(sessionname, configfile)

}

func init() {
	rootCmd.AddCommand(ProjectResumeCmd)
	ProjectCmd.AddCommand(ProjectResumeCmd)
}
