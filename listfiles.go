package tmuxproject

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func ListProject(projectname string) error {

	targets, err := GetSessionPaths(projectname)
	if err != nil {
		return err
	}
	log.Debugf("targets: %d", len(targets))
	_ = ListTargets(targets)
	return nil
}

func ListTargets(targets []string) error {
	for _, target := range targets {
		if !TargetExists(target) {
			log.Errorf("does not exist: %s", target)
			continue
		}
		log.Debugf("target: %s", target)
		// is file a folder?
		fi, err := os.Stat(target)
		if err != nil {
			return err
		}
		mode := fi.Mode()
		if mode.IsRegular() {
			log.Debugf("%s", target)
			fmt.Printf("%s\n", target)
		} else if mode.IsDir() { // folder
			fmt.Printf("%s/\n", target)
			// walk through every file in the folder
			filepath.Walk(target, func(file string, fi os.FileInfo, err error) error {
				// if not a dir, write file content
				if fi.IsDir() {
					log.Debugf("dir %s", file)
				} else {
					log.Debugf("file %s", file)
				}
				return nil
			})
		} else {
			return fmt.Errorf("error: file type not supported")
		}
	}
	return nil
}
