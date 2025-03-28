package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sparkfund/pkg/errors"
)

// Config represents Redis configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Client represents a Redis client
type Client struct {
	client *redis.Client
}

// NewClient creates a new Redis client
func NewClient(cfg *Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.ErrInternalServer(err)
	}

	return &Client{client: client}, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.client.Close()
}

// Set sets a key-value pair
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var err error
	switch v := value.(type) {
	case string:
		err = c.client.Set(ctx, key, v, expiration).Err()
	case []byte:
		err = c.client.Set(ctx, key, v, expiration).Err()
	default:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return errors.ErrInternalServer(err)
		}
		err = c.client.Set(ctx, key, jsonValue, expiration).Err()
	}
	if err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// Get gets a value by key
func (c *Client) Get(ctx context.Context, key string, value interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.ErrNotFound(fmt.Errorf("key not found: %s", key))
		}
		return errors.ErrInternalServer(err)
	}

	switch v := value.(type) {
	case *string:
		*v = val
	case *[]byte:
		*v = []byte(val)
	default:
		if err := json.Unmarshal([]byte(val), value); err != nil {
			return errors.ErrInternalServer(err)
		}
	}
	return nil
}

// Delete deletes a key
func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, errors.ErrInternalServer(err)
	}
	return n > 0, nil
}

// SetNX sets a key-value pair if the key does not exist
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	var err error
	var ok bool
	switch v := value.(type) {
	case string:
		ok, err = c.client.SetNX(ctx, key, v, expiration).Result()
	case []byte:
		ok, err = c.client.SetNX(ctx, key, v, expiration).Result()
	default:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return false, errors.ErrInternalServer(err)
		}
		ok, err = c.client.SetNX(ctx, key, jsonValue, expiration).Result()
	}
	if err != nil {
		return false, errors.ErrInternalServer(err)
	}
	return ok, nil
}

// SetXX sets a key-value pair if the key exists
func (c *Client) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	var err error
	var ok bool
	switch v := value.(type) {
	case string:
		ok, err = c.client.SetXX(ctx, key, v, expiration).Result()
	case []byte:
		ok, err = c.client.SetXX(ctx, key, v, expiration).Result()
	default:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return false, errors.ErrInternalServer(err)
		}
		ok, err = c.client.SetXX(ctx, key, jsonValue, expiration).Result()
	}
	if err != nil {
		return false, errors.ErrInternalServer(err)
	}
	return ok, nil
}

// Incr increments a key
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, errors.ErrInternalServer(err)
	}
	return val, nil
}

// Decr decrements a key
func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, errors.ErrInternalServer(err)
	}
	return val, nil
}

// Expire sets the expiration time for a key
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := c.client.Expire(ctx, key, expiration).Err(); err != nil {
		return errors.ErrInternalServer(err)
	}
	return nil
}

// TTL gets the remaining time to live of a key
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	val, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, errors.ErrInternalServer(err)
	}
	return val, nil
} 