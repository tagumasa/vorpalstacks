package config

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// Store provides configuration storage with support for defaults, environment variables, and persistent storage.
type Store struct {
	*common.BaseStore
	defaults map[string]ConfigEntry
}

// NewStore creates a new configuration store.
//
// Parameters:
//   - store: The underlying storage instance
//
// Returns:
//   - *Store: A new configuration store
func NewStore(store storage.BasicStorage) *Store {
	return &Store{
		BaseStore: common.NewBaseStore(store.Bucket("app_config"), "config"),
		defaults:  loadDefaults(),
	}
}

// Get retrieves a configuration entry by key.
// Lookup order: store, default value. ENV fallback is handled by MergeAndPersist
// during startup, not at runtime.
func (s *Store) Get(key string) (*ConfigEntry, error) {
	var entry ConfigEntry
	if err := s.BaseStore.Get(key, &entry); err == nil {
		entry.Source = ConfigSourceStore
		return &entry, nil
	}

	if def, ok := s.defaults[key]; ok {
		defCopy := def
		defCopy.Source = ConfigSourceDefault
		return &defCopy, nil
	}

	return nil, ErrConfigNotFound
}

// Set sets a configuration value in the store.
// For keys defined in defaults: validates ReadOnly, persists via BaseStore.
// For dynamic keys (not in defaults): creates a synthetic ConfigEntry with
// Type=ConfigTypePort, Category=CategoryPorts, ReadOnly=false.
func (s *Store) Set(key string, value interface{}) error {
	if def, ok := s.defaults[key]; ok {
		if def.ReadOnly {
			return ErrConfigReadOnly
		}
		entry := def
		entry.Value = value
		entry.Source = ConfigSourceStore
		entry.UpdatedAt = time.Now().Unix()
		return s.BaseStore.Put(key, entry)
	}

	entry := ConfigEntry{
		Key:       key,
		Value:     value,
		Type:      ConfigTypePort,
		Source:    ConfigSourceStore,
		Category:  CategoryPorts,
		ReadOnly:  false,
		UpdatedAt: time.Now().Unix(),
	}
	return s.BaseStore.Put(key, entry)
}

// Delete removes a configuration entry from the store.
func (s *Store) Delete(key string) error {
	return s.BaseStore.Delete(key)
}

// Reset resets a configuration entry to its default value.
func (s *Store) Reset(key string) (*ConfigEntry, error) {
	if err := s.BaseStore.Delete(key); err != nil {
		return nil, err
	}
	return s.Get(key)
}

// GetAll retrieves all configuration entries sorted by category then key.
func (s *Store) GetAll() ([]*ConfigEntry, error) {
	entries := make(map[string]*ConfigEntry)

	for key, def := range s.defaults {
		defCopy := def
		defCopy.Source = ConfigSourceDefault
		entries[key] = &defCopy
	}

	err := s.BaseStore.ForEach(func(key string, value []byte) error {
		var entry ConfigEntry
		if err := json.Unmarshal(value, &entry); err != nil {
			return err
		}
		entry.Source = ConfigSourceStore
		entries[key] = &entry
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := make([]*ConfigEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, entry)
	}
	slices.SortFunc(result, func(a, b *ConfigEntry) int {
		if a.Category != b.Category {
			return strings.Compare(string(a.Category), string(b.Category))
		}
		return strings.Compare(a.Key, b.Key)
	})
	return result, nil
}

// ListByCategory retrieves all configuration entries in a category.
func (s *Store) ListByCategory(category ConfigCategory) ([]*ConfigEntry, error) {
	all, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	var result []*ConfigEntry
	for _, entry := range all {
		if entry.Category == category {
			result = append(result, entry)
		}
	}
	return result, nil
}

// GetResourcePort retrieves the port for a specific resource.
func (s *Store) GetResourcePort(servicePortKey, resourceID string) (int, error) {
	resourceKey := servicePortKey + "." + resourceID

	if entry, err := s.Get(resourceKey); err == nil {
		if port, ok := entry.Value.(int); ok {
			return port, nil
		}
		if port, ok := entry.Value.(float64); ok {
			return int(port), nil
		}
	}

	entry, err := s.Get(servicePortKey)
	if err != nil {
		return 0, err
	}
	if port, ok := entry.Value.(int); ok {
		return port, nil
	}
	if port, ok := entry.Value.(float64); ok {
		return int(port), nil
	}
	return 0, ErrConfigInvalid
}

// SetResourcePort sets the port for a specific resource.
func (s *Store) SetResourcePort(servicePortKey, resourceID string, port int) error {
	resourceKey := servicePortKey + "." + resourceID
	return s.Set(resourceKey, port)
}

// DeleteResourcePort removes a resource port allocation.
func (s *Store) DeleteResourcePort(servicePortKey, resourceID string) error {
	resourceKey := servicePortKey + "." + resourceID
	return s.BaseStore.Delete(resourceKey)
}

// ListResourcePorts returns all resource port allocations for a service.
func (s *Store) ListResourcePorts(servicePortKey string) (map[string]int, error) {
	result := make(map[string]int)
	prefix := servicePortKey + "."
	err := s.BaseStore.ForEach(func(key string, value []byte) error {
		if !strings.HasPrefix(key, prefix) {
			return nil
		}
		var entry ConfigEntry
		if err := json.Unmarshal(value, &entry); err != nil {
			return nil
		}
		port := 0
		switch v := entry.Value.(type) {
		case float64:
			port = int(v)
		case int:
			port = v
		default:
			return nil
		}
		resourceID := strings.TrimPrefix(key, prefix)
		result[resourceID] = port
		return nil
	})
	return result, err
}

// GetString retrieves a string configuration value.
func (s *Store) GetString(key string) string {
	entry, err := s.Get(key)
	if err != nil {
		return ""
	}
	if sv, ok := entry.Value.(string); ok {
		return sv
	}
	return ""
}

// GetInt retrieves an integer configuration value.
func (s *Store) GetInt(key string) int {
	entry, err := s.Get(key)
	if err != nil {
		return 0
	}
	switch v := entry.Value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		var i int
		if _, err := parseInt(v, &i); err == nil {
			return i
		}
	}
	return 0
}

// GetBool retrieves a boolean configuration value.
func (s *Store) GetBool(key string) bool {
	entry, err := s.Get(key)
	if err != nil {
		return false
	}
	if b, ok := entry.Value.(bool); ok {
		return b
	}
	return false
}

// GetCategory returns the category for a configuration key.
func (s *Store) GetCategory(key string) ConfigCategory {
	if def, ok := s.defaults[key]; ok {
		return def.Category
	}
	parts := strings.Split(key, ".")
	if len(parts) > 0 {
		switch parts[0] {
		case "server":
			return CategoryServer
		case "aws":
			return CategoryAWS
		case "storage":
			return CategoryStorage
		case "features":
			return CategoryFeatures
		case "endpoints":
			return CategoryEndpoints
		case "ports":
			return CategoryPorts
		case "services":
			return CategoryServices
		case "http":
			return CategoryHTTP
		}
	}
	return ""
}

// ForEach iterates over all raw key-value pairs in the underlying store.
func (s *Store) ForEach(fn func(key string, value []byte) error) error {
	return s.BaseStore.ForEach(fn)
}

// InitDefaults seeds all default values into the store on first startup.
// Only writes entries that don't already exist in the store.
func (s *Store) InitDefaults() error {
	for key, def := range s.defaults {
		var existing ConfigEntry
		if err := s.BaseStore.Get(key, &existing); err != nil {
			entry := def
			entry.Source = ConfigSourceDefault
			entry.UpdatedAt = time.Now().Unix()
			if err := s.BaseStore.Put(key, entry); err != nil {
				return fmt.Errorf("init defaults: write %s: %w", key, err)
			}
		}
	}
	return nil
}

// MergeAndPersist writes a set of configuration values to the store.
// For each key that has an EnvVar in defaults and that env var is set,
// the Source is set to ConfigSourceEnv. Otherwise, ConfigSourceStore.
// This is called once during startup to persist ENV overrides.
func (s *Store) MergeAndPersist(values map[string]interface{}) {
	for key, value := range values {
		def, hasDef := s.defaults[key]

		entry := ConfigEntry{
			Key:       key,
			Value:     value,
			Source:    ConfigSourceStore,
			UpdatedAt: time.Now().Unix(),
		}

		if hasDef {
			entry.Type = def.Type
			entry.Category = def.Category
			entry.Description = def.Description
			entry.ReadOnly = def.ReadOnly
			entry.EnvVar = def.EnvVar

			if def.EnvVar != "" {
				if envVal := os.Getenv(def.EnvVar); envVal != "" {
					entry.Source = ConfigSourceEnv
				}
			}
		} else {
			entry.Type = ConfigTypePort
			entry.Category = CategoryPorts
		}

		if err := s.BaseStore.Put(key, entry); err != nil {
			logs.Warn("MergeAndPersist: failed to write config entry",
				logs.String("key", key),
				logs.Err(err),
			)
		}
	}
}
