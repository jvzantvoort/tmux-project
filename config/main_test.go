package config

import (
	"os"
	"os/user"
	"path"
	"testing"
)

func TestHome(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Fatalf("Failed to get current user: %v", err)
	}

	home := Home()
	if home != usr.HomeDir {
		t.Errorf("Home() = %s, want %s", home, usr.HomeDir)
	}
}

func TestConfigDir(t *testing.T) {
	usr, _ := user.Current()
	expected := path.Join(usr.HomeDir, ".tmux-project")

	configDir := ConfigDir()
	if configDir != expected {
		t.Errorf("ConfigDir() = %s, want %s", configDir, expected)
	}
}

func TestSessionDir(t *testing.T) {
	usr, _ := user.Current()
	expected := path.Join(usr.HomeDir, ".tmux.d")

	sessionDir := SessionDir()
	if sessionDir != expected {
		t.Errorf("SessionDir() = %s, want %s", sessionDir, expected)
	}
}

func TestConfigDirIntegration(t *testing.T) {
	// Test that ConfigDir is a subdirectory of Home
	home := Home()
	configDir := ConfigDir()

	if len(configDir) <= len(home) {
		t.Errorf("ConfigDir should be longer than Home")
	}

	if configDir[:len(home)] != home {
		t.Errorf("ConfigDir should start with Home path")
	}
}

func TestSessionDirIntegration(t *testing.T) {
	// Test that SessionDir is a subdirectory of Home
	home := Home()
	sessionDir := SessionDir()

	if len(sessionDir) <= len(home) {
		t.Errorf("SessionDir should be longer than Home")
	}

	if sessionDir[:len(home)] != home {
		t.Errorf("SessionDir should start with Home path")
	}
}

func TestConfigDirEnvOverride(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Test with custom HOME
	testHome := "/tmp/test-home"
	os.Setenv("HOME", testHome)

	// Note: This test may not work as expected since Home() uses user.Current()
	// which may not respect the HOME environment variable on all systems
	// This is more of a documentation test
	configDir := ConfigDir()
	if configDir == "" {
		t.Error("ConfigDir should not be empty")
	}
}
