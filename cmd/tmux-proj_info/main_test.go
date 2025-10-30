package main

import (
	"os"
	"strings"
	"testing"
)

// TestProjectDef_Initialization verifies ProjectDef struct creation
func TestProjectDef_Initialization(t *testing.T) {
	pd := NewProjectDef("/tmp/project", "/tmp/project/subdir")

	if pd == nil {
		t.Fatal("NewProjectDef returned nil")
	}
	if pd.Path != "/tmp/project/subdir" {
		t.Errorf("Path mismatch: got %s, want /tmp/project/subdir", pd.Path)
	}
	if pd.ProjectDir != "/tmp/project" {
		t.Errorf("ProjectDir mismatch: got %s, want /tmp/project", pd.ProjectDir)
	}
}

// TestProjectDef_AllFieldsAccessible verifies all fields are accessible
func TestProjectDef_AllFieldsAccessible(t *testing.T) {
	pd := &ProjectDef{
		Name:       "test-project",
		Path:       "/tmp/test",
		Branch:     "main",
		AbsPath:    "/tmp/test",
		ProjectDir: "/tmp",
		SubPath:    "test",
		Chapter:    "main",
		Expected:   true,
	}

	// Verify all fields are accessible
	if pd.Name != "test-project" {
		t.Error("Name not accessible")
	}
	if pd.Path != "/tmp/test" {
		t.Error("Path not accessible")
	}
	if pd.Branch != "main" {
		t.Error("Branch not accessible")
	}
	if pd.AbsPath != "/tmp/test" {
		t.Error("AbsPath not accessible")
	}
	if pd.ProjectDir != "/tmp" {
		t.Error("ProjectDir not accessible")
	}
	if pd.SubPath != "test" {
		t.Error("SubPath not accessible")
	}
	if pd.Chapter != "main" {
		t.Error("Chapter not accessible")
	}
	if !pd.Expected {
		t.Error("Expected not accessible")
	}
}

// TestProjectDef_GetFields verifies GetFields method
func TestProjectDef_GetFields(t *testing.T) {
	pd := &ProjectDef{
		Name:   "test",
		Branch: "main",
		Status: make(map[string]int),
	}

	fields := pd.GetFields()

	if len(fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(fields))
	}
	if fields[0] != "test" {
		t.Errorf("First field should be name, got %s", fields[0])
	}
}

// TestProjectDef_GetFieldsWithStatus verifies GetFields with git status
func TestProjectDef_GetFieldsWithStatus(t *testing.T) {
	pd := &ProjectDef{
		Name:   "test-status",
		Branch: "develop",
		Status: map[string]int{
			"M": 2,
			"A": 1,
		},
	}

	fields := pd.GetFields()

	if len(fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(fields))
	}

	// Second field should contain status
	statusField := fields[1]
	if !strings.Contains(statusField, "[") {
		t.Error("Status field should be bracketed")
	}
}

// TestDirDepthMap_EmptyDirectory verifies handling of empty directory
func TestDirDepthMap_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	depthMap, err := DirDepthMap(tmpDir)
	if err != nil {
		t.Fatalf("DirDepthMap failed: %v", err)
	}

	// Empty directory should have empty map (no .git directories)
	if len(depthMap) != 0 {
		t.Errorf("Expected empty map for directory without .git, got %d entries", len(depthMap))
	}
}

// TestDirDepthMap_WithGitDir verifies detection of .git directory
func TestDirDepthMap_WithGitDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a .git directory
	gitDir := tmpDir + "/.git"
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	depthMap, err := DirDepthMap(tmpDir)
	if err != nil {
		t.Fatalf("DirDepthMap failed: %v", err)
	}

	// Should detect the root git directory
	if len(depthMap) == 0 {
		t.Error("Should have detected .git directory")
	}
}

// TestDirDepthMap_NonexistentDirectory verifies error handling
func TestDirDepthMap_NonexistentDirectory(t *testing.T) {
	_, err := DirDepthMap("/nonexistent/directory/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}
}

// TestProjectDef_BranchColors verifies branch coloring logic
func TestProjectDef_BranchColors(t *testing.T) {
	testCases := []struct {
		branch   string
		expected string
	}{
		{"main", "main"},
		{"master", "master"},
		{"develop", "develop"},
		{"feature/test", "feature/test"},
	}

	for _, tc := range testCases {
		pd := &ProjectDef{
			Name:   "test",
			Branch: tc.branch,
			Status: make(map[string]int),
		}

		fields := pd.GetFields()
		// The branch should be in the third field (index 2)
		if len(fields) < 3 {
			t.Errorf("Expected at least 3 fields for branch %s", tc.branch)
			continue
		}

		// Just verify the field contains the branch name (color codes might be present)
		branchField := fields[2]
		if !strings.Contains(branchField, tc.branch) {
			t.Errorf("Branch field should contain %s, got %s", tc.branch, branchField)
		}
	}
}

// TestProjectDef_StatusMap verifies status map handling
func TestProjectDef_StatusMap(t *testing.T) {
	pd := &ProjectDef{
		Name:   "status-test",
		Branch: "main",
	}

	// Initially nil status
	pd.Status = make(map[string]int)
	pd.Status["M"] = 5
	pd.Status["A"] = 3
	pd.Status["D"] = 1

	if len(pd.Status) != 3 {
		t.Errorf("Expected 3 status entries, got %d", len(pd.Status))
	}
	if pd.Status["M"] != 5 {
		t.Errorf("Modified count should be 5, got %d", pd.Status["M"])
	}
	if pd.Status["A"] != 3 {
		t.Errorf("Added count should be 3, got %d", pd.Status["A"])
	}
	if pd.Status["D"] != 1 {
		t.Errorf("Deleted count should be 1, got %d", pd.Status["D"])
	}
}

// TestNewProjectDef_SamePath verifies handling when paths are equal
func TestNewProjectDef_SamePath(t *testing.T) {
	pd := NewProjectDef("/tmp/project", "/tmp/project")

	if pd == nil {
		t.Fatal("NewProjectDef returned nil")
	}
	if pd.Path != "/tmp/project" {
		t.Errorf("Path should be /tmp/project, got %s", pd.Path)
	}
	if pd.ProjectDir != "/tmp/project" {
		t.Errorf("ProjectDir should be /tmp/project, got %s", pd.ProjectDir)
	}
}


