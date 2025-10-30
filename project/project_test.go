package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewProject(t *testing.T) {
	projectName := "test-project"
	proj := NewProject(projectName)

	if proj == nil {
		t.Fatal("NewProject returned nil")
	}

	if proj.Name != projectName {
		t.Errorf("Expected project name %s, got %s", projectName, proj.Name)
	}
}

func TestProjectNameIsValid(t *testing.T) {
	tests := []struct {
		name     string
		projName string
		pattern  string
		expected bool
	}{
		{
			name:     "alphanumeric matches pattern",
			projName: "test123",
			pattern:  "^[a-z0-9]+$",
			expected: true,
		},
		{
			name:     "uppercase fails lowercase pattern",
			projName: "Test123",
			pattern:  "^[a-z0-9]+$",
			expected: false,
		},
		{
			name:     "special chars match when allowed",
			projName: "test-project_123",
			pattern:  "^[a-z0-9_-]+$",
			expected: true,
		},
		{
			name:     "any string matches wildcard",
			projName: "anything goes!",
			pattern:  ".*",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proj := NewProject(tt.projName)
			proj.Pattern = tt.pattern

			result := proj.NameIsValid()
			if result != tt.expected {
				t.Errorf("NameIsValid() = %v, want %v (name: %s, pattern: %s)",
					result, tt.expected, tt.projName, tt.pattern)
			}
		})
	}
}

func TestProjectNameIsValidEmptyPattern(t *testing.T) {
	proj := NewProject("test")
	proj.Pattern = ""

	// Empty pattern should match empty string at start
	result := proj.NameIsValid()
	t.Logf("NameIsValid with empty pattern: %v", result)
}

func TestProjectInjectExternal(t *testing.T) {
	proj := NewProject("test")
	proj.InjectExternal()

	// Check that home directory was set
	if proj.HomeDir == "" {
		t.Error("HomeDir should not be empty after InjectExternal")
	}

	// Check that GOARCH was set
	if proj.GOARCH == "" {
		t.Error("GOARCH should not be empty after InjectExternal")
	}

	// Check that GOOS was set
	if proj.GOOS == "" {
		t.Error("GOOS should not be empty after InjectExternal")
	}

	// GOPATH might be empty on some systems, so we just check it's defined
	t.Logf("GOPATH: %s", proj.GOPATH)

	// Check that USER was set
	if proj.USER == "" {
		t.Error("USER should not be empty after InjectExternal")
	}
}

func TestProjectSetDescription(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "single word",
			input:    []string{"description"},
			expected: "description",
		},
		{
			name:     "multiple words",
			input:    []string{"this", "is", "a", "test"},
			expected: "this is a test",
		},
		{
			name:     "with extra spaces",
			input:    []string{"  test  ", "  description  "},
			expected: "test     description", // joins with space, outer trimmed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proj := NewProject("test")
			proj.SetDescription(tt.input...)

			// strings.Join adds single space, then TrimSpace removes outer whitespace
			if proj.Description != tt.expected {
				t.Errorf("Expected description %q, got %q", tt.expected, proj.Description)
			}
		})
	}
}

func TestProjectSetDescriptionEmpty(t *testing.T) {
	// This test documents behavior when empty description is provided
	// In actual code, it would prompt for input via utils.Ask
	// We can't easily test interactive input, so we skip this
	t.Skip("SetDescription with empty input requires interactive input")
}

func TestProjectNameIsValidPatternCompilation(t *testing.T) {
	proj := NewProject("test")

	// Test various regex patterns
	patterns := []string{
		"^[a-z]+$",
		"^[A-Z][a-z]*$",
		"^\\w+$",
		"^[a-zA-Z0-9_-]+$",
		".*",
	}

	for _, pattern := range patterns {
		proj.Pattern = pattern
		// Should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("NameIsValid panicked with pattern %s: %v", pattern, r)
			}
		}()
		_ = proj.NameIsValid()
	}
}

func TestProjectInjectExternalConsistency(t *testing.T) {
	proj1 := NewProject("test1")
	proj2 := NewProject("test2")

	proj1.InjectExternal()
	proj2.InjectExternal()

	// Both projects should get the same external values
	if proj1.HomeDir != proj2.HomeDir {
		t.Error("HomeDir should be consistent across projects")
	}

	if proj1.GOARCH != proj2.GOARCH {
		t.Error("GOARCH should be consistent across projects")
	}

	if proj1.GOOS != proj2.GOOS {
		t.Error("GOOS should be consistent across projects")
	}

	if proj1.USER != proj2.USER {
		t.Error("USER should be consistent across projects")
	}
}

func TestProjectInjectExternalHomeDir(t *testing.T) {
	proj := NewProject("test")
	proj.InjectExternal()

	// Verify HomeDir is an absolute path
	if !filepath.IsAbs(proj.HomeDir) {
		t.Errorf("HomeDir should be absolute path, got: %s", proj.HomeDir)
	}

	// Verify HomeDir exists
	if _, err := os.Stat(proj.HomeDir); os.IsNotExist(err) {
		t.Errorf("HomeDir should exist: %s", proj.HomeDir)
	}
}

func TestNewProjectEmptyName(t *testing.T) {
	proj := NewProject("")

	if proj == nil {
		t.Fatal("NewProject should not return nil even with empty name")
	}

	if proj.Name != "" {
		t.Error("Project name should be empty string")
	}
}

func TestProjectInjectExternalGOPATH(t *testing.T) {
	proj := NewProject("test")
	proj.InjectExternal()

	// GOPATH can be empty or set
	if proj.GOPATH != "" {
		// If set, it should be an absolute path
		if !filepath.IsAbs(proj.GOPATH) {
			t.Errorf("GOPATH should be absolute if set, got: %s", proj.GOPATH)
		}
	}
}
