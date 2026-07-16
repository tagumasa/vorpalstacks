package appsync

import (
	"encoding/json"
	"time"
)

// ResolverCacheEntry stores a cached resolver result with TTL metadata.
// The cache is scoped per-API and per-resolver-field, keyed by a SHA-256
// hash of the resolver coordinates and resolved caching context paths.
type ResolverCacheEntry struct {
	Result   json.RawMessage `json:"result"`
	CachedAt int64           `json:"cachedAt"` // epoch seconds
	TTL      int64           `json:"ttl"`      // seconds, 1-3600 per AWS spec
}

// resolverCacheBucketName returns the PebbleDB bucket name for resolver cache entries.
func resolverCacheBucketName(region string) string {
	return "appsync-resolver-cache-" + region
}

// GetResolverCacheEntry retrieves a cached resolver result by API ID and cache key.
// Returns an error if the key is not found (Pebble ErrNotFound).
func (s *AppSyncStore) GetResolverCacheEntry(apiId, cacheKey string) (*ResolverCacheEntry, error) {
	key := apiId + "/" + cacheKey
	var entry ResolverCacheEntry
	if err := s.resolverCacheStore.Get(key, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// PutResolverCacheEntry stores a resolver result in the cache.
func (s *AppSyncStore) PutResolverCacheEntry(apiId, cacheKey string, entry *ResolverCacheEntry) error {
	key := apiId + "/" + cacheKey
	return s.resolverCacheStore.Put(key, entry)
}

// FlushResolverCache removes all cached resolver results for a given API.
// Uses prefix scan to delete all entries under {apiId}/.
func (s *AppSyncStore) FlushResolverCache(apiId string) error {
	return s.resolverCacheStore.DeleteByPrefix(apiId + "/")
}

// IsExpired reports whether the cache entry has exceeded its TTL.
func (e *ResolverCacheEntry) IsExpired() bool {
	if e.TTL <= 0 {
		return true
	}
	return time.Now().Unix()-e.CachedAt >= e.TTL
}
