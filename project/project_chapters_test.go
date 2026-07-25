package project

import (
	"strings"
	"testing"
)

func TestChapters_Classify(t *testing.T) {
	chapters := Chapters{
		{Path: "work/*", Chapter: "work"},
		{Path: "vendor/*", Chapter: "vendor"},
	}

	tests := []struct {
		relpath string
		want    string
	}{
		{"", "root"},
		{".", "root"},
		{"work/foo", "work"},
		{"vendor/bar", "vendor"},
		{"work/foo/bar", "rest"}, // "*" does not cross "/"
		{"other", "rest"},
	}

	for _, tt := range tests {
		if got := chapters.Classify(tt.relpath); got != tt.want {
			t.Errorf("Classify(%q) = %q, want %q", tt.relpath, got, tt.want)
		}
	}
}

func TestReadChapters(t *testing.T) {
	yamlDoc := `
- path: "work/*"
  chapter: work
- path: "vendor/*"
  chapter: vendor
`
	chapters, err := ReadChapters(strings.NewReader(yamlDoc))
	if err != nil {
		t.Fatalf("ReadChapters failed: %v", err)
	}

	if len(chapters) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(chapters))
	}

	if got := chapters.Classify("work/foo"); got != "work" {
		t.Errorf("Classify(work/foo) = %q, want work", got)
	}
}

func TestProject_LoadChapters_MissingFile(t *testing.T) {
	proj := Project{Name: "does-not-exist-project"}

	chapters, err := proj.LoadChapters()
	if err != nil {
		t.Fatalf("LoadChapters failed: %v", err)
	}
	if len(chapters) != 0 {
		t.Fatalf("expected no rules, got %d", len(chapters))
	}
	if got := chapters.Classify("anything"); got != "rest" {
		t.Errorf("Classify(anything) = %q, want rest", got)
	}
}
