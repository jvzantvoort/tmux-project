package utils

import (
	"os"
	"testing"
)

func TestGetMode_Valid(t *testing.T) {
	mode, err := GetMode("0644")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != 0644 {
		t.Errorf("expected 0644, got %v", mode)
	}
}

func TestGetMode_Invalid(t *testing.T) {
	_, err := GetMode("notamode")
	if err == nil {
		t.Errorf("expected error for invalid mode")
	}
}

func TestTargetExists(t *testing.T) {
	file := t.TempDir() + "/exists.txt"
	if err := os.WriteFile(file, []byte("hi"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if !TargetExists(file) {
		t.Errorf("expected file to exist")
	}
	if TargetExists(file + ".nope") {
		t.Errorf("expected file to not exist")
	}
}

func TestFileExists(t *testing.T) {
	file := t.TempDir() + "/exists.txt"
	if err := os.WriteFile(file, []byte("hi"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if !FileExists(file) {
		t.Errorf("expected file to exist")
	}
	dir := t.TempDir()
	if FileExists(dir) {
		t.Errorf("expected dir to not be a file")
	}
}
