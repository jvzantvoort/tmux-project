package main

import (
	"os"
)

type Repository struct {
	Path    string
	SubPath string
	Chapter string
	Info    os.FileInfo
}
