package projecttype

// Target defines a structure of a file
type Target struct {
	Name        string `yaml:"name"`
	Destination string `yaml:"destination"`
	Mode        string `yaml:"mode"`
}

type Repo struct {
	Url         string `yaml:"url"`
	Destination string `yaml:"destination"`
	Branch      string `yaml:"branch"`
}

// ProjectTypeConfig defines a structure of a project type
type ProjectTypeConfig struct {
	ProjectType    string   `yaml:"projecttype"`
	ProjectTypeDir string   `yaml:"-"`
	Workdir        string   `yaml:"workdir"`
	Pattern        string   `yaml:"pattern"`
	SetupActions   []string `yaml:"setupactions"`
	Repos          []Repo   `yaml:"repos"`
	Files          []Target `yaml:"targets"`
	ConfigFile     string   `yaml:"-"`
	ConfigDir      string   `yaml:"-"`
}
