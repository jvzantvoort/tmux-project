package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		command     string
		expectError bool
	}{
		{
			name:        "simple echo",
			command:     "echo hello",
			expectError: false,
		},
		{
			name:        "command with output",
			command:     "echo test output",
			expectError: false,
		},
		{
			name:        "pwd command",
			command:     "pwd",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := Exec(tempDir, tt.command)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if stdout == nil {
				t.Error("stdout slice should not be nil")
			}

			if stderr == nil {
				t.Error("stderr slice should not be nil")
			}
		})
	}
}

func TestExecOutput(t *testing.T) {
	tempDir := t.TempDir()

	stdout, stderr, err := Exec(tempDir, "echo test123")

	if err != nil {
		t.Errorf("echo command failed: %v", err)
	}

	if len(stdout) == 0 {
		t.Error("Expected output from echo command")
	}

	found := false
	for _, line := range stdout {
		if strings.Contains(line, "test123") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'test123' in output")
	}

	if len(stderr) > 0 {
		t.Logf("stderr: %v", stderr)
	}
}

func TestExecSilent(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		command     string
		expectError bool
	}{
		{
			name:        "simple echo",
			command:     "echo hello",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := ExecSilent(tempDir, tt.command)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if stdout == nil {
				t.Error("stdout slice should not be nil")
			}

			if stderr == nil {
				t.Error("stderr slice should not be nil")
			}
		})
	}
}

func TestExecWorkingDirectory(t *testing.T) {
	tempDir := t.TempDir()

	stdout, _, err := Exec(tempDir, "pwd")

	if err != nil {
		t.Fatalf("pwd command failed: %v", err)
	}

	if len(stdout) == 0 {
		t.Fatal("Expected output from pwd")
	}

	// The output should contain the temp directory path
	found := false
	for _, line := range stdout {
		if strings.Contains(line, tempDir) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected pwd output to contain %s, got: %v", tempDir, stdout)
	}
}

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

	cwd := "/tmp/test"
	command := "echo test"

	queue.Add(cwd, command)

	if len(queue.Queue) != 1 {
		t.Errorf("Expected queue length 1, got %d", len(queue.Queue))
	}

	item := queue.Queue[0]
	if item.Cwd != cwd {
		t.Errorf("Expected Cwd %s, got %s", cwd, item.Cwd)
	}
	if item.Command != command {
		t.Errorf("Expected Command %s, got %s", command, item.Command)
	}
}

func TestQueueRun(t *testing.T) {
	queue := NewQueue()
	tempDir := t.TempDir()

	queue.Add(tempDir, "echo test1")
	queue.Add(tempDir, "echo test2")

	if err := queue.Run(); err != nil {
		t.Errorf("Queue.Run() error = %v", err)
	}
}

func TestQueueRunInOrder(t *testing.T) {
	queue := NewQueue()
	tempDir := t.TempDir()

	// each step depends on the previous one
	queue.Add(tempDir, "mkdir step1")
	queue.Add(tempDir, "mkdir step1/step2")
	queue.Add(tempDir, "echo done > step1/step2/result")

	if err := queue.Run(); err != nil {
		t.Fatalf("Queue.Run() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "step1", "step2", "result"))
	if err != nil {
		t.Fatalf("result file missing: %v", err)
	}
	if strings.TrimSpace(string(content)) != "done" {
		t.Errorf("result = %q, want %q", content, "done")
	}
}

func TestQueueRunWithFailure(t *testing.T) {
	queue := NewQueue()
	tempDir := t.TempDir()

	queue.Add(tempDir, "false")
	queue.Add(tempDir, "touch should-not-exist")

	if err := queue.Run(); err == nil {
		t.Fatal("Queue.Run() expected an error")
	}

	if _, err := os.Stat(filepath.Join(tempDir, "should-not-exist")); !os.IsNotExist(err) {
		t.Error("Queue.Run() ran an action after a failed one")
	}
}

func TestQueueElementRun(t *testing.T) {
	tempDir := t.TempDir()

	element := QueueElement{
		Cwd:     tempDir,
		Command: "echo 'quoted argument' | tr a-z A-Z",
	}

	if err := element.Run(); err != nil {
		t.Errorf("QueueElement.Run() error = %v", err)
	}
}

func TestExecShell(t *testing.T) {
	stdout, _, err := ExecShell(t.TempDir(), "echo 'hello   world' && echo second")
	if err != nil {
		t.Fatalf("ExecShell() error = %v", err)
	}
	want := []string{"hello   world", "second"}
	if strings.Join(stdout, "|") != strings.Join(want, "|") {
		t.Errorf("ExecShell() stdout = %q, want %q", stdout, want)
	}
}

func TestExecLargeStderr(t *testing.T) {
	// more stderr output than a pipe buffer holds, written before stdout
	done := make(chan error, 1)
	go func() {
		_, stderr, err := ExecSilent(t.TempDir(), "sh -c head${IFS}-c${IFS}300000${IFS}/dev/zero>&2;echo${IFS}ok")
		if err == nil && len(stderr) == 0 {
			err = fmt.Errorf("no stderr captured")
		}
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("ExecSilent() error = %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("ExecSilent() blocked on a large stderr stream")
	}
}

func TestExecCommandParsing(t *testing.T) {
	tempDir := t.TempDir()

	// Test that commands are properly split
	stdout, _, err := Exec(tempDir, "echo hello world")

	if err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	found := false
	for _, line := range stdout {
		if strings.Contains(line, "hello") && strings.Contains(line, "world") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Command arguments not properly parsed")
	}
}

func TestExecSilentErrorHandling(t *testing.T) {
	tempDir := t.TempDir()

	// Test with command that returns non-zero exit code
	_, _, err := ExecSilent(tempDir, "sh -c 'exit 1'")

	if err == nil {
		t.Error("Expected error for non-zero exit code")
	}
}

func TestExecEnvironment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping environment test in short mode")
	}

	tempDir := t.TempDir()

	// Create a test file in tempDir
	testFile := tempDir + "/testfile.txt"
	os.WriteFile(testFile, []byte("test"), 0644)

	// Run ls in the temp directory
	stdout, _, err := Exec(tempDir, "ls")

	if err != nil {
		t.Fatalf("ls command failed: %v", err)
	}

	// Should see our test file
	found := false
	for _, line := range stdout {
		if strings.Contains(line, "testfile.txt") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Command not executed in correct directory")
	}
}
