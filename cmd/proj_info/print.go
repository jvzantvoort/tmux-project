package main

import (
	"fmt"

	"github.com/fatih/color"
)

func printTitle(title string) {

	purple := color.New(TitleColor)
	fmt.Printf("\n  %s:\n\n", purple.Sprint(title))

}

func printInfo(itype, ival string) {
	infNameCol := color.New(InfoNameColor)
	infValCol := color.New(InfoValueColor)
	fmt.Printf("%-24s %s\n", infNameCol.Sprint(itype)+":", infValCol.Sprint(ival))

}
