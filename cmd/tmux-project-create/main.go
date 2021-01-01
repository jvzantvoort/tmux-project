package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	tp "github.com/jvzantvoort/tmux-project"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		DisableLevelTruncation: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func CreateProject(projecttype, projectname string) error {
	configuration := tp.GetProjectTypeConfig(projecttype, projectname)
	tp.DescribeProjectType(configuration)

	tmplvars := tp.NewProjTmplVars(projectname, configuration)
	tmplvars.ProjectDescription = tp.Ask("Description")
	for _, target := range configuration.Files {
		srccontent, _ := tp.LoadFile(target.Name, *tmplvars)
		file, err := os.Create(target.Destination)
		_, err = file.WriteString(srccontent)
		if err != nil {
			return err
		}
		defer file.Close()

		num, err := strconv.Atoi(target.Mode)
		if err != nil {
			return err
		}

		mode, _ := strconv.ParseUint(fmt.Sprintf("%04d", num), 8, 32)
		if err := os.Chmod(target.Destination, os.FileMode(mode)); err != nil {
			return err
		}
	}
	return nil
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
	verbose := false


	flags := flag.NewFlagSet("new", flag.ExitOnError)
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

	if len(projectname) == 0 {
		log.Fatalf("no name provided")
	}


	err := CreateProject(projecttype, projectname)
	if err != nil {
		log.Fatalf("Encountered error: %q", err)
	}

	log.Debugln("End")
}
