package api

import (
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/core/storage"
)

// OperationErrorStore stores operation error definitions.
type OperationErrorStore struct {
	bucket storage.Bucket
}

// NewOperationErrorStore creates a new operation error store.
func NewOperationErrorStore(store storage.BasicStorage) *OperationErrorStore {
	return &OperationErrorStore{
		bucket: store.Bucket("api_operation_errors"),
	}
}

func (s *OperationErrorStore) makeKey(operationID, errorShapeID int64) []byte {
	return []byte(fmt.Sprintf("%d:%d", operationID, errorShapeID))
}

// Put saves an operation error to the store.
func (s *OperationErrorStore) Put(opError *OperationError) error {
	data, err := json.Marshal(opError)
	if err != nil {
		return err
	}
	return s.bucket.Put(s.makeKey(opError.OperationID, opError.ErrorShapeID), data)
}

// GetByOperationID retrieves all errors for an operation.
func (s *OperationErrorStore) GetByOperationID(operationID int64) ([]*OperationError, error) {
	prefix := []byte(fmt.Sprintf("%d:", operationID))
	var errors []*OperationError
	iter := s.bucket.ScanPrefix(prefix)
	defer iter.Close()

	for iter.Next() {
		var e OperationError
		if err := json.Unmarshal(iter.Value(), &e); err != nil {
			return nil, err
		}
		errors = append(errors, &e)
	}

	return errors, iter.Error()
}

// Get retrieves an operation error by operation ID and error shape ID.
func (s *OperationErrorStore) Get(operationID, errorShapeID int64) (*OperationError, error) {
	data, err := s.bucket.Get(s.makeKey(operationID, errorShapeID))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var e OperationError
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// Delete removes an operation error from the store.
func (s *OperationErrorStore) Delete(operationID, errorShapeID int64) error {
	return s.bucket.Delete(s.makeKey(operationID, errorShapeID))
}

// DeleteByOperationID removes all errors for an operation.
func (s *OperationErrorStore) DeleteByOperationID(operationID int64) error {
	prefix := []byte(fmt.Sprintf("%d:", operationID))
	var keys [][]byte
	iter := s.bucket.ScanPrefix(prefix)
	defer iter.Close()

	for iter.Next() {
		key := make([]byte, len(iter.Key()))
		copy(key, iter.Key())
		keys = append(keys, key)
	}

	if err := iter.Error(); err != nil {
		return err
	}

	for _, key := range keys {
		if err := s.bucket.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

// ForEach iterates over all operation errors.
func (s *OperationErrorStore) ForEach(fn func(*OperationError) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		var e OperationError
		if err := json.Unmarshal(v, &e); err != nil {
			return err
		}
		return fn(&e)
	})
}

// Count returns the number of operation errors in the store.
func (s *OperationErrorStore) Count() int {
	return s.bucket.Count()
}
