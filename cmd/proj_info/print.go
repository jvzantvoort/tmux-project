package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func printTitle(title string) {
	css := cases.Title(language.English)
	titlestr := css.String(title)
	purple := color.New(TitleColor)
	fmt.Printf("\n%s\n", purple.Sprint(titlestr))
	fmt.Print(strings.Repeat("-", len(title)))
	fmt.Printf("\n\n")
}
