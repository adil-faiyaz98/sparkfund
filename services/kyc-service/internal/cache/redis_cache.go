package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	prefix string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Prefix   string
}

// DefaultRedisConfig returns default Redis configuration
func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
		Prefix:   "kyc:",
	}
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(cfg RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		prefix: cfg.Prefix,
	}, nil
}

// Get retrieves a value from the cache
func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	data, err := c.client.Get(ctx, c.prefix+key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}

	return json.Unmarshal(data, value)
}

// Set stores a value in the cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.prefix+key, data, ttl).Err()
}

// Delete removes a value from the cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.prefix+key).Err()
}

// Clear removes all values with the prefix from the cache
func (c *RedisCache) Clear(ctx context.Context) error {
	// Get all keys with the prefix
	keys, err := c.client.Keys(ctx, c.prefix+"*").Result()
	if err != nil {
		return err
	}

	// If there are no keys, return nil
	if len(keys) == 0 {
		return nil
	}

	// Delete all keys
	return c.client.Del(ctx, keys...).Err()
}

// GetMulti retrieves multiple values from the cache
func (c *RedisCache) GetMulti(ctx context.Context, keys []string, values interface{}) error {
	// Add prefix to keys
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.prefix + key
	}

	// Get all values
	data, err := c.client.MGet(ctx, prefixedKeys...).Result()
	if err != nil {
		return err
	}

	// Convert to map
	result := make(map[string]interface{})
	for i, key := range keys {
		if data[i] != nil {
			var value interface{}
			if err := json.Unmarshal([]byte(data[i].(string)), &value); err != nil {
				return err
			}
			result[key] = value
		}
	}

	// Marshal and unmarshal to convert to the target type
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return json.Unmarshal(resultBytes, values)
}

// SetMulti stores multiple values in the cache
func (c *RedisCache) SetMulti(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	// Use pipeline for better performance
	pipe := c.client.Pipeline()

	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		pipe.Set(ctx, c.prefix+key, data, ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// DeleteMulti removes multiple values from the cache
func (c *RedisCache) DeleteMulti(ctx context.Context, keys []string) error {
	// Add prefix to keys
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = c.prefix + key
	}

	return c.client.Del(ctx, prefixedKeys...).Err()
}

// Increment increments a counter in the cache
func (c *RedisCache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, c.prefix+key, value).Result()
}

// Decrement decrements a counter in the cache
func (c *RedisCache) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.DecrBy(ctx, c.prefix+key, value).Result()
}

// SetNX sets a value in the cache only if the key does not exist
func (c *RedisCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	return c.client.SetNX(ctx, c.prefix+key, data, ttl).Result()
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}
