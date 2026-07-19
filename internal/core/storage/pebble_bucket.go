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

// makeKey prepends the bucket prefix to the given key.
func (b *PebbleBucket) makeKey(key []byte) []byte {
	return makePrefixedKey(b.prefix, key)
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

	return newPrefixedDBIterator(b.db.NewLazyIterator(start, end), len(b.prefix))
}

// ScanRange returns an iterator for keys within the given range.
// If end is nil, the iterator scans from start to the end of the bucket.
func (b *PebbleBucket) ScanRange(start, end []byte) Iterator {
	lower := b.makeKey(start)
	var upper []byte
	if end != nil {
		upper = b.makeKey(end)
	} else {
		upper = append(append([]byte{}, b.prefix...), 0xFF)
	}

	return newPrefixedDBIterator(b.db.NewLazyIterator(lower, upper), len(b.prefix))
}

// Count returns the number of keys in the bucket.
func (b *PebbleBucket) Count() int {
	start := b.makeKey(nil)
	return b.db.CountPrefix(start)
}

// NewSnapshot creates a point-in-time snapshot of the bucket for
// consistent cross-key reads. The returned Snapshot observes the
// database state at the moment of creation; writes committed after
// creation are invisible. The caller must call Close when finished.
func (b *PebbleBucket) NewSnapshot() Snapshot {
	return &PebbleSnapshot{
		scanner: b.db.NewSnapshotScanner(),
		prefix:  b.prefix,
	}
}

// PebbleSnapshot wraps a pebbledb.SnapshotScanner with bucket prefix
// handling, implementing the storage.Snapshot interface.
type PebbleSnapshot struct {
	scanner *pebbledb.SnapshotScanner
	prefix  []byte
}

func (s *PebbleSnapshot) makeKey(key []byte) []byte {
	return makePrefixedKey(s.prefix, key)
}

// Get retrieves a value by key from the snapshot.
func (s *PebbleSnapshot) Get(key []byte) ([]byte, error) {
	val, err := s.scanner.Get(s.makeKey(key))
	if err != nil {
		if err == pebbledb.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	return val, nil
}

// Has checks whether a key exists in the snapshot.
func (s *PebbleSnapshot) Has(key []byte) bool {
	return s.scanner.Has(s.makeKey(key))
}

// ScanRange returns a prefix-stripped iterator over the given key range
// in the snapshot.
func (s *PebbleSnapshot) ScanRange(start, end []byte) Iterator {
	lower := s.makeKey(start)
	var upper []byte
	if end != nil {
		upper = s.makeKey(end)
	} else {
		upper = append(append([]byte{}, s.prefix...), 0xFF)
	}
	return newPrefixedDBIterator(s.scanner.NewLazyIterator(lower, upper), len(s.prefix))
}

// Close releases the underlying snapshot.
func (s *PebbleSnapshot) Close() {
	s.scanner.Close()
}

// prefixedDBIterator wraps a pebbledb.LazyIterator and strips the bucket prefix from keys.
// It implements the Iterator interface.
type prefixedDBIterator struct {
	lazy      *pebbledb.LazyIterator
	prefixLen int
}

// Next advances the iterator to the next key-value pair.
func (i *prefixedDBIterator) Next() bool {
	return i.lazy.Next()
}

// Key returns the current key with the bucket prefix stripped.
func (i *prefixedDBIterator) Key() []byte {
	key := i.lazy.Key()
	if key == nil {
		return nil
	}
	if len(key) > i.prefixLen {
		origKey := make([]byte, len(key)-i.prefixLen)
		copy(origKey, key[i.prefixLen:])
		return origKey
	}
	return key
}

// Value returns the current value.
func (i *prefixedDBIterator) Value() []byte {
	return i.lazy.Value()
}

// Error returns any accumulated error from the iterator.
func (i *prefixedDBIterator) Error() error {
	return i.lazy.Error()
}

// Close releases resources held by the iterator.
func (i *prefixedDBIterator) Close() {
	i.lazy.Close()
}

func newPrefixedDBIterator(lazy *pebbledb.LazyIterator, prefixLen int) *prefixedDBIterator {
	return &prefixedDBIterator{
		lazy:      lazy,
		prefixLen: prefixLen,
	}
}
