package project

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	errno "github.com/jvzantvoort/tmux-project/errors"
	"github.com/jvzantvoort/tmux-project/utils"
)

// Write json output to an [io.Writer] compatible handle. It returns nil or the
// error of [json.MarshalIndent]
func (proj Project) Write(writer io.Writer) error {
	utils.LogStart()
	defer utils.LogEnd()

	content, err := json.MarshalIndent(proj, "", "  ")
	if err == nil {
		_, err := fmt.Fprintf(writer, "%s\n", string(content))
		if err != nil {
			return err
		}
	}
	return err
}

// Read session content from a [io.Reader] object.
func (proj *Project) Read(reader io.Reader) error {
	utils.LogStart()
	defer utils.LogEnd()

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &proj)
	if err != nil {
		return err
	}
	return nil
}

func (proj Project) ProjectConfigFile() string {
	return filepath.Join(config.SessionDir(), proj.ProjectName+".json")
}

func (proj *Project) Open() error {
	utils.LogStart()
	defer utils.LogEnd()

	projectfile := proj.ProjectConfigFile()
	utils.Debugf("project file: %s", projectfile)

	if _, err := os.Stat(projectfile); os.IsNotExist(err) {
		utils.Debugf("project file not found")
		return errno.ErrProjectNotExist
	}

	filehandle, err := os.Open(projectfile)
	if err != nil {
		utils.Errorf("cannot open project file: %s", err)
		return err
	}

	return proj.Read(filehandle)

}

// Write session configuration to a projectfile
func (proj Project) Save() error {
	utils.LogStart()
	defer utils.LogEnd()

	err := utils.MkdirAll(filepath.Join(proj.ProjectDir, ".tmux-project"))
	if err != nil {
		return err
	}

	projectfile := proj.ProjectConfigFile()
	utils.Debugf("project file: %s", projectfile)

	filehandle, err := os.OpenFile(projectfile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.Errorf("cannot open project file: %s", err)
		return err
	}
	defer filehandle.Close()
	return proj.Write(filehandle)
}

// vim: noexpandtab filetype=go
