package chunk

import (
	"fmt"
	"sync"
)

var defaultIndexManager = &IndexManager{
	indices: make(map[string]*PebbleIndex),
}

// GetDefaultIndexManager returns the default singleton IndexManager instance.
// This manager is used across the application for managing chunk indices.
func GetDefaultIndexManager() *IndexManager {
	return defaultIndexManager
}

// IndexManager manages multiple PebbleIndex instances.
// It provides a thread-safe way to create, retrieve, and close indices.
type IndexManager struct {
	indices map[string]*PebbleIndex
	mu      sync.Mutex
}

// GetIndex returns an existing PebbleIndex for the given path, or creates a new one if it does not exist.
// This method is thread-safe and manages index lifecycle automatically.
func (im *IndexManager) GetIndex(indexPath string) (*PebbleIndex, error) {
	im.mu.Lock()
	defer im.mu.Unlock()

	if idx, ok := im.indices[indexPath]; ok {
		return idx, nil
	}

	idx, err := NewPebbleIndex(&IndexOptions{
		Path: indexPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create index at %s: %w", indexPath, err)
	}

	im.indices[indexPath] = idx
	return idx, nil
}

// CloseAll closes all managed indices and releases their resources.
// After calling this method, the manager will have no open indices.
func (im *IndexManager) CloseAll() {
	im.mu.Lock()
	defer im.mu.Unlock()

	for _, idx := range im.indices {
		idx.Close()
	}
	im.indices = make(map[string]*PebbleIndex)
}

// Count returns the number of indices currently managed by the IndexManager.
func (im *IndexManager) Count() int {
	im.mu.Lock()
	defer im.mu.Unlock()
	return len(im.indices)
}
