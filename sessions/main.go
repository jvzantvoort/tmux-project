package sessions

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
)

type TmuxSession struct {
	Name        string
	Configfile  string
	Environment string
	Workdir     string
	Description string
}

type TmuxSessions struct {
	Sessions []TmuxSession
}

func (tm *TmuxSession) LoadConfig() {
	var err error
	var config_lines []string
	config_lines, err = utils.LoadFile(tm.Configfile)
	if err != nil {
		log.Errorf("%q", err)
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

func (tm TmuxSession) Archive(archivename string) error {

	var buf bytes.Buffer
	targets := tm.TargetPaths()

	log.Debugf("targets: %d", len(targets))

	_ = MakeTarArchive(&buf, targets)

	archivename, _ = utils.Expand(archivename)

	fileToWrite, err := os.OpenFile(archivename, os.O_CREATE|os.O_RDWR, os.FileMode(0600))
	utils.ErrorExit(err)

	_, err = io.Copy(fileToWrite, &buf)
	utils.ErrorExit(err)

	return nil

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

func NewTmuxSessions() *TmuxSessions {
	tmux_sessions := &TmuxSessions{}

	targets, err := os.ReadDir(config.SessionDir())
	if err != nil {
		log.Fatal(err)
	}

	for _, target := range targets {
		target_name := target.Name()

		// we only want the session names
		if !strings.HasSuffix(target_name, ".rc") {
			continue
		}

		// "common" is shared by all others
		if target_name == "common.rc" {
			continue
		}

		target_name = strings.TrimSuffix(target_name, ".rc")

		x := NewTmuxSession(target_name)
		tmux_sessions.Sessions = append(tmux_sessions.Sessions, *x)

	}

	return tmux_sessions

}
