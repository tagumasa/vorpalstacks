package scheduler

import (
	"sync"
	"time"
)

// clientTokenTTL is the time-to-live for idempotency token mappings.
// AWS guarantees idempotency for at least 8 hours; we use 24 hours to
// be conservative (Smithy idempotencyToken trait does not specify TTL).
const clientTokenTTL = 24 * time.Hour

// ClientTokenEntry records the resource created for a given ClientToken.
type ClientTokenEntry struct {
	ResourceArn  string    `json:"resourceArn"`
	ResourceType string    `json:"resourceType"` // "schedule" or "schedule-group"
	CreatedAt    time.Time `json:"createdAt"`
}

// ClientTokenStore provides idempotency token deduplication for
// CreateSchedule and CreateScheduleGroup operations. When a ClientToken
// is reused within the TTL window, the original resource ARN is returned
// without creating a duplicate.
//
// The store uses an in-memory map guarded by a mutex. Entries are
// expired lazily on read. A background goroutine runs periodic cleanup
// to bound memory usage.
type ClientTokenStore struct {
	mu      sync.Mutex
	entries map[string]*ClientTokenEntry
	stopCh  chan struct{}
}

// NewClientTokenStore creates a new ClientTokenStore and starts the
// background cleanup goroutine.
func NewClientTokenStore() *ClientTokenStore {
	s := &ClientTokenStore{
		entries: make(map[string]*ClientTokenEntry),
		stopCh:  make(chan struct{}),
	}
	go s.cleanupLoop()
	return s
}

// LookupOrClaim checks if a ClientToken already maps to a resource.
// If found (and not expired), returns the existing entry and false.
// If not found, claims the token with the given resource details and
// returns the new entry and true.
func (s *ClientTokenStore) LookupOrClaim(token, resourceArn, resourceType string) (*ClientTokenEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for existing entry (lazy expiry).
	if entry, ok := s.entries[token]; ok {
		if time.Since(entry.CreatedAt) < clientTokenTTL {
			return entry, false
		}
		// Expired — remove and allow re-claim.
		delete(s.entries, token)
	}

	// Claim the token.
	entry := &ClientTokenEntry{
		ResourceArn:  resourceArn,
		ResourceType: resourceType,
		CreatedAt:    time.Now().UTC(),
	}
	s.entries[token] = entry
	return entry, true
}

// Release removes a ClientToken entry. Used to roll back a claim when
// resource creation fails after the token was claimed.
func (s *ClientTokenStore) Release(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, token)
}

// cleanupLoop periodically removes expired entries to bound memory.
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
	for token, entry := range s.entries {
		if now.Sub(entry.CreatedAt) >= clientTokenTTL {
			delete(s.entries, token)
		}
	}
}

// Stop shuts down the background cleanup goroutine.
func (s *ClientTokenStore) Stop() {
	close(s.stopCh)
}
