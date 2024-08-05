package main

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/sessions"
	"github.com/jvzantvoort/tmux-project/tmux"
	"github.com/jvzantvoort/tmux-project/utils"
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

	sess := sessions.NewTmuxSessions()
	xsess, err := sess.Find(sessionname)
	utils.ErrorExit(err)
	tmux.Resume(sessionname, xsess.Configfile)

}
