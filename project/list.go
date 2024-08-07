package project

import (
	"fmt"
	"os"
	"strings"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/tmux"
	"github.com/jvzantvoort/tmux-project/utils"
	"github.com/olekukonko/tablewriter"
)

func ListConfigs() []string {
	retv := []string{}
	suffix := ".json"

	inputdir := config.SessionDir()
	targets, err := os.ReadDir(inputdir)
	if err != nil {
		utils.Fatalf("%s", err)
		return retv
	}
	for _, target := range targets {
		target_name := target.Name()

		// we only want the session names
		if strings.HasSuffix(target_name, suffix) {
			retv = append(retv, strings.TrimSuffix(target_name, suffix))
		}
	}
	return retv
}

// PrintFullList prints the list of sessions
func PrintFullList() {
	var err error
	active, _ := tmux.ListActive()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Active", "Description", "Workdir", "Sane"})

	for _, sessionname := range ListConfigs() {
		sessp := NewProject(sessionname)
		err = sessp.Open()
		if err != nil {
			fmt.Printf("%#v\n", err)
			continue
		}
		cols := []string{}
		cols = append(cols, sessionname)
		if utils.StringInSlice(sessionname, active) {
			cols = append(cols, "yes")
		} else {
			cols = append(cols, "")
		}
		cols = append(cols, sessp.ProjectDescription)
		cols = append(cols, sessp.ProjectDir)
		sane := "true"
		for _, target := range sessp.ListFiles() {
			if !utils.TargetExists(target) {
				sane = "false"
			}
		}
		cols = append(cols, sane)
		table.Append(cols)

	}

	table.SetHeaderLine(true)
	table.SetBorder(false)
	table.Render()

}

// PrintShortList prints the list of sessions
func PrintShortList() {
	for _, item := range ListConfigs() {
		fmt.Printf("%s\n", item)
	}
}
