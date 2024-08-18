package project

type Target struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	Content     string `json:"-"`
}

type Repo struct {
	Url         string `json:"url"`
	Destination string `json:"destination"`
	Branch      string `json:"branch"`
}

// Project defines a structure of a project
type Project struct {
	// Stored variables
	ProjectDescription string   `json:"description"`
	ProjectDir         string   `json:"directory"` // Workdir for the project
	ProjectName        string   `json:"name"`
	ProjectType        string   `json:"type"`
	SetupActions       []string `json:"setupactions"`
	Repos              []Repo   `json:"repos"`
	Targets            []Target `json:"targets"`

	// Derived variables
	HomeDir        string `json:"-"`
	ProjectTypeDir string `json:"-"` // Directory where project type files are located
	Pattern        string `json:"-"` // pattern obtained from ProjectType
	GOARCH         string `json:"-"` // target architecture
	GOOS           string `json:"-"` // target operating system
	GOPATH         string `json:"-"` // Go paths
	USER           string `json:"-"` // Username
	Exists         bool   `json:"-"` // project exists
}
