package config

import (
	"os/user"
	"path"
)

type MainConfig struct {
	HomeDir           string
	TmuxDir           string
	ProjTypeConfigDir string
}

func NewMainConfig() *MainConfig {

	v := &MainConfig{}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	v.HomeDir = usr.HomeDir
	v.TmuxDir = path.Join(v.HomeDir, ".bash", "tmux.d")
	v.ProjTypeConfigDir = path.Join(v.HomeDir, ".tmux-project")

	return v

}
