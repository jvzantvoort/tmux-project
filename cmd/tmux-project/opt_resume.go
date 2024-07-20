package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/google/subcommands"
	msg "github.com/jvzantvoort/tmux-project/messages"
	"github.com/jvzantvoort/tmux-project/sessions"
	"github.com/jvzantvoort/tmux-project/tmux"
	"github.com/manifoldco/promptui"
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
	return msg.GetUsage("resume")
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (c *ResumeSubCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.projectname, "projectname", "", "Name of project")
	f.StringVar(&c.projectname, "n", "", "Name of project")
	f.BoolVar(&c.verbose, "v", false, "Verbose logging")
}

func ListSessions() []string {
	retv := []string{}
	tmux := tmux.NewTmux()
	active, _ := tmux.ListActive()

	sess := sessions.NewTmuxSessions()
	for _, sesi := range sess.Sessions {
		state := " "
		if stringInSlice(sesi.Name, active) {
			state = "active"
		}
		message := fmt.Sprintf("%-32s %-6s %s", sesi.Name, state, sesi.Description)

		retv = append(retv, message)
	}

	return retv
}

func (c *ResumeSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")
	//
	if len(c.projectname) == 0 {
		prompt := promptui.Select{
			Label: "Select project",
			Size:  20,
			Items: ListSessions(),
		}
		_, result, err := prompt.Run()

		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		result = strings.Split(result, " ")[0]
		c.projectname = result
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
