package utils

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func InitLogger() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(format string, v ...interface{}) {
	Info.Printf(format, v...)
}

func LogWarn(format string, v ...interface{}) {
	Warn.Printf(format, v...)
}

func LogError(format string, v ...interface{}) {
	Error.Printf(format, v...)
}
