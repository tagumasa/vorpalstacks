package api

import (
	"encoding/json"
	"errors"

	"vorpalstacks/internal/core/storage"
)

// OperationStore stores API operation definitions.
type OperationStore struct {
	bucket storage.Bucket
}

// NewOperationStore creates a new operation store.
func NewOperationStore(store storage.BasicStorage) *OperationStore {
	return &OperationStore{
		bucket: store.Bucket("api_operations"),
	}
}

func (s *OperationStore) makeKey(serviceName, opName string) []byte {
	return []byte(serviceName + ":" + opName)
}

// Get retrieves an operation by service name and operation name.
func (s *OperationStore) Get(serviceName, opName string) (*Operation, error) {
	data, err := s.bucket.Get(s.makeKey(serviceName, opName))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var op Operation
	if err := json.Unmarshal(data, &op); err != nil {
		return nil, err
	}
	return &op, nil
}

// GetByShapeID retrieves an operation by its shape ID.
func (s *OperationStore) GetByShapeID(shapeID string) (*Operation, error) {
	var found *Operation
	err := s.bucket.ForEach(func(k, v []byte) error {
		var op Operation
		if err := json.Unmarshal(v, &op); err != nil {
			return err
		}
		if op.ShapeID == shapeID {
			found = &op
			return errors.New("found")
		}
		return nil
	})
	if err != nil && err.Error() != "found" {
		return nil, err
	}
	return found, nil
}

// Put saves an operation to the store.
func (s *OperationStore) Put(serviceName string, op *Operation) error {
	data, err := json.Marshal(op)
	if err != nil {
		return err
	}
	return s.bucket.Put(s.makeKey(serviceName, op.Name), data)
}

// Delete removes an operation from the store.
func (s *OperationStore) Delete(serviceName, opName string) error {
	return s.bucket.Delete(s.makeKey(serviceName, opName))
}

// ForEach iterates over all operations.
func (s *OperationStore) ForEach(fn func(*Operation) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		var op Operation
		if err := json.Unmarshal(v, &op); err != nil {
			return err
		}
		return fn(&op)
	})
}

// ForEachByService iterates over operations for a specific service.
func (s *OperationStore) ForEachByService(serviceName string, fn func(*Operation) error) error {
	prefix := []byte(serviceName + ":")
	iter := s.bucket.ScanPrefix(prefix)
	defer iter.Close()

	for iter.Next() {
		var op Operation
		if err := json.Unmarshal(iter.Value(), &op); err != nil {
			return err
		}
		if err := fn(&op); err != nil {
			return err
		}
	}
	return iter.Error()
}

// Count returns the number of operations in the store.
func (s *OperationStore) Count() int {
	return s.bucket.Count()
}
