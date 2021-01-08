// Package config provides configuration data globally used
//
//   import (
//     "fmt"
//     "github.com/jvzantvoort/tmux-project/config"
//   )
//
//   mainconfig := NewMainConfig()
//   fmt.Printf("home dir: %s", mainconfig.HomeDir)
//   fmt.Printf("tmux dir: %s", mainconfig.TmuxDir)
//   fmt.Printf("project type config dir: %s", mainconfig.ProjTypeConfigDir)
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
	v.TmuxDir = path.Join(v.HomeDir, ".bash", "tmux.d")
	v.ProjTypeConfigDir = path.Join(v.HomeDir, ".tmux-project")

	return v

}
