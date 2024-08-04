package messages

import (
	"embed"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Content missing godoc.
//
//go:embed usage/* long/* shells/* use/*
var Content embed.FS

func GetContent(folder, name string) string {
	filename := fmt.Sprintf("%s/%s", folder, name)
	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return strings.TrimSuffix(string(msgstr), "\n")

}

func GetUsage(name string) string {
	return GetContent("usage", name)
}

func GetUse(name string) string {
	return GetContent("use", name)
}

func GetLong(name string) string {
	return GetContent("long", name)
}

func GetShell(name string) string {
	return GetContent("shell", name)
}
