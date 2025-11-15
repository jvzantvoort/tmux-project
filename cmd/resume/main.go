package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/project"
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

	proj := project.NewProject(sessionname)
	if err := proj.Open(); err == nil {
		proj.UpdateLastActivity()
	}
	configfile := filepath.Join(config.SessionDir(), proj.Name+".rc")

	tmux.Resume(sessionname, configfile)

}
