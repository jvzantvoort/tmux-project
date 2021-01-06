package tmuxproject

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func ArchiveProject(projectname, archivename string) error {

	var buf bytes.Buffer
	targets, err := GetSessionPaths(projectname)
	if err != nil {
		return err
	}
	log.Debugf("targets: %d", len(targets))

	_ = MakeTarArchive(&buf, targets)

	archivename, _ = ExpandHome(archivename)

	fileToWrite, err := os.OpenFile(archivename, os.O_CREATE|os.O_RDWR, os.FileMode(0600))
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(fileToWrite, &buf); err != nil {
		panic(err)
	}

	return nil
}

func TargetExists(target string) bool {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return false
	}
	return true
}

func MakeTarArchive(buf io.Writer, targets []string) error {

	// tar > gzip > buf
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

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
			filepath.Walk(target, func(file string, fi os.FileInfo, err error) error {
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
			})
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
