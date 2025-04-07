package cache

import (
	"context"
	"errors"
	"time"
)

// ErrCacheMiss is returned when a key is not found in the cache
var ErrCacheMiss = errors.New("cache miss")

// Cache defines the interface for cache implementations
type Cache interface {
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string, value interface{}) error
	
	// Set stores a value in the cache
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	
	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error
	
	// Clear removes all values with the prefix from the cache
	Clear(ctx context.Context) error
	
	// GetMulti retrieves multiple values from the cache
	GetMulti(ctx context.Context, keys []string, values interface{}) error
	
	// SetMulti stores multiple values in the cache
	SetMulti(ctx context.Context, items map[string]interface{}, ttl time.Duration) error
	
	// DeleteMulti removes multiple values from the cache
	DeleteMulti(ctx context.Context, keys []string) error
	
	// Increment increments a counter in the cache
	Increment(ctx context.Context, key string, value int64) (int64, error)
	
	// Decrement decrements a counter in the cache
	Decrement(ctx context.Context, key string, value int64) (int64, error)
	
	// SetNX sets a value in the cache only if the key does not exist
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)
	
	// Close closes the cache connection
	Close() error
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Type      string // "memory", "redis"
	TTL       time.Duration
	RedisHost string
	RedisPort string
	RedisPass string
	RedisDB   int
	Prefix    string
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		Type:      "memory",
		TTL:       5 * time.Minute,
		RedisHost: "localhost",
		RedisPort: "6379",
		RedisPass: "",
		RedisDB:   0,
		Prefix:    "kyc:",
	}
}

// NewCache creates a new cache based on the configuration
func NewCache(cfg CacheConfig) (Cache, error) {
	switch cfg.Type {
	case "memory":
		return NewMemoryCache(MemoryCacheConfig{
			TTL:             cfg.TTL,
			CleanupInterval: cfg.TTL / 2,
			Prefix:          cfg.Prefix,
		}), nil
	case "redis":
		return NewRedisCache(RedisConfig{
			Host:     cfg.RedisHost,
			Port:     cfg.RedisPort,
			Password: cfg.RedisPass,
			DB:       cfg.RedisDB,
			Prefix:   cfg.Prefix,
		})
	default:
		return nil, errors.New("unsupported cache type")
	}
}
