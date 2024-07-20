package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/jvzantvoort/tmux-project/sessions"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

}

func cleanup() {
	if r := recover(); r != nil {
		log.Errorf("Paniced %s", r)
	}
}

func uniqueList(args ...string) []string {
	retv := []string{}
	items_map := make(map[string]bool)

	for _, item := range args {
		items_map[item] = true
	}

	for k := range items_map {
		retv = append(retv, k)
	}

	return retv
}

func PrintHeader(data [][]string) {

	infNameCol := color.New(InfoNameColor)
	infValCol := color.New(InfoValueColor)
	table := tablewriter.NewWriter(os.Stdout)
	// table.SetHeader([]string{"Name", "Value"})
	// table.SetHeaderLine(true)
	table.SetBorder(true)
	for _, slice := range data {
		table.Append([]string{infNameCol.Sprint(slice[0]), infValCol.Sprint(slice[1])})
	}
	fmt.Printf("\n")
	table.Render()
	fmt.Printf("\n")

}

func main() {
	defer cleanup()
	var chapters []string
	verbose := flag.Bool("v", false, "Verbose")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Header
	sessionname := os.Getenv("SESSIONNAME")
	session := sessions.NewTmuxSession(sessionname)

	header := [][]string{
		{"Sessionname", sessionname},
		{"Projectdir", session.Workdir},
		{"Description", session.Description},
	}
	PrintHeader(header)
	brojects := findAllProjects(session.Workdir)

	for _, proj := range brojects {
		chapters = append(chapters, proj.Chapter)
	}

	chapters = uniqueList(chapters...)

	sort.Slice(brojects, func(i, j int) bool { return brojects[i].Name < brojects[j].Name })

	for _, chapter := range chapters {
		printTitle(chapter)
		for _, proj := range brojects {
			if chapter == proj.Chapter {
				proj.PrintLine()
			}
		}
	}
	fmt.Printf("\n")
	fmt.Printf("\n")
}
