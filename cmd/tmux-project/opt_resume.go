package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	tp "github.com/jvzantvoort/tmux-project"
	"github.com/jvzantvoort/tmux-project/sessions"
	"github.com/jvzantvoort/tmux-project/tmux"
	log "github.com/sirupsen/logrus"
)

type ResumeSubCmd struct {
	projectname string
	verbose     bool
}

func (*ResumeSubCmd) Name() string {
	return "resume"
}

func (*ResumeSubCmd) Synopsis() string {
	return "Resume a project"
}

func (*ResumeSubCmd) Usage() string {
	msgstr, err := tp.Asset("messages/usage_resume")
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return string(msgstr)
}

func (c *ResumeSubCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.projectname, "projectname", "", "Name of project")
	f.StringVar(&c.projectname, "n", "", "Name of project")
	f.BoolVar(&c.verbose, "v", false, "Verbose logging")
}

func (c *ResumeSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")
	//
	if len(c.projectname) == 0 {
		log.Fatalf("no name provided")
	}
	_tmux := tmux.NewTmux()
	found := false
	active := false

	sess := sessions.NewTmuxSessions()
	xsess := sessions.TmuxSession{}
	for _, sesi := range sess.Sessions {
		if c.projectname == sesi.Name {
			xsess = sesi
			found = true
		}
	}


	if found {
		if _tmux.SessionExists(xsess.Name) {
			active = true
		}
	} else {
		return subcommands.ExitFailure
	}

	if active {
		_tmux.ResumeSession(xsess)
	} else {
		_tmux.CreateSession(xsess)
	}
	// err := session.Resume(c.archivename)
	// if err != nil {
	// 	log.Fatalf("Encountered error: %q", err)
	// }

	log.Debugln("End")

	return subcommands.ExitSuccess
}
