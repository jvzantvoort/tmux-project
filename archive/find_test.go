package archive

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestFindFiles_Basic tests that FindFiles returns all regular files and no symlinks in a directory.
func TestFindFiles_Basic(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(dir, "file2.txt")
	if err := os.WriteFile(file1, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("world"), 0644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	files, links, err := FindFiles([]string{dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 0 {
		t.Errorf("expected no symlinks, got: %v", links)
	}
	want := []string{file1, file2}
	if !reflect.DeepEqual(sorted(files), sorted(want)) {
		t.Errorf("expected files %v, got %v", want, files)
	}
}

// TestFindFiles_Symlink tests that FindFiles correctly identifies symlinks and their targets.
func TestFindFiles_Symlink(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, []byte("data"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	symlink := filepath.Join(dir, "link.txt")
	if err := os.Symlink(file, symlink); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	files, links, err := FindFiles([]string{dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if links[symlink] != file {
		t.Errorf("expected symlink %s -> %s, got %v", symlink, file, links)
	}
	if len(files) != 1 || files[0] != file {
		t.Errorf("expected file %s, got %v", file, files)
	}
}

// sorted returns a sorted copy of the input string slice (for test comparison).
func sorted(s []string) []string {
	ss := append([]string{}, s...)
	if len(ss) > 1 && ss[0] > ss[1] {
		ss[0], ss[1] = ss[1], ss[0]
	}
	return ss
}
