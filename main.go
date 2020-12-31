package tmuxproject

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"text/template"

	"github.com/jvzantvoort/tmux-project/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ProjTmplVars struct {
	HomeDir            string
	ProjectDescription string
	ProjectDir         string
	ProjectName        string
}

func NewProjTmplVars(projectname string, conf config.ProjectTypeConfig) *ProjTmplVars {

	v := &ProjTmplVars{}
	v.HomeDir = GetHomeDir()
	v.ProjectDir = conf.Workdir
	v.ProjectName = projectname

	return v
}

// buildConfig construct the text from the template definition and arguments.
func (t ProjTmplVars) Parse(templatestring string) string {
	tmpl, err := template.New("prompt").Parse(templatestring)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, t)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func LoadFile(target string, tmplvars ProjTmplVars) (string, error) {
	var retv string
	content, err := ioutil.ReadFile(target)
	if err != nil {
		return "", err
	}
	retv = tmplvars.Parse(string(content))
	return retv, nil
}

func DescribeProjectType(config config.ProjectTypeConfig) {
	log.Debugf("Describe: %s", config.ProjectType)
	log.Debugf("  Workdir: %s", config.Workdir)
	log.Debugf("  Pattern: %s", config.Pattern)

	fileno := len(config.Files)
	actionsno := len(config.SetupActions)

	if fileno > 0 {
		log.Debugf("  Files:")
		for _, act := range config.Files {
			log.Debugf("    - name: %s", act.Name)
			log.Debugf("      destination: %s", act.Destination)
			log.Debugf("      mode: %s", act.Mode)
		}
	}

	if actionsno > 0 {
		log.Debugf("  Actions:")
		for _, act := range config.SetupActions {
			log.Debugf("    - %s", act)
		}

	}

	log.Debugf("Describe: %s, end", config.ProjectType)
}

func GetProjectTypeConfig(configname, projectname string) config.ProjectTypeConfig {
	projtypeconfigdir := GetProjTypeCofigDir()
	tmuxdir := GetTmuxDir()

	var configuration config.ProjectTypeConfig
	viper.SetConfigName("config")
	viper.AddConfigPath(projtypeconfigdir)
	viper.AddConfigPath(MasterConfigDir)

	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	configuration.Workdir, _ = ExpandHome(configuration.Workdir)
	projtmplvars := NewProjTmplVars(projectname, configuration)

	for indx, cfgfile := range configuration.Files {
		// Translate source names
		name := cfgfile.Name
		configuration.Files[indx].Name, err = filepath.Abs(path.Join(projtypeconfigdir, name))
		if err != nil {
			log.Fatal(err)
		}

		// Translate destination names
		dest := cfgfile.Destination
		dest = projtmplvars.Parse(dest)
		configuration.Files[indx].Destination, err = filepath.Abs(path.Join(tmuxdir, dest))
		if err != nil {
			log.Fatal(err)
		}

	}

	for indx, action := range configuration.SetupActions {
		configuration.SetupActions[indx] = projtmplvars.Parse(action)
	}

	return configuration
}

// vim: noexpandtab filetype=go
