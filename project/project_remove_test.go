package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckRemovableDir(t *testing.T) {
	home := "/home/user"
	tests := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{name: "empty directory is skipped", dir: "", wantErr: false},
		{name: "project below home", dir: "/home/user/work/proj", wantErr: false},
		{name: "project outside home", dir: "/srv/projects/proj", wantErr: false},
		{name: "unclean path below home", dir: "/home/user/work/../proj/", wantErr: false},
		{name: "sibling with home as prefix", dir: "/home/user2", wantErr: false},
		{name: "relative path", dir: "work/proj", wantErr: true},
		{name: "root", dir: "/", wantErr: true},
		{name: "home itself", dir: "/home/user", wantErr: true},
		{name: "home with trailing slash", dir: "/home/user/", wantErr: true},
		{name: "parent of home", dir: "/home", wantErr: true},
		{name: "home via dot-dot", dir: "/home/user/work/..", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRemovableDir(tt.dir, home)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckRemovableDir(%q) error = %v, wantErr %v", tt.dir, err, tt.wantErr)
			}
		})
	}
}

func TestRemovePaths(t *testing.T) {
	tmpDir := t.TempDir()

	projDir := filepath.Join(tmpDir, "proj")
	if err := os.MkdirAll(filepath.Join(projDir, "sub"), 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projDir, "sub", "file"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	envFile := filepath.Join(tmpDir, "proj.env")
	if err := os.WriteFile(envFile, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(tmpDir, "other.env")
	if err := os.WriteFile(keep, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := RemovePaths([]string{projDir, envFile}); err != nil {
		t.Fatalf("RemovePaths() error = %v", err)
	}

	for _, gone := range []string{projDir, envFile} {
		if _, err := os.Lstat(gone); !os.IsNotExist(err) {
			t.Errorf("%s still exists", gone)
		}
	}
	if _, err := os.Lstat(keep); err != nil {
		t.Errorf("%s was removed: %v", keep, err)
	}
}

func TestRemovePathsMissing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	if err := RemovePaths([]string{missing}); err != nil {
		t.Errorf("RemovePaths() on missing path error = %v, want nil", err)
	}
}
