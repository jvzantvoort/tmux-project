package utils

import (
	"os"
	"strings"
	"testing"
)

func TestExpand_Home(t *testing.T) {
	home, _ := os.UserHomeDir()
	p, err := Expand("~/testfile")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(p, home) {
		t.Errorf("expected prefix %q, got %q", home, p)
	}
}

func TestExpand_NoTilde(t *testing.T) {
	p, err := Expand("/tmp/testfile")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != "/tmp/testfile" {
		t.Errorf("expected /tmp/testfile, got %q", p)
	}
}

func TestExpand_UserSpecific(t *testing.T) {
	_, err := Expand("~otheruser/file")
	if err == nil {
		t.Errorf("expected error for user-specific home dir")
	}
}
