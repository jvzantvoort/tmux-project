package project

import (
	"fmt"
	"sync"

	"github.com/jvzantvoort/tmux-project/projecttype"
	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
)

var (
	wg sync.WaitGroup
)

func cleanup() {
	if r := recover(); r != nil {
		log.Errorf("Paniced %s", r)
	}
}

func RunSetupAction(workdir, action string) {
	defer wg.Done() // lower counter
	defer cleanup() // handle panics

	stdout_list, stderr_list, eerror := utils.Exec(workdir, action)
	for _, stdout_line := range stdout_list {
		log.Infof("<stdout> %s", stdout_line)
	}
	for _, stderr_line := range stderr_list {
		log.Errorf("<stderr> %s", stderr_line)
	}
	if eerror != nil {
		panic(fmt.Sprintf("Action \"%s\" failed", action))
	}
}

func RunSetupActions(projconf projecttype.ProjectTypeConfig) {
	for _, action := range projconf.SetupActions {
		wg.Add(1)
		go RunSetupAction(projconf.Workdir, action)
	}
	wg.Wait()
}
