package config

import (
    "os"
)

type Config struct {
    Port     string
    Database DatabaseConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        Port: os.Getenv("PORT"),
        Database: DatabaseConfig{
            Host:     os.Getenv("DB_HOST"),
            Port:     os.Getenv("DB_PORT"),
            User:     os.Getenv("DB_USER"),
            Password: os.Getenv("DB_PASSWORD"),
            Name:     os.Getenv("DB_NAME"),
        },
    }
    return cfg, nil
}
