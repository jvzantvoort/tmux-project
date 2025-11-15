// Package project provides custom error types and error checking utilities
// for the tmux-project application.
package project

import (
	"errors"
	"syscall"
)

type syscallErrorType = syscall.Errno

var (
	// ErrProjectNotExist indicates that a requested project does not exist
	ErrProjectNotExist = errors.New("project does not exist")
	// ErrProjectTypeNotDefined indicates that no project type was specified when required
	ErrProjectTypeNotDefined = errors.New("project type not defined")
)

// IsProjectNotExist checks if the error indicates that a project does not exist
func IsProjectNotExist(err error) bool {
	return underlyingErrorIs(err, ErrProjectNotExist)
}

// underlyingError returns the underlying error for known os error types
func underlyingError(err error) error {
	return err
}

// underlyingErrorIs checks if the underlying error matches the target error.
// It unwraps specific error types and handles syscall errors for compatibility.
func underlyingErrorIs(err, target error) bool {
	err = underlyingError(err)
	if err == target {
		return true
	}
	e, ok := err.(syscallErrorType)
	return ok && e.Is(target)
}
