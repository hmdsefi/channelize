package log

import (
	"log"
	"os"
)

const (
	defaultPrefix = ""
)

const (
	Info  Level = "INFO"
	Error Level = "ERROR"
	Warn  Level = "WARN"
	Debug Level = "DEBUG"
)

// Level is an alias type of string that represents the log level.
type Level string

// DefaultLogger is an implementation of channelize log.Logger interface on top
// of the built-in log.Logger. It is used if logger is not specified.
type DefaultLogger struct {
	logger *log.Logger
}

// NewDefaultLogger creates new instance of
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, defaultPrefix, log.LstdFlags),
	}
}

func (l *DefaultLogger) println(level Level, msg string, keyValues ...interface{}) {
	l.logger.Println(append([]interface{}{level, msg}, keyValues...))
}

// Info writes message to the log with info level.
func (l *DefaultLogger) Info(msg string, keyValues ...interface{}) {
	l.println(Info, msg, keyValues)
}

// Error writes message to the log with error level.
func (l *DefaultLogger) Error(msg string, keyValues ...interface{}) {
	l.println(Error, msg, keyValues)
}

// Warn writes message to the log with warn level.
func (l *DefaultLogger) Warn(msg string, keyValues ...interface{}) {
	l.println(Warn, msg, keyValues)
}

// Debug writes message to the log with debug level.
func (l *DefaultLogger) Debug(msg string, keyValues ...interface{}) {
	l.println(Debug, msg, keyValues)
}
