package storage

import (
	"sync"

	"vorpalstacks/internal/core/storage/pebbledb"
)

// ConditionalPebbleBucket provides conditional write operations on a Pebble database.
// It supports optimistic locking patterns where writes are only performed if a
// specified condition (e.g., key does not exist, value matches) is satisfied.
type ConditionalPebbleBucket struct {
	db     *pebbledb.DB
	prefix []byte
	mu     sync.Mutex
}

// NewConditionalPebbleBucket creates a new conditional bucket for optimistic locking operations.
func NewConditionalPebbleBucket(db *pebbledb.DB, prefix []byte) *ConditionalPebbleBucket {
	return &ConditionalPebbleBucket{
		db:     db,
		prefix: prefix,
	}
}

func (b *ConditionalPebbleBucket) makeKey(key []byte) []byte {
	k := make([]byte, len(b.prefix)+len(key))
	copy(k, b.prefix)
	copy(k[len(b.prefix):], key)
	return k
}

// PutIf atomically checks a condition before writing a value.
// Returns whether the write succeeded, the current value, and any error.
func (b *ConditionalPebbleBucket) PutIf(key, value []byte, condition Condition) (ok bool, currentValue []byte, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	fullKey := b.makeKey(key)

	current, err := b.db.Get(fullKey)
	if err != nil && err != pebbledb.ErrKeyNotFound {
		return false, nil, err
	}

	var currentCopy []byte
	if len(current) > 0 {
		currentCopy = make([]byte, len(current))
		copy(currentCopy, current)
	}

	if !condition(currentCopy) {
		return false, currentCopy, nil
	}

	if err := b.db.Set(fullKey, value); err != nil {
		return false, nil, err
	}

	return true, currentCopy, nil
}

// PutReturnOld atomically writes a value and returns the previous value.
// This is useful for implementing atomic increment or get-and-set operations.
func (b *ConditionalPebbleBucket) PutReturnOld(key, value []byte) (oldValue []byte, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	fullKey := b.makeKey(key)

	old, err := b.db.Get(fullKey)
	if err != nil && err != pebbledb.ErrKeyNotFound {
		return nil, err
	}

	var oldCopy []byte
	if len(old) > 0 {
		oldCopy = make([]byte, len(old))
		copy(oldCopy, old)
	}

	if err := b.db.Set(fullKey, value); err != nil {
		return nil, err
	}

	return oldCopy, nil
}

// DeleteIf atomically checks a condition before deleting a key.
// Returns whether the delete succeeded, the old value, and any error.
func (b *ConditionalPebbleBucket) DeleteIf(key []byte, condition Condition) (ok bool, oldValue []byte, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	fullKey := b.makeKey(key)

	old, err := b.db.Get(fullKey)
	if err != nil && err != pebbledb.ErrKeyNotFound {
		return false, nil, err
	}

	var oldCopy []byte
	if len(old) > 0 {
		oldCopy = make([]byte, len(old))
		copy(oldCopy, old)
	}

	if !condition(oldCopy) {
		return false, oldCopy, nil
	}

	if err := b.db.Delete(fullKey); err != nil {
		return false, nil, err
	}

	return true, oldCopy, nil
}

// ConditionalCheckFailedError is returned when a conditional write operation fails
// because the condition was not satisfied. It contains the key, current value, and
// the underlying error if any.
type ConditionalCheckFailedError struct {
	Key          []byte
	CurrentValue []byte
	Err          error
}

// Error returns the error message for the conditional check failure.
func (e *ConditionalCheckFailedError) Error() string {
	if e.Err != nil {
		return "conditional check failed: " + e.Err.Error()
	}
	return "conditional check failed"
}

// Unwrap returns the underlying error, if any.
func (e *ConditionalCheckFailedError) Unwrap() error { return e.Err }

// NewConditionalCheckFailedError creates a new error for failed conditional checks.
func NewConditionalCheckFailedError(key, currentValue []byte) *ConditionalCheckFailedError {
	return &ConditionalCheckFailedError{
		Key:          key,
		CurrentValue: currentValue,
	}
}
