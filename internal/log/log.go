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

func Info(msg string) {
	fmt.Fprintf(os.Stderr, "%s%sINFO:%s %s\n", bold, colorBlue, colorReset, msg)
}

func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "%s%sWARN:%s %s\n", bold, colorYellow, colorReset, msg)
}

func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s%sERROR:%s %s\n", bold, colorRed, colorReset, msg)
}

func Success(msg string) {
	fmt.Fprintf(os.Stderr, "%s%sSUCCESS:%s %s\n", bold, colorGreen, colorReset, msg)
}
