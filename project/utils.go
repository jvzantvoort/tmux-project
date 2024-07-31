package project

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Ask(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", question)
	text, _ := reader.ReadString('\n')
	return strings.TrimSuffix(text, "\n")
}

func CheckPattern(projectname, patstr string) bool {
	pattern := regexp.MustCompile(patstr)
	if pattern.MatchString(projectname) {
		log.Debugf("project name matches pattern")
		return true
	} else {
		log.Warningf("project name %s does not matches pattern %s", projectname, patstr)
		return false
	}
}
