package projecttype

import (
	"errors"
)

var (
	// ErrFileNoExists indicates that a required file does not exist
	ErrFileNoExists = errors.New("file does not exist")
	// ErrProjectNotExists indicates that the requested project type does not exist
	ErrProjectNotExists = errors.New("project Type does not exist")
	// ErrProjectNameEmpty indicates that no project name was provided when required
	ErrProjectNameEmpty = errors.New("no project name provided")
)
