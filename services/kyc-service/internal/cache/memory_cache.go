package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// MemoryCache implements the Cache interface using in-memory storage
type MemoryCache struct {
	items             map[string]item
	mu                sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	prefix            string
	stopCleanup       chan bool
}

// MemoryCacheConfig holds memory cache configuration
type MemoryCacheConfig struct {
	TTL             time.Duration
	CleanupInterval time.Duration
	Prefix          string
}

// item represents a cached item
type item struct {
	Value      []byte
	Expiration int64
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache(cfg MemoryCacheConfig) *MemoryCache {
	cache := &MemoryCache{
		items:             make(map[string]item),
		defaultExpiration: cfg.TTL,
		cleanupInterval:   cfg.CleanupInterval,
		prefix:            cfg.Prefix,
		stopCleanup:       make(chan bool),
	}

	// Start cleanup goroutine if cleanup interval is greater than 0
	if cfg.CleanupInterval > 0 {
		go cache.startCleanup()
	}

	return cache
}

// Get retrieves a value from the cache
func (c *MemoryCache) Get(ctx context.Context, key string, value interface{}) error {
	c.mu.RLock()
	item, found := c.items[c.prefix+key]
	c.mu.RUnlock()

	if !found {
		return ErrCacheMiss
	}

	// Check if item has expired
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return ErrCacheMiss
	}

	return json.Unmarshal(item.Value, value)
}

// Set stores a value in the cache
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Use default expiration if ttl is 0
	if ttl == 0 {
		ttl = c.defaultExpiration
	}

	// Marshal value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Calculate expiration
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	// Store item in cache
	c.mu.Lock()
	c.items[c.prefix+key] = item{
		Value:      data,
		Expiration: expiration,
	}
	c.mu.Unlock()

	return nil
}

// Delete removes a value from the cache
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	delete(c.items, c.prefix+key)
	c.mu.Unlock()
	return nil
}

// Clear removes all values with the prefix from the cache
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	for key := range c.items {
		if len(key) >= len(c.prefix) && key[:len(c.prefix)] == c.prefix {
			delete(c.items, key)
		}
	}
	c.mu.Unlock()
	return nil
}

// GetMulti retrieves multiple values from the cache
func (c *MemoryCache) GetMulti(ctx context.Context, keys []string, values interface{}) error {
	// Get current time once for all items
	now := time.Now().UnixNano()

	// Create result map
	result := make(map[string]interface{})

	c.mu.RLock()
	for _, key := range keys {
		item, found := c.items[c.prefix+key]
		if found && (item.Expiration == 0 || now < item.Expiration) {
			var value interface{}
			if err := json.Unmarshal(item.Value, &value); err != nil {
				c.mu.RUnlock()
				return err
			}
			result[key] = value
		}
	}
	c.mu.RUnlock()

	// Marshal and unmarshal to convert to the target type
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return json.Unmarshal(resultBytes, values)
}

// SetMulti stores multiple values in the cache
func (c *MemoryCache) SetMulti(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	// Use default expiration if ttl is 0
	if ttl == 0 {
		ttl = c.defaultExpiration
	}

	// Calculate expiration
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	c.mu.Lock()
	for key, value := range items {
		// Marshal value to JSON
		data, err := json.Marshal(value)
		if err != nil {
			c.mu.Unlock()
			return err
		}

		// Store item in cache
		c.items[c.prefix+key] = item{
			Value:      data,
			Expiration: expiration,
		}
	}
	c.mu.Unlock()

	return nil
}

// DeleteMulti removes multiple values from the cache
func (c *MemoryCache) DeleteMulti(ctx context.Context, keys []string) error {
	c.mu.Lock()
	for _, key := range keys {
		delete(c.items, c.prefix+key)
	}
	c.mu.Unlock()
	return nil
}

// Increment increments a counter in the cache
func (c *MemoryCache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get current value
	item, found := c.items[c.prefix+key]
	if !found {
		// Key doesn't exist, create it
		data, err := json.Marshal(value)
		if err != nil {
			return 0, err
		}

		c.items[c.prefix+key] = item{
			Value:      data,
			Expiration: 0,
		}

		return value, nil
	}

	// Check if item has expired
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		// Key has expired, create it
		data, err := json.Marshal(value)
		if err != nil {
			return 0, err
		}

		c.items[c.prefix+key] = item{
			Value:      data,
			Expiration: item.Expiration,
		}

		return value, nil
	}

	// Unmarshal current value
	var currentValue int64
	if err := json.Unmarshal(item.Value, &currentValue); err != nil {
		return 0, err
	}

	// Increment value
	newValue := currentValue + value

	// Marshal new value
	data, err := json.Marshal(newValue)
	if err != nil {
		return 0, err
	}

	// Store new value
	c.items[c.prefix+key] = item{
		Value:      data,
		Expiration: item.Expiration,
	}

	return newValue, nil
}

// Decrement decrements a counter in the cache
func (c *MemoryCache) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return c.Increment(ctx, key, -value)
}

// SetNX sets a value in the cache only if the key does not exist
func (c *MemoryCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if key exists and is not expired
	item, found := c.items[c.prefix+key]
	if found && (item.Expiration == 0 || time.Now().UnixNano() < item.Expiration) {
		return false, nil
	}

	// Use default expiration if ttl is 0
	if ttl == 0 {
		ttl = c.defaultExpiration
	}

	// Marshal value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	// Calculate expiration
	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	// Store item in cache
	c.items[c.prefix+key] = item{
		Value:      data,
		Expiration: expiration,
	}

	return true, nil
}

// startCleanup starts the cleanup process
func (c *MemoryCache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopCleanup:
			return
		}
	}
}

// deleteExpired deletes expired items from the cache
func (c *MemoryCache) deleteExpired() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	for key, item := range c.items {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.items, key)
		}
	}
	c.mu.Unlock()
}

// Close stops the cleanup process
func (c *MemoryCache) Close() error {
	close(c.stopCleanup)
	return nil
}
