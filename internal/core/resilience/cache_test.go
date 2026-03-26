package resilience

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache := NewCache(&CacheConfig{
		TTL:             time.Minute,
		MaxSize:         100,
		CleanupInterval: time.Hour,
	})
	defer cache.Stop()

	t.Run("Set and Get", func(t *testing.T) {
		cache.Set("key1", "value1")
		val, ok := cache.Get("key1")
		assert.True(t, ok)
		assert.Equal(t, "value1", val)
	})

	t.Run("Get non-existent key", func(t *testing.T) {
		_, ok := cache.Get("non-existent")
		assert.False(t, ok)
	})

	t.Run("Delete", func(t *testing.T) {
		cache.Set("key2", "value2")
		cache.Delete("key2")
		_, ok := cache.Get("key2")
		assert.False(t, ok)
	})

	t.Run("Clear", func(t *testing.T) {
		cache.Set("key3", "value3")
		cache.Set("key4", "value4")
		cache.Clear()
		_, ok := cache.Get("key3")
		assert.False(t, ok)
		_, ok = cache.Get("key4")
		assert.False(t, ok)
	})

	t.Run("Update updates existing key", func(t *testing.T) {
		cache.Set("key5", "old-value")
		cache.Set("key5", "new-value")
		val, ok := cache.Get("key5")
		assert.True(t, ok)
		assert.Equal(t, "new-value", val)
	})
}

func TestCacheItemExpired(t *testing.T) {
	t.Run("Not expired", func(t *testing.T) {
		item := &CacheItem{
			Value:      "test",
			Expiration: time.Now().Add(time.Hour),
		}
		assert.False(t, item.Expired())
	})

	t.Run("Expired", func(t *testing.T) {
		item := &CacheItem{
			Value:      "test",
			Expiration: time.Now().Add(-time.Hour),
		}
		assert.True(t, item.Expired())
	})

	t.Run("Zero expiration", func(t *testing.T) {
		item := &CacheItem{
			Value:      "test",
			Expiration: time.Time{},
		}
		assert.False(t, item.Expired())
	})
}

func TestDefaultCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()
	assert.NotNil(t, config)
	assert.Equal(t, DefaultCacheTTL, config.TTL)
	assert.Equal(t, DefaultCacheMaxSize, config.MaxSize)
	assert.Equal(t, DefaultCacheCleanupInterval, config.CleanupInterval)
}

func TestCacheEviction(t *testing.T) {
	cache := NewCache(&CacheConfig{
		TTL:             time.Minute,
		MaxSize:         3,
		CleanupInterval: time.Hour,
	})
	defer cache.Stop()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	cache.Set("key4", "value4")

	time.Sleep(10 * time.Millisecond)
}
