package project

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/sessions"
	"github.com/jvzantvoort/tmux-project/tmux"
	"github.com/jvzantvoort/tmux-project/utils"
	"github.com/olekukonko/tablewriter"
)

func ListActive() []string {
	tmux := tmux.NewTmux()
	retv, _ := tmux.ListActive()
	return retv
}

func PrintFullList() {
	active := ListActive()
	sessiondata, err := sessions.ListTmuxConfigs(config.SessionDir())
	utils.ErrorExit(err)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Active", "Description", "Workdir", "Sane"})

	for _, entries := range sessiondata {
		if entries[0] == "default" {
			continue
		}

		// create cols
		cols := []string{entries[0], ""}
		cols = append(cols, entries[1:]...)

		if utils.StringInSlice(entries[0], active) {
			cols[1] = "yes"
		}

		table.Append(cols)
	}

	table.SetHeaderLine(true)
	table.SetBorder(false)
	table.Render()
}

func PrintShortList() {
	sessiondata, err := sessions.ListTmuxConfigs(config.SessionDir())
	utils.ErrorExit(err)

	for _, item := range sessiondata {
		fmt.Printf("%s\n", item[0])

	}
}
