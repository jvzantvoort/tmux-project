package utils

import (
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/messages"
)

func SetupSessionDir() error {
	basedir := config.SessionDir()
	commonfile := filepath.Join(basedir, "common.rc")
	tos_file := filepath.Join(basedir, "tmux_opt_source")

	if FileExists(commonfile) {
		return nil
	}

	err := MkdirAll(config.SessionDir())
	LogIfError(err)
	content := messages.GetConfig("common.rc")
	tos_content := messages.GetConfig("tmux_opt_source")

	rc_fh, err := os.OpenFile(commonfile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Errorf("cannot open common.rc file for writing: %s", err)
		return err
	}
	defer rc_fh.Close()
	_, err = rc_fh.WriteString(content)
	if err != nil {
		return err
	}

	tos_fh, err := os.OpenFile(tos_file, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		Errorf("cannot open tmux_opt_source file for writing: %s", err)
		return err
	}
	defer tos_fh.Close()
	_, err = tos_fh.WriteString(tos_content)

	return err
}
