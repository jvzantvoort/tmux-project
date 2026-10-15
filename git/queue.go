package git

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/jvzantvoort/tmux-project/utils"
)

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

// Run executes all queued git operations concurrently and waits for them
// to finish. It returns the errors of all failed operations joined.
func (q *Queue) Run() error {
	var wg sync.WaitGroup
	errs := make([]error, len(q.Queue))

	for indx, item := range q.Queue {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[indx] = item.Run()
		}()
	}
	wg.Wait()

	return errors.Join(errs...)
}

// NewQueue creates a new empty git operation queue
func NewQueue() *Queue {
	return &Queue{}
}

// Run executes a single git clone and, when a branch is set, checks out
// that branch in the new clone.
func (e QueueElement) Run() error {
	destination := filepath.Join(e.Basedir, e.Destination)

	obj := NewGitCmd(e.Basedir)
	if err := obj.Clone(e.Url, destination); err != nil {
		return fmt.Errorf("clone %s to %s failed: %w", e.Url, destination, err)
	}
	utils.Infof("cloned %s", e.Url)

	if len(e.Branch) == 0 {
		return nil
	}

	// Clone has pointed obj at the new clone
	if err := obj.Checkout(e.Branch); err != nil {
		return fmt.Errorf("checkout of branch %s in %s failed: %w", e.Branch, destination, err)
	}
	utils.Infof("set branch %s", e.Branch)
	return nil
}
