package projecttype

import (
	"embed"
	"fmt"
	"os"
	"path"
	"runtime"
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
}

// prefix returns a prefix for logging and messages based on function name.
func (ptc ProjectTypeConfig) prefix() string {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return elements[len(elements)-1]
}

func (ptc *ProjectTypeConfig) readConfig(projtypeconfigdir string) {
	// Setup logging
	log_prefix := ptc.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	filepath := path.Join(projtypeconfigdir, "config.yml")
	log.Debugf("filepath: %s", filepath)

	yamlFile, err := os.ReadFile(filepath)
	utils.ErrorExit(err)

	err = yaml.Unmarshal(yamlFile, &ptc)
	utils.ErrorExit(err)

}

// Describe describe
func (ptc ProjectTypeConfig) Describe() {
	// Setup logging
	log_prefix := ptc.prefix()
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

func (ptc ProjectTypeConfig) Write(boxname, target string) error {
	// Setup logging
	log_prefix := ptc.prefix()
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

func (ptc ProjectTypeConfig) Exists(targetpath string) bool {
	// Setup logging
	log_prefix := ptc.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)
	_, err := os.Stat(targetpath)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (ptc ProjectTypeConfig) UpdateConfigFile(target string) error {
	// Setup logging
	log_prefix := ptc.prefix()
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

func (ptc ProjectTypeConfig) MkdirAll(targetpath string) error {
	// Setup logging
	log_prefix := ptc.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)
	log.Debugf("mkdir %s, start", targetpath)
	defer log.Debugf("mkdir %s, end", targetpath)

	if target, err := os.Stat(targetpath); !os.IsNotExist(err) {
		if !target.IsDir() {
			return fmt.Errorf("mkdir %s, exists but is not a directory", targetpath)
		}
		log.Debugf("mkdir %s, already exists", targetpath)
		return nil
	}

	err := os.MkdirAll(targetpath, os.FileMode(int(0755)))

	if err != nil {
		return fmt.Errorf("mkdir %s, failed: %s", targetpath, err)
	}
	return nil

}

func (ptc *ProjectTypeConfig) Init(projtypeconfigdir, projecttype string) error {
	// Setup logging
	log_prefix := ptc.prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	log.Debugf("Init Start: %s", projecttype)
	projtypeconfigdir = path.Join(projtypeconfigdir, projecttype)

	ptc.ProjectType = projecttype
	ptc.ProjectTypeDir = projtypeconfigdir

	if err := ptc.MkdirAll(ptc.ProjectTypeDir); err != nil {
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

// NewProjectTypeConfig read the relevant configfile and return
// ProjectTypeConfig object with relevant data.
func NewProjectTypeConfig(projecttype string) ProjectTypeConfig {

	// Load main configuration targets
	projtypeconfigdir := path.Join(config.ConfigDir(), projecttype)

	log.Debugf("project type config dir: %s", projtypeconfigdir)
	log.Debugf("tmux dir: %s", config.SessionDir())

	v := ProjectTypeConfig{}
	v.readConfig(projtypeconfigdir)

	var err error
	v.Workdir, err = utils.Expand(v.Workdir)
	if err != nil {
		log.Errorf("%q", err)
	}

	log.Debugf("config >> %#v", v)

	return v
}

func CreateProjectType(projecttype string) error {
	var pt ProjectTypeConfig
	pt.Init(config.ConfigDir(), projecttype)
	return nil
}
