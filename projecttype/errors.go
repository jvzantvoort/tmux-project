package projecttype

import (
	"errors"
)

var (
	ErrFileNoExists     = errors.New("file does not exist")
	ErrProjectNotExists = errors.New("project Type does not exist")
	ErrProjectNameEmpty = errors.New("no project name provided")
)
