package pebbledb

import (
	"context"
	"time"
)

// TTLManager provides a convenient interface for managing time-to-live on keys.
// It wraps database operations with TTL functionality.
type TTLManager struct {
	db *DB
}

// NewTTLManager creates a new TTL manager for the given database.
// This provides a convenient interface for managing time-to-live on keys.
func NewTTLManager(db *DB) *TTLManager {
	return &TTLManager{db: db}
}

// SetWithExpiry sets a key-value pair with a time-to-live expiry.
func (m *TTLManager) SetWithExpiry(key, value []byte, ttl time.Duration) error {
	return m.db.SetWithTTL(key, value, ttl)
}

// Get retrieves a value and reports whether it exists.
func (m *TTLManager) Get(key []byte) ([]byte, bool, error) {
	val, err := m.db.Get(key)
	if err != nil {
		if err == ErrKeyNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return val, true, nil
}

// Delete removes a key from the database.
func (m *TTLManager) Delete(key []byte) error {
	return m.db.Delete(key)
}

// Run starts the TTL cleanup background task.
func (m *TTLManager) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.db.opts.TTL.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			m.db.cleanExpiredKeys()
		}
	}
}

// CleanupNow immediately cleans up expired keys.
func (m *TTLManager) CleanupNow() {
	m.db.cleanExpiredKeys()
}

// TTLManager returns a TTL manager for this database.
func (d *DB) TTLManager() *TTLManager {
	return NewTTLManager(d)
}
