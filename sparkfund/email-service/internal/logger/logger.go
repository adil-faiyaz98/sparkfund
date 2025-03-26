package logger

import (
	"os" // Import the os package

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger initializes the logger with the specified environment
func NewLogger(env string) (*zap.Logger, error) {
	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
		//Consider setting output paths in production
		//config.OutputPaths = []string{"/var/log/myapp/app.log", "stderr"}
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
