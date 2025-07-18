package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// printTitle formats and prints a title with a specific color.
// It uses the TitleColor constant for the color formatting.
// The title is converted to title case using the cases package for better readability.
// The title is printed with a line of dashes below it for visual separation.
// This function is useful for creating visually distinct sections in the output.
// It can be used to highlight important sections or headings in the command output.
// The title is printed in a consistent format, making it easy to identify different parts of the output.
// The function does not return any value, it directly prints to the standard output.
// It is typically used in conjunction with other output functions to create a structured and readable command output
func printTitle(title string) {
	css := cases.Title(language.English)
	titlestr := css.String(title)
	purple := color.New(TitleColor)
	fmt.Printf("\n%s\n", purple.Sprint(titlestr))
	fmt.Print(strings.Repeat("-", len(title)))
	fmt.Printf("\n\n")
}
