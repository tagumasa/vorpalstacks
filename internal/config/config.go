// Package config provides global configuration access for vorpalstacks.
package config

import (
	"sync"

	"vorpalstacks/internal/core/storage"
	storeconfig "vorpalstacks/internal/store/config"
)

var (
	globalStore *storeconfig.Store
	mu          sync.RWMutex
)

// Initialise initializes the global configuration store.
//
// Parameters:
//   - store: The storage instance to use for configuration
func Initialise(store storage.BasicStorage) {
	mu.Lock()
	defer mu.Unlock()
	globalStore = storeconfig.NewStore(store)
}

// GetStore returns the global configuration store.
//
// Returns:
//   - *storeconfig.Store: The configuration store instance
func GetStore() *storeconfig.Store {
	mu.RLock()
	defer mu.RUnlock()
	return globalStore
}

// Get retrieves a configuration entry by key.
//
// Parameters:
//   - key: The configuration key
//
// Returns:
//   - *storeconfig.ConfigEntry: The configuration entry
//   - error: An error if retrieval fails
func Get(key string) (*storeconfig.ConfigEntry, error) {
	mu.RLock()
	defer mu.RUnlock()
	if globalStore == nil {
		return nil, storeconfig.ErrConfigNotInitialised
	}
	return globalStore.Get(key)
}

// Set sets a configuration value.
//
// Parameters:
//   - key: The configuration key
//   - value: The value to set
//
// Returns:
//   - error: An error if setting fails
func Set(key string, value interface{}) error {
	mu.Lock()
	defer mu.Unlock()
	if globalStore == nil {
		return storeconfig.ErrConfigNotInitialised
	}
	return globalStore.Set(key, value)
}

// GetString retrieves a string configuration value.
//
// Parameters:
//   - key: The configuration key
//
// Returns:
//   - string: The string value, or empty string if not initialised
func GetString(key string) string {
	mu.RLock()
	defer mu.RUnlock()
	if globalStore == nil {
		return ""
	}
	return globalStore.GetString(key)
}

// GetInt retrieves an integer configuration value.
//
// Parameters:
//   - key: The configuration key
//
// Returns:
//   - int: The integer value, or 0 if not initialised
func GetInt(key string) int {
	mu.RLock()
	defer mu.RUnlock()
	if globalStore == nil {
		return 0
	}
	return globalStore.GetInt(key)
}

// GetBool retrieves a boolean configuration value.
//
// Parameters:
//   - key: The configuration key
//
// Returns:
//   - bool: The boolean value, or false if not initialised
func GetBool(key string) bool {
	mu.RLock()
	defer mu.RUnlock()
	if globalStore == nil {
		return false
	}
	return globalStore.GetBool(key)
}

// GetAll retrieves all configuration entries.
//
// Returns:
//   - []*storeconfig.ConfigEntry: All configuration entries
//   - error: An error if retrieval fails
func GetAll() ([]*storeconfig.ConfigEntry, error) {
	mu.RLock()
	defer mu.RUnlock()
	if globalStore == nil {
		return nil, storeconfig.ErrConfigNotInitialised
	}
	return globalStore.GetAll()
}

// ListByCategory retrieves all configuration entries in a category.
//
// Parameters:
//   - category: The configuration category to filter by
//
// Returns:
//   - []*storeconfig.ConfigEntry: Configuration entries in the category
//   - error: An error if retrieval fails
func ListByCategory(category storeconfig.ConfigCategory) ([]*storeconfig.ConfigEntry, error) {
	mu.RLock()
	defer mu.RUnlock()
	if globalStore == nil {
		return nil, storeconfig.ErrConfigNotInitialised
	}
	return globalStore.ListByCategory(category)
}

// Reset resets a configuration entry to its default value.
//
// Parameters:
//   - key: The configuration key to reset
//
// Returns:
//   - *storeconfig.ConfigEntry: The reset configuration entry
//   - error: An error if reset fails
func Reset(key string) (*storeconfig.ConfigEntry, error) {
	mu.Lock()
	defer mu.Unlock()
	if globalStore == nil {
		return nil, storeconfig.ErrConfigNotInitialised
	}
	return globalStore.Reset(key)
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
func GetResourcePort(servicePortKey, resourceID string) (int, error) {
	mu.RLock()
	defer mu.RUnlock()
	if globalStore == nil {
		return 0, storeconfig.ErrConfigNotInitialised
	}
	return globalStore.GetResourcePort(servicePortKey, resourceID)
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
func SetResourcePort(servicePortKey, resourceID string, port int) error {
	mu.Lock()
	defer mu.Unlock()
	if globalStore == nil {
		return storeconfig.ErrConfigNotInitialised
	}
	return globalStore.SetResourcePort(servicePortKey, resourceID, port)
}

// ServerPort returns the configured HTTP server port.
//
// Returns:
//   - int: The server port
func ServerPort() int {
	return GetInt("server.port")
}

// GRPCWebPort returns the configured gRPC-web server port.
//
// Returns:
//   - int: The gRPC-web server port
func GRPCWebPort() int {
	return GetInt("server.grpc_web_port")
}

// BindAddr returns the configured server bind address.
//
// Returns:
//   - string: The bind address
func BindAddr() string {
	return GetString("server.bind_addr")
}

// AWSAccountID returns the configured AWS account ID.
//
// Returns:
//   - string: The AWS account ID
func AWSAccountID() string {
	return GetString("aws.account_id")
}

// AWSRegion returns the configured AWS region.
//
// Returns:
//   - string: The AWS region
func AWSRegion() string {
	return GetString("aws.region")
}

// DataPath returns the configured data storage path.
//
// Returns:
//   - string: The data path
func DataPath() string {
	return GetString("storage.data_path")
}

// BaseURL returns the configured base URL for the service.
//
// Returns:
//   - string: The base URL
func BaseURL() string {
	return GetString("endpoints.base_url")
}

// AuditEnabled returns whether audit logging is enabled.
//
// Returns:
//   - bool: True if audit logging is enabled
func AuditEnabled() bool {
	return GetBool("features.audit_enabled")
}

// SignatureVerificationEnabled returns whether AWS signature verification is enabled.
//
// Returns:
//   - bool: True if signature verification is enabled
func SignatureVerificationEnabled() bool {
	return GetBool("features.signature_verification")
}

// Route53DNSEnabled returns whether Route53 DNS management is enabled.
//
// Returns:
//   - bool: True if Route53 DNS is enabled
func Route53DNSEnabled() bool {
	return GetBool("features.route53_dns")
}
