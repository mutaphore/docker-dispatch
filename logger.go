package main

import (
	"bytes"
	"fmt"
	"log"
)

const (
	LogLevel1 = "Debug"
	LogLevel2 = "Info"
	LogLevel3 = "Silent"
)

type Logger struct {
	logLevel string
	logger   *log.Logger
}

func NewLogger(logLevel string) *Logger {
	switch logLevel {
	default:
		logLevel = LogLevel3
	case LogLevel1:
		logLevel = LogLevel1
	case LogLevel2:
		logLevel = LogLevel2
	case LogLevel3:
		logLevel = LogLevel3
	}
	var buf bytes.Buffer
	return &Logger{
		logLevel: logLevel,
		logger:   log.New(&buf, "logger: ", log.Ldate|log.Ltime),
	}
}

func (l *Logger) debug(msg string) {
	if l.logLevel == LogLevel1 || l.logLevel == LogLevel2 {
		fmt.Print(msg)
	}
}

func (l *Logger) info(msg string) {
	if l.logLevel == LogLevel2 {
		fmt.Print(msg)
	}
}

func (l *Logger) fatal(msg string) {
	log.Fatal(msg)
}

func (l *Logger) failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
