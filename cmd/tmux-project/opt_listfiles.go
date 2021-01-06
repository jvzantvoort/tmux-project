package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	tp "github.com/jvzantvoort/tmux-project"
	log "github.com/sirupsen/logrus"
)

type ListFilesSubCmd struct {
	projectname string
	verbose     bool
}

func (*ListFilesSubCmd) Name() string {
	return "listfiles"
}

func (*ListFilesSubCmd) Synopsis() string {
	return "Archive a project"
}

func (*ListFilesSubCmd) Usage() string {
	msgstr, err := tp.Asset("messages/usage_listfiles")
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return string(msgstr)
}

func (c *ListFilesSubCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.projectname, "projectname", "", "Name of project")
	f.StringVar(&c.projectname, "n", "", "Name of project")
	f.BoolVar(&c.verbose, "v", false, "Verbose logging")
}

func (c *ListFilesSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")
	//
	if len(c.projectname) == 0 {
		log.Fatalf("no name provided")
	}

	err := tp.ListProject(c.projectname)
	if err != nil {
		log.Fatalf("Encountered error: %q", err)
	}

	log.Debugln("End")

	return subcommands.ExitSuccess
}
