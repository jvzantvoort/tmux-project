package tmuxproject

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func PrintFullList() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Workdir"})
	for _, target := range ListTmuxConfigs() {
		var cols []string
		cols = append(cols, target.Name)
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
