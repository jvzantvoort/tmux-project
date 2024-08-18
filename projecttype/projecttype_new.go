package projecttype

import (
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/utils"
)

// New read the relevant configfile and return
// ProjectTypeConfig object with relevant data.
func New(projecttype string) (ProjectTypeConfig, error) {
	utils.LogStart()
	defer utils.LogEnd()

	utils.LogArgument("projecttype", projecttype)

	var rerr error // the one we return
	var err error

	// Create the object
	retv := ProjectTypeConfig{}

	if len(projecttype) == 0 {
		rerr = ErrProjectNameEmpty
		return retv, rerr
	}

	retv.ProjectType = projecttype

	retv.ProjectTypeDir = filepath.Join(config.ConfigDir(), retv.ProjectType)
	utils.LogVariable("retv.ProjectTypeDir", retv.ProjectTypeDir)

	retv.ConfigDir = filepath.Join(config.ConfigDir(), retv.ProjectType)
	utils.LogVariable("retv.ConfigDir", retv.ConfigDir)

	retv.ConfigFile = filepath.Join(retv.ConfigDir, "config.yml")
	utils.LogVariable("retv.ConfigDir", retv.ConfigDir)

	err = retv.Open()

	if err != nil {
		utils.Debugf("Project Type %s does not exist", retv.ProjectType)
		rerr = ErrProjectNotExists
	}

	retv.Directory, err = utils.Expand(retv.Directory)
	if err != nil {
		utils.Errorf("Failed to get workdir %s", err)
	}

	return retv, rerr
}
