package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init initializes the logger
func Init() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	if log == nil {
		Init()
	}
	return log
}

// WithFields creates a new entry with fields
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// Info logs an info message
func Info(msg string, fields map[string]interface{}) {
	WithFields(fields).Info(msg)
}

// Error logs an error message
func Error(err error, msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	WithFields(fields).Error(msg)
}

// Debug logs a debug message
func Debug(msg string, fields map[string]interface{}) {
	WithFields(fields).Debug(msg)
}

// Warn logs a warning message
func Warn(msg string, fields map[string]interface{}) {
	WithFields(fields).Warn(msg)
}

// Fatal logs a fatal message and exits
func Fatal(err error, msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	WithFields(fields).Fatal(msg)
}
