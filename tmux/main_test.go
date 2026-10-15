package tmux

import (
	"os"
	"reflect"
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

func TestSessionName(t *testing.T) {
	tests := map[string]string{
		"plain":    "plain",
		"foo.bar":  "foo_bar",
		"a:b.c":    "a_b_c",
		"JIRA-123": "JIRA-123",
	}
	for input, want := range tests {
		if got := SessionName(input); got != want {
			t.Errorf("SessionName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestResumeCommands(t *testing.T) {
	tests := []struct {
		name       string
		exists     bool
		insideTmux bool
		directory  string
		want       [][]string
	}{
		{
			name: "outside tmux, new session",
			want: [][]string{{"-f", "p.rc", "new-session", "-s", "p", "-c", "/work/p"}},
		},
		{
			name:   "outside tmux, existing session",
			exists: true,
			want:   [][]string{{"-f", "p.rc", "attach-session", "-d", "-t", "=p"}},
		},
		{
			name:       "inside tmux, new session",
			insideTmux: true,
			want: [][]string{
				{"-f", "p.rc", "new-session", "-d", "-s", "p", "-c", "/work/p"},
				{"switch-client", "-t", "=p"},
			},
		},
		{
			name:       "inside tmux, existing session",
			exists:     true,
			insideTmux: true,
			want:       [][]string{{"switch-client", "-t", "=p"}},
		},
		{
			name:      "unknown project has no start directory",
			directory: "-",
			want:      [][]string{{"-f", "p.rc", "new-session", "-s", "p"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			directory := "/work/p"
			if tt.directory == "-" {
				directory = ""
			}
			got := resumeCommands("p", "p.rc", directory, tt.exists, tt.insideTmux)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resumeCommands() = %q, want %q", got, tt.want)
			}
		})
	}
}
