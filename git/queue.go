package git

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/jvzantvoort/tmux-project/utils"
)

var (
	wg sync.WaitGroup
)

func cleanup() {
	if r := recover(); r != nil {
		utils.Errorf("Paniced %s", r)
	}
}

// QueueElement represents a git repository clone operation with destination and branch
type QueueElement struct {
	Url         string
	Basedir     string
	Destination string
	Branch      string
}

// Queue manages concurrent git repository clone operations
type Queue struct {
	Queue []QueueElement
}

// Add adds a repository clone operation to the queue
func (q *Queue) Add(url, basedir, destination, branch string) {
	q.Queue = append(q.Queue, QueueElement{
		Url:         url,
		Basedir:     basedir,
		Destination: destination,
		Branch:      branch,
	})
}

// Run executes all queued git operations concurrently
func (q *Queue) Run() {
	for _, item := range q.Queue {
		wg.Add(1)
		go item.Run()
	}
	wg.Wait()
}

// NewQueue creates a new empty git operation queue
func NewQueue() *Queue {
	return &Queue{}
}

// Run executes a single git clone and checkout operation
func (e QueueElement) Run() {
	defer wg.Done()
	defer cleanup()

	obj := NewGitCmd(e.Basedir)

	err := obj.Clone(e.Url, filepath.Join(e.Basedir, e.Destination))
	if err != nil {
		panic(fmt.Sprintf("Action \"%s\"/\"%s\" failed", e.Url, e.Destination))
	}
	utils.Infof("cloned %s", e.Url)

	err = obj.Checkout(e.Branch)
	if err != nil {
		panic(fmt.Sprintf("Action \"%s\"/\"%s\" failed", e.Url, e.Destination))
	}
	utils.Infof("set branch %s", e.Branch)
}
