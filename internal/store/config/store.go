package config

import (
	"encoding/json"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/serviceports"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// Store is the Pebble-backed configuration store, seeded with defaults on first init.
type Store struct {
	*common.BaseStore
	defaults map[string]ConfigEntry
	initOnce sync.Once
}

// NewStore creates a new configuration store backed by the given storage.
func NewStore(store storage.BasicStorage) *Store {
	return &Store{
		BaseStore: common.NewBaseStore(store.Bucket("app_config"), "config"),
		defaults:  loadDefaults(),
	}
}

// Initialise seeds default values for keys not yet in Pebble, then applies
// environment variable overrides. Must be called once after NewStore.
func (s *Store) Initialise() {
	s.initOnce.Do(func() {
		s.seedDefaults()
		s.applyEnvOverrides()
	})
}

// seedDefaults writes every default key that does not already exist in Pebble.
func (s *Store) seedDefaults() {
	for key, def := range s.defaults {
		if s.BaseStore.Exists(key) {
			continue
		}
		entry := def
		entry.UpdatedAt = time.Now().Unix()
		if err := s.BaseStore.Put(key, entry); err != nil {
			logs.Warn("seedDefaults: failed to write default entry",
				logs.String("key", key),
				logs.Err(err),
			)
		}
	}
}

// applyEnvOverrides overwrites Pebble entries for keys whose EnvVar is set
// in the current process environment. After this call the server never reads
// ENV vars again — Pebble is the single source of truth.
func (s *Store) applyEnvOverrides() {
	for key, def := range s.defaults {
		if def.EnvVar == "" {
			continue
		}
		envVal := os.Getenv(def.EnvVar)
		if envVal == "" {
			continue
		}

		var value interface{}
		switch def.Type {
		case ConfigTypeBool:
			parsed, err := strconv.ParseBool(envVal)
			if err != nil {
				logs.Warn("applyEnvOverrides: invalid bool value, skipping",
					logs.String("key", key),
					logs.String("env", def.EnvVar),
					logs.String("value", envVal),
				)
				continue
			}
			value = parsed
		case ConfigTypeInt, ConfigTypePort:
			parsed, err := strconv.Atoi(envVal)
			if err != nil {
				logs.Warn("applyEnvOverrides: invalid int value, skipping",
					logs.String("key", key),
					logs.String("env", def.EnvVar),
					logs.String("value", envVal),
				)
				continue
			}
			value = parsed
		default:
			value = envVal
		}

		entry := def
		entry.Value = value
		entry.UpdatedAt = time.Now().Unix()
		if err := s.BaseStore.Put(key, entry); err != nil {
			logs.Warn("applyEnvOverrides: failed to write config entry",
				logs.String("key", key),
				logs.Err(err),
			)
		}
	}
}

// Get retrieves a configuration entry from Pebble by key.
func (s *Store) Get(key string) (*ConfigEntry, error) {
	var entry ConfigEntry
	if err := s.BaseStore.Get(key, &entry); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}
	return &entry, nil
}

// Set persists a configuration value to Pebble.
// For keys defined in defaults: validates ReadOnly and copies schema metadata.
// For dynamic keys (resource ports): creates a synthetic ConfigEntry.
func (s *Store) Set(key string, value interface{}) error {
	if def, ok := s.defaults[key]; ok {
		if def.ReadOnly {
			return ErrConfigReadOnly
		}
		entry := def
		entry.Value = value
		entry.UpdatedAt = time.Now().Unix()
		return s.BaseStore.Put(key, entry)
	}

	entry := ConfigEntry{
		Key:       key,
		Value:     value,
		Type:      ConfigTypePort,
		Category:  CategoryPorts,
		ReadOnly:  false,
		UpdatedAt: time.Now().Unix(),
	}
	return s.BaseStore.Put(key, entry)
}

// Delete removes a configuration entry by key.
func (s *Store) Delete(key string) error {
	return s.BaseStore.Delete(key)
}

// Reset deletes the Pebble entry and re-seeds the default value.
func (s *Store) Reset(key string) (*ConfigEntry, error) {
	def, ok := s.defaults[key]
	if !ok {
		return nil, ErrConfigNotFound
	}
	if err := s.BaseStore.Delete(key); err != nil {
		return nil, err
	}
	entry := def
	entry.UpdatedAt = time.Now().Unix()
	if err := s.BaseStore.Put(key, entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// GetAll retrieves all configuration entries from Pebble, sorted by category then key.
func (s *Store) GetAll() ([]*ConfigEntry, error) {
	var result []*ConfigEntry
	err := s.BaseStore.ForEach(func(key string, value []byte) error {
		var entry ConfigEntry
		if err := json.Unmarshal(value, &entry); err != nil {
			return err
		}
		result = append(result, &entry)
		return nil
	})
	if err != nil {
		return nil, err
	}
	slices.SortFunc(result, func(a, b *ConfigEntry) int {
		if a.Category != b.Category {
			return strings.Compare(string(a.Category), string(b.Category))
		}
		return strings.Compare(a.Key, b.Key)
	})
	return result, nil
}

// ListByCategory returns all configuration entries in the given category.
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

// GetResourcePort returns the allocated port for a resource, or 0 if unassigned.
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

// SetResourcePort persists the port allocation for a resource.
func (s *Store) SetResourcePort(servicePortKey, resourceID string, port int) error {
	resourceKey := servicePortKey + "." + resourceID
	return s.Set(resourceKey, port)
}

// DeleteResourcePort removes a port allocation for a resource.
func (s *Store) DeleteResourcePort(servicePortKey, resourceID string) error {
	resourceKey := servicePortKey + "." + resourceID
	return s.BaseStore.Delete(resourceKey)
}

// ListResourcePorts returns all resource-to-port mappings for a service port key.
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

// ListAllResourcePorts returns all port allocations across all service keys.
// Only includes entries whose value is an integer within the dynamic port range.
func (s *Store) ListAllResourcePorts() (map[string]map[string]int, error) {
	result := make(map[string]map[string]int)
	err := s.BaseStore.ForEach(func(key string, value []byte) error {
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
		if port < serviceports.DynamicRangeStart || port > serviceports.DynamicRangeEnd {
			return nil
		}
		dot := strings.Index(key, ".")
		if dot < 0 {
			return nil
		}
		svcKey := key[:dot]
		resID := key[dot+1:]
		if result[svcKey] == nil {
			result[svcKey] = make(map[string]int)
		}
		result[svcKey][resID] = port
		return nil
	})
	return result, err
}

// GetString retrieves a configuration value as a string.
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

// GetInt retrieves a configuration value as an integer.
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

// GetBool retrieves a configuration value as a boolean.
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

// GetCategory returns the category of a configuration key.
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

// ForEach iterates over all raw configuration key-value pairs.
func (s *Store) ForEach(fn func(key string, value []byte) error) error {
	return s.BaseStore.ForEach(fn)
}

// GetSchema returns the configuration schema derived from the defaults map.
func (s *Store) GetSchema() []ConfigSchema {
	schema := make([]ConfigSchema, 0, len(s.defaults))
	for _, entry := range s.defaults {
		schema = append(schema, ConfigSchema{
			Key:         entry.Key,
			Type:        entry.Type,
			Default:     entry.Value,
			Description: entry.Description,
			ReadOnly:    entry.ReadOnly,
			EnvVar:      entry.EnvVar,
			Category:    entry.Category,
		})
	}
	return schema
}

// GetKeysByCategory returns all configuration keys in a category.
func (s *Store) GetKeysByCategory(category ConfigCategory) []string {
	keys := make([]string, 0)
	for key, entry := range s.defaults {
		if entry.Category == category {
			keys = append(keys, key)
		}
	}
	return keys
}
