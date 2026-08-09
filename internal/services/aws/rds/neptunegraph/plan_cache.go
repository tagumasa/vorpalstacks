package neptunegraph

import (
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// planCacheEntry holds a cached query plan with an expiry timestamp.
type planCacheEntry struct {
	key       string
	plan      interface{}
	expiresAt time.Time
}

// queryPlanCache is an LRU cache for Cypher query plans with a 5-minute
// TTL and a capacity of 1000 entries, matching the AWS Neptune Analytics
// specification for the planCache feature.
type queryPlanCache struct {
	mu       sync.Mutex
	items    map[string]*list.Element
	lru      *list.List
	capacity int
	ttl      time.Duration
}

// newQueryPlanCache creates a plan cache with the specified capacity and TTL.
func newQueryPlanCache(capacity int, ttl time.Duration) *queryPlanCache {
	return &queryPlanCache{
		items:    make(map[string]*list.Element, capacity),
		lru:      list.New(),
		capacity: capacity,
		ttl:      ttl,
	}
}

// planCacheKey computes a deterministic cache key from graph ID, query
// string, and stringified parameters. The key is prefixed with the graph
// ID so that purgeByGraph can efficiently evict all entries for a deleted
// or reset graph. Map keys are sorted before hashing to guarantee a
// stable digest regardless of Go's randomised map iteration order.
func planCacheKey(graphID, query string, params map[string]any) string {
	h := sha256.New()
	h.Write([]byte(graphID))
	h.Write([]byte{0})
	h.Write([]byte(query))
	if len(params) > 0 {
		h.Write([]byte{0})
		keys := make([]string, 0, len(params))
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var sb strings.Builder
		for _, k := range keys {
			sb.WriteString(k)
			sb.WriteByte('=')
			sb.WriteString(strings.TrimSpace(strings.Trim(stringifyValue(params[k]), `"`)))
			sb.WriteByte(';')
		}
		h.Write([]byte(sb.String()))
	}
	return graphID + ":" + hex.EncodeToString(h.Sum(nil))
}

// get returns a cached plan if it exists and has not expired. Expired
// entries are evicted on access.
func (c *queryPlanCache) get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}
	entry := elem.Value.(*planCacheEntry)
	if time.Now().After(entry.expiresAt) {
		c.lru.Remove(elem)
		delete(c.items, key)
		return nil, false
	}
	c.lru.MoveToFront(elem)
	return entry.plan, true
}

// put stores a plan in the cache, evicting the least-recently-used entry
// if the cache is at capacity.
func (c *queryPlanCache) put(key string, plan interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok {
		entry := elem.Value.(*planCacheEntry)
		entry.plan = plan
		entry.expiresAt = time.Now().Add(c.ttl)
		c.lru.MoveToFront(elem)
		return
	}
	entry := &planCacheEntry{
		key:       key,
		plan:      plan,
		expiresAt: time.Now().Add(c.ttl),
	}
	elem := c.lru.PushFront(entry)
	c.items[key] = elem
	if c.lru.Len() > c.capacity {
		oldest := c.lru.Back()
		if oldest != nil {
			c.lru.Remove(oldest)
			delete(c.items, oldest.Value.(*planCacheEntry).key)
		}
	}
}

// purge removes all entries for a specific graph ID. This is called when
// a graph is deleted or reset.
func (c *queryPlanCache) purgeByGraph(graphID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	prefix := graphID + ":"
	for key, elem := range c.items {
		if strings.HasPrefix(key, prefix) {
			c.lru.Remove(elem)
			delete(c.items, key)
		}
	}
}

// stringifyValue converts a parameter value to a string for cache key
// generation.
func stringifyValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strings.TrimRight(strings.TrimRight(
			strconv.FormatFloat(val, 'f', -1, 64), "0"), ".")
	case int:
		return strconv.Itoa(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
