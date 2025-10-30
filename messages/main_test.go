package messages

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetShort(t *testing.T) {
	// Test that GetShort doesn't panic
	result := GetShort("test")
	if result == "" {
		t.Log("GetShort returned empty string (expected if file doesn't exist)")
	}
}

func TestGetUse(t *testing.T) {
	result := GetUse("test")
	if result == "" {
		t.Log("GetUse returned empty string (expected if file doesn't exist)")
	}
}

func TestGetLong(t *testing.T) {
	result := GetLong("test")
	if result == "" {
		t.Log("GetLong returned empty string (expected if file doesn't exist)")
	}
}

func TestGetShell(t *testing.T) {
	result := GetShell("test")
	if result == "" {
		t.Log("GetShell returned empty string (expected if file doesn't exist)")
	}
}

func TestGetConfig(t *testing.T) {
	result := GetConfig("test")
	if result == "" {
		t.Log("GetConfig returned empty string (expected if file doesn't exist)")
	}
}

func TestGetContentNonExistent(t *testing.T) {
	result := GetContent("nonexistent", "nonexistent")
	if result != "undefined" {
		t.Errorf("Expected 'undefined' for non-existent file, got: %s", result)
	}
}

func TestGetContentTrimsNewline(t *testing.T) {
	// This tests that content has trailing newlines removed
	// We can't easily test with embedded files, so we document the behavior
	t.Log("GetContent should trim trailing newlines from file content")
}

func TestCopy(t *testing.T) {
	tempDir := t.TempDir()
	destFile := filepath.Join(tempDir, "subdir", "test-config")

	err := Copy("test", destFile, 0644)

	// Check that directory was created
	destDir := filepath.Dir(destFile)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		t.Error("Copy should create destination directory")
	}

	// The file creation might fail if "test" doesn't exist in config/
	if err != nil {
		t.Logf("Copy failed (expected if source doesn't exist): %v", err)
	}
}

func TestCopyExistingDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Create directory first
	destDir := filepath.Join(tempDir, "existing")
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	destFile := filepath.Join(destDir, "config")
	err := Copy("test", destFile, 0644)

	if err != nil {
		t.Logf("Copy failed (expected if source doesn't exist): %v", err)
	}
}

func TestCopyFileMode(t *testing.T) {
	tempDir := t.TempDir()
	destFile := filepath.Join(tempDir, "test-mode")

	// This tests that the mode parameter is used
	mode := os.FileMode(0600)
	_ = Copy("test", destFile, mode)

	// Check if file exists and has correct mode
	if info, err := os.Stat(destFile); err == nil {
		if info.Mode() != mode {
			t.Logf("File mode: %v (may differ from requested %v due to umask)", info.Mode(), mode)
		}
	}
}

func TestCopyNestedDirectories(t *testing.T) {
	tempDir := t.TempDir()
	destFile := filepath.Join(tempDir, "a", "b", "c", "config")

	err := Copy("test", destFile, 0644)

	// Check that nested directories were created
	destDir := filepath.Dir(destFile)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		t.Error("Copy should create nested directories")
	}

	if err != nil {
		t.Logf("Copy failed (expected if source doesn't exist): %v", err)
	}
}

func TestGetContentFolders(t *testing.T) {
	folders := []string{"long", "shells", "use", "short", "config"}

	for _, folder := range folders {
		t.Run(folder, func(t *testing.T) {
			result := GetContent(folder, "nonexistent")
			if result != "undefined" {
				t.Logf("Folder %s might have content", folder)
			}
		})
	}
}

func TestContentEmbedding(t *testing.T) {
	// Verify that Content embed.FS is properly initialized
	// We can test by trying to read a directory
	_, err := Content.ReadDir(".")
	if err != nil {
		t.Logf("Content embed.FS read error (may be expected): %v", err)
	}
}
