// Package storage provides storage functionality for vorpalstacks.
package storage

import (
	"encoding/json"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage/pebbledb"
)

type idempotencyEntry struct {
	Token     string `json:"token"`
	Result    []byte `json:"result,omitempty"`
	CreatedAt int64  `json:"created_at"`
	Completed bool   `json:"completed"`
}

// PebbleIdempotencyStore stores idempotency tokens in a Pebble database.
// It is used to ensure that operations can be safely retried without
// duplicate execution.
type PebbleIdempotencyStore struct {
	db     *pebbledb.DB
	prefix []byte
	mu     sync.RWMutex
}

// NewPebbleIdempotencyStore creates a new Pebble idempotency store.
func NewPebbleIdempotencyStore(db *pebbledb.DB) *PebbleIdempotencyStore {
	return &PebbleIdempotencyStore{
		db:     db,
		prefix: []byte(IdempotencyPrefix),
	}
}

func (s *PebbleIdempotencyStore) makeKey(token string) []byte {
	return append(s.prefix, []byte(token)...)
}

// CheckAndStore checks if an idempotency token exists and stores it if not.
// Returns true if this is a new token, false if it already exists.
func (s *PebbleIdempotencyStore) CheckAndStore(token string, ttl time.Duration) (isNew bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.makeKey(token)

	data, err := s.db.Get(key)
	if err != nil && err != pebbledb.ErrKeyNotFound {
		return false, err
	}

	if err == nil {
		var entry idempotencyEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			return false, err
		}
		return false, nil
	}

	newEntry := idempotencyEntry{
		Token:     token,
		CreatedAt: time.Now().Unix(),
		Completed: false,
	}

	entryData, err := json.Marshal(newEntry)
	if err != nil {
		return false, err
	}

	if err := s.db.SetWithTTL(key, entryData, ttl); err != nil {
		return false, err
	}

	return true, nil
}

// StoreResult stores the result of an idempotent operation.
// This should be called after the operation completes successfully.
func (s *PebbleIdempotencyStore) StoreResult(token string, result []byte, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.makeKey(token)

	data, err := s.db.Get(key)
	if err != nil {
		return err
	}

	var entry idempotencyEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return err
	}

	entry.Result = result
	entry.Completed = true

	entryData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return s.db.SetWithTTL(key, entryData, ttl)
}

// GetResult retrieves the stored result for an idempotency token.
// Returns the result, whether it exists and is complete, and any error.
func (s *PebbleIdempotencyStore) GetResult(token string) ([]byte, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.makeKey(token)

	data, err := s.db.Get(key)
	if err != nil {
		if err == pebbledb.ErrKeyNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}

	var entry idempotencyEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false, err
	}

	if !entry.Completed {
		return nil, false, nil
	}

	return entry.Result, true, nil
}

// Delete removes an idempotency token from the store.
func (s *PebbleIdempotencyStore) Delete(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.makeKey(token)
	return s.db.Delete(key)
}
