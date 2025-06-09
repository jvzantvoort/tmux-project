package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/utils"
)

// TarArchive encapsulates the tar.gz creation logic
type TarArchive struct {
	OutputFile string
	Targets    []string
	Links      map[string]string
}

// NewTarArchive creates a new TarArchive
func NewTarArchive(OutputFile string) *TarArchive {
	retv := &TarArchive{OutputFile: OutputFile}
	retv.Links = make(map[string]string)
	return retv
}

// AddFile adds a file to the list of files to be included in the tar.gz archive
func (t *TarArchive) AddFile(filePath string) {
	t.Targets = append(t.Targets, filePath)
}

// AddSymlink adds a symbolic link to the tar.gz archive
func (t *TarArchive) AddSymlink(linkName, target string) {
	t.Links[linkName] = target
}

// AddFiles adds a list of files to the tar.gz archive
func (t *TarArchive) AddFiles(paths []string) {
	targets, links, err := FindFiles(paths)
	utils.ErrorExit(err)
	for _, target := range targets {
		t.AddFile(target)
	}

	for linkName, target := range links {
		t.AddSymlink(linkName, target)
	}
}

// CreateArchive creates the tar.gz archive with the added files
func (t *TarArchive) CreateArchive() error {
	utils.LogStart()
	defer utils.LogEnd()
	// Create the output file
	outFile, err := os.Create(t.OutputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outFile.Close()

	// Create a new gzip writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	// Create a new tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Add files to the tar archive
	for _, file := range t.Targets {
		if err := t.addFileToTar(tarWriter, file); err != nil {
			return fmt.Errorf("error adding file to tar: %w", err)
		}
	}

	// Add symbolic links to the tar archive
	for linkName, target := range t.Links {
		if err := t.addSymlinkToTar(tarWriter, linkName, target); err != nil {
			return fmt.Errorf("error adding symlink to tar: %w", err)
		}
	}

	return nil
}

// addFileToTar adds a file to the given tar writer
func (t *TarArchive) addFileToTar(tarWriter *tar.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar header from the file info
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// must provide real name
	// (see https://golang.org/src/archive/tar/common.go?#L626)
	header.Name = filepath.ToSlash(filePath)
	utils.Debugf("header.Name: %s", header.Name)

	// Write the tar header
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy the file data to the tar writer
	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return err
	}

	return nil
}

// addSymlinkToTar adds a symbolic link to the given tar writer
func (t *TarArchive) addSymlinkToTar(tarWriter *tar.Writer, linkName, target string) error {
	// Create a tar header for the symlink
	header := &tar.Header{
		Name:     linkName,
		Mode:     0777, // Permissions for the symlink
		Typeflag: tar.TypeSymlink,
		Linkname: target,
	}

	// Write the tar header
	err := tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	return nil
}
