package storage

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cockroachdb/pebble/v2"
	"vorpalstacks/internal/core/storage/pebbledb"
)

// PebbleStorage provides a storage implementation backed by PebbleDB.
// It implements the Storage interface and supports key-value operations,
// transactions, TTL, encryption, and various bucket types.
type PebbleStorage struct {
	db               *pebbledb.DB
	mu               sync.RWMutex
	closed           bool
	idempotencyStore *PebbleIdempotencyStore
	lockManager      *PebbleLockManager
}

// NewPebbleStorage creates a new PebbleStorage instance with the given configuration.
// If cfg is nil, default configuration is used. The storage path defaults to "./data".
func NewPebbleStorage(cfg *Config) (*PebbleStorage, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	if cfg.Path == "" {
		cfg.Path = "./data"
	}

	opts := []pebbledb.Option{pebbledb.WithPath(cfg.Path)}

	if len(cfg.EncryptionKey) > 0 {
		opts = append(opts, pebbledb.WithEncryption(cfg.EncryptionKey))
	}

	ttlEnabled := cfg.TTLEnabled
	ttlInterval := cfg.TTLCheckInterval
	if ttlInterval == 0 {
		ttlInterval = 5 * time.Minute
	}
	opts = append(opts, pebbledb.WithTTL(ttlEnabled, ttlInterval, 24*time.Hour))

	db, err := pebbledb.Open(opts...)
	if err != nil {
		return nil, err
	}

	storage := &PebbleStorage{
		db: db,
	}
	storage.idempotencyStore = NewPebbleIdempotencyStore(storage.db)
	storage.lockManager = NewPebbleLockManager(storage.db, nil)
	return storage, nil
}

// Close closes the storage and releases all resources.
func (s *PebbleStorage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	return s.db.Close()
}

// Bucket returns a bucket with the given name for key-value operations.
func (s *PebbleStorage) Bucket(name string) Bucket {
	return &PebbleBucket{
		db:     s.db,
		prefix: []byte(BucketPrefix + name + ":"),
	}
}

// CreateBucket creates a new bucket with the given name.
func (s *PebbleStorage) CreateBucket(name string) error {
	return nil
}

// DeleteBucket deletes a bucket and all its contents.
func (s *PebbleStorage) DeleteBucket(name string) error {
	prefix := []byte(BucketPrefix + name + ":")
	return s.db.DeleteRange(prefix, append(prefix, 0xFF))
}

// ListBuckets returns a list of all bucket names.
func (s *PebbleStorage) ListBuckets() []string {
	prefixSet := make(map[string]struct{})
	prefix := []byte(BucketPrefix)
	iter, err := s.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
	})
	if err != nil {
		return nil
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		if len(key) <= len(prefix) {
			continue
		}
		remainder := key[len(prefix):]
		for i, b := range remainder {
			if b == ':' {
				prefixSet[string(remainder[:i])] = struct{}{}
				break
			}
		}
	}
	if err := iter.Error(); err != nil {
		return nil
	}

	buckets := make([]string, 0, len(prefixSet))
	for name := range prefixSet {
		buckets = append(buckets, name)
	}
	return buckets
}

// View executes a read-only transaction.
// The provided function receives a Transaction that can be used for read operations.
func (s *PebbleStorage) View(ctx context.Context, fn func(Transaction) error) error {
	return s.db.View(ctx, func(txn *pebbledb.Txn) error {
		return fn(&PebbleDBTransaction{txn: txn})
	})
}

// Update executes a read-write transaction.
// The provided function receives a Transaction that can be used for read and write operations.
func (s *PebbleStorage) Update(ctx context.Context, fn func(Transaction) error) error {
	return s.db.Update(ctx, func(txn *pebbledb.Txn) error {
		return fn(&PebbleDBTransaction{txn: txn})
	})
}

// PebbleDBTransaction adapts pebbledb.Txn to the Transaction interface.
type PebbleDBTransaction struct {
	txn *pebbledb.Txn
}

// Bucket returns a bucket with the given name within the transaction.
func (t *PebbleDBTransaction) Bucket(name string) Bucket {
	prefix := []byte(BucketPrefix + name + ":")
	return &PebbleDBTransactionBucket{
		txn:     t.txn,
		prefix:  prefix,
		prefixL: len(prefix),
	}
}

// CreateBucket creates a new bucket within the transaction.
func (t *PebbleDBTransaction) CreateBucket(name string) error {
	return nil
}

// DeleteBucket deletes a bucket within the transaction.
func (t *PebbleDBTransaction) DeleteBucket(name string) error {
	prefix := []byte(BucketPrefix + name + ":")
	return t.txn.DeleteRange(prefix, append(prefix, 0xFF))
}

// Commit commits the transaction.
func (t *PebbleDBTransaction) Commit() error {
	return nil
}

// Rollback aborts the transaction.
func (t *PebbleDBTransaction) Rollback() error {
	return nil
}

// PebbleDBTransactionBucket provides a bucket interface for pebbledb.Txn.
type PebbleDBTransactionBucket struct {
	txn     *pebbledb.Txn
	prefix  []byte
	prefixL int
}

func (b *PebbleDBTransactionBucket) makeKey(key []byte) []byte {
	k := make([]byte, len(b.prefix)+len(key))
	copy(k, b.prefix)
	copy(k[len(b.prefix):], key)
	return k
}

// Get retrieves a value by key from the transaction bucket.
func (b *PebbleDBTransactionBucket) Get(key []byte) ([]byte, error) {
	val, err := b.txn.Get(b.makeKey(key))
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

// Put stores a key-value pair in the transaction bucket.
func (b *PebbleDBTransactionBucket) Put(key, value []byte) error {
	return b.txn.Set(b.makeKey(key), value)
}

// Delete removes a key from the transaction bucket.
func (b *PebbleDBTransactionBucket) Delete(key []byte) error {
	return b.txn.Delete(b.makeKey(key))
}

// Has checks whether a key exists in the transaction bucket.
func (b *PebbleDBTransactionBucket) Has(key []byte) bool {
	_, err := b.txn.Get(b.makeKey(key))
	return err == nil
}

// ForEach iterates over all key-value pairs in the transaction bucket.
func (b *PebbleDBTransactionBucket) ForEach(fn func(k, v []byte) error) error {
	iter, err := b.txn.NewIter(&pebble.IterOptions{
		LowerBound: b.prefix,
		UpperBound: append(append([]byte{}, b.prefix...), 0xFF),
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
		if len(key) > b.prefixL {
			origKey := make([]byte, len(key)-b.prefixL)
			copy(origKey, key[b.prefixL:])
			if err := fn(origKey, val); err != nil {
				return err
			}
		}
	}
	return iter.Error()
}

// ScanPrefix returns an iterator for keys with the given prefix in the transaction bucket.
func (b *PebbleDBTransactionBucket) ScanPrefix(prefix []byte) Iterator {
	start := make([]byte, len(b.prefix)+len(prefix))
	copy(start, b.prefix)
	copy(start[len(b.prefix):], prefix)

	end := make([]byte, len(b.prefix)+len(prefix)+1)
	copy(end, b.prefix)
	copy(end[len(b.prefix):], prefix)
	end[len(b.prefix)+len(prefix)] = 0xFF

	return newTxnPebbleDBIterator(b.txn, start, end, b.prefixL)
}

// ScanRange returns an iterator for keys within the given range in the transaction bucket.
func (b *PebbleDBTransactionBucket) ScanRange(start, end []byte) Iterator {
	lower := b.makeKey(start)
	upper := b.makeKey(end)

	return newTxnPebbleDBIterator(b.txn, lower, upper, b.prefixL)
}

// Count returns the number of keys in the transaction bucket.
func (b *PebbleDBTransactionBucket) Count() int {
	count := 0
	_ = b.ForEach(func(k, v []byte) error {
		count++
		return nil
	})
	return count
}

// Backup writes a backup of the database to the given writer.
// Note: Use BackupToDir for production backups instead.
func (s *PebbleStorage) Backup(w io.Writer) error {
	return fmt.Errorf("Backup to writer not implemented, use BackupToDir instead")
}

// Compact triggers a manual compaction of the database.
func (s *PebbleStorage) Compact() error {
	return s.db.Compact(context.Background())
}

// Stats returns statistics about the storage.
func (s *PebbleStorage) Stats() Stats {
	metrics := s.db.Metrics()
	return Stats{
		KeyCount:    int64(metrics.Keys.RangeKeySetsCount + metrics.Keys.TombstoneCount),
		SizeBytes:   int64(metrics.DiskSpaceUsage()),
		Compactions: int64(metrics.Compact.Count),
		ReadAmp:     float64(metrics.ReadAmp()),
	}
}

// DB returns the underlying CockroachDB Pebble instance.
func (s *PebbleStorage) DB() *pebble.DB {
	return s.db.DB()
}

// PebbleDB returns the underlying pebbledb instance.
func (s *PebbleStorage) PebbleDB() *pebbledb.DB {
	return s.db
}

// IdempotencyStore returns the idempotency store for this storage.
func (s *PebbleStorage) IdempotencyStore() IdempotencyStore {
	return s.idempotencyStore
}

// TwoPhaseTransaction returns a two-phase transaction manager.
func (s *PebbleStorage) TwoPhaseTransaction() TwoPhaseTransaction {
	return NewTwoPhaseTransaction(s)
}

// MultiItemTransaction returns a multi-item transaction manager.
func (s *PebbleStorage) MultiItemTransaction() MultiItemTransaction {
	return NewMultiItemTransaction(s)
}

// ConditionalBucket returns a bucket with conditional write operations.
func (s *PebbleStorage) ConditionalBucket(name string) ConditionalBucket {
	return NewConditionalPebbleBucket(s.db, []byte(ConditionalPrefix+name+":"))
}

// VersionedBucket returns a bucket with versioning support.
func (s *PebbleStorage) VersionedBucket(name string) VersionedBucket {
	return NewPebbleVersionedBucket(s.db, []byte(VersionedPrefix+name+":"))
}

// LockManager returns the distributed lock manager for this storage.
func (s *PebbleStorage) LockManager() LockManager {
	return s.lockManager
}

// CleanupExpired removes expired entries from the storage.
func (s *PebbleStorage) CleanupExpired() error {
	return s.idempotencyStore.Cleanup()
}

// BackupToDir creates a checkpoint backup of the database to the specified directory.
func (s *PebbleStorage) BackupToDir(dir string) error {
	return s.db.Checkpoint(dir)
}

// SetWithTTL sets a key-value pair with a time-to-live expiry.
func (s *PebbleStorage) SetWithTTL(key, value []byte, ttl time.Duration) error {
	return s.db.SetWithTTL(key, value, ttl)
}

// GetTTL retrieves a value with TTL checking.
func (s *PebbleStorage) GetTTL(key []byte) ([]byte, error) {
	return s.db.Get(key)
}

// SetJSONWithTTL encodes a value as JSON and stores it with a time-to-live expiry.
func (s *PebbleStorage) SetJSONWithTTL(key []byte, value any, ttl time.Duration) error {
	return s.db.SetJSONWithTTL(key, value, ttl)
}

// GetJSON retrieves a JSON-encoded value and decodes it into dest.
func (s *PebbleStorage) GetJSON(key []byte, dest any) error {
	return s.db.GetJSON(key, dest)
}
