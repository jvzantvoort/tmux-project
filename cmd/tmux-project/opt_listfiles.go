package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	msg "github.com/jvzantvoort/tmux-project/messages"

	"github.com/jvzantvoort/tmux-project/sessions"
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
	return msg.GetUsage("listfiles")
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
	session := sessions.NewTmuxSession(c.projectname)
	for _, ink := range session.TargetPaths() {
		fmt.Printf("%s\n", ink)
	}

	log.Debugln("End")

	return subcommands.ExitSuccess
}
