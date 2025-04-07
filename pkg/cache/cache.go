package cache

import (
	"sync"
	"time"
)

// Item represents a cached item
type Item struct {
	Value      interface{}
	Expiration int64
}

// Cache is a simple in-memory cache
type Cache struct {
	items             map[string]Item
	mu                sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	stopCleanup       chan bool
}

// Config holds cache configuration
type Config struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

// DefaultConfig returns default cache configuration
func DefaultConfig() Config {
	return Config{
		DefaultExpiration: 5 * time.Minute,
		CleanupInterval:   10 * time.Minute,
	}
}

// New creates a new cache with the given default expiration and cleanup interval
func New(cfg Config) *Cache {
	cache := &Cache{
		items:             make(map[string]Item),
		defaultExpiration: cfg.DefaultExpiration,
		cleanupInterval:   cfg.CleanupInterval,
		stopCleanup:       make(chan bool),
	}

	// Start cleanup goroutine if cleanup interval is greater than 0
	if cfg.CleanupInterval > 0 {
		go cache.startCleanup()
	}

	return cache
}

// Set adds an item to the cache with the default expiration
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithExpiration(key, value, c.defaultExpiration)
}

// SetWithExpiration adds an item to the cache with a custom expiration
func (c *Cache) SetWithExpiration(key string, value interface{}, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		// 0 means use default expiration
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.mu.Lock()
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
	c.mu.Unlock()
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	// Return nil, false if item not found
	if !found {
		return nil, false
	}

	// Return nil, false if item has expired
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Value, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Flush removes all items from the cache
func (c *Cache) Flush() {
	c.mu.Lock()
	c.items = make(map[string]Item)
	c.mu.Unlock()
}

// startCleanup starts the cleanup process
func (c *Cache) startCleanup() {
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
func (c *Cache) deleteExpired() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// Stop stops the cleanup process
func (c *Cache) Stop() {
	c.stopCleanup <- true
}
