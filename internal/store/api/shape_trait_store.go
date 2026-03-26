package api

import (
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/core/storage"
)

// ShapeTraitStore provides storage operations for shape trait annotations.
// Traits are metadata annotations attached to API shapes.
type ShapeTraitStore struct {
	bucket storage.Bucket
}

// NewShapeTraitStore creates a new ShapeTraitStore with the provided storage backend.
func NewShapeTraitStore(store storage.BasicStorage) *ShapeTraitStore {
	return &ShapeTraitStore{
		bucket: store.Bucket("api_shape_traits"),
	}
}

func (s *ShapeTraitStore) makeKey(shapeID int64, traitName string) []byte {
	return []byte(fmt.Sprintf("%d:%s", shapeID, traitName))
}

// Put stores a shape trait in the store.
func (s *ShapeTraitStore) Put(trait *ShapeTrait) error {
	data, err := json.Marshal(trait)
	if err != nil {
		return err
	}
	return s.bucket.Put(s.makeKey(trait.ShapeID, trait.TraitName), data)
}

// GetByShapeID retrieves all traits attached to a specific shape.
func (s *ShapeTraitStore) GetByShapeID(shapeID int64) ([]*ShapeTrait, error) {
	prefix := []byte(fmt.Sprintf("%d:", shapeID))
	var traits []*ShapeTrait
	iter := s.bucket.ScanPrefix(prefix)
	defer iter.Close()

	for iter.Next() {
		var t ShapeTrait
		if err := json.Unmarshal(iter.Value(), &t); err != nil {
			return nil, err
		}
		traits = append(traits, &t)
	}

	return traits, iter.Error()
}

// GetTraitValue retrieves the value of a specific trait on a shape.
func (s *ShapeTraitStore) GetTraitValue(shapeID int64, traitName string) (string, error) {
	data, err := s.bucket.Get(s.makeKey(shapeID, traitName))
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	var t ShapeTrait
	if err := json.Unmarshal(data, &t); err != nil {
		return "", err
	}
	return t.TraitValue, nil
}

// Delete removes a specific trait from a shape.
func (s *ShapeTraitStore) Delete(shapeID int64, traitName string) error {
	return s.bucket.Delete(s.makeKey(shapeID, traitName))
}

// DeleteByShapeID removes all traits attached to a specific shape.
func (s *ShapeTraitStore) DeleteByShapeID(shapeID int64) error {
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

// ForEach iterates over all shape traits and calls the provided function for each one.
func (s *ShapeTraitStore) ForEach(fn func(*ShapeTrait) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		var t ShapeTrait
		if err := json.Unmarshal(v, &t); err != nil {
			return err
		}
		return fn(&t)
	})
}

// Count returns the number of shape traits in the store.
func (s *ShapeTraitStore) Count() int {
	return s.bucket.Count()
}
