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
	return filepath.Join(config.SessionDir(), proj.Name+".json")
}

func (proj *Project) Open() error {
	utils.LogStart()
	defer utils.LogEnd()

	configfile := proj.ProjectConfigFile()
	utils.Debugf("project file: %s", configfile)

	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		utils.Debugf("project file not found")
		return errno.ErrProjectNotExist
	}

	filehandle, err := os.Open(configfile)
	if err != nil {
		utils.Errorf("cannot open project file for reading: %s", err)
		return err
	}

	return proj.Read(filehandle)

}

// Write session configuration to a projectfile
func (proj Project) Save() error {
	utils.LogStart()
	defer utils.LogEnd()

	err := utils.SetupSessionDir(false)
	if err != nil {
		return err
	}

	configfile := proj.ProjectConfigFile()
	utils.Debugf("project file: %s", configfile)

	filehandle, err := os.OpenFile(configfile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.Errorf("cannot open project file for writing: %s", err)
		return err
	}
	defer filehandle.Close()
	return proj.Write(filehandle)
}

// vim: noexpandtab filetype=go
