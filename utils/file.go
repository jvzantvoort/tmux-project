package utils

import (
	"io/fs"
	"os"
	"strconv"
)

func GetMode(instr string) (fs.FileMode, error) {
	var retv fs.FileMode
	mode, err := strconv.ParseUint(instr, 8, 32)
	if err != nil {
		return retv, err
	}
	retv = os.FileMode(mode)
	return retv, err

}

//
// func ReplaceBytesToBytes(content []byte, term, replacement string) []byte {
// 	return []byte(strings.Replace(string(content), term, replacement, -1))
// }

func TargetExists(target string) bool {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return false
	}
	return true
}

func FileExists(targetpath string) bool {
	if !TargetExists(targetpath) {
		Debugf("FileExists[%s]: does not exist", targetpath)
		return false
	}

	// is file a folder?
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
