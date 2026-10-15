package main

import (
	"fmt"
	"strings"

	errno "github.com/jvzantvoort/tmux-project/errors"
	"github.com/jvzantvoort/tmux-project/project"
	"github.com/jvzantvoort/tmux-project/tmux"
	"github.com/jvzantvoort/tmux-project/utils"
)

// ProjectRemove removes the project directory, its target files and its
// session files. Unless assumeYes is set the user has to confirm first.
// With withArchive the project is archived before anything is removed.
func ProjectRemove(projectname string, withArchive, assumeYes bool) error {
	utils.LogStart()
	defer utils.LogEnd()

	proj := project.NewProject(projectname)
	if err := proj.RefreshStruct(); err != nil {
		if errno.IsProjectNotExist(err) || err == errno.ErrProjectTypeNotDefined {
			return fmt.Errorf("project %s does not exist", projectname)
		}
		return err
	}

	active, _ := tmux.ListActive()
	if utils.StringInSlice(tmux.SessionName(proj.Name), active) {
		return fmt.Errorf("tmux session %s is still running, close it first", proj.Name)
	}

	targets := proj.RemovalTargets()
	if len(targets) == 0 {
		return fmt.Errorf("nothing to remove for project %s", projectname)
	}

	fmt.Printf("The following will be removed:\n\n")
	for _, target := range targets {
		fmt.Printf("  %s\n", target)
	}
	fmt.Println()

	if !assumeYes {
		answer := strings.ToLower(strings.TrimSpace(utils.Ask("Remove project " + proj.Name + "? [y/N]")))
		if answer != "y" && answer != "yes" {
			return fmt.Errorf("aborted, nothing removed")
		}
	}

	if withArchive {
		if err := CreateArchive(projectname, ""); err != nil {
			return fmt.Errorf("archive failed, nothing removed: %w", err)
		}
	}

	if err := proj.Remove(); err != nil {
		return err
	}
	utils.Infof("project %s removed", proj.Name)
	return nil
}
