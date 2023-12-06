package logger

import "fmt"

const (
	infoColor    = "\033[1;34m%+v\033[0m"
	noticeColor  = "\033[1;36m%+v\033[0m"
	warningColor = "\033[1;33m%+v\033[0m"
	errorColor   = "\033[1;31m%+v\033[0m"
	debugColor   = "\033[0;36m%+v\033[0m"
)

func Log(v ...any) {
	fmt.Printf("%+v", v)
	fmt.Printf("\n")
}

func Info(v ...any) {
	fmt.Printf(infoColor, v)
	fmt.Printf("\n")
}

func Notice(v ...any) {
	fmt.Printf(noticeColor, v)
	fmt.Printf("\n")
}

func Warn(v ...any) {
	fmt.Printf(warningColor, v)
	fmt.Printf("\n")
}

func Error(v ...any) {
	fmt.Printf(errorColor, v)
	fmt.Printf("\n")
}

func Debug(v ...any) {
	fmt.Printf(debugColor, v)
	fmt.Printf("\n")
}
