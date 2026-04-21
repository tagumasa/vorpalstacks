// Package config provides runtime configuration storage for vorpalstacks.
package config

import (
	"errors"
)

// ConfigType represents the type of a configuration value.
type ConfigType string

const (
	// ConfigTypeString represents a string configuration value.
	ConfigTypeString ConfigType = "STRING"
	// ConfigTypeInt represents an integer configuration value.
	ConfigTypeInt ConfigType = "INT"
	// ConfigTypeBool represents a boolean configuration value.
	ConfigTypeBool ConfigType = "BOOL"
	// ConfigTypePort represents a port configuration value.
	ConfigTypePort ConfigType = "PORT"
	// ConfigTypeURL represents a URL configuration value.
	ConfigTypeURL ConfigType = "URL"
)

// ConfigSource represents the source of a configuration value.
type ConfigSource string

const (
	// ConfigSourceEnv indicates the value comes from an environment variable.
	ConfigSourceEnv ConfigSource = "ENV"
	// ConfigSourceStore indicates the value comes from persistent storage.
	ConfigSourceStore ConfigSource = "STORE"
	// ConfigSourceDefault indicates the value is the default.
	ConfigSourceDefault ConfigSource = "DEFAULT"
)

// ConfigCategory represents a configuration category.
type ConfigCategory string

const (
	// CategoryServer contains server-related configuration keys.
	CategoryServer ConfigCategory = "server"
	// CategoryAWS contains AWS-related configuration keys.
	CategoryAWS ConfigCategory = "aws"
	// CategoryStorage contains storage-related configuration keys.
	CategoryStorage ConfigCategory = "storage"
	// CategoryFeatures contains feature flag configuration keys.
	CategoryFeatures ConfigCategory = "features"
	// CategoryEndpoints contains endpoint configuration keys.
	CategoryEndpoints ConfigCategory = "endpoints"
	// CategoryPorts contains port configuration keys.
	CategoryPorts ConfigCategory = "ports"
	// CategoryHTTP contains HTTP-related configuration keys (e.g. CORS).
	CategoryHTTP ConfigCategory = "http"
)

// ConfigEntry represents a single configuration entry.
type ConfigEntry struct {
	Key         string         `json:"key"`
	Value       interface{}    `json:"value"`
	Type        ConfigType     `json:"type"`
	Source      ConfigSource   `json:"source"`
	Description string         `json:"description"`
	ReadOnly    bool           `json:"read_only"`
	EnvVar      string         `json:"env_var,omitempty"`
	Category    ConfigCategory `json:"category"`
	UpdatedAt   int64          `json:"updated_at"`
}

// ConfigSchema represents the schema for a configuration entry.
type ConfigSchema struct {
	Key         string         `json:"key"`
	Type        ConfigType     `json:"type"`
	Default     interface{}    `json:"default"`
	Description string         `json:"description"`
	ReadOnly    bool           `json:"read_only"`
	EnvVar      string         `json:"env_var"`
	Category    ConfigCategory `json:"category"`
}

var (
	// ErrConfigNotFound is returned when a configuration key is not found.
	ErrConfigNotFound = errors.New("configuration not found")
	// ErrConfigReadOnly is returned when attempting to modify a read-only configuration.
	ErrConfigReadOnly = errors.New("configuration is read-only")
	// ErrConfigInvalid is returned when a configuration value is invalid.
	ErrConfigInvalid = errors.New("invalid configuration value")
	// ErrConfigNotInitialised is returned when the config store has not been initialised.
	ErrConfigNotInitialised = errors.New("config store not initialised")
)
