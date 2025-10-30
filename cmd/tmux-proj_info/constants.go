// Constants for color formatting in the tmux project info command
package main

import (
	"github.com/fatih/color"
)

const (
	TitleColor         color.Attribute = color.FgMagenta // Color for titles
	InfoNameColor      color.Attribute = color.Bold      // Color for info names
	InfoValueColor     color.Attribute = color.FgYellow  // Color for info values
	BranchDefaultColor color.Attribute = color.FgBlue    // Color for default branch
	BranchChangedColor color.Attribute = color.FgYellow  // Color for changed branch
)
