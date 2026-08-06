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
type ClientTokenEntry struct {
	ResourceArn  string    `json:"resourceArn"`
	ResourceType string    `json:"resourceType"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ClientTokenStore provides idempotency token deduplication for
// CreateSchedule and CreateScheduleGroup operations. Entries are persisted
// to Pebble so that idempotency survives server restarts.
type ClientTokenStore struct {
	mu     sync.Mutex
	store  *common.BaseStore
	stopCh chan struct{}
}

// NewClientTokenStore creates a new ClientTokenStore backed by the given
// BaseStore (a Pebble bucket). Existing entries are loaded on construction
// so that previously claimed tokens are honoured after a server restart.
func NewClientTokenStore(store *common.BaseStore) *ClientTokenStore {
	s := &ClientTokenStore{
		store:  store,
		stopCh: make(chan struct{}),
	}
	s.loadExisting()
	go s.cleanupLoop()
	return s
}

// loadExisting scans the Pebble bucket and removes expired entries left
// over from a previous server run. Non-expired entries are kept so that
// idempotency is preserved across restarts.
func (s *ClientTokenStore) loadExisting() {
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

// LookupOrClaim checks if a ClientToken already maps to a resource.
// If found (and not expired), returns the existing entry and false.
// If not found, claims the token and persists it to Pebble.
func (s *ClientTokenStore) LookupOrClaim(token, resourceArn, resourceType string) (*ClientTokenEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := tokenKeyPrefix + token

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
		ResourceArn:  resourceArn,
		ResourceType: resourceType,
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.store.PutRaw(key, mustMarshal(entry)); err != nil {
		logs.Error("Failed to persist ClientToken entry",
			logs.String("token", token),
			logs.Err(err))
	}

	return entry, true
}

// Release removes a ClientToken entry from Pebble. Used to roll back a
// claim when resource creation fails after the token was claimed.
func (s *ClientTokenStore) Release(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.store.Delete(tokenKeyPrefix + token)
}

// cleanupLoop periodically removes expired entries from Pebble.
func (s *ClientTokenStore) cleanupLoop() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.cleanupExpired()
		case <-s.stopCh:
			return
		}
	}
}

// cleanupExpired removes all entries older than clientTokenTTL.
func (s *ClientTokenStore) cleanupExpired() {
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

// Stop shuts down the background cleanup goroutine.
func (s *ClientTokenStore) Stop() {
	close(s.stopCh)
}

// mustMarshal serialises a ClientTokenEntry to JSON. Panics are impossible
// because ClientTokenEntry contains only primitive types.
func mustMarshal(entry *ClientTokenEntry) []byte {
	b, _ := json.Marshal(entry)
	return b
}
