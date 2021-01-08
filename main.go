package tmuxproject

import (
	"path"
	"path/filepath"
	"regexp"

	"github.com/jvzantvoort/tmux-project/config"
	pt "github.com/jvzantvoort/tmux-project/projecttype"
	log "github.com/sirupsen/logrus"
)

var (
	mainconfig = config.NewMainConfig()
)

// NewProjectConfig derives from ProjectTypeConfig and returns an updated
// object with translated values.
func NewProjectConfig(projecttype, projectname string) pt.ProjectTypeConfig {

	ptc := pt.NewProjectTypeConfig(projecttype)

	projtypeconfigdir := path.Join(mainconfig.ProjTypeConfigDir, projecttype)
	projtmplvars := NewProjTmplVars(projectname, ptc)

	ptc.Workdir = projtmplvars.Parse(ptc.Workdir)
	pattern := regexp.MustCompile(ptc.Pattern)
	if pattern.MatchString(projectname) {
		log.Debugf("project name matches pattern")
	} else {
		log.Warningf("project name %s does not matches pattern %s", projectname, ptc.Pattern)
	}

	var err error
	for indx, cfgfile := range ptc.Files {
		// Translate source names
		name := cfgfile.Name
		ptc.Files[indx].Name, err = filepath.Abs(path.Join(projtypeconfigdir, name))
		if err != nil {
			log.Fatal(err)
		}

		// Translate destination names
		dest := cfgfile.Destination
		dest = projtmplvars.Parse(dest)
		ptc.Files[indx].Destination, err = filepath.Abs(path.Join(mainconfig.TmuxDir, dest))
		if err != nil {
			log.Fatal(err)
		}

	}

	for indx, action := range ptc.SetupActions {
		ptc.SetupActions[indx] = projtmplvars.Parse(action)
	}

	return ptc
}

// vim: noexpandtab filetype=go
