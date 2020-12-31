package config

type ProjectTypeFiles struct {
	Name        string
	Destination string
	Mode        string
}

type ProjectTypeConfig struct {
	ProjectType  string
	Workdir      string
	Pattern      string
	SetupActions []string
	Files        []ProjectTypeFiles
}
