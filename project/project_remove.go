package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	"github.com/jvzantvoort/tmux-project/utils"
)

// sessionFileSuffixes are the per-project files kept in the session
// directory. They are removed even when they are not listed as targets,
// e.g. for projects created before targets were stored.
var sessionFileSuffixes = []string{".json", ".rc", ".env", ".chapters.yml"}

// RemovalTargets returns the paths that Remove deletes: the project
// directory, the target files and the session files of the project.
// Paths that do not exist are left out.
func (proj Project) RemovalTargets() []string {
	utils.LogStart()
	defer utils.LogEnd()

	candidates := proj.ListFiles()
	for _, suffix := range sessionFileSuffixes {
		candidates = append(candidates, filepath.Join(config.SessionDir(), proj.Name+suffix))
	}

	retv := []string{}
	seen := map[string]bool{}
	for _, target := range candidates {
		if len(target) == 0 {
			continue
		}
		target = filepath.Clean(target)
		if seen[target] {
			continue
		}
		seen[target] = true
		if _, err := os.Lstat(target); err != nil {
			continue
		}
		retv = append(retv, target)
	}
	return retv
}

// CheckRemovableDir returns an error when dir is not a sane project
// directory to remove recursively: it must be absolute and must not be
// the filesystem root, the home directory or one of its parents.
func CheckRemovableDir(dir, home string) error {
	if len(dir) == 0 {
		return nil
	}
	if !filepath.IsAbs(dir) {
		return fmt.Errorf("refusing to remove relative directory %q", dir)
	}
	dir = filepath.Clean(dir)
	if dir == filepath.Dir(dir) {
		return fmt.Errorf("refusing to remove root directory %q", dir)
	}
	if len(home) != 0 {
		// home is dir itself or lives below it
		rel, err := filepath.Rel(dir, filepath.Clean(home))
		if err == nil && (rel == "." || filepath.IsLocal(rel)) {
			return fmt.Errorf("refusing to remove %q: it contains the home directory", dir)
		}
	}
	return nil
}

// Remove deletes the project directory, its target files and its session
// files. It keeps going after a failure and returns the first error.
func (proj Project) Remove() error {
	utils.LogStart()
	defer utils.LogEnd()

	home, _ := os.UserHomeDir()
	if err := CheckRemovableDir(proj.Directory, home); err != nil {
		return err
	}

	return RemovePaths(proj.RemovalTargets())
}

// RemovePaths removes each path recursively, logging every removal. It
// keeps going after a failure and returns the first error.
func RemovePaths(paths []string) error {
	var rerr error
	for _, target := range paths {
		if err := os.RemoveAll(target); err != nil {
			utils.Errorf("failed to remove %s: %s", target, err)
			if rerr == nil {
				rerr = err
			}
			continue
		}
		utils.Infof("removed %s", target)
	}
	return rerr
}
