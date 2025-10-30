package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMkdirAll(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test", "nested", "directory")

	err := MkdirAll(testPath)
	if err != nil {
		t.Errorf("MkdirAll failed: %v", err)
	}

	// Verify directory was created
	info, err := os.Stat(testPath)
	if err != nil {
		t.Errorf("Directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}
}

func TestMkdirAllExisting(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "existing")

	// Create directory first
	if err := os.Mkdir(testPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Should not error on existing directory
	err := MkdirAll(testPath)
	if err != nil {
		t.Errorf("MkdirAll should not error on existing directory: %v", err)
	}
}

func TestMkdirAllEmpty(t *testing.T) {
	err := MkdirAll("")

	if err == nil {
		t.Error("MkdirAll should return error for empty path")
	}

	if err.Error() != "mkdir called with empty directory" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestMkdirAllFileExists(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "file.txt")

	// Create a file
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Try to create directory with same name
	err := MkdirAll(testFile)

	if err == nil {
		t.Error("MkdirAll should return error when target is a file")
	}

	if err.Error() != "target exists but is not a directory" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestMkdirAllNestedPath(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "a", "b", "c", "d", "e", "f")

	err := MkdirAll(testPath)
	if err != nil {
		t.Errorf("MkdirAll failed for deeply nested path: %v", err)
	}

	// Verify all intermediate directories were created
	checkPath := filepath.Join(tempDir, "a", "b", "c")
	if _, err := os.Stat(checkPath); os.IsNotExist(err) {
		t.Error("Intermediate directories were not created")
	}

	// Verify final directory exists
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Error("Final directory was not created")
	}
}

func TestMkdirAllPermissions(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "testperms")

	err := MkdirAll(testPath)
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	info, err := os.Stat(testPath)
	if err != nil {
		t.Fatalf("Failed to stat directory: %v", err)
	}

	// Check that permissions are 0755
	expectedPerm := os.FileMode(0755)
	actualPerm := info.Mode().Perm()

	if actualPerm != expectedPerm {
		t.Errorf("Expected permissions %v, got %v", expectedPerm, actualPerm)
	}
}

func TestMkdirAllIdempotent(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "idempotent")

	// Call multiple times
	for i := 0; i < 3; i++ {
		err := MkdirAll(testPath)
		if err != nil {
			t.Errorf("MkdirAll call %d failed: %v", i+1, err)
		}
	}

	// Verify directory still exists and is valid
	info, err := os.Stat(testPath)
	if err != nil {
		t.Error("Directory does not exist after multiple calls")
	}

	if !info.IsDir() {
		t.Error("Path is not a directory")
	}
}

func TestMkdirAllRelativePath(t *testing.T) {
	// Save current directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create relative path
	testPath := "relative/nested/path"

	err = MkdirAll(testPath)
	if err != nil {
		t.Errorf("MkdirAll failed for relative path: %v", err)
	}

	// Verify directory was created
	fullPath := filepath.Join(tempDir, testPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("Directory was not created at expected location")
	}
}

func TestMkdirAllSingleDirectory(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "single")

	err := MkdirAll(testPath)
	if err != nil {
		t.Errorf("MkdirAll failed for single directory: %v", err)
	}

	info, err := os.Stat(testPath)
	if err != nil {
		t.Error("Single directory was not created")
	}

	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}
}
