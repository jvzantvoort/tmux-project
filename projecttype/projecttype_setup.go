package projecttype

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jvzantvoort/tmux-project/utils"
)

func (ptc *ProjectTypeConfig) Setup() error {

	// write basic files
	targets := []string{"config.yml", "default.rc", "default.env"}

	if err := utils.MkdirAll(ptc.ProjectTypeDir); err != nil {
		return err
	}

	for _, target := range targets {
		fpath := path.Join(ptc.ProjectTypeDir, target)
		err := ptc.Copy(target, fpath)
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}
		err = ptc.UpdateConfigFile(fpath)
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}
	}

	return nil
}

func (ptc ProjectTypeConfig) UpdateConfigFile(target string) error {
	utils.LogStart()
	defer utils.LogEnd()

	read, err := os.ReadFile(target)
	if err != nil {
		return err
	}

	content := string(read)
	ncontent := strings.Replace(content, "PROJECTTYPE", ptc.ProjectType, -1)
	if content == ncontent {
		return nil
	} else {
		content = ncontent
	}

	err = os.WriteFile(target, []byte(content), 0)
	if err != nil {
		return err
	}
	return nil
}
