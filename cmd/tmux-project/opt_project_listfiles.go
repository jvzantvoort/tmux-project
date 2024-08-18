package main

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/project"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ProjectListfileCmd represents the list command
var ProjectListfileCmd = &cobra.Command{
	Use:   messages.GetUse("project/listfiles"),
	Short: messages.GetShort("project/listfiles"),
	Long:  messages.GetLong("project/listfiles"),
	Run:   handleProjectListfileCmd,
}

func handleProjectListfileCmd(cmd *cobra.Command, args []string) {
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
	ProjectName := args[0]

	project := project.NewProject(ProjectName)
	cobra.CheckErr(project.RefreshStruct())

	for _, ink := range project.ListFiles() {
		fmt.Printf("%s\n", ink)
	}

}

func init() {
	ProjectCmd.AddCommand(ProjectListfileCmd)
}
