package projecttype

import "testing"

func TestErrorVars(t *testing.T) {
	if ErrFileNoExists.Error() != "file does not exist" {
		t.Errorf("unexpected error string: %s", ErrFileNoExists)
	}
	if ErrProjectNotExists.Error() != "project Type does not exist" {
		t.Errorf("unexpected error string: %s", ErrProjectNotExists)
	}
	if ErrProjectNameEmpty.Error() != "no project name provided" {
		t.Errorf("unexpected error string: %s", ErrProjectNameEmpty)
	}
}
