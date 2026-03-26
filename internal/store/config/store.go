package config

import (
	"encoding/json"
	"os"
	"strings"
	"time"

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
//
// The lookup order is: store, environment variable, default value.
//
// Parameters:
//   - key: The configuration key
//
// Returns:
//   - *ConfigEntry: The configuration entry
//   - error: An error if the key is not found
func (s *Store) Get(key string) (*ConfigEntry, error) {
	var entry ConfigEntry
	if err := s.BaseStore.Get(key, &entry); err == nil {
		entry.Source = ConfigSourceStore
		return &entry, nil
	}

	if envVal := s.getFromEnv(key); envVal != nil {
		return envVal, nil
	}

	if def, ok := s.defaults[key]; ok {
		defCopy := def
		defCopy.Source = ConfigSourceDefault
		return &defCopy, nil
	}

	return nil, ErrConfigNotFound
}

func (s *Store) getFromEnv(key string) *ConfigEntry {
	def, ok := s.defaults[key]
	if !ok || def.EnvVar == "" {
		return nil
	}

	envVal := os.Getenv(def.EnvVar)
	if envVal == "" {
		return nil
	}

	entry := def
	entry.Value = s.parseEnvValue(envVal, def.Type)
	entry.Source = ConfigSourceEnv
	return &entry
}

func (s *Store) parseEnvValue(val string, typ ConfigType) interface{} {
	switch typ {
	case ConfigTypeInt, ConfigTypePort:
		var i int
		if _, err := parseInt(val, &i); err == nil {
			return i
		}
		return 0
	case ConfigTypeBool:
		return val == "true" || val == "1" || val == "yes"
	default:
		return val
	}
}

// Set sets a configuration value in the store.
//
// Parameters:
//   - key: The configuration key
//   - value: The value to set
//
// Returns:
//   - error: An error if setting fails or the key is read-only
func (s *Store) Set(key string, value interface{}) error {
	entry, err := s.Get(key)
	if err != nil {
		return err
	}
	if entry.ReadOnly {
		return ErrConfigReadOnly
	}

	entry.Value = value
	entry.Source = ConfigSourceStore
	entry.UpdatedAt = time.Now().Unix()

	return s.BaseStore.Put(key, entry)
}

// Delete removes a configuration entry from the store.
//
// Parameters:
//   - key: The configuration key to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *Store) Delete(key string) error {
	return s.BaseStore.Delete(key)
}

// Reset resets a configuration entry to its default value.
//
// Parameters:
//   - key: The configuration key to reset
//
// Returns:
//   - *ConfigEntry: The reset configuration entry
//   - error: An error if reset fails
func (s *Store) Reset(key string) (*ConfigEntry, error) {
	if err := s.BaseStore.Delete(key); err != nil {
		return nil, err
	}
	return s.Get(key)
}

// GetAll retrieves all configuration entries.
//
// Returns:
//   - []*ConfigEntry: All configuration entries
//   - error: An error if retrieval fails
func (s *Store) GetAll() ([]*ConfigEntry, error) {
	entries := make(map[string]*ConfigEntry)

	for key, def := range s.defaults {
		defCopy := def
		defCopy.Source = ConfigSourceDefault
		entries[key] = &defCopy

		if envVal := s.getFromEnv(key); envVal != nil {
			entries[key] = envVal
		}
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
	return result, nil
}

// ListByCategory retrieves all configuration entries in a category.
//
// Parameters:
//   - category: The category to filter by
//
// Returns:
//   - []*ConfigEntry: Configuration entries in the category
//   - error: An error if retrieval fails
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
//
// Parameters:
//   - servicePortKey: The service port configuration key
//   - resourceID: The resource identifier
//
// Returns:
//   - int: The resource port
//   - error: An error if retrieval fails
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
//
// Parameters:
//   - servicePortKey: The service port configuration key
//   - resourceID: The resource identifier
//   - port: The port to set
//
// Returns:
//   - error: An error if setting fails
func (s *Store) SetResourcePort(servicePortKey, resourceID string, port int) error {
	resourceKey := servicePortKey + "." + resourceID
	return s.Set(resourceKey, port)
}

// GetString retrieves a string configuration value.
//
// Parameters:
//   - key: The configuration key
//
// Returns:
//   - string: The string value, or empty string if not found
func (s *Store) GetString(key string) string {
	entry, err := s.Get(key)
	if err != nil {
		return ""
	}
	if s, ok := entry.Value.(string); ok {
		return s
	}
	return ""
}

// GetInt retrieves an integer configuration value.
//
// Parameters:
//   - key: The configuration key
//
// Returns:
//   - int: The integer value, or 0 if not found
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
//
// Parameters:
//   - key: The configuration key
//
// Returns:
//   - bool: The boolean value, or false if not found
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
//
// Parameters:
//   - key: The configuration key
//
// Returns:
//   - ConfigCategory: The category of the configuration key
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
		}
	}
	return ""
}
