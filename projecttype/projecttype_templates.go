package projecttype

import (
	"embed"
	"fmt"
	"os"

	"github.com/jvzantvoort/tmux-project/utils"
)

//go:embed templates/*
var Content embed.FS

func (ptc ProjectTypeConfig) Copy(srcFile, destFile string) error {
	utils.LogStart()
	defer utils.LogEnd()

	filename := fmt.Sprintf("templates/%s", srcFile)

	content, err := Content.ReadFile(filename)
	if err != nil {
		utils.Errorf("Error: %s", err)
		content = []byte("undefined")
	}
	file, _ := os.Create(destFile)
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
