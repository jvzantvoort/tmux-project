package project

type ProjectTarget struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	Content     string `json:"-"`
}

type Project struct {
	HomeDir            string          `json:"homedir"`
	ProjectDescription string          `json:"description"`
	ProjectDir         string          `json:"directory"` // Workdir for the project
	ProjectName        string          `json:"name"`
	ProjectType        string          `json:"type"`
	SetupActions       []string        `json:"setupactions"`
	Targets            []ProjectTarget `json:"targets"`
	ProjectTypeDir     string          `json:"-"`
	Pattern            string          `json:"-"` // pattern obtained from ProjectType
	GOARCH             string          `json:"-"` // target architecture
	GOOS               string          `json:"-"` // target operating system
	GOPATH             string          `json:"-"` // Go paths
	USER               string          `json:"-"`
}
