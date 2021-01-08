package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	tp "github.com/jvzantvoort/tmux-project"
	"github.com/jvzantvoort/tmux-project/projecttype"
	log "github.com/sirupsen/logrus"
)

type InitProjSubCmd struct {
	projecttype string
	force       bool
	verbose     bool
}

func (*InitProjSubCmd) Name() string {
	return "init"
}

func (*InitProjSubCmd) Synopsis() string {
	return "Initialize a new project type"
}

func (*InitProjSubCmd) Usage() string {
	msgstr, err := tp.Asset("messages/usage_projectinit")
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return string(msgstr)
}

func (c *InitProjSubCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.projecttype, "projecttype", "default", "Type of project")
	f.StringVar(&c.projecttype, "t", "default", "Type of project")
	f.BoolVar(&c.force, "f", false, "Force (re)creation")
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
	} else if c.projecttype == "default" {
		if !c.force {
			log.Fatalf("Cannot overwrite default")
		}
	}
	log.Debugf("type: %s", c.projecttype)
	projecttype.CreateProjectType(c.projecttype)
	log.Debugln("End")

	return subcommands.ExitSuccess
}
