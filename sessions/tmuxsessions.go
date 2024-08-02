package sessions

import (
	"os"
	"strings"

	"github.com/jvzantvoort/tmux-project/config"
	log "github.com/sirupsen/logrus"
)

type TmuxSessions struct {
	Sessions []TmuxSession
}

func NewTmuxSessions() *TmuxSessions {
	tmux_sessions := &TmuxSessions{}

	// list ~/.tmux.d
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
