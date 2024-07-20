// Package config provides configuration data globally used
//
package config

import (
	"os/user"
	"path"
	"path/filepath"
)

type MainConfig struct {
	HomeDir           string
	TmuxDir           string
	ProjTypeConfigDir string
}

// ExpandHome expand the tilde in a given path.
func (m MainConfig) ExpandHome(pathstr string) (string, error) {
	if len(pathstr) == 0 {
		return pathstr, nil
	}

	if pathstr[0] != '~' {
		return pathstr, nil
	}

	return filepath.Join(m.HomeDir, pathstr[1:]), nil

}

func NewMainConfig() *MainConfig {

	v := &MainConfig{}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	v.HomeDir = usr.HomeDir
	v.TmuxDir = path.Join(v.HomeDir, ".tmux.d")
	v.ProjTypeConfigDir = path.Join(v.HomeDir, ".tmux-project")

	return v

}
