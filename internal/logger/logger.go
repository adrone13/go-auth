package logger

import (
	"log"
)

/*
TODO:
User log standard package instead of fmt to add timestamp, src file etc.

https://stackoverflow.com/questions/60721367/how-to-log-error-without-exit-the-program-in-golang
*/
const (
	infoColor    = "\033[1;34m%+v\033[0m"
	noticeColor  = "\033[1;36m%+v\033[0m"
	warningColor = "\033[1;33m%+v\033[0m"
	errorColor   = "\033[1;31m%+v\033[0m"
	debugColor   = "\033[0;36m%+v\033[0m"
)

func Log(v any) {
	log.Printf("%+v", v)
}

func Info(v any) {
	log.Printf(infoColor, v)
}

func Warn(v any) {
	log.Printf(warningColor, v)
}

func Error(v any) {
	log.Printf(errorColor, v)
}

func Debug(v any) {
	log.Printf(debugColor, v)
}
