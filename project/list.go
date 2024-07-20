package project

import (
	"fmt"
	"os"
	"strings"

	"github.com/jvzantvoort/tmux-project/sessions"
	"github.com/jvzantvoort/tmux-project/utils"
	"github.com/jvzantvoort/tmux-project/tmux"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

type ListTable struct {
	Name        string
	Description string
	Workdir     string
	Sane        bool
}

func PrintFullList() {
	tmux := tmux.NewTmux()
	active, _ := tmux.ListActive()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Active", "Description", "Workdir", "Sane"})

	for _, target := range ListTmuxConfigs() {
		if target.Name == "default" {
			continue
		}
		var cols []string
		cols = append(cols, target.Name)

		if utils.StringInSlice(target.Name, active) {
			cols = append(cols, "yes")
		} else {
			cols = append(cols, "")
		}
		cols = append(cols, target.Description)
		cols = append(cols, target.Workdir)
		if target.Sane {
			cols = append(cols, "true")
		} else {
			cols = append(cols, "false")

		}
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

func ListTmuxConfigs() []ListTable {
	var retv []ListTable
	targets, err := os.ReadDir(mainconfig.TmuxDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, target := range targets {
		target_name := target.Name()

		// we only want the session names
		if !strings.HasSuffix(target_name, ".rc") {
			continue
		}

		// "common" is shared by all others
		if target_name == "common.rc" {
			continue
		}

		target_name = strings.TrimSuffix(target_name, ".rc")

		session := sessions.NewTmuxSession(target_name)

		t := ListTable{}
		t.Name = session.Name
		t.Description = session.Description
		t.Workdir = session.Workdir
		t.Sane = session.IsSane()
		retv = append(retv, t)
	}
	return retv
}
