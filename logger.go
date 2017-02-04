package main

import (
	"log"
	"strings"
)

type Logger struct {
	logLevel int
}

func (l *Logger) debug(format string, v ...interface{}) {
	if l.logLevel < 4 {
		return
	}

	logMessage("DEBUG", format, v)
}

func (l *Logger) info(format string, v ...interface{}) {
	if l.logLevel < 3 {
		return
	}

	logMessage("INFO", format, v)
}

func (l *Logger) warning(format string, v ...interface{}) {
	if l.logLevel < 2 {
		return
	}

	logMessage("WARNING", format, v)
}

func (l *Logger) error(format string, v ...interface{}) {
	if l.logLevel < 1 {
		return
	}

	logMessage("ERROR", format, v)
}

func logMessage(level string, format string, v interface{}) {
	vInterface := v.([]interface{})
	if len(vInterface) > 0 {
		log.Printf("["+level+"]\t"+format, vInterface...)
		return
	}

	log.Printf("[%s]\t"+format, level)
}

func NewLogger() *Logger {
	logLevel := 0

	switch strings.ToLower(config.Loglevel) {
	case "debug":
		logLevel = 4
	case "info":
		logLevel = 3
	case "warning":
		logLevel = 2
	case "error":
		logLevel = 1
	default:
		logLevel = 2
	}

	return &Logger{logLevel}
}
