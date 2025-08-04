package projecttype

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/utils"
	"github.com/olekukonko/tablewriter"
)

func ListProjectTypeConfigs() error {

	ptconfigs := []ProjectTypeConfig{}
	inputdir := config.ConfigDir()

	if err := utils.MkdirAll(inputdir); err != nil {
		return err
	}

	targets, err := os.ReadDir(inputdir)
	if err != nil {
		utils.Fatalf("%s", err)
		return err
	}
	for _, target := range targets {
		target_name := target.Name()
		obj, err := New(target_name)
		utils.LogIfError(err)
		ptconfigs = append(ptconfigs, obj)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Name", "SetupActions", "Targets", "Repos", "ConfigDir", "Pattern"})

	for _, item := range ptconfigs {
		cols := []string{}
		cols = append(cols, item.ProjectType)
		cols = append(cols, fmt.Sprintf("%d", len(item.SetupActions)))
		cols = append(cols, fmt.Sprintf("%d", len(item.Targets))) // FIXME
		cols = append(cols, fmt.Sprintf("%d", len(item.Repos)))
		cols = append(cols, item.ConfigDir)
		cols = append(cols, item.Pattern)
		table.Append(cols)
	}
	table.Render()

	return nil

}
