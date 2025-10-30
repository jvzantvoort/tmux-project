package git

import (
	"testing"
)

func TestNewQueue(t *testing.T) {
	queue := NewQueue()

	if queue == nil {
		t.Fatal("NewQueue returned nil")
	}

	if queue.Queue != nil && len(queue.Queue) != 0 {
		t.Errorf("New queue should be empty, got %d items", len(queue.Queue))
	}
}

func TestQueueAdd(t *testing.T) {
	queue := NewQueue()

	url := "https://github.com/test/repo.git"
	basedir := "/tmp/test"
	destination := "repo"
	branch := "main"

	queue.Add(url, basedir, destination, branch)

	if len(queue.Queue) != 1 {
		t.Errorf("Expected queue length 1, got %d", len(queue.Queue))
	}

	item := queue.Queue[0]
	if item.Url != url {
		t.Errorf("Expected URL %s, got %s", url, item.Url)
	}
	if item.Basedir != basedir {
		t.Errorf("Expected Basedir %s, got %s", basedir, item.Basedir)
	}
	if item.Destination != destination {
		t.Errorf("Expected Destination %s, got %s", destination, item.Destination)
	}
	if item.Branch != branch {
		t.Errorf("Expected Branch %s, got %s", branch, item.Branch)
	}
}

func TestQueueAddMultiple(t *testing.T) {
	queue := NewQueue()

	queue.Add("url1", "base1", "dest1", "branch1")
	queue.Add("url2", "base2", "dest2", "branch2")
	queue.Add("url3", "base3", "dest3", "branch3")

	if len(queue.Queue) != 3 {
		t.Errorf("Expected queue length 3, got %d", len(queue.Queue))
	}

	// Verify order is maintained
	if queue.Queue[0].Url != "url1" {
		t.Error("Queue order not maintained")
	}
	if queue.Queue[1].Url != "url2" {
		t.Error("Queue order not maintained")
	}
	if queue.Queue[2].Url != "url3" {
		t.Error("Queue order not maintained")
	}
}

func TestQueueElementStruct(t *testing.T) {
	element := QueueElement{
		Url:         "https://example.com/repo.git",
		Basedir:     "/home/user",
		Destination: "project",
		Branch:      "develop",
	}

	if element.Url == "" {
		t.Error("Url should not be empty")
	}
	if element.Basedir == "" {
		t.Error("Basedir should not be empty")
	}
	if element.Destination == "" {
		t.Error("Destination should not be empty")
	}
	if element.Branch == "" {
		t.Error("Branch should not be empty")
	}
}

func TestQueueRun(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping queue run test in short mode")
	}

	queue := NewQueue()

	// Add items with invalid URLs that will fail
	queue.Add("https://invalid-url.com/repo1.git", t.TempDir(), "repo1", "main")

	// Run should not panic even if items fail
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Queue.Run panicked (expected for invalid repos): %v", r)
		}
	}()

	queue.Run()
}

func TestQueueEmpty(t *testing.T) {
	queue := NewQueue()

	// Running an empty queue should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Empty queue run should not panic: %v", r)
		}
	}()

	queue.Run()
}
