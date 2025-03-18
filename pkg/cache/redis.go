package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache represents a Redis cache client
type Cache struct {
	client *redis.Client
}

// Config holds Redis configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewCache creates a new Redis cache client
func NewCache(cfg Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{client: client}, nil
}

// Set stores a value in the cache with the given key and expiration
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set cache key: %w", err)
	}

	return nil
}

// Get retrieves a value from the cache by key
func (c *Cache) Get(ctx context.Context, key string, value interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get cache key: %w", err)
	}

	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a value from the cache by key
func (c *Cache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cache key: %w", err)
	}

	return nil
}

// SetNX sets a value in the cache only if the key does not exist
func (c *Cache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	ok, err := c.client.SetNX(ctx, key, data, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set cache key: %w", err)
	}

	return ok, nil
}

// GetWithTTL retrieves a value and its TTL from the cache by key
func (c *Cache) GetWithTTL(ctx context.Context, key string, value interface{}) (time.Duration, error) {
	pipe := c.client.Pipeline()
	getCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)

	if _, err := pipe.Exec(ctx); err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("key not found: %s", key)
		}
		return 0, fmt.Errorf("failed to get cache key: %w", err)
	}

	data, err := getCmd.Bytes()
	if err != nil {
		return 0, fmt.Errorf("failed to get cache value: %w", err)
	}

	if err := json.Unmarshal(data, value); err != nil {
		return 0, fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return ttlCmd.Val(), nil
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// Example usage:
// cache, err := cache.NewCache(cache.Config{
//     Host: "localhost",
//     Port: 6379,
// })
// if err != nil {
//     log.Fatal(err)
// }
// defer cache.Close()
//
// err = cache.Set(ctx, "user:123", user, 24*time.Hour)
// var user User
// err = cache.Get(ctx, "user:123", &user)
