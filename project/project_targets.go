package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/utils"
)

func (proj Project) CalcDestination(instr string) string {
	retv := proj.Parse(instr)

	if !filepath.IsAbs(retv) {
		retv = filepath.Join(config.SessionDir(), retv)
	}
	return retv
}

func (proj Project) ProcessTarget(element *Target) error {

	utils.LogStart()
	defer utils.LogEnd()

	element.Name = proj.Parse(element.Name)

	mode, err := utils.GetMode(element.Mode)
	if err != nil {
		return err
	}

	// source_file := filepath.Join(proj.ProjectTypeDir, element.Name)
	dest_file := proj.CalcDestination(element.Destination)
	content := proj.Parse(element.Content)

	filehandle, err := os.Create(dest_file)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", dest_file, err)
	}
	defer func() {
		if cerr := filehandle.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	_, err = filehandle.WriteString(content)
	if err != nil {
		return err
	}
	if err := os.Chmod(dest_file, mode); err != nil {
		return err
	}

	return nil
}
