package pebbledb

import (
	"time"
)

// Options contains configuration options for the Pebble database.
type Options struct {
	Path string

	Encryption EncryptionOptions

	TTL TTLOptions

	SyncWrites bool
}

// EncryptionOptions contains configuration for database encryption.
type EncryptionOptions struct {
	Enabled bool
	Key     []byte
}

// TTLOptions contains configuration for time-to-live expiry of database entries.
type TTLOptions struct {
	Enabled       bool
	CheckInterval time.Duration
	DefaultTTL    time.Duration
}

// DefaultOptions returns the default options for the Pebble database.
func DefaultOptions() *Options {
	return &Options{
		Path: "./data",
		Encryption: EncryptionOptions{
			Enabled: false,
		},
		TTL: TTLOptions{
			Enabled:       true,
			CheckInterval: 5 * time.Minute,
			DefaultTTL:    24 * time.Hour,
		},
	}
}

// Option is a function that modifies Options.
type Option func(*Options)

// WithPath sets the path for the Pebble database.
func WithPath(path string) Option {
	return func(o *Options) {
		o.Path = path
	}
}

// WithEncryption enables encryption with the specified key.
func WithEncryption(key []byte) Option {
	return func(o *Options) {
		o.Encryption.Enabled = len(key) > 0
		o.Encryption.Key = key
	}
}

// WithTTL configures time-to-live settings for database entries.
func WithTTL(enabled bool, interval, defaultTTL time.Duration) Option {
	return func(o *Options) {
		o.TTL.Enabled = enabled
		o.TTL.CheckInterval = interval
		o.TTL.DefaultTTL = defaultTTL
	}
}
