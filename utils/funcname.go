package utils

import (
	"fmt"
	"runtime"
	"strings"
)

// prefix returns a prefix for logging and messages based on function name.
func FunctionName(indent ...int) string {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	retv := elements[len(elements)-1]
	if len(indent) != 0 {
		retv = fmt.Sprintf("%s%s", strings.Repeat(" ", indent[0]), retv)
	}
	return retv
}
