package sessions

import (
	"os"
	"path"
	"strings"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/utils"
)

type TmuxSession struct {
	Name        string
	Configfile  string
	Environment string
	Workdir     string
	Description string
}

func (tm *TmuxSession) LoadConfig() {
	var err error
	var config_lines []string
	config_lines, err = utils.LoadFile(tm.Configfile)
	if err != nil {
		utils.Errorf("%q", err)
	}

	if len(config_lines) != 0 {
		config_matches := utils.GetMatches(`^#\s+DESCRIPTION\:\s*(?P<description>.*)\s*$`, config_lines)
		env_matches := utils.GetMatches(`^#\s+WORKDIR\:\s*(?P<workdir>.*)\s*$`, config_lines)

		tm.Description = strings.TrimSuffix(config_matches["description"], "\n")
		tm.Workdir = strings.TrimSuffix(env_matches["workdir"], "\n")
		tm.Workdir, _ = utils.Expand(tm.Workdir)
	}
}

func (tm TmuxSession) TargetPaths() (targets []string) {
	targets = append(targets, tm.Configfile)
	targets = append(targets, tm.Environment)
	targets = append(targets, tm.Workdir+"/")
	return
}

func (tm TmuxSession) IsSane() bool {

	_, err := os.Stat(tm.Workdir)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func NewTmuxSession(sessionname string) *TmuxSession {
	tm := &TmuxSession{Name: sessionname}
	tm.Configfile = path.Join(config.SessionDir(), tm.Name+".rc")
	tm.Environment = path.Join(config.SessionDir(), tm.Name+".env")

	tm.LoadConfig()

	return tm
}
