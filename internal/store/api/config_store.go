// Package api provides storage implementations for API metadata including services, operations, and shapes.
package api

import (
	"encoding/json"
	"sync"

	"vorpalstacks/internal/core/storage"
)

// ConfigStore stores service configuration.
type ConfigStore struct {
	bucket storage.Bucket
	cache  sync.Map
}

// NewConfigStore creates a new config store.
func NewConfigStore(store storage.BasicStorage) *ConfigStore {
	return &ConfigStore{
		bucket: store.Bucket("service_config"),
	}
}

// LoadAll loads all service configurations into the cache.
func (s *ConfigStore) LoadAll() error {
	err := s.bucket.ForEach(func(k, v []byte) error {
		var cfg ServiceConfig
		if err := json.Unmarshal(v, &cfg); err != nil {
			return err
		}
		s.cache.Store(cfg.ServiceName, &cfg)
		return nil
	})
	return err
}

// Get retrieves a service configuration by name.
func (s *ConfigStore) Get(serviceName string) (*ServiceConfig, error) {
	if v, ok := s.cache.Load(serviceName); ok {
		cfg := v.(*ServiceConfig)
		return &ServiceConfig{
			ServiceName: cfg.ServiceName,
			Mode:        cfg.Mode,
			Error:       cfg.Error,
			Enabled:     cfg.Enabled,
		}, nil
	}

	data, err := s.bucket.Get([]byte(serviceName))
	if err != nil {
		return nil, err
	}
	if data == nil {
		cfg := NewServiceConfig(serviceName, ServiceModeFallback)
		s.cache.Store(serviceName, cfg)
		return &ServiceConfig{
			ServiceName: cfg.ServiceName,
			Mode:        cfg.Mode,
			Error:       cfg.Error,
			Enabled:     cfg.Enabled,
		}, nil
	}
	var cfg ServiceConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	s.cache.Store(serviceName, &cfg)
	return &ServiceConfig{
		ServiceName: cfg.ServiceName,
		Mode:        cfg.Mode,
		Error:       cfg.Error,
		Enabled:     cfg.Enabled,
	}, nil
}

// Put saves a service configuration to the store.
func (s *ConfigStore) Put(cfg *ServiceConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := s.bucket.Put([]byte(cfg.ServiceName), data); err != nil {
		return err
	}
	s.cache.Store(cfg.ServiceName, cfg)
	return nil
}

// Delete removes a service configuration from the store.
func (s *ConfigStore) Delete(serviceName string) error {
	if err := s.bucket.Delete([]byte(serviceName)); err != nil {
		return err
	}
	s.cache.Delete(serviceName)
	return nil
}

// ForEach iterates over all service configurations.
func (s *ConfigStore) ForEach(fn func(*ServiceConfig) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		var cfg ServiceConfig
		if err := json.Unmarshal(v, &cfg); err != nil {
			return err
		}
		return fn(&cfg)
	})
}

// SetMode updates the service mode for a configuration.
func (s *ConfigStore) SetMode(serviceName string, mode ServiceMode) error {
	cfg, err := s.Get(serviceName)
	if err != nil {
		return err
	}
	cfg.Mode = mode
	return s.Put(cfg)
}

// SetError sets the error configuration for a service.
func (s *ConfigStore) SetError(serviceName string, errCfg *ErrorConfig) error {
	cfg, err := s.Get(serviceName)
	if err != nil {
		return err
	}
	cfg.Error = errCfg
	cfg.Mode = ServiceModeErrorInjection
	return s.Put(cfg)
}

// SetEnabled enables or disables a service.
func (s *ConfigStore) SetEnabled(serviceName string, enabled bool) error {
	cfg, err := s.Get(serviceName)
	if err != nil {
		return err
	}
	cfg.Enabled = enabled
	return s.Put(cfg)
}

// Count returns the number of service configurations.
func (s *ConfigStore) Count() int {
	return s.bucket.Count()
}
