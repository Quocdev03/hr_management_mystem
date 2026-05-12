package loggers

import (
	"log"
	"os"
)

// Logger wrapper đơn giản

var (
	inforLogger = log.New(os.Stdout, "[INFO]  ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger  = log.New(os.Stdout, "[WARN]  ", log.Ldate|log.Ltime)
)

func Infor(format string, v ...interface{}) {
	inforLogger.Printf(format, v...)
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
