package main

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/utils"
)

func main() {

	// Header
	sessionname := os.Getenv("SESSIONNAME")
	proj_obj := project.NewProject(sessionname)
	err := proj_obj.RefreshStruct()
	utils.ErrorExit(err)
	fmt.Printf("%s\n", proj_obj.Description)

}
