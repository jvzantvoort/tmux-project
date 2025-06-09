package utils

import (
	"os"
	"strings"
	"testing"
)

func TestAsk_SimulateInput(t *testing.T) {
	// Save original stdin
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	input := "answer\n"
	if _, err := w.Write([]byte(input)); err != nil {
		t.Fatalf("failed to write to pipe: %v", err)
	}
	w.Close()

	result := Ask("What is your name?")
	if result != strings.TrimSpace(input) {
		t.Errorf("expected %q, got %q", strings.TrimSpace(input), result)
	}
}
