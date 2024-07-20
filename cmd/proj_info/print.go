package main

import (
	"fmt"

	"github.com/fatih/color"
)

func printTitle(title string) {

	purple := color.New(TitleColor)
	fmt.Printf("\n  %s:\n\n", purple.Sprint(title))

}

