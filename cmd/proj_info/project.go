package main

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/jvzantvoort/tmux-project/utils"
)

func walkAllProjects(projdir string) ([]string, error) {
	var retv []string
	utils.LogStart()
	defer utils.LogEnd()

	err := filepath.Walk(projdir, func(file string, fi os.FileInfo, inerr error) error {
		err := inerr
		if err != nil {
			utils.Errorf("this passed an error: %q", err)
		}
		if fi.IsDir() && fi.Name() == ".git" {
			retv = append(retv, filepath.Dir(file))
		}
		return nil
	})

	return retv, err
}

func findAllProjects(projdir string) []ProjectDef {
	var data []ProjectDef
	var retv []ProjectDef
	utils.LogStart()
	defer utils.LogEnd()

	var wg sync.WaitGroup

	dirnames, _ := walkAllProjects(projdir)

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
