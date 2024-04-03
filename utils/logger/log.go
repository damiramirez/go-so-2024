package utils

import (
	"io"
	"log"
	"os"
)

// Sets up a logger writing to stdout and a file.
func ConfigureLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	return log.New(mw, "", log.LstdFlags)
}

// Closes the log file.
func CloseLogger(logFile *os.File) {
	logFile.Close()
}