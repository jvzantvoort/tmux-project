package projecttype

import (
	"bytes"
	"os"
	"testing"
)

func TestProjectTypeConfig_Content(t *testing.T) {
	dir := t.TempDir()
	file := "testfile.yml"
	content := "hello: world"
	if err := os.WriteFile(dir+"/"+file, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	ptc := ProjectTypeConfig{ConfigDir: dir}
	got, err := ptc.Content(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != content {
		t.Errorf("expected %q, got %q", content, got)
	}
}

func TestProjectTypeConfig_ReadWrite(t *testing.T) {
	ptc := &ProjectTypeConfig{
		ProjectType:  "test",
		Description:  "desc",
		Directory:    "dir",
		Pattern:      "*.go",
		SetupActions: []string{"echo hi"},
		Repos:        []Repo{{Url: "url", Destination: "dest", Branch: "main"}},
		Targets:      []Target{{Name: "n", Destination: "d", Mode: "0644"}},
	}
	buf := &bytes.Buffer{}
	err := ptc.Write(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ptc2 := &ProjectTypeConfig{}
	err = ptc2.Read(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ptc2.ProjectType != ptc.ProjectType {
		t.Errorf("expected %q, got %q", ptc.ProjectType, ptc2.ProjectType)
	}
}
