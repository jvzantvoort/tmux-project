package projecttype

import (
	"embed"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/jvzantvoort/tmux-project/config"
)

var (
	mainconfig = config.NewMainConfig()
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

func (ptc *ProjectTypeConfig) readConfig(projtypeconfigdir string) {

	viper.SetConfigName("config")
	viper.AddConfigPath(projtypeconfigdir)
	// viper.AddConfigPath(tp.MasterConfigDir)

	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
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

func (ptc *ProjectTypeConfig) Init(projtypeconfigdir, projecttype string) error {

	log.Debugf("Init Start: %s", projecttype)
	projtypeconfigdir = path.Join(projtypeconfigdir, projecttype)

	ptc.ProjectType = projecttype
	ptc.ProjectTypeDir = projtypeconfigdir

	if ptc.Exists(ptc.ProjectTypeDir) {
		return fmt.Errorf("directory already exists: %s", ptc.ProjectTypeDir)
	}

	if err := os.MkdirAll(ptc.ProjectTypeDir, os.FileMode(int(0755))); err != nil {
		return fmt.Errorf("directory cannot be created: %s", ptc.ProjectTypeDir)
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
	projtypeconfigdir := path.Join(mainconfig.ProjTypeConfigDir, projecttype)

	log.Debugf("project type config dir: %s", projtypeconfigdir)
	log.Debugf("tmux dir: %s", mainconfig.TmuxDir)

	v := ProjectTypeConfig{}
	v.readConfig(projtypeconfigdir)

	var err error
	v.Workdir, err = mainconfig.ExpandHome(v.Workdir)
	if err != nil {
		log.Errorf("%q", err)
	}

	return v
}

func CreateProjectType(projecttype string) error {
	var pt ProjectTypeConfig
	pt.Init(mainconfig.ProjTypeConfigDir, projecttype)
	return nil
}
