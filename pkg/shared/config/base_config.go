// Package config provides configuration utilities.
//
// Deprecated: This package is being migrated to github.com/adil-faiyaz98/sparkfund/pkg/config.
// Please use that package for new code.
package config

import (
	"time"

	"github.com/spf13/viper"
)

// BaseConfig provides common configuration structure
type BaseConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Security SecurityConfig `mapstructure:"security"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Feature  FeatureConfig  `mapstructure:"feature"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type SecurityConfig struct {
	JWT            JWTConfig  `mapstructure:"jwt"`
	CORS           CORSConfig `mapstructure:"cors"`
	RateLimit      RateLimit  `mapstructure:"rate_limit"`
	TrustedProxies []string   `mapstructure:"trusted_proxies"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
	Issuer     string        `mapstructure:"issuer"`
}

// LoadBaseConfig provides common configuration loading logic
func LoadBaseConfig(serviceName string) (*BaseConfig, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Load from files
	v.SetConfigName("config")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config BaseConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Load environment-specific overrides
	loadEnvOverrides(v)

	return &config, nil
}
