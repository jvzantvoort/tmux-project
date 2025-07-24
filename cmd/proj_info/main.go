package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/utils"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Set up the logger
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

// cleanup handles any panic that occurs during execution, logging the error.
// It is deferred to ensure that it runs even if a panic occurs.
func cleanup() {
	if r := recover(); r != nil {
		log.Errorf("Paniced %s", r)
	}
}

// uniqueList returns a slice of unique strings from the provided arguments.
// It uses a map to track unique items and returns them in a slice.
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

// PrintHeader formats and prints the header information for the project.
// It uses color formatting for the names and values, making it visually distinct.
// The header includes the session name, project directory, description, and project type.
// It uses the tablewriter package to create a nicely formatted table output.
// The header is printed to the standard output.
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

// main is the entry point for the proj_info command.
// It retrieves the project information based on the current tmux session name.
// The session name is expected to be set in the environment variable SESSIONNAME.
// If the session name is not set, it will print an error message and exit.
// It also allows for verbose output and depth control via command line flags.
// The verbose flag enables detailed logging, and the depth flag controls the maximum depth for searching projects.
// The project information is printed in a formatted table, including the session name, project directory, description, and type.
// It also finds all projects in the specified directory and prints their information grouped by chapter.
// The output is formatted using the tablewriter package for better readability.
func main() {
	defer cleanup()
	var chapters []string
	verbose := flag.Bool("v", false, "Verbose")
	depth := flag.Int("d", 1, "Max depth in search")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Header
	sessionname := os.Getenv("SESSIONNAME")
	proj_obj := project.NewProject(sessionname)
	err := proj_obj.RefreshStruct()
	utils.ErrorExit(err)

	header := [][]string{
		{"Sessionname", proj_obj.Name},
		{"Projectdir", proj_obj.Directory},
		{"Description", proj_obj.Description},
		{"Type", proj_obj.ProjectType},
	}
	PrintHeader(header)
	brojects := findAllProjects(proj_obj.Directory, *depth)

	for _, proj := range brojects {
		chapters = append(chapters, proj.Chapter)
	}

	chapters = uniqueList(chapters...)

	sort.Slice(brojects, func(i, j int) bool { return brojects[i].Name < brojects[j].Name })

	for _, chapter := range chapters {
		printTitle(chapter)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Status", "Branch"})
		table.SetBorder(false)
		for _, proj := range brojects {
			if chapter == proj.Chapter {
				table.Append(proj.GetFields())
			}
		}
		table.Render()
	}
	fmt.Printf("\n")
}
