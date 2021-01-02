package tmuxproject

import (
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetTmuxConfigFile(sessionname string) string {
	return path.Join(mainconfig.TmuxDir, sessionname+".rc")
}

func GetTmuxEnvFile(sessionname string) string {
	return path.Join(mainconfig.TmuxDir, sessionname+".env")
}

func GetDescription(sessionname string) (string, error) {
	target := GetTmuxConfigFile(sessionname)

	lines, err := LoadStringLines(target)
	if err != nil {
		return "", err
	}

	matches := GetMatches(`^#\s+DESCRIPTION\:\s*(?P<description>.*)\s*$`, lines)
	return strings.TrimSuffix(matches["description"], "\n"), nil

}

func GetWorkdir(sessionname string) (string, error) {
	target := GetTmuxConfigFile(sessionname)

	lines, err := LoadStringLines(target)
	if err != nil {
		return "", err
	}

	matches := GetMatches(`^#\s+WORKDIR\:\s*(?P<workdir>.*)\s*$`, lines)
	return strings.TrimSuffix(matches["workdir"], "\n"), nil

}

func GetSessionPaths(sessionname string) ([]string, error) {
	var targets []string

	rcfile := GetTmuxConfigFile(sessionname)
	envfile := GetTmuxEnvFile(sessionname)

	workdir, err := GetWorkdir(sessionname)
	if err != nil {
		return targets, err
	}

	targets = append(targets, workdir)
	targets = append(targets, rcfile)
	targets = append(targets, envfile)

	for indx, target := range targets {
		targetn, nerr := ExpandHome(target)
		if nerr != nil {
			log.Fatalf("%q", nerr)
		}
		targets[indx] = targetn
	}
	return targets, nil

}
