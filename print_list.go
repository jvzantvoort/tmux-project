package tmuxproject

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/tmux"
	"github.com/olekukonko/tablewriter"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func PrintFullList() {
	tmux := tmux.NewTmux()
	active, _ := tmux.ListActive()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Active", "Description", "Workdir"})

	for _, target := range ListTmuxConfigs() {
		var cols []string
		cols = append(cols, target.Name)

		if stringInSlice(target.Name, active) {
			cols = append(cols, "yes")
		} else {
			cols = append(cols, "")
		}
		cols = append(cols, target.Description)
		cols = append(cols, target.Workdir)
		table.Append(cols)
	}
	table.SetHeaderLine(true)
	table.SetBorder(false)
	table.Render()
}

func PrintShortList() {
	for _, target := range ListTmuxConfigs() {
		fmt.Println(target.Name)
	}
}
