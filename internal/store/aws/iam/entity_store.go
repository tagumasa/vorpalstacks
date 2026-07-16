package iam

import (
	"encoding/json"
	"strings"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

type entityStore[T any] struct {
	*common.BaseStore
	kl common.KeyLocker
}

func newEntityStore[T any](store storage.BasicStorage, bucketName string) entityStore[T] {
	return entityStore[T]{
		BaseStore: common.NewBaseStore(store.Bucket(bucketName), "iam"),
	}
}

// Get retrieves an entity by its storage key.
func (s *entityStore[T]) Get(key string) (*T, error) {
	var entity T
	if err := s.BaseStore.Get(key, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

// Delete removes an entity by its storage key.
func (s *entityStore[T]) Delete(key string) error {
	return s.BaseStore.Delete(key)
}

// Exists checks whether an entity with the given storage key exists.
func (s *entityStore[T]) Exists(key string) bool {
	return s.BaseStore.Exists(key)
}

// Count returns the total number of entities in the store.
func (s *entityStore[T]) Count() int {
	return s.BaseStore.Count()
}

func getEntityByID[T any](store *common.BaseStore, id string, getID func(*T) string, op string, notFound error) (*T, error) {
	var found *T
	err := store.ForEach(func(key string, value []byte) error {
		var entity T
		if err := json.Unmarshal(value, &entity); err != nil {
			return err
		}
		if getID(&entity) == id && found == nil {
			found = &entity
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError(op, err)
	}
	if found == nil {
		return nil, NewStoreError(op, notFound)
	}
	return found, nil
}

func listEntitiesWithPathPrefix[T any](store *common.BaseStore, pathPrefix, marker string, maxItems int, getPath func(*T) string) ([]*T, bool, string, error) {
	var filter common.FilterFunc[T]
	if pathPrefix != "" {
		filter = func(e *T) bool { return strings.HasPrefix(getPath(e), pathPrefix) }
	}
	result, err := common.List[T](store, common.ListOptions{Marker: marker, MaxItems: maxItems}, filter)
	if err != nil {
		return nil, false, "", err
	}
	return result.Items, result.IsTruncated, result.NextMarker, nil
}

func listEntitiesByPrefix[T any](store *common.BaseStore, prefix string, getField func(*T) string) ([]*T, error) {
	items, err := common.ListAll[T](store)
	if err != nil {
		return nil, err
	}
	if prefix == "" {
		return items, nil
	}
	var filtered []*T
	for _, item := range items {
		if strings.HasPrefix(getField(item), prefix) {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}
