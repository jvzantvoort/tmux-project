package sessions

import (
	"os"
	"sort"
	"strings"

	"github.com/jvzantvoort/tmux-project/config"
	log "github.com/sirupsen/logrus"
)

type TmuxSessions struct {
	Sessions []TmuxSession
}

// Define a type for a slice of slices of strings
type SliceOfStringSlices [][]string

// Implement the sort.Interface for the type
func (s SliceOfStringSlices) Len() int {
	return len(s)
}

func (s SliceOfStringSlices) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SliceOfStringSlices) Less(i, j int) bool {
	// Compare the first elements of the sub-slices (strings)
	return s[i][0] < s[j][0]
}

func NewTmuxSessions() *TmuxSessions {
	retv := &TmuxSessions{}

	sess, err := ListSessions(config.SessionDir())
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range sess {
		retv.Sessions = append(retv.Sessions, item)
	}
	return retv

}

func ListSessions(inputdir string) (map[string]TmuxSession, error) {
	retv := make(map[string]TmuxSession)
	targets, err := os.ReadDir(inputdir)
	if err != nil {
		log.Fatal(err)
		return retv, err
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

		retv[target_name] = *NewTmuxSession(target_name)

	}
	return retv, nil

}

func ListTmuxConfigs(inputdir string) ([][]string, error) {
	retv := [][]string{}

	sessiondata, err := ListSessions(inputdir)

	for _, session := range sessiondata {
		sane := "false"
		if session.IsSane() {
			sane = "true"
		}
		cols := []string{session.Name, session.Description, session.Workdir, sane}
		retv = append(retv, cols)
	}
	sort.Sort(SliceOfStringSlices(retv))

	return retv, err
}
