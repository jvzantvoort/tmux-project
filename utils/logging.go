// Package utils provides common utility functions for the tmux-project application,
// including logging helpers, file operations, and command execution utilities.
package utils

import (
	"fmt"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// LogIfError logs an error message if the provided message is not nil.
// It includes the calling function name in the log message.
func LogIfError(msg interface{}) {
	if msg == nil {
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	log.Errorf("%s: return not nil: %s", elements[len(elements)-1], msg)
}

// LogStart logs the start of a function execution with the function name
func LogStart() {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	log.Debugf("%s: start", elements[len(elements)-1])

}

// LogEnd logs the end of a function execution with the function name
func LogEnd() {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	log.Debugf("%s: end", elements[len(elements)-1])
}

// LogArgument logs a function argument with its name, type, and value
func LogArgument(name, input interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	log.Debugf("%s: %s", elements[len(elements)-1], fmt.Sprintf(" argument: %s(%T)\n", name, input))
	log.Debugf("%s: %s", elements[len(elements)-1], fmt.Sprintf("    value: %#v\n", input))

}

// LogVariable logs a variable with its name, type, and value
func LogVariable(name, input interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	log.Debugf("%s: %s", elements[len(elements)-1], fmt.Sprintf(" variable: %s(%T)\n", name, input))
	log.Debugf("%s: %s", elements[len(elements)-1], fmt.Sprintf("    value: %#v\n", input))

}

// Debugf logs a debug message prefixed with the calling function name
func Debugf(format string, args ...interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	msg := fmt.Sprintf(format, args...)
	log.Debugf("%s: %s", elements[len(elements)-1], msg)

}

// Infof logs an info message prefixed with the calling function name
func Infof(format string, args ...interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	msg := fmt.Sprintf(format, args...)
	log.Infof("%s: %s", elements[len(elements)-1], msg)

}

// Errorf logs an error message prefixed with the calling function name
func Errorf(format string, args ...interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	msg := fmt.Sprintf(format, args...)
	log.Errorf("%s: %s", elements[len(elements)-1], msg)

}

// Warningf logs a warning message prefixed with the calling function name
func Warningf(format string, args ...interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	msg := fmt.Sprintf(format, args...)
	log.Warningf("%s: %s", elements[len(elements)-1], msg)

}

// Fatalf logs a fatal error message prefixed with the calling function name and exits
func Fatalf(format string, args ...interface{}) {

	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	msg := fmt.Sprintf(format, args...)
	log.Fatalf("%s: %s", elements[len(elements)-1], msg)

}
