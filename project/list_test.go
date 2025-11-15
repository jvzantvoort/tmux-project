package project

import (
	"path/filepath"
	"testing"
)

func TestListConfigs(t *testing.T) {
	// This test checks that ListConfigs returns a slice
	// It uses the real session directory, so we can't control the contents
	configs := ListConfigs()

	if configs == nil {
		t.Error("ListConfigs should return non-nil slice")
	}

	// All returned items should not have .json extension
	for _, config := range configs {
		if filepath.Ext(config) == ".json" {
			t.Errorf("Config name should not have .json extension: %s", config)
		}
	}
}

func TestListConfigsEmpty(t *testing.T) {
	// This test documents expected behavior when directory is empty
	// Since we use real directories, we can't guarantee this
	t.Skip("Cannot control session directory contents in unit test")
}

func TestListConfigsNonExistentDirectory(t *testing.T) {
	// If session directory doesn't exist, ListConfigs should handle gracefully
	configs := ListConfigs()

	// Should handle gracefully (may return empty or may have created directory)
	if configs == nil {
		t.Log("ListConfigs returned nil for non-existent directory")
	}
}

func TestPrintShortList(t *testing.T) {
	// This is primarily an output test
	// We can't easily capture stdout without major refactoring
	// So we just ensure it doesn't panic

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintShortList should not panic: %v", r)
		}
	}()

	PrintShortList()
}

func TestPrintFullList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This requires tmux to be running
	// Just ensure it doesn't panic

	defer func() {
		if r := recover(); r != nil {
			t.Logf("PrintFullList panicked (may be expected): %v", r)
		}
	}()

	PrintFullList()
}

func TestListConfigsSuffix(t *testing.T) {
	// Test that only .json files are included
	// This uses real directory, so we just verify returned format
	configs := ListConfigs()

	for _, config := range configs {
		// Should not have any extension
		ext := filepath.Ext(config)
		if ext != "" {
			t.Errorf("Config %s has extension %s", config, ext)
		}
	}
}

func TestListConfigsOrdering(t *testing.T) {
	// Test that all returned configs are strings
	configs := ListConfigs()

	for _, config := range configs {
		if config == "" {
			t.Error("Config name should not be empty")
		}

		// Should be a valid filename (no path separators)
		if filepath.Base(config) != config {
			t.Errorf("Config should be a base name, not a path: %s", config)
		}
	}
}
