// Package config provides configuration data globally used
package config

import (
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/jvzantvoort/tmux-project/utils"
)

type MainConfig struct {
	HomeDir           string
	TmuxDir           string
	ProjTypeConfigDir string
}

func NewMainConfig() *MainConfig {

	v := &MainConfig{}

	home, err := homedir.Dir()
	utils.ErrorExit(err)
	v.HomeDir = home
	v.TmuxDir = path.Join(v.HomeDir, ".tmux.d")
	v.ProjTypeConfigDir = path.Join(v.HomeDir, ".tmux-project")

	return v

}
