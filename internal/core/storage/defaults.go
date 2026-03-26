// Package storage provides storage functionality for vorpalstacks.
package storage

import "time"

// Default constants for storage operations.
const (
	// DefaultLockMaxRetries is the default maximum number of lock retry attempts.
	DefaultLockMaxRetries = 100
	// DefaultVersionMaxRetries is the default maximum number of version retry attempts.
	DefaultVersionMaxRetries = 100
	// DefaultLockRetryDelay is the default delay between lock retry attempts.
	DefaultLockRetryDelay = 50 * time.Millisecond

	// BucketPrefix is the prefix used for bucket keys.
	BucketPrefix = "__bucket__:"
	// LockPrefix is the prefix used for lock keys.
	LockPrefix = "__locks__:"
	// IndexPrefix is the prefix used for index keys.
	IndexPrefix = "__indexes__:"
	// VersionedPrefix is the prefix used for versioned keys.
	VersionedPrefix = "__versioned__:"
	// ConditionalPrefix is the prefix used for conditional keys.
	ConditionalPrefix = "__cond__:"
	// IdempotencyPrefix is the prefix used for idempotency keys.
	IdempotencyPrefix = "__idempotency__:"
	// MetaPrefix is the prefix used for metadata keys.
	MetaPrefix = "__meta_"
	// CounterPrefix is the prefix used for counter keys.
	CounterPrefix = "__counter__"
)
