package projecttype

import (
	"github.com/jvzantvoort/tmux-project/config"
)

var (
	mainconfig = config.NewMainConfig()
)

func CreateProjectType(projecttype string) error {
	var pt ProjectTypeConfig
	pt.Init(mainconfig.ProjTypeConfigDir, projecttype)
	return nil
}
