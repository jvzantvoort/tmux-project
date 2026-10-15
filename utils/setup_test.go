package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupSessionDir_NoExec(t *testing.T) {
	// This should not create files, just print
	err := SetupSessionDir(true)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestSetupSessionDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := SetupSessionDir(false); err != nil {
		t.Fatalf("SetupSessionDir() error = %v", err)
	}

	for name, executable := range map[string]bool{"common.rc": false, "tmux_opt_source": true} {
		info, err := os.Stat(filepath.Join(home, ".tmux.d", name))
		if err != nil {
			t.Fatalf("%s not created: %v", name, err)
		}
		if got := info.Mode().Perm()&0100 != 0; got != executable {
			t.Errorf("%s executable = %v, want %v", name, got, executable)
		}
	}

	// an existing common.rc is left alone
	common := filepath.Join(home, ".tmux.d", "common.rc")
	if err := os.WriteFile(common, []byte("custom"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := SetupSessionDir(false); err != nil {
		t.Fatalf("SetupSessionDir() second run error = %v", err)
	}
	if content, _ := os.ReadFile(common); string(content) != "custom" {
		t.Errorf("SetupSessionDir() overwrote common.rc: %q", content)
	}
}
