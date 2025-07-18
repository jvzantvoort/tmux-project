// proj_descr/main.go
package main

import (
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/utils"
)

// main is the entry point for the proj_descr command, which prints the project description.
// It retrieves the project description based on the current tmux session name.
// The session name is expected to be set in the environment variable SESSIONNAME.
// If the session name is not set, it will print an error message and exit.
func main() {

	// Header
	sessionname := os.Getenv("SESSIONNAME")
	proj_obj := project.NewProject(sessionname)
	err := proj_obj.RefreshStruct()
	utils.ErrorExit(err)
	fmt.Printf("%s\n", proj_obj.Description)

}
