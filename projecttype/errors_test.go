package projecttype

import "testing"

func TestErrorVars(t *testing.T) {
	const errMsgFmt = "unexpected error string: %s"
	if ErrFileNoExists.Error() != "file does not exist" {
		t.Errorf(errMsgFmt, ErrFileNoExists)
	}
	if ErrProjectNotExists.Error() != "project Type does not exist" {
		t.Errorf(errMsgFmt, ErrProjectNotExists)
	}
	if ErrProjectNameEmpty.Error() != "no project name provided" {
		t.Errorf(errMsgFmt, ErrProjectNameEmpty)
	}
}
