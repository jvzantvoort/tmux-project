package sessions

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var tw *tar.Writer

func Walker(file string, fi os.FileInfo, inerr error) error {
	err := inerr
	if err != nil {
		log.Errorf("this passed an error: %q", err)
	}
	// generate tar header
	header, err := tar.FileInfoHeader(fi, file)
	if err != nil {
		return err
	}

	// must provide real name
	// (see https://golang.org/src/archive/tar/common.go?#L626)
	header.Name = filepath.ToSlash(file)

	// write header
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	// if not a dir, write file content
	if !fi.IsDir() {
		data, err := os.Open(file)
		if err != nil {
			return err
		}
		if _, err := io.Copy(tw, data); err != nil {
			return err
		}
		data.Close()
	}
	return nil

}

func MakeTarArchive(buf io.Writer, targets []string) error {

	// tar > gzip > buf
	zr := gzip.NewWriter(buf)
	tw = tar.NewWriter(zr)

	for _, target := range targets {
		if !TargetExists(target) {
			log.Errorf("does not exist: %s", target)
			continue
		}
		log.Debugf("target: %s", target)
		// is file a folder?
		fi, err := os.Stat(target)
		if err != nil {
			return err
		}
		mode := fi.Mode()
		if mode.IsRegular() {
			// get header
			header, err := tar.FileInfoHeader(fi, target)
			if err != nil {
				return err
			}
			header.Name = filepath.ToSlash(target)
			// write header
			if err := tw.WriteHeader(header); err != nil {
				return err
			}
			// get content
			data, err := os.Open(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
			data.Close()
		} else if mode.IsDir() { // folder

			// walk through every file in the folder
			filepath.Walk(target, Walker)
		} else {
			return fmt.Errorf("error: file type not supported")
		}
	}

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}
	return nil
}
