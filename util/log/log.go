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

func Info(input string) {
	fmt.Printf(infoColor, input)
	fmt.Println("")
}

func Error(input string) {
	fmt.Printf(errorColor, input)
	fmt.Println("")
}

func Notice(input string) {
	fmt.Printf(noticeColor, input)
	fmt.Println("")
}

func Warning(input string) {
	fmt.Printf(warningColor, input)
	fmt.Println("")
}

func Debug(input string) {
	fmt.Printf(debugColor, input)
	fmt.Println("")
}
