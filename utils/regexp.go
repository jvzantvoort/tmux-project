package utils

import (
	"regexp"
)

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
