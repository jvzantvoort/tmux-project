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
	"os"
	"os/user"
	"path/filepath"
)

// Home returns the home directory of the user. $HOME takes precedence, as
// it does for os.UserHomeDir, so tests and wrappers can relocate it.
func Home() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

// ConfigDir returns the directory holding the project type configurations
func ConfigDir() string {
	return filepath.Join(Home(), ".tmux-project")
}

// SessionDir returns the directory holding the project and tmux files
func SessionDir() string {
	return filepath.Join(Home(), ".tmux.d")
}
