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
	w.Write([]byte(input))
	w.Close()

	result := Ask("What is your name?")
	if result != strings.TrimSpace(input) {
		t.Errorf("expected %q, got %q", strings.TrimSpace(input), result)
	}
}
