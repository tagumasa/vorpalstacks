package api

// Package api provides API model data store implementations for vorpalstacks.

import (
	"encoding/json"

	"vorpalstacks/internal/core/storage"
)

// ShapeStore provides storage operations for API shape definitions.
// Shapes represent data structures in AWS API models.
type ShapeStore struct {
	bucket storage.Bucket
}

// NewShapeStore creates a new ShapeStore with the provided storage backend.
func NewShapeStore(store storage.BasicStorage) *ShapeStore {
	return &ShapeStore{
		bucket: store.Bucket("api_shapes"),
	}
}

// Get retrieves a shape by its shape ID.
func (s *ShapeStore) Get(shapeID string) (*Shape, error) {
	data, err := s.bucket.Get([]byte(shapeID))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var shape Shape
	if err := json.Unmarshal(data, &shape); err != nil {
		return nil, err
	}
	return &shape, nil
}

// GetByID retrieves a shape by its numeric ID.
func (s *ShapeStore) GetByID(id int64) (*Shape, error) {
	var found *Shape
	err := s.bucket.ForEach(func(k, v []byte) error {
		var shape Shape
		if err := json.Unmarshal(v, &shape); err != nil {
			return err
		}
		if shape.ID == id {
			found = &shape
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return found, nil
}

// Put stores a shape in the store.
func (s *ShapeStore) Put(shape *Shape) error {
	data, err := json.Marshal(shape)
	if err != nil {
		return err
	}
	return s.bucket.Put([]byte(shape.ShapeID), data)
}

// Delete removes a shape from the store by its shape ID.
func (s *ShapeStore) Delete(shapeID string) error {
	return s.bucket.Delete([]byte(shapeID))
}

// ForEach iterates over all shapes and calls the provided function for each one.
func (s *ShapeStore) ForEach(fn func(*Shape) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		var shape Shape
		if err := json.Unmarshal(v, &shape); err != nil {
			return err
		}
		return fn(&shape)
	})
}

// ForEachByService iterates over shapes belonging to a specific service and calls the provided function for each one.
func (s *ShapeStore) ForEachByService(serviceID int64, fn func(*Shape) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		var shape Shape
		if err := json.Unmarshal(v, &shape); err != nil {
			return err
		}
		if shape.ServiceID != nil && *shape.ServiceID == serviceID {
			return fn(&shape)
		}
		return nil
	})
}

// Count returns the number of shapes in the store.
func (s *ShapeStore) Count() int {
	return s.bucket.Count()
}
