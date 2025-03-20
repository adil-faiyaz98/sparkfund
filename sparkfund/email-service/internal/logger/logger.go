package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the logger with the specified environment
func Init(env string) error {
	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// GetLogger returns the logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		// Initialize with development config if not initialized
		Init("development")
	}
	return log
} 