package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	tp "github.com/jvzantvoort/tmux-project"
	msg "github.com/jvzantvoort/tmux-project/messages"
	log "github.com/sirupsen/logrus"
)

type ListSubCmd struct {
	projecttype string
	projectname string
	printfull   bool
	verbose     bool
}

func (*ListSubCmd) Name() string {
	return "list"
}

func (*ListSubCmd) Synopsis() string {
	return "List projects"
}

func (*ListSubCmd) Usage() string {
	return msg.GetUsage("list")
}

func (c *ListSubCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.projectname, "projectname", "", "Name of project")
	f.StringVar(&c.projectname, "n", "", "Name of project")
	f.BoolVar(&c.printfull, "f", false, "Print full")
	f.BoolVar(&c.verbose, "v", false, "Verbose logging")
}

func (c *ListSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")

	if c.printfull {
		tp.PrintFullList()
	} else {
		tp.PrintShortList()
	}

	log.Debugln("End")

	return subcommands.ExitSuccess
}
