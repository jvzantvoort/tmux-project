package project

import (
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
)

func (proj Project) ProcessProjectTarget(element *ProjectTarget) error {
	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)

	element.Name = proj.Parse(element.Name)
	element.Destination = proj.Parse(element.Destination)

	mode, err := utils.GetMode(element.Mode)
	if err != nil {
		return err
	}

	// source_file := filepath.Join(proj.ProjectTypeDir, element.Name)
	dest_file := filepath.Join(config.SessionDir(), element.Destination)
	content := proj.Parse(element.Content)

	filehandle, _ := os.Create(dest_file)
	_, err = filehandle.WriteString(content)
	if err != nil {
		return err
	}
	defer filehandle.Close()
	if err := os.Chmod(dest_file, mode); err != nil {
		return err
	}

	return nil
}
