package tmuxproject

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	mainconfig = config.NewMainConfig()
)

func NewProjectTypeConfig(projecttype, projectname string) config.ProjectTypeConfig {
	projtypeconfigdir := path.Join(mainconfig.ProjTypeConfigDir, projecttype)
	log.Debugf("project type config dir: %s", projtypeconfigdir)
	log.Debugf("tmux dir: %s", mainconfig.TmuxDir)

	var configuration config.ProjectTypeConfig
	viper.SetConfigName("config")
	viper.AddConfigPath(projtypeconfigdir)

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
	configuration.Workdir = projtmplvars.Parse(configuration.Workdir)

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
		configuration.Files[indx].Destination, err = filepath.Abs(path.Join(mainconfig.TmuxDir, dest))
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
