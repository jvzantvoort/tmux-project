package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	tp "github.com/jvzantvoort/tmux-project"
	msg "github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/sessions"
	log "github.com/sirupsen/logrus"
)

type EditSubCmd struct {
	projecttype string
	projectname string
	verbose     bool
}

func (*EditSubCmd) Name() string {
	return "edit"
}

func (*EditSubCmd) Synopsis() string {
	return "Edit a projects tmux configuration"
}

func (*EditSubCmd) Usage() string {
	return msg.GetUsage("edit")
}

func (c *EditSubCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.projectname, "projectname", "", "Name of project")
	f.StringVar(&c.projectname, "n", "", "Name of project")
	f.BoolVar(&c.verbose, "v", false, "Verbose logging")
}

func (c *EditSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")
	//
	if len(c.projectname) == 0 {
		log.Fatalf("no name provided")
	}

	session := sessions.NewTmuxSession(c.projectname)

	tp.Edit("-O", session.Configfile, session.Environment)

	log.Debugln("End")

	return subcommands.ExitSuccess
}
