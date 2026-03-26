// Package api provides storage implementations for API metadata including services, operations, and shapes.
package api

import (
	"encoding/json"

	"vorpalstacks/internal/core/storage"
)

// ServiceStore provides storage operations for AWS service definitions.
// Services are stored by their name and can be retrieved by name or shape ID.
type ServiceStore struct {
	bucket storage.Bucket
}

// NewServiceStore creates a new ServiceStore with the provided storage backend.
func NewServiceStore(store storage.BasicStorage) *ServiceStore {
	return &ServiceStore{
		bucket: store.Bucket("api_services"),
	}
}

// Get retrieves a service by its name.
func (s *ServiceStore) Get(name string) (*Service, error) {
	data, err := s.bucket.Get([]byte(name))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var svc Service
	if err := json.Unmarshal(data, &svc); err != nil {
		return nil, err
	}
	return &svc, nil
}

// GetByShapeID retrieves a service by its shape ID.
// This is useful when looking up a service associated with a specific shape.
func (s *ServiceStore) GetByShapeID(shapeID string) (*Service, error) {
	var found *Service
	err := s.bucket.ForEach(func(k, v []byte) error {
		var svc Service
		if err := json.Unmarshal(v, &svc); err != nil {
			return err
		}
		if svc.ShapeID == shapeID {
			found = &svc
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return found, nil
}

// Put stores a service in the store.
func (s *ServiceStore) Put(svc *Service) error {
	data, err := json.Marshal(svc)
	if err != nil {
		return err
	}
	return s.bucket.Put([]byte(svc.Name), data)
}

// Delete removes a service from the store by its name.
func (s *ServiceStore) Delete(name string) error {
	return s.bucket.Delete([]byte(name))
}

// ForEach iterates over all services and calls the provided function for each one.
func (s *ServiceStore) ForEach(fn func(*Service) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		var svc Service
		if err := json.Unmarshal(v, &svc); err != nil {
			return err
		}
		return fn(&svc)
	})
}

// Count returns the number of services in the store.
func (s *ServiceStore) Count() int {
	return s.bucket.Count()
}
