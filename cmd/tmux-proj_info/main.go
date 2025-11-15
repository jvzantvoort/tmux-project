package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

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
func PrintHeader(data [][]string) {
	infNameCol := color.New(InfoNameColor)
	infValCol := color.New(InfoValueColor)
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Name", "Value"})
	for _, slice := range data {
		table.Append([]string{infNameCol.Sprint(slice[0]), infValCol.Sprint(slice[1])})
	}
	fmt.Printf("\n")
	table.Render()
	fmt.Printf("\n")
}

// getCurrentTmuxSession returns the current tmux session name.
// It tries multiple methods to determine the session name:
// 1. Check SESSIONNAME environment variable
// 2. Check TMUX_PANE and query tmux for session name
// 3. Run tmux display-message to get session name
func getCurrentTmuxSession() (string, error) {
	// Try SESSIONNAME environment variable first
	if sessionname := os.Getenv("SESSIONNAME"); sessionname != "" {
		return sessionname, nil
	}

	// Try using TMUX_PANE
	if pane := os.Getenv("TMUX_PANE"); pane != "" {
		cmd := exec.Command("tmux", "display-message", "-p", "-t", pane, "#{session_name}")
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output)), nil
		}
	}

	// Try getting current session directly
	cmd := exec.Command("tmux", "display-message", "-p", "#{session_name}")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a tmux session or unable to determine session name")
	}

	return strings.TrimSpace(string(output)), nil
}

// main is the entry point for the tmux-proj_info command.
// This version is designed to work in tmux display-popup which runs in a separate shell.
// It detects the current tmux session automatically instead of relying on SESSIONNAME env var.
func main() {
	defer cleanup()
	var chapters []string
	verbose := flag.Bool("v", false, "Verbose")
	depth := flag.Int("d", 1, "Max depth in search")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Get current tmux session name (works in popup)
	sessionname, err := getCurrentTmuxSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		fmt.Fprintf(os.Stderr, "\nThis command must be run from within a tmux session.\n")
		os.Exit(1)
	}

	log.Debugf("Detected session: %s", sessionname)

	// Load project information
	proj_obj := project.NewProject(sessionname)
	err = proj_obj.RefreshStruct()
	utils.ErrorExit(err)
	lastActivityStr := "undef"

	if at, ok := proj_obj.TimeSinceLastActivity(); ok == nil {
		lastActivityStr = project.FormatDuration(at)
		lastActivityStr += " ago"
	}

	header := [][]string{
		{"Sessionname", proj_obj.Name},
		{"Projectdir", proj_obj.Directory},
		{"Description", proj_obj.Description},
		{"Type", proj_obj.ProjectType},
		{"Access", lastActivityStr},
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
		table.Header([]string{"Name", "Status", "Branch"})
		for _, proj := range brojects {
			if chapter == proj.Chapter {
				table.Append(proj.GetFields())
			}
		}
		table.Render()
	}
	fmt.Printf("\n")
}
