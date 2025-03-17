package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the service
type Config struct {
	DatabaseURL string
	GrpcPort    int
	HttpPort    int
	LogLevel    string
	Environment string
}

// Load reads the configuration from environment variables
func Load() (*Config, error) {
	grpcPort, err := strconv.Atoi(getEnv("GRPC_PORT", "50051"))
	if err != nil {
		return nil, err
	}

	httpPort, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/money_pulse?sslmode=disable"),
		GrpcPort:    grpcPort,
		HttpPort:    httpPort,
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
