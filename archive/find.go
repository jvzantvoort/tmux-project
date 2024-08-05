package archive

import (
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/utils"
)

// FindFiles walks through the given targets and returns a slice of all found non-directory targets.
func FindFiles(targets []string) ([]string, map[string]string, error) {
	utils.LogStart()
	defer utils.LogEnd()

	var files []string
	links := make(map[string]string)

	for _, target := range targets {
		err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.Mode()&os.ModeSymlink != 0 {
				linkTarget, err := os.Readlink(path)
				if err != nil {
					return err
				}
				links[path] = linkTarget
			} else if !info.IsDir() {
				utils.Debugf("add target: %s", path)
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
	}

	return files, links, nil
}
