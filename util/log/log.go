package log

import (
	"fmt"
)

const (
	infoColor    = "\033[1;34m%s\033[0m"
	noticeColor  = "\033[1;36m%s\033[0m"
	warningColor = "\033[1;33m%s\033[0m"
	errorColor   = "\033[1;31m%s\033[0m"
	debugColor   = "\033[0;36m%s\033[0m"
)

// Info logs info message
func Info(input string) {
	fmt.Printf(infoColor, input)
	fmt.Println("")
}

// Error logs error message
func Error(input string) {
	fmt.Printf(errorColor, input)
	fmt.Println("")
}

// Notice logs notice message
func Notice(input string) {
	fmt.Printf(noticeColor, input)
	fmt.Println("")
}

// Warning logs warning message
func Warning(input string) {
	fmt.Printf(warningColor, input)
	fmt.Println("")
}

// Debug logs debug message
func Debug(input string) {
	fmt.Printf(debugColor, input)
	fmt.Println("")
}
