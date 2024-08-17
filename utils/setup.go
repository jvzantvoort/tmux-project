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

	session_dir := config.SessionDir()
	commonfile := filepath.Join(session_dir, "common.rc")
	tos_file := filepath.Join(session_dir, "tmux_opt_source")

	if FileExists(commonfile) {
		Debugf("common.rc file exists")
		return nil
	}

	if noexec {
		fmt.Printf("create %s file\n", commonfile)
		fmt.Printf("create %s file\n", tos_file)
		return nil
	}

	if err := messages.Copy("common.rc", commonfile, 0644); err != nil {
		return err
	}

	if err := messages.Copy("tmux_opt_source", tos_file, 0755); err != nil {
		return err
	}

	return nil
}
