package utils

import (
	"fmt"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/messages"
)

// SetupSessionDir setup the ~/.tmux.d directory
func SetupSessionDir(noexec bool) error {
	LogStart()
	defer LogEnd()

	sessionDir := config.SessionDir()
	commonfile := filepath.Join(sessionDir, "common.rc")
	tmuxOptSourceFile := filepath.Join(sessionDir, "tmux_opt_source")

	if FileExists(commonfile) {
		Debugf("common.rc file exists")
		return nil
	}

	if noexec {
		fmt.Printf("create %s file\n", commonfile)
		fmt.Printf("create %s file\n", tmuxOptSourceFile)
		return nil
	}

	if err := messages.Copy("common.rc", commonfile, 0644); err != nil {
		return err
	}

	if err := messages.Copy("tmux_opt_source", tmuxOptSourceFile, 0755); err != nil {
		return err
	}

	return nil
}
