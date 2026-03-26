// Package storage provides storage functionality for vorpalstacks.
package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// BackendType specifies the type of storage backend to use.
type BackendType string

// BackendType constants define supported storage backends.
const (
	BackendPebble BackendType = "pebble"
)

// Open opens a storage backend at the given path with the specified options.
func Open(path string, opts ...Option) (Storage, error) {
	cfg := &Config{
		Path:           path,
		CacheSizeBytes: 256 << 20,
		MaxOpenFiles:   1000,
		BytesPerSync:   1 << 20,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.Path), 0755); err != nil { // #nosec G301
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return NewPebbleStorage(cfg)
}

// OpenBackend opens a storage backend of the specified type at the given path.
func OpenBackend(backend BackendType, path string, opts ...Option) (Storage, error) {
	switch backend {
	case BackendPebble:
		return Open(path, opts...)
	default:
		return nil, fmt.Errorf("unknown storage backend: %s", backend)
	}
}
