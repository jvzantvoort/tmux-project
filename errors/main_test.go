package project

import (
	"fmt"
	"testing"
)

func TestIsProjectNotExist(t *testing.T) {
	if !IsProjectNotExist(ErrProjectNotExist) {
		t.Errorf("Test for ErrProjectNotExist failed")
	}
}

func TestUnderlyingError(t *testing.T) {
	err := fmt.Errorf("error")
	retv := underlyingError(err)
	if err != retv {
		t.Errorf("Expected %q but got %q", err, retv)
	}
}

