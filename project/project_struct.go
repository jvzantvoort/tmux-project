package project

// Target defines a structure of a file
type Target struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	Content     string `json:"-"`
}

// Repo defines the struct of a (git) repository
type Repo struct {
	Url         string `json:"url"`
	Destination string `json:"destination"`
	Branch      string `json:"branch"`
}

// Project defines a structure of a project
type Project struct {
	Name           string   `json:"name"`
	ProjectType    string   `json:"type"`
	Description    string   `json:"description"`
	Directory      string   `json:"directory"`
	SetupActions   []string `json:"setupactions"`
	Repos          []Repo   `json:"repos"`
	Targets        []Target `json:"targets"`
	Root           string   `json:"root"` // Root type of the project directory or git url
	HomeDir        string   `json:"-"`
	ProjectTypeDir string   `json:"-"` // Directory where project type files are located
	Pattern        string   `json:"-"` // pattern obtained from ProjectType
	GOARCH         string   `json:"-"` // target architecture
	GOOS           string   `json:"-"` // target operating system
	GOPATH         string   `json:"-"` // Go paths
	USER           string   `json:"-"` // Username
	Exists         bool     `json:"-"` // project exists
}
