package project

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/jvzantvoort/tmux-project/config"
	"gopkg.in/yaml.v2"
)

// ChapterRule maps a glob pattern, relative to the project root, to a
// chapter name. Patterns are matched with path.Match: "*" matches any
// sequence of non-separator characters, so it stays within one path
// segment and does not match across "/".
type ChapterRule struct {
	Path    string `yaml:"path"`
	Chapter string `yaml:"chapter"`
}

// Chapters is an ordered list of ChapterRule, evaluated first match wins.
type Chapters []ChapterRule

// Classify returns the chapter for relpath. The project root itself
// ("" or ".") is always chapter "root". If no rule matches, it falls
// back to chapter "rest".
func (c Chapters) Classify(relpath string) string {
	if relpath == "" || relpath == "." {
		return "root"
	}
	for _, rule := range c {
		if ok, _ := path.Match(rule.Path, relpath); ok {
			return rule.Chapter
		}
	}
	return "rest"
}

// ReadChapters parses chapter classification rules from reader.
func ReadChapters(reader io.Reader) (Chapters, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var retv Chapters
	if err := yaml.Unmarshal(data, &retv); err != nil {
		return nil, err
	}
	return retv, nil
}

// ChaptersFile returns the full path to the project's chapter
// classification file.
func (proj Project) ChaptersFile() string {
	return filepath.Join(config.SessionDir(), proj.Name+".chapters.yml")
}

// LoadChapters loads the chapter classification rules for the project.
// A missing file is not an error; it yields an empty rule set so callers
// fall back to the default root/rest classification.
func (proj Project) LoadChapters() (Chapters, error) {
	file := proj.ChaptersFile()

	filehandle, err := os.Open(file) // #nosec G304 - controlled config file path
	if err != nil {
		if os.IsNotExist(err) {
			return Chapters{}, nil
		}
		return nil, err
	}
	defer func() {
		_ = filehandle.Close()
	}()

	return ReadChapters(filehandle)
}
