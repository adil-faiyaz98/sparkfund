package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init initializes the logger with the specified level
func Init(level string) {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

// WithFields returns a new logger with the specified fields
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return log.WithFields(fields)
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Panic logs a panic message and panics
func Panic(args ...interface{}) {
	log.Panic(args...)
} 