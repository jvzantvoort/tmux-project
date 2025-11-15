package project

import (
	"fmt"
	"os"
	"sort"
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
		// If directory doesn't exist, create it and return empty list
		if os.IsNotExist(err) {
			if err := os.MkdirAll(inputdir, 0755); err != nil {
				utils.Fatalf("failed to create session directory: %s", err)
			}
			return retv
		}
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

func GenerateListOfMaps() []map[string]string {
	var err error
	active, _ := tmux.ListActive()
	data := []map[string]string{}

	for _, sessionname := range ListConfigs() {
		sessp := NewProject(sessionname)
		err = sessp.Open()
		if err != nil {
			utils.Errorf("%#v", err)
			continue
		}

		dict := map[string]string{}
		dict["Name"] = sessionname
		dict["Type"] = sessp.ProjectType
		dict["Description"] = sessp.Description
		dict["Directory"] = sessp.Directory
		dict["Status"] = sessp.Status

		if utils.StringInSlice(sessionname, active) {
			dict["Active"] = "yes"
		} else {
			dict["Active"] = ""
		}

		sane := "true"
		for _, target := range sessp.ListFiles() {
			if !utils.TargetExists(target) {
				sane = "false"
			}
		}
		dict["Sane"] = sane

		data = append(data, dict)
	}
	return data

}

func indexOf(target string, list []string) int {
	for i, v := range list {
		if v == target {
			return i
		}
	}
	return -1 // not found
}

func resortSlice(PrimaryColumn, SecondaryColumn int, rows [][]string) [][]string {
	sort.Slice(rows, func(i, j int) bool {
		name1 := rows[i][PrimaryColumn]
		name2 := rows[j][PrimaryColumn]
		type1 := rows[i][SecondaryColumn]
		type2 := rows[j][SecondaryColumn]

		if type1 != type2 {
			return type1 < type2
		}
		return name1 < name2
	})
	return rows
}

func GenerateFullList(orderField int, fields ...string) [][]string {
	retv := [][]string{}
	for _, row := range GenerateListOfMaps() {
		tmplist := []string{}
		for _, field := range fields {
			if value, ok := row[field]; ok {
				tmplist = append(tmplist, value)
			} else {
				tmplist = append(tmplist, "")
			}
		}
		retv = append(retv, tmplist)
	}
	return retv
}

// PrintFullList prints the list of sessions
func PrintFullList() {
	legend := []string{"Name", "Type", "Status", "Active", "Description", "Directory"}

	SecondaryColumn := indexOf("Type", legend)

	rows := GenerateFullList(SecondaryColumn, legend...)
	rows = resortSlice(0, SecondaryColumn, rows)

	outputtable := tablewriter.NewWriter(os.Stdout)
	outputtable.Header(legend)
	outputtable.Bulk(rows)
	outputtable.Render()

}

// PrintFullList prints the list of sessions
func PrintSanityList() {

	legend := []string{"Name", "Type", "Status", "Sane", "Description", "Directory"}

	SecondaryColumn := indexOf("Type", legend)

	rows := GenerateFullList(SecondaryColumn, legend...)
	rows = resortSlice(0, SecondaryColumn, rows)

	outputtable := tablewriter.NewWriter(os.Stdout)
	outputtable.Header(legend)
	outputtable.Bulk(rows)
	outputtable.Render()
}

// PrintShortList prints the list of sessions
func PrintShortList() {
	for _, item := range ListConfigs() {
		fmt.Printf("%s\n", item)
	}
}
