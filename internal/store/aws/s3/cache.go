package s3

import (
	"github.com/dgraph-io/ristretto/v2"
)

const (
	// VersioningCacheNumCounters is the number of counters for the versioning cache.
	VersioningCacheNumCounters = 1000
	// VersioningCacheMaxCost is the maximum cost of the versioning cache.
	VersioningCacheMaxCost = 1 << 20
	// VersioningCacheBufferItems is the buffer size for the versioning cache.
	VersioningCacheBufferItems = 64
)

// VersioningCache caches versioning status for S3 buckets.
type VersioningCache struct {
	cache *ristretto.Cache[string, bool]
}

// NewVersioningCache creates a new versioning cache instance.
func NewVersioningCache() (*VersioningCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, bool]{
		NumCounters: VersioningCacheNumCounters,
		MaxCost:     VersioningCacheMaxCost,
		BufferItems: VersioningCacheBufferItems,
	})
	if err != nil {
		return nil, err
	}
	return &VersioningCache{cache: cache}, nil
}

// Get retrieves the versioning status for a bucket from the cache.
func (c *VersioningCache) Get(bucket string) (bool, bool) {
	return c.cache.Get(bucket)
}

// Set stores the versioning status for a bucket in the cache.
func (c *VersioningCache) Set(bucket string, enabled bool) {
	c.cache.Set(bucket, enabled, 1)
}

// Delete removes a bucket's versioning status from the cache.
func (c *VersioningCache) Delete(bucket string) {
	c.cache.Del(bucket)
}

// Close closes the versioning cache and releases resources.
func (c *VersioningCache) Close() {
	c.cache.Close()
}
