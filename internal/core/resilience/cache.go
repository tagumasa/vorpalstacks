// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import (
	"fmt"
	"sync"
	"time"
)

// CacheItem represents a single item stored in the cache.
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// Expired returns true if the cache item has expired.
func (ci *CacheItem) Expired() bool {
	return !ci.Expiration.IsZero() && time.Now().After(ci.Expiration)
}

// CacheConfig defines the configuration for a cache.
type CacheConfig struct {
	TTL             time.Duration
	MaxSize         int
	CleanupInterval time.Duration
}

// DefaultCacheConfig returns the default cache configuration.
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		TTL:             5 * time.Minute,
		MaxSize:         1000,
		CleanupInterval: 1 * time.Minute,
	}
}

// Cache provides an in-memory cache with TTL support.
type Cache struct {
	items    map[string]*CacheItem
	mu       sync.RWMutex
	config   *CacheConfig
	stopChan chan struct{}
	stopOnce sync.Once
}

// NewCache creates a new cache with the given configuration.
func NewCache(config *CacheConfig) *Cache {
	if config == nil {
		config = DefaultCacheConfig()
	}

	cache := &Cache{
		items:    make(map[string]*CacheItem),
		config:   config,
		stopChan: make(chan struct{}),
	}

	go cache.cleanup()

	return cache
}

// Get retrieves a value from the cache by key.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	if item.Expired() {
		return nil, false
	}

	return item.Value, true
}

// Set stores a value in the cache with the default TTL.
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.config.TTL)
}

// SetWithTTL stores a value in the cache with a specific TTL.
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) >= c.config.MaxSize {
		c.evictOne()
	}

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	c.items[key] = &CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Delete removes a value from the cache by key.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all items from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*CacheItem)
}

// Size returns the number of items in the cache.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// Keys returns all keys in the cache.
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}

	return keys
}

func (c *Cache) evictOne() {
	var oldestKey string
	var oldestTime time.Time
	var hasTTL bool

	for key, item := range c.items {
		if item.Expiration.IsZero() {
			delete(c.items, key)
			return
		}

		if !hasTTL || item.Expiration.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.Expiration
			hasTTL = true
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

func (c *Cache) cleanup() {
	if c.config.CleanupInterval <= 0 {
		return
	}

	ticker := time.NewTicker(c.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stopChan:
			return
		}
	}
}

func (c *Cache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if item.Expired() {
			delete(c.items, key)
		}
	}
}

// Stop stops the cache cleanup goroutine.
func (c *Cache) Stop() {
	c.stopOnce.Do(func() {
		close(c.stopChan)
	})
}

// Stats returns the current statistics of the cache.
func (c *Cache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	expiredCount := 0
	for _, item := range c.items {
		if item.Expired() {
			expiredCount++
		}
	}

	return CacheStats{
		Size:         len(c.items),
		MaxSize:      c.config.MaxSize,
		ExpiredCount: expiredCount,
		TTL:          c.config.TTL,
	}
}

// CacheStats represents statistics about the cache.
type CacheStats struct {
	Size         int
	MaxSize      int
	ExpiredCount int
	TTL          time.Duration
}

// String returns the string representation of the CacheStats.
func (s CacheStats) String() string {
	return fmt.Sprintf("Size: %d/%d, Expired: %d, TTL: %s", s.Size, s.MaxSize, s.ExpiredCount, s.TTL)
}
