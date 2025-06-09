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

// QueueElement represents an item in the queue
type QueueElement struct {
	Url         string
	Basedir     string
	Destination string
	Branch      string
}

// Queue represents a queue of items
type Queue struct {
	Queue []QueueElement
}

// Add adds an item to the queue
func (q *Queue) Add(url, basedir, destination, branch string) {
	q.Queue = append(q.Queue, QueueElement{
		Url:         url,
		Basedir:     basedir,
		Destination: destination,
		Branch:      branch,
	})
}

// Run runs the queue, waiting for all items to finish
func (q *Queue) Run() {
	for _, item := range q.Queue {
		wg.Add(1)
		go item.Run()
	}
	wg.Wait()
}

// NewQueue creates a new queue
func NewQueue() *Queue {
	return &Queue{}
}

// Run runs an item
func (e QueueElement) Run() {
	defer wg.Done() // lower counter
	defer cleanup() // handle panics

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
