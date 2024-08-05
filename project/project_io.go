package project

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
)

// Write json output to an [io.Writer] compatible handle. It returns nil or the
// error of [json.MarshalIndent]
func (proj Project) Write(writer io.Writer) error {
	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)
	content, err := json.MarshalIndent(proj, "", "  ")
	if err == nil {
		charno, err := fmt.Fprintf(writer, "%s\n", string(content))
		log.Debugf("%s:  chars: %d\n", charno)
		if err != nil {
			return err
		}
	}
	return err
}

// Read session content from a [io.Reader] object.
func (proj *Project) Read(reader io.Reader) error {
	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)
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

// Write session configuration to a projectfile
func (proj Project) Save() error {
	functionname := utils.FunctionName(2)
	log.Debugf("%s: start", functionname)
	defer log.Debugf("%s: end", functionname)

	err := utils.MkdirAll(filepath.Join(proj.ProjectDir, ".tmux-project"))
	if err != nil {
		return err
	}

	projectfile := filepath.Join(proj.ProjectDir, ".tmux-project", "project.json")

	filehandle, err := os.OpenFile(projectfile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer filehandle.Close()
	return proj.Write(filehandle)
}

// vim: noexpandtab filetype=go
