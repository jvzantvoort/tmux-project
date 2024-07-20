package utils

import (
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Which(command string) string {
	Path := strings.Split(os.Getenv("PATH"), ":")
	var retv string
	for _, dirname := range Path {
		fullpath := path.Join(dirname, command)
		log.Debugf("test path: %s", fullpath)

		_, err := os.Stat(fullpath)
		if err == nil {
			if os.IsNotExist(err) {
				continue
			}
			retv = fullpath
		}
	}
	return retv

}
