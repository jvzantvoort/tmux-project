package tmux

import (
	"os"
	"testing"
)

func TestListActive(t *testing.T) {
	// This test will only work if tmux is available
	sessions, err := ListActive()

	// Check if tmux is available
	if err != nil && os.Getenv("TMUX") == "" {
		t.Skip("tmux not available or not running in tmux session")
	}

	if sessions == nil {
		t.Error("ListActive should return non-nil slice even on error")
	}

	// Sessions can be empty if no tmux sessions are running
	t.Logf("Found %d active sessions", len(sessions))
}

func TestListActiveFormat(t *testing.T) {
	sessions, err := ListActive()

	if err != nil && os.Getenv("TMUX") == "" {
		t.Skip("tmux not available")
	}

	// Check that session names don't contain quotes
	for _, session := range sessions {
		if len(session) > 0 && (session[0] == '"' || session[len(session)-1] == '"') {
			t.Errorf("Session name should not contain quotes: %s", session)
		}
	}
}

// TestNew and TestAttach are difficult to test without actually running tmux
// These are integration tests that would require a running tmux server

func TestNewWithInvalidConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test documents expected behavior
	// New() will call utils.ErrorExit on failure, which exits the program
	// In a real scenario, this would need to be refactored to return errors
	t.Log("New() and Attach() call utils.ErrorExit on failure")
}

func TestAttachWithInvalidConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Log("Attach() calls utils.ErrorExit on failure")
}

func TestResumeLogic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Resume logic:
	// 1. Get list of active sessions
	// 2. If session exists in active list, attach
	// 3. Otherwise, create new session

	t.Log("Resume() combines ListActive(), Attach(), and New()")
}

func TestListActiveEmpty(t *testing.T) {
	// When no sessions are running, should return empty slice
	sessions, err := ListActive()

	if err != nil && os.Getenv("TMUX") == "" {
		t.Skip("tmux not available")
	}

	if sessions == nil {
		t.Error("ListActive should return empty slice, not nil")
	}
}

func TestListActiveNoDuplicates(t *testing.T) {
	sessions, err := ListActive()

	if err != nil && os.Getenv("TMUX") == "" {
		t.Skip("tmux not available")
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, session := range sessions {
		if seen[session] {
			t.Errorf("Duplicate session found: %s", session)
		}
		seen[session] = true
	}
}
