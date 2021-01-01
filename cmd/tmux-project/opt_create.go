package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	tp "github.com/jvzantvoort/tmux-project"
	log "github.com/sirupsen/logrus"
)

type CreateSubCmd struct {
	projecttype string
	projectname string
	verbose     bool
}

func (*CreateSubCmd) Name() string {
	return "create"
}

func (*CreateSubCmd) Synopsis() string {
	return "Create a new project"
}

func (*CreateSubCmd) Usage() string {
	msgstr, err := tp.Asset("messages/usage_create")
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return string(msgstr)
}

func (c *CreateSubCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.projecttype, "projecttype", "default", "Type of project")
	f.StringVar(&c.projecttype, "t", "default", "Type of project")
	f.StringVar(&c.projectname, "projectname", "", "Name of project")
	f.StringVar(&c.projectname, "n", "", "Name of project")
	f.BoolVar(&c.verbose, "v", false, "Verbose logging")
}

func (c *CreateSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")
	//
	if len(c.projectname) == 0 {
		log.Fatalf("no name provided")
	}
	err := tp.CreateProject(c.projecttype, c.projectname)
	if err != nil {
		log.Fatalf("Encountered error: %q", err)
	}
	//
	log.Debugln("End")

	return subcommands.ExitSuccess
}
