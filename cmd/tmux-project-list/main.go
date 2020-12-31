package main

import (
	"flag"
	"fmt"

	// 	"fmt"
	"os"
	// 	"strconv"

	tp "github.com/jvzantvoort/tmux-project"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {

	defer func() {
		if panicname := recover(); panicname != nil {
			log.Println("Found error", panicname)
			os.Exit(1)
		}
	}()

	projecttype := "default"
	projectname := ""


	flags := flag.NewFlagSet("list", flag.ExitOnError)
	flags.StringVar(&projecttype, "projecttype", projecttype, "Type of project")
	flags.StringVar(&projecttype, "t", projecttype, "Type of project")
	flags.StringVar(&projectname, "projectname", projectname, "Name of project")
	flags.StringVar(&projectname, "n", projectname, "Name of project")
	flags.BoolVar(&verbose, "v", false, "Verbose logging")
	flags.Parse(os.Args[1:])

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")

	for _, target := range tp.ListTmuxConfigs() {
		fmt.Println(target)

	}

	log.Debugln("End")
}
