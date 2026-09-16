package scheduler

import (
	"encoding/json"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
)

// clientTokenTTL is the time-to-live for idempotency token mappings.
// AWS guarantees idempotency for at least 8 hours; we use 24 hours to
// be conservative (Smithy idempotencyToken trait does not specify TTL).
const clientTokenTTL = 24 * time.Hour

// tokenKeyPrefix is the Pebble key prefix for persisted ClientToken entries.
const tokenKeyPrefix = "token:"

// ClientTokenEntry records the resource created for a given ClientToken.
// The claim's operation and resource scope live in the storage key alone
// (clientTokenKey); the entry carries only what a replay needs to answer.
type ClientTokenEntry struct {
	ResourceArn string    `json:"resourceArn"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ClientTokenStore provides idempotency token deduplication for the
// schedule create, update and delete operations. Entries are keyed by the
// operation, the resource ARN and the token, and persisted to Pebble so
// that idempotency survives server restarts.
type ClientTokenStore struct {
	mu       sync.Mutex
	store    *common.BaseStore
	stopCh   chan struct{}
	stopOnce sync.Once
}

// NewClientTokenStore creates a new ClientTokenStore backed by the given
// BaseStore (a Pebble bucket). Existing entries are reaped on construction
// so that expired leftovers from a previous server run are removed while
// non-expired entries are kept, preserving idempotency across restarts.
func NewClientTokenStore(store *common.BaseStore) *ClientTokenStore {
	s := &ClientTokenStore{
		store:  store,
		stopCh: make(chan struct{}),
	}
	s.reapExpired()
	go s.cleanupLoop()
	return s
}

// reapExpired scans the token bucket and removes every entry older than
// clientTokenTTL. It runs once on construction and hourly from cleanupLoop.
// The mutex keeps the scan-delete atomic with LookupOrClaim's claim write:
// without it the reaper could read an expired value, let a concurrent claim
// overwrite the key, and then delete that fresh claim — losing the
// idempotency mapping a client retry depends on.
func (s *ClientTokenStore) reapExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	_ = s.store.ScanPrefix(tokenKeyPrefix, func(key string, value []byte) error {
		var entry ClientTokenEntry
		if err := json.Unmarshal(value, &entry); err != nil {
			return nil
		}
		if now.Sub(entry.CreatedAt) >= clientTokenTTL {
			_ = s.store.Delete(key)
		}
		return nil
	})
}

// clientTokenKey scopes an idempotency claim to one operation on one
// resource. A token reused by a different request (another operation, or
// the same operation on another resource) claims a separate key instead of
// replaying the first request's outcome.
func clientTokenKey(resourceType, resourceArn, token string) string {
	return tokenKeyPrefix + resourceType + ":" + resourceArn + ":" + token
}

// LookupOrClaim checks if a ClientToken already maps to a resource within
// the same operation and resource scope. If found (and not expired), returns
// the existing entry and false. If not found, claims the token and persists
// it to Pebble.
func (s *ClientTokenStore) LookupOrClaim(token, resourceArn, resourceType string) (*ClientTokenEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := clientTokenKey(resourceType, resourceArn, token)

	// Check Pebble for existing entry.
	data, err := s.store.GetRaw(key)
	if err == nil && data != nil {
		var entry ClientTokenEntry
		if json.Unmarshal(data, &entry) == nil {
			if time.Since(entry.CreatedAt) < clientTokenTTL {
				return &entry, false
			}
		}
	}

	// Claim the token.
	entry := &ClientTokenEntry{
		ResourceArn: resourceArn,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.store.PutRaw(key, mustMarshal(entry)); err != nil {
		logs.Error("Failed to persist ClientToken entry",
			logs.String("token", token),
			logs.Err(err))
	}

	return entry, true
}

// Release removes a ClientToken entry from Pebble. Used to roll back a
// claim when resource creation fails after the token was claimed. The scope
// arguments must match the LookupOrClaim call that claimed the token.
func (s *ClientTokenStore) Release(token, resourceArn, resourceType string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.store.Delete(clientTokenKey(resourceType, resourceArn, token))
}

// cleanupLoop periodically reaps expired entries from Pebble.
func (s *ClientTokenStore) cleanupLoop() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.reapExpired()
		case <-s.stopCh:
			return
		}
	}
}

// Stop shuts down the background cleanup goroutine. It is idempotent:
// StopEngine closes every cached store without evicting the cache, so a
// repeated shutdown must not close an already-closed channel.
func (s *ClientTokenStore) Stop() {
	s.stopOnce.Do(func() { close(s.stopCh) })
}

// mustMarshal serialises a ClientTokenEntry to JSON. Panics are impossible
// because ClientTokenEntry contains only primitive types.
func mustMarshal(entry *ClientTokenEntry) []byte {
	b, _ := json.Marshal(entry)
	return b
}
