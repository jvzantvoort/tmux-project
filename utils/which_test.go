package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWhich_FindsExecutable(t *testing.T) {
	// Create a temp dir and fake executable
	dir := t.TempDir()
	exe := "mycmd"
	exePath := filepath.Join(dir, exe)
	os.WriteFile(exePath, []byte("#!/bin/sh\necho hi\n"), 0755)

	// Prepend temp dir to PATH
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", dir+":"+oldPath)
	defer os.Setenv("PATH", oldPath)

	found := Which(exe)
	if found != exePath {
		t.Errorf("expected %q, got %q", exePath, found)
	}
}

func TestWhich_NotFound(t *testing.T) {
	name := "definitelynotfoundcmd"
	found := Which(name)
	if found != name {
		t.Errorf("expected fallback to input, got %q", found)
	}
}
