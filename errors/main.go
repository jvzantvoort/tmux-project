package project

import (
	"errors"
	"syscall"
)

type syscallErrorType = syscall.Errno

var (
	ErrProjectNotExist       = errors.New("project does not exist")
	ErrProjectTypeNotDefined = errors.New("project type not defined")
)

func IsProjectNotExist(err error) bool {
	return underlyingErrorIs(err, ErrProjectNotExist)
}

// underlyingError returns the underlying error for known os error types.
func underlyingError(err error) error {
	// switch err := err.(type) {
	// case *PathError:
	// 	return err.Err
	// case *LinkError:
	// 	return err.Err
	// case *SyscallError:
	// 	return err.Err
	// }
	return err
}

func underlyingErrorIs(err, target error) bool {
	// Note that this function is not errors.Is:
	// underlyingError only unwraps the specific error-wrapping types
	// that it historically did, not all errors implementing Unwrap().
	err = underlyingError(err)
	if err == target {
		return true

	}
	// To preserve prior behavior, only examine syscall errors.
	e, ok := err.(syscallErrorType)
	return ok && e.Is(target)
}
