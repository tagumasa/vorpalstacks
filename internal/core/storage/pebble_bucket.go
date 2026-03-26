// Package storage provides storage functionality for vorpalstacks.
package storage

import (
	"vorpalstacks/internal/core/storage/pebbledb"
)

// PebbleBucket provides a namespaced key-value bucket interface for Pebble.
// It allows logical grouping of keys with a common prefix.
type PebbleBucket struct {
	db     *pebbledb.DB
	prefix []byte
}

func (b *PebbleBucket) makeKey(key []byte) []byte {
	k := make([]byte, len(b.prefix)+len(key))
	copy(k, b.prefix)
	copy(k[len(b.prefix):], key)
	return k
}

// Get retrieves a value by key from the bucket.
func (b *PebbleBucket) Get(key []byte) ([]byte, error) {
	val, err := b.db.Get(b.makeKey(key))
	if err != nil {
		if err == pebbledb.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	result := make([]byte, len(val))
	copy(result, val)
	return result, nil
}

// Put stores a key-value pair in the bucket.
func (b *PebbleBucket) Put(key, value []byte) error {
	return b.db.Set(b.makeKey(key), value)
}

// Delete removes a key from the bucket.
func (b *PebbleBucket) Delete(key []byte) error {
	return b.db.Delete(b.makeKey(key))
}

// Has checks whether a key exists in the bucket.
func (b *PebbleBucket) Has(key []byte) bool {
	exists, err := b.db.Has(b.makeKey(key))
	if err != nil {
		return false
	}
	return exists
}

// ForEach iterates over all key-value pairs in the bucket.
func (b *PebbleBucket) ForEach(fn func(k, v []byte) error) error {
	start := b.makeKey(nil)
	end := make([]byte, len(b.prefix)+1)
	copy(end, b.prefix)
	end[len(b.prefix)] = 0xFF

	return b.db.ScanRange(start, end, func(k, v []byte) error {
		origKey := make([]byte, len(k)-len(b.prefix))
		copy(origKey, k[len(b.prefix):])
		return fn(origKey, v)
	})
}

// ScanPrefix returns an iterator for keys with the given prefix.
func (b *PebbleBucket) ScanPrefix(prefix []byte) Iterator {
	start := b.makeKey(prefix)
	end := make([]byte, len(b.prefix)+len(prefix)+1)
	copy(end, b.prefix)
	copy(end[len(b.prefix):], prefix)
	end[len(b.prefix)+len(prefix)] = 0xFF

	return &PebbleDBIterator{
		db:        b.db,
		start:     start,
		end:       end,
		prefixLen: len(b.prefix),
	}
}

// ScanRange returns an iterator for keys within the given range.
func (b *PebbleBucket) ScanRange(start, end []byte) Iterator {
	lower := b.makeKey(start)
	upper := b.makeKey(end)

	return &PebbleDBIterator{
		db:        b.db,
		start:     lower,
		end:       upper,
		prefixLen: len(b.prefix),
	}
}

// Count returns the number of keys in the bucket.
func (b *PebbleBucket) Count() int {
	start := b.makeKey(nil)
	end := make([]byte, len(b.prefix)+1)
	copy(end, b.prefix)
	end[len(b.prefix)] = 0xFF
	return b.db.CountPrefix(start)
}

// PebbleDBIterator provides an iterator implementation for pebbledb.
type PebbleDBIterator struct {
	db        *pebbledb.DB
	start     []byte
	end       []byte
	keys      [][]byte
	values    [][]byte
	current   int
	err       error
	prefixLen int
	closed    bool
}

// Next advances the iterator to the next key-value pair.
// Returns false if there are no more items or an error occurred.
func (i *PebbleDBIterator) Next() bool {
	if i.err != nil || i.closed {
		return false
	}
	if i.keys == nil {
		i.keys = make([][]byte, 0)
		i.values = make([][]byte, 0)
		i.current = -1

		err := i.db.ScanRange(i.start, i.end, func(k, v []byte) error {
			keyCopy := make([]byte, len(k))
			copy(keyCopy, k)
			valCopy := make([]byte, len(v))
			copy(valCopy, v)
			i.keys = append(i.keys, keyCopy)
			i.values = append(i.values, valCopy)
			return nil
		})
		if err != nil {
			i.err = err
			return false
		}
	}
	i.current++
	return i.current < len(i.keys)
}

// Key returns the key at the current iterator position.
func (i *PebbleDBIterator) Key() []byte {
	if i.current < 0 || i.current >= len(i.keys) {
		return nil
	}
	key := i.keys[i.current]
	if len(key) > i.prefixLen {
		origKey := make([]byte, len(key)-i.prefixLen)
		copy(origKey, key[i.prefixLen:])
		return origKey
	}
	return key
}

// Value returns the value at the current iterator position.
func (i *PebbleDBIterator) Value() []byte {
	if i.current < 0 || i.current >= len(i.values) {
		return nil
	}
	result := make([]byte, len(i.values[i.current]))
	copy(result, i.values[i.current])
	return result
}

// Error returns the error that occurred during iteration, if any.
func (i *PebbleDBIterator) Error() error {
	return i.err
}

// Close closes the iterator and releases resources.
func (i *PebbleDBIterator) Close() {
	i.closed = true
}
