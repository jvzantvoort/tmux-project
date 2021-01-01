package main

import (
	"flag"
	"os"

	tp "github.com/jvzantvoort/tmux-project"
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
	log.SetLevel(log.DebugLevel)
}

func main() {

	defer func() {
		if panicname := recover(); panicname != nil {
			log.Println("Found error", panicname)
			os.Exit(1)
		}
	}()

	projectname := ""
	archivename := ""
	verbose := false

	flags := flag.NewFlagSet("archive", flag.ExitOnError)
	flags.StringVar(&projectname, "projectname", projectname, "Name of project")
	flags.StringVar(&projectname, "n", projectname, "Name of project")
	flags.StringVar(&archivename, "archivename", archivename, "Archive file")
	flags.StringVar(&archivename, "a", archivename, "Archive file")
	flags.BoolVar(&verbose, "v", false, "Verbose logging")
	flags.Parse(os.Args[1:])

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")

	if len(projectname) == 0 {
		log.Fatalf("no name provided")
	}

	if archivename == "" {
		archivename, _ = tp.GetWorkdir(projectname)
		archivename = archivename + ".tar.gz"
	}

	err := tp.ArchiveProject(projectname, archivename)
	if err != nil {
		log.Fatalf("Encountered error: %q", err)
	}

	log.Debugln("End")
}
