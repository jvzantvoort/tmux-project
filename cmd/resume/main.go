package main

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/sessions"
	"github.com/jvzantvoort/tmux-project/tmux"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE:\n\n\t%s [ls|<session>]\n\n", os.Args[0])
		return
	}
	sessionname := os.Args[1]
	if sessionname == "ls" {
		project.PrintFullList()
		return
	}

	_tmux := tmux.NewTmux()
	found := false
	active := false

	sess := sessions.NewTmuxSessions()
	xsess := sessions.TmuxSession{}
	for _, sesi := range sess.Sessions {
		if sessionname == sesi.Name {
			xsess = sesi
			found = true
		}
	}

	if found {
		if _tmux.SessionExists(xsess.Name) {
			active = true
		}
	} else {
		os.Exit(1)
	}

	if active {
		_tmux.ResumeSession(xsess)
	} else {
		_tmux.CreateSession(xsess)
	}

}
