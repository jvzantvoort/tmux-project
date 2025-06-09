package utils

import (
	"testing"
)

func TestSetupSessionDir_NoExec(t *testing.T) {
	// This should not create files, just print
	err := SetupSessionDir(true)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Note: Full test for SetupSessionDir with file creation would require
// mocking config.SessionDir and messages.Copy, which is not trivial without
// refactoring for testability. This test covers the noexec branch.
