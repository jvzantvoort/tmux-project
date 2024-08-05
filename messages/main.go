package messages

import (
	"embed"
	"fmt"
	"strings"

	"github.com/jvzantvoort/tmux-project/utils"
)

// Content missing godoc.
//
//go:embed long/* shells/* use/* short/*
var Content embed.FS

func GetContent(folder, name string) string {
	filename := fmt.Sprintf("%s/%s", folder, name)

	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		utils.Errorf("%s", err)
		msgstr = []byte("undefined")
	}
	return strings.TrimSuffix(string(msgstr), "\n")

}

func GetShort(name string) string {
	return GetContent("short", name)
}

func GetUse(name string) string {
	return GetContent("use", name)
}

func GetLong(name string) string {
	return GetContent("long", name)
}

func GetShell(name string) string {
	return GetContent("shells", name)
}
