package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

const (
	CONSOLE_WIDTH  int    = 80
	CONSOLE_INDENT int    = 2
	VBORDER        string = "|"
	HBORDER        string = "-"
	CBORDER        string = "+"
)

// ErrorExit prints the msg with the prefix 'Error:' and exits with error code 1. If the msg is nil, it does nothing.
func ErrorExit(msg interface{}) {
	if msg != nil {
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}

func Abort(format string, input ...interface{}) {
	msg := fmt.Sprintf(format, input...)

	if len(msg) == 0 {
		os.Exit(1)
	}
	fmt.Println("")
	fmt.Println(box_header())

	msg = strings.TrimSpace(msg)
	msg = wordwrap.WrapString(msg, box_widthu())
	for _, instr := range strings.Split(msg, "\n") {
		fmt.Println(center_text(instr))
	}

	fmt.Println(box_footer())
	fmt.Println("")

	os.Exit(1)
}

func box_width() int {
	return int(CONSOLE_WIDTH - (CONSOLE_INDENT * 2))
}

func box_widthu() uint {
	return uint(box_width())
}

func box_indent() int {
	return int(CONSOLE_INDENT)
}

func inner_box_width() int {
	return int(box_width() - (box_indent() * 2))
}

func get_inner_string(instr string) string {
	instrlen := len(instr)
	restlen := box_width() - instrlen                // get the spaces left in the line
	restlen -= (box_indent() * 2)                    // remove the indent chars
	instr = strings.Repeat(" ", (restlen/2)) + instr // prefix half of the spaces
	return instr
}

func box_header() string {
	width := box_width() - (len(CBORDER) * 2)

	retv := ""
	retv += strings.Repeat(" ", int(box_indent()))
	retv += CBORDER
	retv += strings.Repeat("-", width)
	retv += CBORDER
	retv += "\n" + center_text("")
	return retv
}

func box_footer() string {
	width := box_width() - (len(CBORDER) * 2)
	retv := ""
	retv += center_text("") + "\n"
	retv += strings.Repeat(" ", int(box_indent()))
	retv += CBORDER
	retv += strings.Repeat("-", width)
	retv += CBORDER

	return retv
}

func center_text(instr string) string {
	vborder := VBORDER
	instr = get_inner_string(instr)
	stringfmt := fmt.Sprintf("%s%s %%-%ds %s", strings.Repeat(" ", box_indent()), vborder, inner_box_width(), vborder)

	return fmt.Sprintf(stringfmt, instr)
}
