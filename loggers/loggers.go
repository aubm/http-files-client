package loggers

import (
	"log"
	"os"
)

var (
	// Info logs information data
	Info *log.Logger
	// Error logs errors
	Error *log.Logger
)

// InitLoggers inits Info and Error loggers
func InitLoggers() *os.File {
	logFileName := "logs.txt"
	out, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open " + logFileName)
	}

	flag := log.Ldate | log.Ltime
	Info = log.New(out, "INFO: ", flag)
	Error = log.New(out, "ERROR: ", flag)

	return out
}
