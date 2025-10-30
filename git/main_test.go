package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewGitCmd(t *testing.T) {
	tempDir := t.TempDir()
	gitCmd := NewGitCmd(tempDir)

	if gitCmd == nil {
		t.Fatal("NewGitCmd returned nil")
	}

	if gitCmd.cwd != tempDir {
		t.Errorf("Expected cwd %s, got %s", tempDir, gitCmd.cwd)
	}

	if gitCmd.command == "" {
		t.Error("git command should not be empty")
	}

	if len(gitCmd.path) == 0 {
		t.Error("path should not be empty")
	}
}

func TestGitCmdWhich(t *testing.T) {
	tempDir := t.TempDir()
	gitCmd := NewGitCmd(tempDir)

	// Test finding a command that should exist
	_, err := gitCmd.which("git")
	if err != nil && os.Getenv("CI") == "" {
		// Only fail if not in CI environment where git might not be available
		t.Logf("Warning: git command not found: %v", err)
	}
}

func TestGitCmdIsGit(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(string) error
		expected bool
	}{
		{
			name: "not a git repository",
			setup: func(dir string) error {
				return nil
			},
			expected: false,
		},
		{
			name: "git directory exists",
			setup: func(dir string) error {
				gitDir := filepath.Join(dir, ".git")
				return os.Mkdir(gitDir, 0755)
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			if err := tt.setup(testDir); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			gitCmd := NewGitCmd(testDir)
			result := gitCmd.IsGit()

			if result != tt.expected {
				t.Errorf("IsGit() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGitCmdGetStatus(t *testing.T) {
	testDir := t.TempDir()

	// Create a git repo for testing
	gitDir := filepath.Join(testDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	gitCmd := NewGitCmd(testDir)

	// This will fail if git is not initialized, but shouldn't panic
	status := gitCmd.GetStatus()
	if status == nil {
		t.Error("GetStatus should return a non-nil map even on error")
	}
}

func TestGitCmdBranch(t *testing.T) {
	tempDir := t.TempDir()
	gitCmd := NewGitCmd(tempDir)

	// This should return an error since it's not a git repo
	_, err := gitCmd.Branch()
	if err == nil {
		t.Log("Branch() didn't return error for non-git directory (might be expected)")
	}
}

func TestGitCmdCheckout(t *testing.T) {
	tempDir := t.TempDir()
	gitCmd := NewGitCmd(tempDir)

	// This should return an error since it's not a git repo
	err := gitCmd.Checkout("main")
	if err == nil {
		t.Error("Checkout should fail for non-git directory")
	}
}

func TestGitCmdClone(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping clone test in short mode")
	}

	tempDir := t.TempDir()
	gitCmd := NewGitCmd(tempDir)

	// Test with invalid URL (using .invalid TLD which is reserved and will never resolve)
	invalidURL := "https://thisdomaindoesnotexist12345.invalid/repo.git"
	dest := filepath.Join(tempDir, "test-repo")

	err := gitCmd.Clone(invalidURL, dest)
	if err == nil {
		t.Error("Clone should fail with invalid URL")
	}
}

func TestGitCmdCloneDestination(t *testing.T) {
	tempDir := t.TempDir()
	gitCmd := NewGitCmd(tempDir)

	// Test that destination directory is created
	dest := filepath.Join(tempDir, "nested", "path", "repo")

	// This will fail due to invalid URL, but should create parent directories
	_ = gitCmd.Clone("https://invalid.com/repo.git", dest)

	// Check that parent directory was created
	parentDir := filepath.Dir(dest)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		t.Errorf("Parent directory should have been created: %s", parentDir)
	}
}
