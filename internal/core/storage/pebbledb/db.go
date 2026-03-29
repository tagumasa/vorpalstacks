package pebbledb

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cockroachdb/pebble/v2"
)

var (
	// ErrClosed is returned when attempting to operate on a database that has been closed.
	ErrClosed = errors.New("database is closed")
	// ErrKeyNotFound is returned when a requested key does not exist in the database.
	ErrKeyNotFound = errors.New("key not found")
	// ErrInvalidTTLKey is returned when a key has an invalid TTL format.
	ErrInvalidTTLKey = errors.New("invalid TTL key format")
)

const (
	ttlMagicPrefix = 0x54
	ttlMetaSize    = 9
	ttlMetaVersion = 1
)

// TTLValue represents a value with a time-to-live expiration.
// ExpiresAt is the Unix timestamp (in seconds) after which the value should be considered expired.
// Data contains the actual stored data.
type TTLValue struct {
	ExpiresAt int64
	Data      []byte
}

// Marshal serialises the TTL value to binary format.
func (v *TTLValue) Marshal() ([]byte, error) {
	if v.ExpiresAt < 0 {
		return nil, errors.New("invalid TTL value: ExpiresAt cannot be negative")
	}
	buf := make([]byte, 1+8+len(v.Data))
	buf[0] = ttlMagicPrefix
	binary.BigEndian.PutUint64(buf[1:9], uint64(v.ExpiresAt))
	copy(buf[9:], v.Data)
	return buf, nil
}

// UnmarshalTTLValue deserialises a TTL-encoded value from binary data.
// Returns an error if the data is malformed or contains an invalid magic prefix.
func UnmarshalTTLValue(data []byte) (*TTLValue, error) {
	if len(data) < 9 {
		return nil, errors.New("invalid TTL value: too short")
	}
	if data[0] != ttlMagicPrefix {
		return nil, errors.New("not a TTL value")
	}
	u := binary.BigEndian.Uint64(data[1:9])
	if u > math.MaxInt64 {
		return nil, errors.New("invalid TTL value: overflow")
	}
	return &TTLValue{
		ExpiresAt: int64(u),
		Data:      data[9:],
	}, nil
}

// TTLKeyMeta contains metadata for a TTL-encoded key.
// It is used internally for tracking expiration information.
type TTLKeyMeta struct {
	ExpiresAt int64
	Value     []byte
}

// DB provides a high-level interface to a Pebble database.
// It supports encryption, TTL, and various data access patterns.
type DB struct {
	db             *pebble.DB
	encryptor      Encryptor
	opts           *Options
	closed         bool
	mu             sync.RWMutex
	ttlStop        chan struct{}
	pendingDeletes sync.WaitGroup
}

type pebbleLogger struct{}

// Infof logs an informational message from the Pebble database.
// It formats the message with the given arguments and logs it with a PEBBLE INFO prefix.
func (l *pebbleLogger) Infof(format string, args ...interface{}) {
	slog.Info(fmt.Sprintf(format, args...), "source", "pebble")
}

// Fatalf logs a fatal error message from the Pebble database.
// It formats the message with the given arguments and logs it with a PEBBLE FATAL prefix.
// This is typically called when the database encounters a critical error from which it cannot recover.
func (l *pebbleLogger) Fatalf(format string, args ...interface{}) {
	slog.Error(fmt.Sprintf(format, args...), "source", "pebble", "severity", "critical")
}

// Errorf logs an error message from the Pebble database.
// It formats the message with the given arguments and logs it with a PEBBLE ERROR prefix.
func (l *pebbleLogger) Errorf(format string, args ...interface{}) {
	slog.Error(fmt.Sprintf(format, args...), "source", "pebble")
}

// Open opens a new Pebble database with the given options.
// If no options are provided, default options are used.
// The database is created if it does not already exist.
func Open(opts ...Option) (*DB, error) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	if err := os.MkdirAll(options.Path, 0755); err != nil { // nosec G301
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	pebbleOpts := &pebble.Options{
		Logger: &pebbleLogger{},
	}
	db, err := pebble.Open(options.Path, pebbleOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to open pebble database: %w", err)
	}

	var encryptor Encryptor
	if options.Encryption.Enabled {
		encryptor, err = NewAESGCMEncryptor(options.Encryption.Key)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to create encryptor: %w", err)
		}
	} else {
		encryptor = NewNoopEncryptor()
	}

	pdb := &DB{
		db:        db,
		encryptor: encryptor,
		opts:      options,
		ttlStop:   make(chan struct{}),
	}

	if options.TTL.Enabled {
		go pdb.ttlGC()
	}

	return pdb, nil
}

// Close closes the database and releases all associated resources.
// It is safe to call Close multiple times; subsequent calls will return nil.
func (d *DB) Close() error {
	d.mu.Lock()

	if d.closed {
		d.mu.Unlock()
		return nil
	}

	d.closed = true

	if d.ttlStop != nil {
		close(d.ttlStop)
	}

	d.mu.Unlock()
	d.pendingDeletes.Wait()
	return d.db.Close()
}

// Get retrieves the value associated with the given key.
// Returns ErrKeyNotFound if the key does not exist.
// The returned slice is a copy and may be safely modified.
func (d *DB) Get(key []byte) ([]byte, error) {
	d.mu.RLock()
	if d.closed {
		d.mu.RUnlock()
		return nil, ErrClosed
	}

	val, closer, err := d.db.Get(key)
	if err != nil {
		d.mu.RUnlock()
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	defer closer.Close()

	result := make([]byte, len(val))
	copy(result, val)

	decrypted, err := d.encryptor.Decrypt(result)
	if err != nil {
		d.mu.RUnlock()
		return nil, fmt.Errorf("failed to decrypt value: %w", err)
	}

	if d.opts.TTL.Enabled {
		ttlVal, err := UnmarshalTTLValue(decrypted)
		if err == nil {
			if ttlVal.ExpiresAt > 0 && time.Now().Unix() > ttlVal.ExpiresAt {
				d.mu.RUnlock()
				d.mu.Lock()
				if !d.closed {
					d.pendingDeletes.Add(1)
					go func() {
						defer d.pendingDeletes.Done()
						_ = d.Delete(key)
					}()
				}
				d.mu.Unlock()
				return nil, ErrKeyNotFound
			}
			d.mu.RUnlock()
			return ttlVal.Data, nil
		}
	}

	d.mu.RUnlock()
	return decrypted, nil
}

// Set stores a key-value pair in the database with no expiration.
// This is equivalent to calling SetWithTTL with a zero duration.
func (d *DB) Set(key, value []byte) error {
	return d.SetWithTTL(key, value, 0)
}

// SetWithTTL stores a key-value pair with a time-to-live expiration.
// If ttl is zero, the value is stored without expiration.
// Returns an error if the database is closed.
func (d *DB) SetWithTTL(key, value []byte, ttl time.Duration) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return ErrClosed
	}

	var data []byte
	if ttl > 0 && d.opts.TTL.Enabled {
		expiresAt := time.Now().Add(ttl).Unix()
		ttlValue := &TTLValue{
			ExpiresAt: expiresAt,
			Data:      value,
		}
		var err error
		data, err = ttlValue.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal TTL value: %w", err)
		}
	} else {
		data = value
	}

	encrypted, err := d.encryptor.Encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	return d.db.Set(key, encrypted, d.syncFlag())
}

// Delete removes the specified key from the database.
// Returns nil if the key did not exist.
func (d *DB) Delete(key []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return ErrClosed
	}

	return d.db.Delete(key, d.syncFlag())
}

// DeleteRange removes all keys in the range [start, end).
// Both start and end are inclusive lower bounds and exclusive upper bounds.
func (d *DB) DeleteRange(start, end []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return ErrClosed
	}

	return d.db.DeleteRange(start, end, d.syncFlag())
}

// NewIter returns an iterator over the database.
// The caller must close the iterator when finished.
func (d *DB) NewIter(opts *pebble.IterOptions) (*pebble.Iterator, error) {
	d.mu.RLock()
	if d.closed {
		d.mu.RUnlock()
		return nil, ErrClosed
	}
	d.mu.RUnlock()
	return d.db.NewIter(opts)
}

// SetRaw stores a key-value pair without encryption.
// Use this for internal keys that should not be encrypted.
func (d *DB) SetRaw(key, value []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return ErrClosed
	}

	return d.db.Set(key, value, d.syncFlag())
}

// GetRaw retrieves a value without decryption.
// Returns the raw bytes stored in the database.
func (d *DB) GetRaw(key []byte) ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.closed {
		return nil, ErrClosed
	}

	val, closer, err := d.db.Get(key)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	defer closer.Close()

	result := make([]byte, len(val))
	copy(result, val)
	return result, nil
}

// Compact triggers a compaction of the entire database.
// This can improve read performance by eliminating obsolete files.
func (d *DB) Compact(ctx context.Context) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.closed {
		return ErrClosed
	}

	start := []byte{0x00}
	end := []byte{0xFF}
	return d.db.Compact(ctx, start, end, true)
}

// Metrics returns current database metrics including compactions, memtable sizes, and disk usage.
// Returns nil if the database is closed.
func (d *DB) Metrics() *pebble.Metrics {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.closed {
		return nil
	}

	return d.db.Metrics()
}

// Checkpoint creates a checkpoint of the database to the specified directory.
func (d *DB) Checkpoint(dir string) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.closed {
		return ErrClosed
	}

	return d.db.Checkpoint(dir)
}

func (d *DB) ttlGC() {
	ticker := time.NewTicker(d.opts.TTL.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-d.ttlStop:
			return
		case <-ticker.C:
			d.cleanExpiredKeys()
		}
	}
}

func (d *DB) cleanExpiredKeys() {
	d.mu.Lock()
	if d.closed {
		d.mu.Unlock()
		return
	}

	now := time.Now().Unix()

	keysToDelete := [][]byte{}

	iter, err := d.db.NewIter(nil)
	if err != nil {
		d.mu.Unlock()
		return
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()

		val, err := iter.ValueAndErr()
		if err != nil {
			continue
		}

		decrypted, err := d.encryptor.Decrypt(val)
		if err != nil {
			continue
		}

		ttlVal, err := UnmarshalTTLValue(decrypted)
		if err != nil {
			continue
		}

		if ttlVal.ExpiresAt > 0 && ttlVal.ExpiresAt <= now {
			keyCopy := make([]byte, len(key))
			copy(keyCopy, key)
			keysToDelete = append(keysToDelete, keyCopy)
		}
	}
	d.mu.Unlock()

	if len(keysToDelete) > 0 {
		batch := d.db.NewBatch()
		defer batch.Close()

		for _, key := range keysToDelete {
			_ = batch.Delete(key, nil) // intentional: batch errors are caught by Commit
		}

		if err := batch.Commit(pebble.Sync); err != nil {
			slog.Warn("failed to commit batch deletion of expired keys", "error", err, "count", len(keysToDelete))
		}
	}
}

// GetJSON retrieves a value and unmarshals it into the destination object.
// The destination must be a pointer.
func (d *DB) GetJSON(key []byte, dest any) error {
	data, err := d.Get(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// SetJSON marshals the value as JSON and stores it.
// This is a convenience method for storing structured data.
func (d *DB) SetJSON(key []byte, value any) error {
	return d.SetJSONWithTTL(key, value, 0)
}

// SetJSONWithTTL marshals the value as JSON and stores it with a time-to-live expiration.
func (d *DB) SetJSONWithTTL(key []byte, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return d.SetWithTTL(key, data, ttl)
}

// ListKeys returns all keys that begin with the given prefix.
// The keys are returned as a slice of byte slices.
func (d *DB) ListKeys(prefix []byte) ([][]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.closed {
		return nil, ErrClosed
	}

	var keys [][]byte

	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: append(prefix, 0xFF),
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)
		keys = append(keys, keyCopy)
	}

	return keys, nil
}

// Has checks if a key exists in the database.
// Returns true if the key exists, false otherwise.
func (d *DB) Has(key []byte) (bool, error) {
	_, err := d.Get(key)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ForEach iterates over all key-value pairs in the database.
// The iteration stops if the callback returns an error.
func (d *DB) ForEach(fn func(k, v []byte) error) error {
	iter, err := d.NewIter(nil)
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		val, err := iter.ValueAndErr()
		if err != nil {
			return err
		}

		decrypted, err := d.encryptor.Decrypt(val)
		if err != nil {
			return err
		}

		if d.opts.TTL.Enabled {
			ttlVal, err := UnmarshalTTLValue(decrypted)
			if err == nil && ttlVal.ExpiresAt > 0 {
				if time.Now().Unix() > ttlVal.ExpiresAt {
					continue
				}
				decrypted = ttlVal.Data
			}
		}

		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)

		if err := fn(keyCopy, decrypted); err != nil {
			return err
		}
	}

	return iter.Error()
}

// ScanPrefix iterates over all key-value pairs with the given prefix.
// The iteration stops if the callback returns an error.
func (d *DB) ScanPrefix(prefix []byte, fn func(k, v []byte) error) error {
	upperBound := make([]byte, len(prefix)+1)
	copy(upperBound, prefix)
	upperBound[len(prefix)] = 0xFF

	iter, err := d.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: upperBound,
	})
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		val, err := iter.ValueAndErr()
		if err != nil {
			return err
		}

		decrypted, err := d.encryptor.Decrypt(val)
		if err != nil {
			return err
		}

		if d.opts.TTL.Enabled {
			ttlVal, err := UnmarshalTTLValue(decrypted)
			if err == nil && ttlVal.ExpiresAt > 0 {
				if time.Now().Unix() > ttlVal.ExpiresAt {
					continue
				}
				decrypted = ttlVal.Data
			}
		}

		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)

		if err := fn(keyCopy, decrypted); err != nil {
			return err
		}
	}

	return iter.Error()
}

// ScanRange iterates over all key-value pairs in the range [start, end).
// The iteration stops if the callback returns an error.
func (d *DB) ScanRange(start, end []byte, fn func(k, v []byte) error) error {
	iter, err := d.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		val, err := iter.ValueAndErr()
		if err != nil {
			return err
		}

		decrypted, err := d.encryptor.Decrypt(val)
		if err != nil {
			return err
		}

		if d.opts.TTL.Enabled {
			ttlVal, err := UnmarshalTTLValue(decrypted)
			if err == nil && ttlVal.ExpiresAt > 0 {
				if time.Now().Unix() > ttlVal.ExpiresAt {
					continue
				}
				decrypted = ttlVal.Data
			}
		}

		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)

		if err := fn(keyCopy, decrypted); err != nil {
			return err
		}
	}

	return iter.Error()
}

// Count returns the number of keys in the database.
func (d *DB) Count() int {
	count := 0
	err := d.ForEach(func(k, v []byte) error {
		count++
		return nil
	})
	if err != nil {
		return 0
	}
	return count
}

// CountPrefix returns the number of keys with the given prefix.
func (d *DB) CountPrefix(prefix []byte) int {
	count := 0
	_ = d.ScanPrefix(prefix, func(k, v []byte) error {
		count++
		return nil
	})
	return count
}

// DB returns the underlying Pebble database instance.
// This allows direct access to Pebble-specific operations.
func (d *DB) DB() *pebble.DB {
	return d.db
}

// Encryptor returns the encryptor used by this database.
// This allows inspection of the encryption configuration.
func (d *DB) Encryptor() Encryptor {
	return d.encryptor
}

func (d *DB) syncFlag() *pebble.WriteOptions {
	if d.opts.SyncWrites {
		return pebble.Sync
	}
	return pebble.NoSync
}

// DataDir returns the path to the data directory for a given base path.
// This is a utility function for constructing database paths.
func DataDir(basePath string) string {
	return filepath.Join(basePath, "data")
}

// LazyIterator provides lazy, on-demand iteration over a key range.
// Unlike ScanRange which materialises all results upfront, this iterator
// wraps a pebble.Iterator and decrypts values one at a time.
type LazyIterator struct {
	iter      *pebble.Iterator
	encryptor Encryptor
	ttlOpts   *TTLOptions
	key       []byte
	value     []byte
	err       error
	closed    bool
	first     bool
}

// NewLazyIterator creates a lazy iterator over the given range.
// The caller must call Close when finished.
func (d *DB) NewLazyIterator(start, end []byte) *LazyIterator {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.closed {
		return &LazyIterator{err: ErrClosed}
	}

	iter, err := d.db.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	if err != nil {
		return &LazyIterator{err: err}
	}

	var ttlOpts *TTLOptions
	if d.opts.TTL.Enabled {
		ttlOpts = &d.opts.TTL
	}

	li := &LazyIterator{
		iter:      iter,
		encryptor: d.encryptor,
		ttlOpts:   ttlOpts,
	}
	li.first = true
	return li
}

func (li *LazyIterator) advance() {
	for li.iter.Valid() {
		key := li.iter.Key()
		val, err := li.iter.ValueAndErr()
		if err != nil {
			li.err = err
			return
		}

		decrypted, err := li.encryptor.Decrypt(val)
		if err != nil {
			li.err = err
			return
		}

		if li.ttlOpts != nil {
			ttlVal, err := UnmarshalTTLValue(decrypted)
			if err == nil && ttlVal.ExpiresAt > 0 {
				if time.Now().Unix() > ttlVal.ExpiresAt {
					li.iter.Next()
					continue
				}
				decrypted = ttlVal.Data
			}
		}

		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)
		li.key = keyCopy
		li.value = decrypted
		return
	}
	li.key = nil
	li.value = nil
}

// Next advances the iterator to the next key-value pair.
func (li *LazyIterator) Next() bool {
	if li.err != nil || li.closed {
		return false
	}
	if li.first {
		li.first = false
		li.iter.First()
	} else {
		li.iter.Next()
	}
	li.advance()
	return li.key != nil
}

// Key returns the current key.
func (li *LazyIterator) Key() []byte {
	return li.key
}

// Value returns the current decrypted value.
func (li *LazyIterator) Value() []byte {
	return li.value
}

// Error returns any error encountered during iteration.
func (li *LazyIterator) Error() error {
	return li.err
}

// Close releases the underlying iterator resources.
func (li *LazyIterator) Close() {
	if li.closed {
		return
	}
	li.closed = true
	if li.iter != nil {
		li.iter.Close()
	}
}
