package tmuxproject

import (
	"path"
)

func GetDescription(sessionname string) (string, error) {
	tmuxdir := GetTmuxDir()
	target := path.Join(tmuxdir, sessionname+".rc")

	lines, err := LoadStringLines(target)
	if err != nil {
		return "", err
	}

	matches := GetMatches(`^#\s+DESCRIPTION\:\s*(?P<description>.*)\s*$`, lines)
	return matches["description"], nil

}

func GetWorkdir(sessionname string) (string, error) {
	tmuxdir := GetTmuxDir()
	target := path.Join(tmuxdir, sessionname+".rc")

	lines, err := LoadStringLines(target)
	if err != nil {
		return "", err
	}

	matches := GetMatches(`^#\s+WORKDIR\:\s*(?P<workdir>.*)\s*$`, lines)
	return matches["workdir"], nil

}
