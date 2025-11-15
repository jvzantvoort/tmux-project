package version

import (
	"testing"
	"time"
)

func TestGetVersion(t *testing.T) {
	// Test with default values
	Version = ""
	Commit = ""
	BuildTime = ""

	info := GetVersion()

	if info.Version != "dev" {
		t.Errorf("Expected version 'dev', got '%s'", info.Version)
	}
}

func TestGetVersionWithValues(t *testing.T) {
	// Test with set values
	Version = "v1.2.3"
	Commit = "abc123"
	BuildTime = "2024-01-01T12:00:00Z"

	info := GetVersion()

	if info.Version != "v1.2.3" {
		t.Errorf("Expected version 'v1.2.3', got '%s'", info.Version)
	}

	if info.Commit != "abc123" {
		t.Errorf("Expected commit 'abc123', got '%s'", info.Commit)
	}

	if info.BuildTime != "2024-01-01T12:00:00Z" {
		t.Errorf("Expected build time '2024-01-01T12:00:00Z', got '%s'", info.BuildTime)
	}

	// Check if build date was parsed
	if info.BuildDate.IsZero() {
		t.Error("Expected BuildDate to be parsed, but it's zero")
	}

	expectedDate := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	if !info.BuildDate.Equal(expectedDate) {
		t.Errorf("Expected BuildDate %v, got %v", expectedDate, info.BuildDate)
	}
}

func TestInfoString(t *testing.T) {
	tests := []struct {
		name     string
		info     Info
		expected string
	}{
		{
			name:     "version only",
			info:     Info{Version: "v1.0.0"},
			expected: "v1.0.0",
		},
		{
			name:     "version with commit",
			info:     Info{Version: "v1.0.0", Commit: "abc123"},
			expected: "v1.0.0 (commit: abc123)",
		},
		{
			name:     "version with commit and build time",
			info:     Info{Version: "v1.0.0", Commit: "abc123", BuildTime: "2024-01-01T12:00:00Z"},
			expected: "v1.0.0 (commit: abc123, built: 2024-01-01T12:00:00Z)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.info.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestInfoShort(t *testing.T) {
	info := Info{
		Version:   "v1.2.3",
		Commit:    "abc123",
		BuildTime: "2024-01-01T12:00:00Z",
	}

	if info.Short() != "v1.2.3" {
		t.Errorf("Expected short version 'v1.2.3', got '%s'", info.Short())
	}
}

func TestIsRelease(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{"release version", "v1.0.0", true},
		{"dev version", "dev", false},
		{"devel version", "(devel)", false},
		{"empty version", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := Info{Version: tt.version}
			result := info.IsRelease()
			if result != tt.expected {
				t.Errorf("Expected IsRelease() = %v for version '%s', got %v",
					tt.expected, tt.version, result)
			}
		})
	}
}
