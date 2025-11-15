// Package projecttype handles project type templates and configurations.
// It provides functionality for creating, loading, and managing reusable
// project templates with pre-defined structures and settings.
package projecttype

// Target defines a file template that will be created for a project
type Target struct {
	Name        string `yaml:"name"`
	Destination string `yaml:"destination"`
	Mode        string `yaml:"mode"`
}

// Repo defines a git repository to be cloned as part of project initialization
type Repo struct {
	Url         string `yaml:"url"`
	Destination string `yaml:"destination"`
	Branch      string `yaml:"branch"`
}

// ProjectTypeConfig represents a reusable project template configuration
type ProjectTypeConfig struct {
	ProjectType    string   `yaml:"type"`
	Description    string   `yaml:"description"`
	Directory      string   `yaml:"directory"`
	Pattern        string   `yaml:"pattern"`
	SetupActions   []string `yaml:"setupactions"`
	Repos          []Repo   `yaml:"repos"`
	Targets        []Target `yaml:"targets"`
	ProjectTypeDir string   `yaml:"-"`
	ConfigFile     string   `yaml:"-"`
	ConfigDir      string   `yaml:"-"`
	Root           string   `yaml:"root"`
}
