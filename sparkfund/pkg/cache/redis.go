package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config holds Redis configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Client represents a Redis client
type Client struct {
	rdb *redis.Client
}

// NewClient creates a new Redis client
func NewClient(config *Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

// Set stores a value in Redis
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.rdb.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

// Get retrieves a value from Redis
func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a value from Redis
func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete value: %w", err)
	}

	return nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.rdb.Close()
}
