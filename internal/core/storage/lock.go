// Package storage provides storage functionality for vorpalstacks.
package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"vorpalstacks/internal/core/storage/pebbledb"

	"github.com/cockroachdb/pebble/v2"
)

type lockEntry struct {
	Token     string `json:"token"`
	Mode      int    `json:"mode"`
	Version   uint64 `json:"version"`
	ExpiresAt int64  `json:"expires_at"`
}

// PebbleLockManager provides distributed locking capabilities using a Pebble database.
// It supports shared and exclusive locks with automatic expiration.
type PebbleLockManager struct {
	db             *pebbledb.DB
	prefix         []byte
	maxNonConflict int
}

// NewPebbleLockManager creates a new lock manager backed by Pebble.
func NewPebbleLockManager(db *pebbledb.DB, writeOpts *pebble.WriteOptions) *PebbleLockManager {
	return &PebbleLockManager{
		db:             db,
		prefix:         []byte(LockPrefix),
		maxNonConflict: 10,
	}
}

func (m *PebbleLockManager) makeKey(key []byte) []byte {
	k := make([]byte, len(m.prefix)+len(key))
	copy(k, m.prefix)
	copy(k[len(m.prefix):], key)
	return k
}

func (m *PebbleLockManager) generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("crypto/rand failed: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (m *PebbleLockManager) getLockEntry(fullKey []byte) (*lockEntry, error) {
	data, err := m.db.Get(fullKey)
	if err != nil {
		if errors.Is(err, pebbledb.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var entry lockEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func (m *PebbleLockManager) tryAcquire(key []byte, mode LockMode, ttl time.Duration) (*LockHandle, error) {
	fullKey := m.makeKey(key)

	for i := 0; i < DefaultLockMaxRetries; i++ {
		existing, err := m.getLockEntry(fullKey)
		if err != nil {
			return nil, err
		}

		var existingVersion uint64
		if existing != nil {
			existingVersion = existing.Version
			return nil, &LockConflictError{
				Key: key,
				ExistingHandle: &LockHandle{
					Key:       key,
					Token:     existing.Token,
					Mode:      LockMode(existing.Mode),
					ExpiresAt: existing.ExpiresAt,
				},
			}
		}

		newVersion := existingVersion + 1
		token, err := m.generateToken()
		if err != nil {
			return nil, fmt.Errorf("failed to generate lock token: %w", err)
		}

		now := time.Now().Unix()
		newEntry := lockEntry{
			Token:     token,
			Mode:      int(mode),
			Version:   newVersion,
			ExpiresAt: now + int64(ttl.Seconds()),
		}

		entryData, err := json.Marshal(newEntry)
		if err != nil {
			return nil, err
		}

		batch := m.db.NewBatch()
		if err := batch.Set(fullKey, entryData); err != nil {
			batch.Close()
			return nil, err
		}
		if err := batch.Commit(pebble.Sync); err != nil {
			batch.Close()
			return nil, err
		}
		batch.Close()

		verifyEntry, err := m.getLockEntry(fullKey)
		if err != nil {
			continue
		}
		if verifyEntry != nil && verifyEntry.Token == newEntry.Token && verifyEntry.Version == newVersion {
			return &LockHandle{
				Key:       key,
				Token:     newEntry.Token,
				Mode:      mode,
				ExpiresAt: newEntry.ExpiresAt,
			}, nil
		}
	}

	return nil, &LockConflictError{Key: key}
}

// TryLock attempts to acquire a lock without blocking.
// Returns a lock handle if successful, or a conflict error if the lock is held.
func (m *PebbleLockManager) TryLock(key []byte, mode LockMode, ttl time.Duration) (*LockHandle, error) {
	return m.tryAcquire(key, mode, ttl)
}

// Lock acquires a lock, blocking until it is available.
// The lock is automatically released when the context is cancelled or times out.
func (m *PebbleLockManager) Lock(ctx context.Context, key []byte, mode LockMode, ttl time.Duration) (*LockHandle, error) {
	nonConflictErrors := 0
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		handle, err := m.tryAcquire(key, mode, ttl)
		if err == nil {
			return handle, nil
		}

		var conflictErr *LockConflictError
		if errors.As(err, &conflictErr) {
			timer := time.NewTimer(DefaultLockRetryDelay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
			continue
		}

		nonConflictErrors++
		if nonConflictErrors >= m.maxNonConflict {
			return nil, fmt.Errorf("lock acquisition failed after %d non-conflict errors: %w", nonConflictErrors, err)
		}
		return nil, err
	}
}

// Unlock releases a previously acquired lock.
// The handle token must match or a token mismatch error is returned.
func (m *PebbleLockManager) Unlock(handle *LockHandle) error {
	fullKey := m.makeKey(handle.Key)

	for i := 0; i < DefaultLockMaxRetries; i++ {
		existing, err := m.getLockEntry(fullKey)
		if err != nil {
			return err
		}
		if existing == nil {
			return nil
		}

		if existing.Token != handle.Token {
			return &LockTokenMismatchError{Key: handle.Key}
		}

		batch := m.db.NewBatch()
		if err := batch.Delete(fullKey); err != nil {
			batch.Close()
			return err
		}
		if err := batch.Commit(pebble.Sync); err != nil {
			batch.Close()
			return err
		}
		batch.Close()

		verifyEntry, err := m.getLockEntry(fullKey)
		if err != nil {
			return err
		}
		if verifyEntry == nil {
			return nil
		}
		if verifyEntry.Token != handle.Token {
			return nil
		}
	}

	return nil
}

// Extend extends the TTL of an existing lock.
// The handle token must match or a token mismatch error is returned.
func (m *PebbleLockManager) Extend(handle *LockHandle, ttl time.Duration) error {
	fullKey := m.makeKey(handle.Key)

	for i := 0; i < DefaultLockMaxRetries; i++ {
		existing, err := m.getLockEntry(fullKey)
		if err != nil {
			return err
		}
		if existing == nil {
			return &LockNotFoundError{Key: handle.Key}
		}

		if existing.Token != handle.Token {
			return &LockTokenMismatchError{Key: handle.Key}
		}

		now := time.Now().Unix()
		newEntry := lockEntry{
			Token:     existing.Token,
			Mode:      existing.Mode,
			Version:   existing.Version + 1,
			ExpiresAt: now + int64(ttl.Seconds()),
		}

		entryData, err := json.Marshal(newEntry)
		if err != nil {
			return err
		}

		batch := m.db.NewBatch()
		if err := batch.Set(fullKey, entryData); err != nil {
			batch.Close()
			return err
		}
		if err := batch.Commit(pebble.Sync); err != nil {
			batch.Close()
			return err
		}
		batch.Close()

		verifyEntry, err := m.getLockEntry(fullKey)
		if err != nil {
			continue
		}
		if verifyEntry != nil && verifyEntry.Token == handle.Token && verifyEntry.Version == newEntry.Version {
			handle.ExpiresAt = newEntry.ExpiresAt
			return nil
		}
	}

	return &LockConflictError{Key: handle.Key}
}

// IsLocked checks if a key is currently locked.
// Returns whether it is locked, the lock handle if locked, and any error.
func (m *PebbleLockManager) IsLocked(key []byte) (bool, *LockHandle, error) {
	fullKey := m.makeKey(key)

	entry, err := m.getLockEntry(fullKey)
	if err != nil {
		return false, nil, err
	}
	if entry == nil {
		return false, nil, nil
	}

	return true, &LockHandle{
		Key:       key,
		Token:     entry.Token,
		Mode:      LockMode(entry.Mode),
		ExpiresAt: entry.ExpiresAt,
	}, nil
}

// LockConflictError is returned when attempting to acquire a lock that is already held.
type LockConflictError struct {
	Key            []byte
	ExistingHandle *LockHandle
	Err            error
}

// Error returns the error message for the lock conflict.
func (e *LockConflictError) Error() string {
	if e.Err != nil {
		return "lock conflict: " + e.Err.Error()
	}
	return "lock conflict: key is already locked"
}

// Unwrap returns the underlying error, if any.
func (e *LockConflictError) Unwrap() error { return e.Err }

// LockTokenMismatchError is returned when a lock token does not match the existing lock.
type LockTokenMismatchError struct {
	Key []byte
	Err error
}

// Error returns the error message for the lock token mismatch.
func (e *LockTokenMismatchError) Error() string {
	if e.Err != nil {
		return "lock token mismatch: " + e.Err.Error()
	}
	return "lock token mismatch"
}

// Unwrap returns the underlying error, if any.
func (e *LockTokenMismatchError) Unwrap() error { return e.Err }

// LockNotFoundError is returned when attempting to operate on a lock that does not exist.
type LockNotFoundError struct {
	Key []byte
	Err error
}

// Error returns the error message when the lock is not found.
func (e *LockNotFoundError) Error() string {
	if e.Err != nil {
		return "lock not found: " + e.Err.Error()
	}
	return "lock not found"
}

// Unwrap returns the underlying error, if any.
func (e *LockNotFoundError) Unwrap() error { return e.Err }
