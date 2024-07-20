package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	msg "github.com/jvzantvoort/tmux-project/messages"
	log "github.com/sirupsen/logrus"
)

type ShellProfileCmd struct {
	shellname string
	verbose   bool
}

func (*ShellProfileCmd) Name() string {
	return "shell"
}

func (*ShellProfileCmd) Synopsis() string {
	return "Edit a projects tmux configuration"
}

func (*ShellProfileCmd) Usage() string {
	return msg.GetUsage("shell")
}

func (c *ShellProfileCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.shellname, "shellname", "bash", "Name of the shell profile to provide")
	f.StringVar(&c.shellname, "s", "bash", "Name of the shell profile to provide")
	f.BoolVar(&c.verbose, "v", false, "Verbose logging")
}

func (c *ShellProfileCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")

	fmt.Print(msg.GetShell(c.shellname))

	log.Debugln("End")

	return subcommands.ExitSuccess
}
