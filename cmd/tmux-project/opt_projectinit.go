package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	tp "github.com/jvzantvoort/tmux-project"
	log "github.com/sirupsen/logrus"
)

type InitProjSubCmd struct {
	projecttype string
	verbose     bool
}

func (*InitProjSubCmd) Name() string {
	return "init"
}

func (*InitProjSubCmd) Synopsis() string {
	return "Initialize a new project type"
}

func (*InitProjSubCmd) Usage() string {
	return `print [-capitalize] <some text>:
	    Print args to stdout.
	    `
}

func (c *InitProjSubCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.projecttype, "projecttype", "default", "Type of project")
	f.StringVar(&c.projecttype, "t", "default", "Type of project")
	f.BoolVar(&c.verbose, "v", false, "Verbose logging")
}

func (c *InitProjSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")
	//
	if len(c.projecttype) == 0 {
		log.Fatalf("no type provided")
	}
	tp.CreateProjectType(c.projecttype)
	log.Debugln("End")

	return subcommands.ExitSuccess
}
