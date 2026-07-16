// Package storage provides core storage abstractions and implementations for vorpalstacks.
package storage

import (
	"fmt"
	"path/filepath"
	"sync"

	"vorpalstacks/internal/common/defaults"
)

// RegionStorageManager manages storage instances for different AWS regions.
type RegionStorageManager struct {
	baseDir  string
	config   *Config
	storages map[string]Storage
	mu       sync.RWMutex
}

// NewRegionStorageManager creates a new region storage manager.
//
// Parameters:
//   - cfg: The storage configuration
//
// Returns:
//   - *RegionStorageManager: A new region storage manager
//   - error: An error if creation fails
func NewRegionStorageManager(cfg *Config) (*RegionStorageManager, error) {
	if cfg == nil {
		cfg = &Config{Path: "./data"}
	}
	return &RegionStorageManager{
		baseDir:  cfg.Path,
		config:   cfg,
		storages: make(map[string]Storage),
	}, nil
}

// GetStorage retrieves or creates a storage instance for the specified region.
//
// Parameters:
//   - region: The AWS region
//
// Returns:
//   - Storage: The storage instance for the region
//   - error: An error if retrieval or creation fails
func (m *RegionStorageManager) GetStorage(region string) (Storage, error) {
	if region == "" {
		region = defaults.DefaultRegion
	}

	m.mu.RLock()
	if s, ok := m.storages[region]; ok {
		m.mu.RUnlock()
		return s, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if s, ok := m.storages[region]; ok {
		return s, nil
	}

	regionPath := filepath.Join(m.baseDir, region)
	regionCfg := &Config{
		Path:           regionPath,
		CacheSizeBytes: m.config.CacheSizeBytes,
		MaxOpenFiles:   m.config.MaxOpenFiles,
		BytesPerSync:   m.config.BytesPerSync,
	}

	s, err := NewPebbleStorage(regionCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage for region %s: %w", region, err)
	}

	m.storages[region] = s
	return s, nil
}

// GetGlobalStorage retrieves or creates the global storage instance.
func (m *RegionStorageManager) GetGlobalStorage() (Storage, error) {
	return m.GetStorage("global")
}

// Close closes all storage instances managed by this manager.
func (m *RegionStorageManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for region, s := range m.storages {
		if err := s.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close storage for region %s: %w", region, err)
		}
	}
	m.storages = make(map[string]Storage)
	return lastErr
}

// ListRegions returns a list of all regions that have storage instances.
func (m *RegionStorageManager) ListRegions() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	regions := make([]string, 0, len(m.storages))
	for r := range m.storages {
		regions = append(regions, r)
	}
	return regions
}

// GetActiveRegions returns a list of all active (non-global) regions.
func (m *RegionStorageManager) GetActiveRegions() []string {
	m.mu.RLock()
	regions := make([]string, 0, len(m.storages))
	for r := range m.storages {
		if r != "global" {
			regions = append(regions, r)
		}
	}
	m.mu.RUnlock()

	if len(regions) == 0 {
		return []string{defaults.DefaultRegion}
	}
	return regions
}
