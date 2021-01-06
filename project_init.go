package tmuxproject

import (
	"github.com/jvzantvoort/tmux-project/config"
)

func CreateProjectType(projecttype string) error {
	var config config.ProjectTypeConfig
	config.Init(mainconfig.ProjTypeConfigDir, projecttype)
	return nil
}
