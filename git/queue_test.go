package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestNewQueue(t *testing.T) {
	queue := NewQueue()

	if queue == nil {
		t.Fatal("NewQueue returned nil")
	}

	if len(queue.Queue) != 0 {
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
	queue.Add("https://invalid-url.invalid/repo1.git", t.TempDir(), "repo1", "main")

	if err := queue.Run(); err == nil {
		t.Error("Queue.Run() expected an error for an invalid repository")
	}
}

func TestQueueEmpty(t *testing.T) {
	queue := NewQueue()

	if err := queue.Run(); err != nil {
		t.Errorf("Queue.Run() on empty queue error = %v", err)
	}
}

// gitRun runs git in dir for test setup
func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func TestQueueRunChecksOutBranchInClone(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// source repository with a "feature" branch next to the default one
	source := filepath.Join(t.TempDir(), "source")
	if err := os.MkdirAll(source, 0750); err != nil {
		t.Fatal(err)
	}
	gitRun(t, source, "init", "-q", "-b", "main")
	gitRun(t, source, "commit", "-q", "--allow-empty", "-m", "initial")
	gitRun(t, source, "branch", "feature")

	projdir := t.TempDir()

	queue := NewQueue()
	queue.Add(source, projdir, "with-branch", "feature")
	queue.Add(source, projdir, "without-branch", "")

	if err := queue.Run(); err != nil {
		t.Fatalf("Queue.Run() error = %v", err)
	}

	for dest, want := range map[string]string{"with-branch": "feature", "without-branch": "main"} {
		branch, err := NewGitCmd(filepath.Join(projdir, dest)).Branch()
		if err != nil {
			t.Fatalf("Branch() of %s error = %v", dest, err)
		}
		if branch != want {
			t.Errorf("%s is on branch %q, want %q", dest, branch, want)
		}
	}
}
