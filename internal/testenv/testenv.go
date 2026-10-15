// Package testenv provides helpers for tests.
package testenv

import (
	"fmt"
	"os"
	"testing"
)

// IsolateHome runs the tests of a package with $HOME pointing at a fresh
// temporary directory, so tests never read or write the real ~/.tmux.d
// and ~/.tmux-project. Use it from TestMain:
//
//	func TestMain(m *testing.M) { os.Exit(testenv.IsolateHome(m)) }
func IsolateHome(m *testing.M) int {
	home, err := os.MkdirTemp("", "tmux-project-test-home-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot create test home: %s\n", err)
		return 1
	}
	defer func() { _ = os.RemoveAll(home) }()

	if err := os.Setenv("HOME", home); err != nil {
		fmt.Fprintf(os.Stderr, "cannot set HOME: %s\n", err)
		return 1
	}
	return m.Run()
}
