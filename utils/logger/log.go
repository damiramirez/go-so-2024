package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
)

type Logger struct {
	Logger      *log.Logger
	LogFile     *os.File
	environment string
}

// Log formats and logs a message at the given severity level.
func (l *Logger) Log(message string, level LogLevel) {
	if l.environment == "prod" && level == "DEBUG" {
		return
	}

	format := fmt.Sprintf("[%s]: %s", level, message)
	l.Logger.Println(format)
}

func (l *Logger) CloseLogger() {
	l.LogFile.Close()
}

// Initializes a Logger to write to both stdout and a specified file.
func ConfigureLogger(filename, env string) Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logger := Logger{Logger: log.New(mw, "", log.LstdFlags), LogFile: logFile, environment: env}
	return logger
}
