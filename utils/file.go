package utils

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	log "github.com/sirupsen/logrus"
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

func ReplaceBytesToBytes(content []byte, term, replacement string) []byte {
	return []byte(strings.Replace(string(content), term, replacement, -1))
}

func TargetExists(target string) bool {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return false
	}
	return true
}

func FileExists(targetpath string) bool {
	if !TargetExists(targetpath) {
		log.Debugf("FileExists[%s]: does not exist", targetpath)
		return false
	}

	// is file a folder?
	fi, err := os.Stat(targetpath)
	if err != nil {
		log.Debugf("FileExists[%s]: Cannot be identified", targetpath)
		return false
	}

	mode := fi.Mode()
	if mode.IsRegular() {
		return true
	} else {
		log.Debugf("FileExists[%s]: is not a regular file", targetpath)
		return false
	}
}

func LoadFile(targetpath string) ([]string, error) {
	var retv []string

	if !FileExists(targetpath) {
		return retv, fmt.Errorf("file %s does not exists", targetpath)
	}

	content, err := os.ReadFile(targetpath)
	if err != nil {
		return retv, err
	}

	retv = append(retv, strings.Split(string(content), "\n")...)
	return retv, nil

}
