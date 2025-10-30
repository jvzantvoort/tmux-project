package utils

import (
	"os"
	"strings"
	"testing"
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

	if queue.Queue != nil && len(queue.Queue) != 0 {
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
	if testing.Short() {
		t.Skip("Skipping queue run test in short mode")
	}

	queue := NewQueue()
	tempDir := t.TempDir()

	queue.Add(tempDir, "echo test1")
	queue.Add(tempDir, "echo test2")

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Queue.Run should not panic: %v", r)
		}
	}()

	queue.Run()
}

func TestQueueRunWithFailure(t *testing.T) {
	t.Skip("Queue.Run with failure calls ErrorExit which terminates the program")
}

func TestQueueElementRun(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping element run test in short mode")
	}

	tempDir := t.TempDir()

	element := QueueElement{
		Cwd:     tempDir,
		Command: "echo test",
	}

	// Should not panic for successful command
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("QueueElement.Run should not panic for successful command: %v", r)
		}
	}()

	element.Run()
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
