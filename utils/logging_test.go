package utils

import (
	"testing"
)

func TestLogIfError(t *testing.T) {
	// Should not panic or log for nil
	LogIfError(nil)
	// Should log for non-nil (no assertion, just coverage)
	LogIfError("an error")
}

func TestLogStartEnd(t *testing.T) {
	// Should not panic
	LogStart()
	LogEnd()
}
