package utils

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func MkdirAll(targetpath string) error {
	if len(targetpath) == 0 {
		return fmt.Errorf("mkdir called with empty directory")
	}
	prefix := fmt.Sprintf("mkdir %s", targetpath)

	log.Debugf("%s, start", prefix)
	defer log.Debugf("%s, end", prefix)

	if target, err := os.Stat(targetpath); !os.IsNotExist(err) {
		if !target.IsDir() {
			return fmt.Errorf("%s, exists but is not a directory", prefix)
		}
		log.Debugf("%s, already exists", prefix)
		return nil
	}

	err := os.MkdirAll(targetpath, os.FileMode(int(0755)))

	if err != nil {
		return fmt.Errorf("%s, failed: %s", prefix, err)
	}
	return nil

}
