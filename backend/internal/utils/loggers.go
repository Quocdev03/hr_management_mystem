package utils

import (
	"log"
	"os"
)

// Logger wrapper đơn giản

var (
	infoLogger = log.New(os.Stdout, "[APP-INFO]  ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stdout, "[APP-WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger  = log.New(os.Stderr, "[APP-ERROR]  ", log.Ldate|log.Ltime)
)

func Info(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func Warn(format string, v ...interface{}) {
	warnLogger.Printf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	errorLogger.Fatalf(format, v...)
}
