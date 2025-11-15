package project

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jvzantvoort/tmux-project/config"
	errno "github.com/jvzantvoort/tmux-project/errors"
	"github.com/jvzantvoort/tmux-project/utils"
)

// Write serializes the project to JSON format and writes it to the provided io.Writer.
// It returns nil on success or an error from json.MarshalIndent.
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

// Read deserializes project data from an io.Reader and populates the project struct.
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

// ProjectConfigFile returns the full path to the project's configuration file.
func (proj Project) ProjectConfigFile() string {
	return filepath.Join(config.SessionDir(), proj.Name+".json")
}

// Open reads and loads an existing project configuration from disk.
// Returns ErrProjectNotExist if the project file doesn't exist.
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

// Save writes the project configuration to disk in the session directory.
func (proj Project) Save() error {
	utils.LogStart()
	defer utils.LogEnd()

	err := utils.SetupSessionDir(false)
	if err != nil {
		return err
	}

	configfile := proj.ProjectConfigFile()
	utils.Debugf("project file: %s", configfile)

	filehandle, err := os.OpenFile(configfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		utils.Errorf("cannot open project file for writing: %s", err)
		return err
	}
	defer filehandle.Close()
	return proj.Write(filehandle)
}

// UpdateLastActivity updates the LastActivity timestamp to the current time and saves the project.
func (proj *Project) UpdateLastActivity() error {
	utils.LogStart()
	defer utils.LogEnd()

	proj.LastActivity = time.Now().Format(time.RFC3339)
	return proj.Save()
}
