package projecttype

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/utils"

	"gopkg.in/yaml.v2"
)

//go:embed templates/*
var Content embed.FS

// ProjectTypeFile defines a structure of a file
type ProjectTypeFile struct {
	Name        string `yaml:"name"`
	Destination string `yaml:"destination"`
	Mode        string `yaml:"mode"`
}

type RepoListItem struct {
	Name      string `yaml:"name"`
	RepoNames string `yaml:"reponames"`
}

// ProjectTypeConfig defines a structure of a project type
type ProjectTypeConfig struct {
	ProjectType    string `yaml:"projecttype"`
	ProjectTypeDir string
	Workdir        string            `yaml:"workdir"`
	Pattern        string            `yaml:"pattern"`
	SetupActions   []string          `yaml:"setupactions"`
	RepoListItems  []RepoListItem    `yaml:"repolist"`
	Files          []ProjectTypeFile `yaml:"files"`
	ConfigFile     string            `yaml:"-"`
	ConfigDir      string            `yaml:"-"`
}

func (ptc *ProjectTypeConfig) readConfig() error {
	// Setup logging
	log_prefix := utils.FunctionName()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	if !utils.FileExists(ptc.ConfigFile) {
		return fmt.Errorf("ConfigFile: %s does not exist", ptc.ConfigFile)
	}

	yamlFile, err := os.ReadFile(ptc.ConfigFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &ptc)
	return err
}

// Describe describe
func (ptc ProjectTypeConfig) Describe() {
	// Setup logging
	log_prefix := utils.FunctionName()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

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

func (ptc *ProjectTypeConfig) SetupProjectTypeConfig() error {

	if err := utils.MkdirAll(ptc.ProjectTypeDir); err != nil {
		return err
	}

	// write basic files
	targets := []string{"config.yml", "default.rc", "default.env"}
	for _, target := range targets {
		fpath := path.Join(ptc.ProjectTypeDir, target)
		err := ptc.Write(target, fpath)
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

func (ptc ProjectTypeConfig) Content(target string) (string, error) {
	ipath := filepath.Join(ptc.ConfigDir, target)

	content, err := os.ReadFile(ipath)
	if err != nil {
		return "", err
	}
	return string(content), err
}

func (ptc ProjectTypeConfig) Write(boxname, target string) error {
	log_prefix := utils.FunctionName()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)
	filename := fmt.Sprintf("templates/%s", boxname)
	content, err := Content.ReadFile(filename)
	if err != nil {
		log.Error(err)
		content = []byte("undefined")
	}
	file, _ := os.Create(target)
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func (ptc ProjectTypeConfig) UpdateConfigFile(target string) error {
	log_prefix := utils.FunctionName()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

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

// NewProjectTypeConfig read the relevant configfile and return
// ProjectTypeConfig object with relevant data.
func NewProjectTypeConfig(projecttype string) ProjectTypeConfig {
	var err error

	retv := ProjectTypeConfig{ProjectType: projecttype}

	retv.ProjectTypeDir = filepath.Join(config.ConfigDir(), retv.ProjectType)
	retv.ConfigDir = filepath.Join(config.ConfigDir(), retv.ProjectType)
	retv.ConfigFile = filepath.Join(retv.ConfigDir, "config.yml")

	err = retv.readConfig()

	if err != nil {
		log.Debugf("Project Type %s does not exist", retv.ProjectType)
	}

	retv.Workdir, err = utils.Expand(retv.Workdir)
	if err != nil {
		log.Errorf("Failed to get workdir %s", err)
	}

	return retv
}
