// Package storage provides storage functionality for vorpalstacks.
package storage

import (
	"context"
	"io"
	"time"
)

// Closer is an interface for objects that can be closed.
type Closer interface {
	Close() error
}

// WriteOptions contains options for write operations.
type WriteOptions struct {
	Sync bool
}

// DefaultWriteOptions contains the default options for write operations.
var DefaultWriteOptions = &WriteOptions{Sync: true}

// IterOptions contains options for iteration.
type IterOptions struct {
	LowerBound []byte
	UpperBound []byte
}

// BasicStorage provides basic storage operations.
type BasicStorage interface {
	Close() error
	Bucket(name string) Bucket
	CreateBucket(name string) error
	DeleteBucket(name string) error
	ListBuckets() []string
}

// TransactionalStorage provides transaction support.
type TransactionalStorage interface {
	BasicStorage
	View(ctx context.Context, fn func(Transaction) error) error
	Update(ctx context.Context, fn func(Transaction) error) error
}

// TransactionalStorageWith2PC provides transaction support with two-phase commit.
type TransactionalStorageWith2PC interface {
	TransactionalStorage
	TwoPhaseTransaction() TwoPhaseTransaction
}

// MaintenanceStorage provides maintenance operations.
type MaintenanceStorage interface {
	TransactionalStorage
	Backup(w io.Writer) error
	Compact() error
	Stats() Stats
}

// FeatureStorage provides advanced feature accessors.
type FeatureStorage interface {
	IdempotencyStore() IdempotencyStore
	ConditionalBucket(name string) ConditionalBucket
	VersionedBucket(name string) VersionedBucket
	LockManager() LockManager
	TwoPhaseTransaction() TwoPhaseTransaction
	MultiItemTransaction() MultiItemTransaction
}

// Storage defines the full interface for a storage system.
// Deprecated: Use BasicStorage, TransactionalStorage, or TransactionalStorageWith2PC instead.
// This interface is kept only for test utilities and backward compatibility.
type Storage interface {
	MaintenanceStorage
	FeatureStorage
}

// Bucket provides key-value storage operations.
type Bucket interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error
	Has(key []byte) bool
	ForEach(fn func(k, v []byte) error) error
	ScanPrefix(prefix []byte) Iterator
	ScanRange(start, end []byte) Iterator
	Count() int
}

// Transaction provides a transactional view of storage.
// Commit and rollback are handled automatically by View()/Update().
type Transaction interface {
	Bucket(name string) Bucket
	CreateBucket(name string) error
	DeleteBucket(name string) error
}

// ConditionalBucket provides conditional write operations.
type ConditionalBucket interface {
	PutIf(key, value []byte, condition Condition) (ok bool, currentValue []byte, err error)
	PutReturnOld(key, value []byte) (oldValue []byte, err error)
	DeleteIf(key []byte, condition Condition) (ok bool, oldValue []byte, err error)
}

// Condition is a function that checks a value and returns true if the condition is met.
type Condition func(currentValue []byte) bool

// ConditionNotExists returns a condition that checks whether a key does not exist.
func ConditionNotExists() Condition {
	return func(currentValue []byte) bool {
		return len(currentValue) == 0
	}
}

// ConditionExists returns a condition that checks whether a key exists.
func ConditionExists() Condition {
	return func(currentValue []byte) bool {
		return len(currentValue) > 0
	}
}

// ConditionEquals returns a condition that checks whether the current value
// equals the expected value.
func ConditionEquals(expected []byte) Condition {
	return func(currentValue []byte) bool {
		if len(currentValue) != len(expected) {
			return false
		}
		for i, b := range currentValue {
			if b != expected[i] {
				return false
			}
		}
		return true
	}
}

// ConditionNotEquals returns a condition that checks whether the current value
// does not equal the expected value.
func ConditionNotEquals(expected []byte) Condition {
	return func(currentValue []byte) bool {
		return !ConditionEquals(expected)(currentValue)
	}
}

// IdempotencyStore provides idempotency token storage.
type IdempotencyStore interface {
	CheckAndStore(token string, ttl time.Duration) (isNew bool, err error)
	StoreResult(token string, result []byte, ttl time.Duration) error
	GetResult(token string) ([]byte, bool, error)
	Delete(token string) error
}

// WriteResult contains the result of a write operation.
type WriteResult struct {
	Value    []byte
	Version  uint64
	Created  bool
	Updated  bool
	Deleted  bool
	OldValue []byte
}

// TwoPhaseTransaction provides two-phase commit support.
type TwoPhaseTransaction interface {
	AddValidator(v Validator)
	AddExecutor(e Executor)
	Commit(ctx context.Context) error
	Clear()
	ValidatorCount() int
	ExecutorCount() int
}

// Validator validates a transaction before commit.
type Validator interface {
	Validate(ctx context.Context, txn Transaction) error
}

// Executor executes operations within a transaction.
type Executor interface {
	Execute(ctx context.Context, txn Transaction) error
}

// ValidatorFunc is a function that implements Validator.
type ValidatorFunc func(ctx context.Context, txn Transaction) error

// Validate implements the Validator interface.
func (f ValidatorFunc) Validate(ctx context.Context, txn Transaction) error {
	return f(ctx, txn)
}

// ExecutorFunc is a function that implements Executor.
type ExecutorFunc func(ctx context.Context, txn Transaction) error

// Execute implements the Executor interface.
func (f ExecutorFunc) Execute(ctx context.Context, txn Transaction) error {
	return f(ctx, txn)
}

// VersionedValue represents a value with version information.
type VersionedValue struct {
	Value     []byte
	Version   uint64
	Timestamp int64
	Deleted   bool
}

// VersionListOptions contains options for listing versions.
type VersionListOptions struct {
	StartVersion   uint64
	Limit          int
	IncludeDeleted bool
	Reverse        bool
}

// VersionedBucket provides versioned key-value storage.
type VersionedBucket interface {
	GetVersion(key []byte, version uint64) (*VersionedValue, error)
	GetLatest(key []byte) (*VersionedValue, error)
	PutWithVersion(key, value []byte) (*VersionedValue, error)
	DeleteWithVersion(key []byte) (*VersionedValue, error)
	ListVersions(key []byte, opts VersionListOptions) ([]*VersionedValue, error)
	PurgeVersions(key []byte, keepVersions int) error
}

// LockMode represents the type of lock.
type LockMode int

// LockMode constants define the types of locks available.
const (
	LockModeShared LockMode = iota
	LockModeExclusive
)

// LockHandle represents a lock handle.
type LockHandle struct {
	Key       []byte
	Token     string
	Mode      LockMode
	ExpiresAt int64
}

// LockManager provides distributed locking support.
type LockManager interface {
	TryLock(key []byte, mode LockMode, ttl time.Duration) (*LockHandle, error)
	Lock(ctx context.Context, key []byte, mode LockMode, ttl time.Duration) (*LockHandle, error)
	Unlock(handle *LockHandle) error
	Extend(handle *LockHandle, ttl time.Duration) error
	IsLocked(key []byte) (bool, *LockHandle, error)
}

// MultiItemTransaction provides multi-item transaction support.
type MultiItemTransaction interface {
	Put(bucket string, key, value []byte)
	Delete(bucket string, key []byte)
	ConditionCheck(bucket string, key []byte, condition Condition)
	Commit(ctx context.Context) error
	Clear()
	OperationCount() int
	ConditionCount() int
}

// Iterator provides iteration over key-value pairs.
type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
	Error() error
	Close()
}

// Stats contains storage statistics.
type Stats struct {
	KeyCount    int64
	BucketCount int
	SizeBytes   int64
	Compactions int64
	WriteAmp    float64
	ReadAmp     float64
}

// KeyValueStore provides basic key-value storage operations.
type KeyValueStore interface {
	Get(ctx context.Context, key []byte) ([]byte, error)
	Put(ctx context.Context, key, value []byte) error
	Delete(ctx context.Context, key []byte) error
	Has(ctx context.Context, key []byte) (bool, error)
	ScanPrefix(ctx context.Context, prefix []byte) (map[string][]byte, error)
}

// Config contains storage configuration options.
type Config struct {
	Path           string
	CacheSizeBytes int64
	MaxOpenFiles   int
	BytesPerSync   int

	// Encryption settings
	EncryptionKey []byte // 16, 24, or 32 bytes for AES-128, AES-192, or AES-256

	// TTL settings
	TTLEnabled       bool
	TTLCheckInterval time.Duration
}

// Option configures a storage Config.
type Option func(*Config)

// WithCacheSize sets the cache size in bytes for the storage configuration.
func WithCacheSize(size int64) Option {
	return func(c *Config) {
		c.CacheSizeBytes = size
	}
}

// WithMaxOpenFiles sets the maximum number of open files for the storage configuration.
func WithMaxOpenFiles(n int) Option {
	return func(c *Config) {
		c.MaxOpenFiles = n
	}
}

// WithBytesPerSync sets the number of bytes per sync operation for the storage configuration.
func WithBytesPerSync(n int) Option {
	return func(c *Config) {
		c.BytesPerSync = n
	}
}

// WithEncryptionKey enables encryption with the specified key.
// Key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
func WithEncryptionKey(key []byte) Option {
	return func(c *Config) {
		c.EncryptionKey = key
	}
}

// WithTTL enables TTL with the specified check interval.
func WithTTL(enabled bool, checkInterval time.Duration) Option {
	return func(c *Config) {
		c.TTLEnabled = enabled
		c.TTLCheckInterval = checkInterval
	}
}
