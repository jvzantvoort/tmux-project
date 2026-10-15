package config

import (
	"path/filepath"
	"testing"
)

func TestHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if got := Home(); got != home {
		t.Errorf("Home() = %s, want %s", got, home)
	}
}

func TestConfigDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	want := filepath.Join(home, ".tmux-project")
	if got := ConfigDir(); got != want {
		t.Errorf("ConfigDir() = %s, want %s", got, want)
	}
}

func TestSessionDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	want := filepath.Join(home, ".tmux.d")
	if got := SessionDir(); got != want {
		t.Errorf("SessionDir() = %s, want %s", got, want)
	}
}

func TestHomeFallback(t *testing.T) {
	// without $HOME the home directory comes from the user database
	t.Setenv("HOME", "")

	if got := Home(); got == "" {
		t.Error("Home() returned an empty string without $HOME")
	}
}
