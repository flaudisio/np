// Package log writes colored log messages to stderr.
//
// Prefixes use ANSI bold+color: blue (Info), yellow (Warn),
// red (Error), green (Success).
package log

import (
	"fmt"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	bold        = "\033[1m"
)

// Info writes a blue bold "INFO:" message to stderr.
func Info(msg string) {
	fmt.Fprintf(os.Stderr, "%s%sINFO:%s %s\n", bold, colorBlue, colorReset, msg)
}

// Warn writes a yellow bold "WARN:" message to stderr.
func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "%s%sWARN:%s %s\n", bold, colorYellow, colorReset, msg)
}

// Error writes a red bold "ERROR:" message to stderr.
func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s%sERROR:%s %s\n", bold, colorRed, colorReset, msg)
}

// Success writes a green bold "SUCCESS:" message to stderr.
func Success(msg string) {
	fmt.Fprintf(os.Stderr, "%s%sSUCCESS:%s %s\n", bold, colorGreen, colorReset, msg)
}
