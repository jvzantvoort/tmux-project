package main

import (
//	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sort"

	"github.com/jvzantvoort/tmux-project/utils"
)

// DirDepthMap returns a map of directory names to their depth relative to projdir.
func DirDepthMap(projdir string) (map[string]int, error) {
	utils.LogStart()
	defer utils.LogEnd()
	retv := make(map[string]int)

	reldepth := strings.Count(projdir, string(os.PathSeparator))
	utils.Debugf("relative depth %d", reldepth)

	err := filepath.Walk(projdir, func(path string, fi os.FileInfo, inerr error) error {
		err := inerr
		if err != nil {
			utils.Errorf("this passed an error: %q", err)
		}

		curdepth := strings.Count(path, string(os.PathSeparator)) - reldepth - 1 // -1 to not count the project directory itself

		if fi.IsDir() && fi.Name() == ".git" {
			fpath := filepath.Dir(path)
			fpath = strings.TrimPrefix(fpath, projdir)
			fpath = strings.TrimPrefix(fpath, string(os.PathSeparator)) // Remove leading separator
			retv[fpath] = curdepth
		}

		return nil

	})
	if err != nil {
		return nil, err
	}
	return retv, nil
}


// walkAllProjects walks through the project directory and returns a list of all projects found.
// It limits the depth of the search to avoid going too deep into subdirectories.
// The function returns a slice of strings representing the paths to the projects found.
// It uses the filepath.Walk function to traverse the directory structure.
// The maximum depth is controlled by the maxDepth parameter, which is compared against the current depth
// relative to the project directory.
// If the current depth exceeds maxDepth, it skips further exploration of that directory.
// The function returns an error if any issues occur during the directory traversal.
func walkAllProjects(projdir string, maxDepth int) ([]string, error) {
	var retv []string
	utils.LogStart()
	defer utils.LogEnd()
	data, err := DirDepthMap(projdir)

	for k, v := range data {
		md := maxDepth
		if strings.HasPrefix(k, string("work/")) {
			md = maxDepth + 1 // work directories can be one level deeper
		}

		if v <= md {
			retv = append(retv, filepath.Join(projdir, k))
		}
	}

	// sort the list of projects by path
	sort.Strings(retv)

	reldepth := strings.Count(projdir, string(os.PathSeparator))
	utils.Debugf("relative depth %d", reldepth)
	utils.Debugf("maxium depth %d", maxDepth)

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
