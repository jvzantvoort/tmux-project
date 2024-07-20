package messages

import (
	"embed"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Content missing godoc.
//
//go:embed usage/* long/* shells/*
var Content embed.FS

func GetUsage(name string) string {
	filename := fmt.Sprintf("usage/%s", name)
	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return strings.TrimSuffix(string(msgstr), "\n")
}

func GetLong(name string) string {
	filename := fmt.Sprintf("long/%s", name)
	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return string(msgstr)
}

func GetShell(name string) string {
	filename := fmt.Sprintf("shells/%s", name)
	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Error(err)
		msgstr = []byte("# undefined")
	}
	return string(msgstr)
}
