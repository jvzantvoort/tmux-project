package sessions

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func TargetExists(target string) bool {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return false
	}
	return true
}

func FileExists(targetpath string) bool {
	if !TargetExists(targetpath) {
		log.Debugf("FileExists[%s]: does not exist", targetpath)
		return false
	}

	// is file a folder?
	fi, err := os.Stat(targetpath)
	if err != nil {
		log.Debugf("FileExists[%s]: Cannot be identified", targetpath)
		return false
	}

	mode := fi.Mode()
	if mode.IsRegular() {
		return true
	} else {
		log.Debugf("FileExists[%s]: is not a regular file", targetpath)
		return false
	}
}

func LoadFile(targetpath string) ([]string, error) {
	var retv []string

	if !FileExists(targetpath) {
		return retv, fmt.Errorf("file %s does not exists", targetpath)
	}

	content, err := os.ReadFile(targetpath)
	if err != nil {
		return retv, err
	}

	retv = append(retv, strings.Split(string(content), "\n")...)
	return retv, nil

}

// GetMatches parses each line and sticks the findings in a map
func GetMatches(regEx string, lines []string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)

	paramsMap = make(map[string]string)

	for _, line := range lines {
		match := compRegEx.FindStringSubmatch(line)
		for i, name := range compRegEx.SubexpNames() {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}
		}
	}
	return
}
