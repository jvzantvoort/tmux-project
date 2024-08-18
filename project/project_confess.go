package project

import (
	"fmt"

	"github.com/jvzantvoort/tmux-project/utils"
)

func (proj Project) Confess() {
	content := proj.Parse(GetScriptContent("confess"))
	fmt.Print(content)
}

func (projt Target) Confess() {
	utils.Debugf("  %-32s %s", "Name", projt.Name)
	utils.Debugf("      %-28s %s", "Name", projt.Destination)
	utils.Debugf("      %-28s %s", "Name", projt.Mode)
}
