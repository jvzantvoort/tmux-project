package projecttype

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/utils"

	"gopkg.in/yaml.v2"
)

// Read content from a project type file
func (ptc ProjectTypeConfig) Content(target string) (string, error) {
	ipath := filepath.Join(ptc.ConfigDir, target)

	content, err := os.ReadFile(ipath)
	if err != nil {
		return "", err
	}

	return string(content), err
}

// Read session configuration from a filehandle
func (ptc *ProjectTypeConfig) Read(reader io.Reader) error {
	utils.LogStart()
	defer utils.LogEnd()

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &ptc)
	utils.LogIfError(err)
	return err
}

// Write json output to an [io.Writer] compatible handle. It returns nil or the
// error of [json.MarshalIndent]
func (ptc ProjectTypeConfig) Write(writer io.Writer) error {
	utils.LogStart()
	defer utils.LogEnd()

	content, err := yaml.Marshal(ptc)
	if err == nil {
		_, err := fmt.Fprintf(writer, "%s\n", string(content))
		if err != nil {
			return err
		}
	}
	return err
}

// Open session configuration from a configfile
func (ptc *ProjectTypeConfig) Open() error {
	utils.LogStart()
	defer utils.LogEnd()

	configfile := ptc.ConfigFile

	utils.Debugf("project type file: %s", configfile)

	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}

	filehandle, err := os.Open(configfile)
	if err != nil {
		utils.Errorf("cannot open project type file for reading: %s", err)
		return fmt.Errorf("read error")
	}

	return ptc.Read(filehandle)

}

// Write session configuration to a configfile
func (ptc ProjectTypeConfig) Save() error {
	utils.LogStart()
	defer utils.LogEnd()

	err := utils.SetupSessionDir(false)
	if err != nil {
		return err
	}

	configfile := ptc.ConfigFile
	utils.Debugf("project type file: %s", configfile)

	filehandle, err := os.OpenFile(configfile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.Errorf("cannot open project type file for writing: %s", err)
		return err
	}
	defer filehandle.Close()
	return ptc.Write(filehandle)
}

// vim: noexpandtab filetype=go
