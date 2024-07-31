// Package config provides configuration data globally used
//
// Example:
//
//	import (
//	  "fmt"
//	  "github.com/jvzantvoort/tmux-project/config"
//	)
//
//	fmt.Printf("home dir: %s", config.Home())
//	fmt.Printf("tmux dir: %s", config.SessionDir())
//	fmt.Printf("project type config dir: %s", mainconfig.ConfigDir())
package config

import (
	"path"

	"github.com/jvzantvoort/tmux-project/utils"
)

func Home() string {
	return utils.HomeDir()
}

func ConfigDir() string {
	return path.Join(Home(), ".tmux-project")
}

func SessionDir() string {
	return path.Join(Home(), ".tmux.d")
}
