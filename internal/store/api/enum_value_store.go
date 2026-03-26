package api

import (
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/core/storage"
)

// EnumValueStore stores API enum values.
type EnumValueStore struct {
	bucket storage.Bucket
}

// NewEnumValueStore creates a new enum value store.
func NewEnumValueStore(store storage.BasicStorage) *EnumValueStore {
	return &EnumValueStore{
		bucket: store.Bucket("api_enum_values"),
	}
}

func (s *EnumValueStore) makeKey(shapeID int64, name string) []byte {
	return []byte(fmt.Sprintf("%d:%s", shapeID, name))
}

// Put saves an enum value to the store.
func (s *EnumValueStore) Put(ev *ShapeEnumValue) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	return s.bucket.Put(s.makeKey(ev.ShapeID, ev.Name), data)
}

// GetByShapeID retrieves all enum values for a shape.
func (s *EnumValueStore) GetByShapeID(shapeID int64) ([]*ShapeEnumValue, error) {
	prefix := []byte(fmt.Sprintf("%d:", shapeID))
	var values []*ShapeEnumValue
	iter := s.bucket.ScanPrefix(prefix)
	defer iter.Close()

	for iter.Next() {
		var ev ShapeEnumValue
		if err := json.Unmarshal(iter.Value(), &ev); err != nil {
			return nil, err
		}
		values = append(values, &ev)
	}

	return values, iter.Error()
}

// Get retrieves an enum value by shape ID and name.
func (s *EnumValueStore) Get(shapeID int64, name string) (*ShapeEnumValue, error) {
	data, err := s.bucket.Get(s.makeKey(shapeID, name))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var ev ShapeEnumValue
	if err := json.Unmarshal(data, &ev); err != nil {
		return nil, err
	}
	return &ev, nil
}

// Delete removes an enum value from the store.
func (s *EnumValueStore) Delete(shapeID int64, name string) error {
	return s.bucket.Delete(s.makeKey(shapeID, name))
}

// DeleteByShapeID removes all enum values for a shape.
func (s *EnumValueStore) DeleteByShapeID(shapeID int64) error {
	prefix := []byte(fmt.Sprintf("%d:", shapeID))
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

// ForEach iterates over all enum values.
func (s *EnumValueStore) ForEach(fn func(*ShapeEnumValue) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		var ev ShapeEnumValue
		if err := json.Unmarshal(v, &ev); err != nil {
			return err
		}
		return fn(&ev)
	})
}

// Count returns the number of enum values in the store.
func (s *EnumValueStore) Count() int {
	return s.bucket.Count()
}
