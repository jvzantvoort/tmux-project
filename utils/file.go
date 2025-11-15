package utils

import (
	"io/fs"
	"os"
	"strconv"
)

// GetMode converts an octal string to a file mode
func GetMode(instr string) (fs.FileMode, error) {
	var retv fs.FileMode
	mode, err := strconv.ParseUint(instr, 8, 32)
	if err != nil {
		return retv, err
	}
	retv = os.FileMode(mode)
	return retv, err

}

// TargetExists checks if a file or directory exists at the given path
func TargetExists(target string) bool {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return false
	}
	return true
}

// FileExists checks if a regular file exists at the given path
func FileExists(targetpath string) bool {
	if !TargetExists(targetpath) {
		Debugf("FileExists[%s]: does not exist", targetpath)
		return false
	}

	fi, err := os.Stat(targetpath)
	if err != nil {
		Debugf("FileExists[%s]: Cannot be identified", targetpath)
		return false
	}

	mode := fi.Mode()
	if mode.IsRegular() {
		return true
	} else {
		Debugf("FileExists[%s]: is not a regular file", targetpath)
		return false
	}
}
