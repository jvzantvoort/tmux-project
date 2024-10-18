package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jvzantvoort/tmux-project/utils"
)

func walkAllProjects(projdir string, maxDepth int) ([]string, error) {
	var retv []string
	utils.LogStart()
	defer utils.LogEnd()

	reldepth := strings.Count(projdir, string(os.PathSeparator))
	utils.Debugf("relative depth %d", reldepth)
	utils.Debugf("maxium depth %d", maxDepth)

	err := filepath.Walk(projdir, func(file string, fi os.FileInfo, inerr error) error {
		err := inerr
		if err != nil {
			utils.Errorf("this passed an error: %q", err)
		}

		curdepth := strings.Count(file, string(os.PathSeparator)) - reldepth

		if curdepth > maxDepth {
			utils.Debugf("depth %d/%d is too high for %s", curdepth, maxDepth, file)
			return fs.SkipDir
		} else {
			utils.Debugf("depth %d/%d is good for %s", curdepth, maxDepth, file)
		}

		if fi.IsDir() && fi.Name() == ".git" {
			retv = append(retv, filepath.Dir(file))
		}

		return nil
	})

	return retv, err
}

func findAllProjects(projdir string, depth int) []ProjectDef {
	var data []ProjectDef
	var retv []ProjectDef
	utils.LogStart()
	defer utils.LogEnd()

	var wg sync.WaitGroup

	dirnames, _ := walkAllProjects(projdir, depth)

	for _, dirname := range dirnames {
		repos := NewProjectDef(projdir, dirname)
		data = append(data, *repos)
	}

	for _, obj := range data {
		wg.Add(1)

		go func(obj ProjectDef) {
			defer wg.Done()
			obj.Init()
			retv = append(retv, obj)
		}(obj)

	}
	// Wait for all goroutines to finish
	wg.Wait()

	return retv

}
