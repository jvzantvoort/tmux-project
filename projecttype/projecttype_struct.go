package projecttype

// Target defines a structure of a file
type Target struct {
	Name        string `yaml:"name"`
	Destination string `yaml:"destination"`
	Mode        string `yaml:"mode"`
}

// Repo defines the struct of a (git) repository
type Repo struct {
	Url         string `yaml:"url"`
	Destination string `yaml:"destination"`
	Branch      string `yaml:"branch"`
}

// ProjectTypeConfig defines a structure of a project type
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
