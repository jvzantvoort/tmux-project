package main

import (
	"flag"
	"os"
	"sort"

	"github.com/jvzantvoort/tmux-project/sessions"
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

func main() {
	defer cleanup()
	var chapters []string
	verbose := flag.Bool("v", false, "Verbose")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	sessionname := os.Getenv("SESSIONNAME")
	session := sessions.NewTmuxSession(sessionname)
	printInfo("Sessionname", sessionname)
	printInfo("Projectdir", session.Workdir)
	printInfo("Description", session.Description)

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
}
