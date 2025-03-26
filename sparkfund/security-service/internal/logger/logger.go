package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger with appropriate configuration based on environment
func NewLogger(env string) (*zap.Logger, error) {
	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Override config with environment variables if they exist
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		lvl, err := zapcore.ParseLevel(level)
		if err == nil {
			config.Level = zap.NewAtomicLevelAt(lvl)
		}
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log, nil
}
