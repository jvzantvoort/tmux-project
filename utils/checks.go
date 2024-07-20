package utils

import (
	"fmt"
	"os"
)

// ErrorExit prints the msg with the prefix 'Error:' and exits with error code 1. If the msg is nil, it does nothing.
func ErrorExit(msg interface{}) {
	if msg != nil {
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}
