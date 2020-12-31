package config

import ()

type ProjectTypeFiles struct {
	Name        string
	Destination string
	Mode        string
}

type Configuration struct {
	ProjectType  string
	Workdir      string
	Pattern      string
	SetupActions []string
	Files        []ProjectTypeFiles
}
