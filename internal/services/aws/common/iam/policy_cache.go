// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"time"

	"vorpalstacks/internal/services/aws/iam/policy"

	"github.com/dgraph-io/ristretto/v2"
)

// PolicyCache caches IAM policies for performance.
type PolicyCache struct {
	cache *ristretto.Cache[string, *policy.Document]
}

// NewPolicyCache creates a new policy cache.
func NewPolicyCache() (*PolicyCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, *policy.Document]{
		NumCounters: DefaultCacheNumCounters,
		MaxCost:     DefaultCacheMaxCost,
		BufferItems: DefaultCacheBufferItems,
	})
	if err != nil {
		return nil, err
	}
	return &PolicyCache{cache: cache}, nil
}

// Get retrieves a policy from the cache.
func (c *PolicyCache) Get(roleName string) (*policy.Document, bool) {
	return c.cache.Get(roleName)
}

// Set stores a policy in the cache.
func (c *PolicyCache) Set(roleName string, doc *policy.Document) bool {
	return c.cache.Set(roleName, doc, 1)
}

// SetWithTTL stores a policy in the cache with a time-to-live.
func (c *PolicyCache) SetWithTTL(roleName string, doc *policy.Document, ttl time.Duration) bool {
	return c.cache.SetWithTTL(roleName, doc, 1, ttl)
}

// Delete removes a policy from the cache.
func (c *PolicyCache) Delete(roleName string) {
	c.cache.Del(roleName)
}

// Close closes the policy cache.
func (c *PolicyCache) Close() {
	c.cache.Close()
}
