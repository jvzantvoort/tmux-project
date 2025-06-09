package projecttype

import (
	"os"
	"testing"
)

func TestListProjectTypeConfigs_EmptyDir(t *testing.T) {
	// Set config dir to a temp dir
	old := os.Getenv("CONFIGDIR")
	tmp := t.TempDir()
	os.Setenv("CONFIGDIR", tmp)
	defer os.Setenv("CONFIGDIR", old)

	err := ListProjectTypeConfigs()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Note: This test only covers the empty directory case. More comprehensive tests would require
// refactoring to inject dependencies and avoid side effects (stdout, filesystem).
