package config

import (
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ProjectTypeFile defines a structure of a file
type ProjectTypeFile struct {
	Name        string `yaml:"name"`
	Destination string `yaml:"destination"`
	Mode        string `yaml:"mode"`
}

// ProjectTypeConfig defines a structure of a project type
type ProjectTypeConfig struct {
	ProjectType    string `yaml:"projecttype"`
	ProjectTypeDir string
	Workdir        string            `yaml:"workdir"`
	Pattern        string            `yaml:"pattern"`
	SetupActions   []string          `yaml:"setupactions"`
	Files          []ProjectTypeFile `yaml:"files"`
}

func (ptc *ProjectTypeConfig) readConfig(projtypeconfigdir string) {

	viper.SetConfigName("config")
	viper.AddConfigPath(projtypeconfigdir)
	// viper.AddConfigPath(tp.MasterConfigDir)

	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	err = viper.Unmarshal(&ptc)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

// Describe describe
func (ptc ProjectTypeConfig) Describe() {
	log.Debugf("Describe: %s", ptc.ProjectType)
	log.Debugf("  Workdir: %s", ptc.Workdir)
	log.Debugf("  Pattern: %s", ptc.Pattern)

	fileno := len(ptc.Files)
	actionsno := len(ptc.SetupActions)

	if fileno > 0 {
		log.Debugf("  Files:")
		for _, act := range ptc.Files {
			log.Debugf("    - name: %s", act.Name)
			log.Debugf("      destination: %s", act.Destination)
			log.Debugf("      mode: %s", act.Mode)
		}
	}

	if actionsno > 0 {
		log.Debugf("  Actions:")
		for _, act := range ptc.SetupActions {
			log.Debugf("    - %s", act)
		}

	}

	log.Debugf("Describe: %s, end", ptc.ProjectType)

}

func (ptc ProjectTypeConfig) Write(boxname, target string) error {
	content, _ := Asset("templates/" + boxname)
	file, err := os.Create(target)
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func (ptc ProjectTypeConfig) Exists(targetpath string) bool {
	_, err := os.Stat(targetpath)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (ptc *ProjectTypeConfig) Init(projtypeconfigdir, projecttype string) error {

	projtypeconfigdir = path.Join(projtypeconfigdir, projecttype)

	if ptc.Exists(projtypeconfigdir) {
		return fmt.Errorf("Directory already exists: %s", projtypeconfigdir)
	}

	if err := os.MkdirAll(projtypeconfigdir, os.FileMode(int(0755))); err != nil {
		return fmt.Errorf("Directory cannot be created: %s", projtypeconfigdir)
	}

	targets := []string{"config.yml", "default.rc", "default.env"}
	for _, target := range targets {
		ptc.Write(target, path.Join(projtypeconfigdir, target))
	}
	return nil
}
