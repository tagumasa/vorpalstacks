// Package common provides common AWS store utilities for vorpalstacks.
package common

import (
	"encoding/json"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
)

// DefaultMaxItems is the default maximum number of items for paginated list operations.
const DefaultMaxItems = 100

// AbsoluteMaxItems is the hard upper limit for paginated list operations.
const AbsoluteMaxItems = 1000

// BaseStore provides common storage operations for AWS services.
type BaseStore struct {
	bucket  storage.Bucket
	service string
}

// NewBaseStore creates a new base store with the given bucket and service name.
func NewBaseStore(bucket storage.Bucket, service string) *BaseStore {
	return &BaseStore{bucket: bucket, service: service}
}

// Bucket returns the underlying storage bucket.
func (s *BaseStore) Bucket() storage.Bucket {
	return s.bucket
}

// Service returns the service name for this store.
func (s *BaseStore) Service() string {
	return s.service
}

// Get retrieves an item from the store by key.
func (s *BaseStore) Get(key string, dest interface{}) error {
	data, err := s.bucket.Get([]byte(key))
	if err != nil {
		return NewStoreErrorWithKey(s.service, "get", key, err)
	}
	if data == nil {
		return NewStoreErrorWithKey(s.service, "get", key, ErrNotFound)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return NewStoreErrorWithKey(s.service, "get", key, err)
	}
	return nil
}

// Put stores an item in the store with the given key.
func (s *BaseStore) Put(key string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return NewStoreErrorWithKey(s.service, "put", key, err)
	}
	if err := s.bucket.Put([]byte(key), bytes); err != nil {
		return NewStoreErrorWithKey(s.service, "put", key, err)
	}
	return nil
}

// Delete removes an item from the store by key.
func (s *BaseStore) Delete(key string) error {
	if err := s.bucket.Delete([]byte(key)); err != nil {
		return NewStoreErrorWithKey(s.service, "delete", key, err)
	}
	return nil
}

// Exists checks if an item exists in the store by key.
func (s *BaseStore) Exists(key string) bool {
	return s.bucket.Has([]byte(key))
}

// GetProto retrieves an item from the store and unmarshals it as protobuf.
func (s *BaseStore) GetProto(key string, dest proto.Message) error {
	data, err := s.bucket.Get([]byte(key))
	if err != nil {
		return NewStoreErrorWithKey(s.service, "get_proto", key, err)
	}
	if data == nil {
		return NewStoreErrorWithKey(s.service, "get_proto", key, ErrNotFound)
	}
	if err := proto.Unmarshal(data, dest); err != nil {
		return NewStoreErrorWithKey(s.service, "get_proto", key, err)
	}
	return nil
}

// PutProto stores an item in the store after marshaling it as protobuf.
func (s *BaseStore) PutProto(key string, src proto.Message) error {
	bytes, err := proto.Marshal(src)
	if err != nil {
		return NewStoreErrorWithKey(s.service, "put_proto", key, err)
	}
	if err := s.bucket.Put([]byte(key), bytes); err != nil {
		return NewStoreErrorWithKey(s.service, "put_proto", key, err)
	}
	return nil
}

// Count returns the number of items in the store.
func (s *BaseStore) Count() int {
	return s.bucket.Count()
}

// ForEach iterates over all items in the store.
func (s *BaseStore) ForEach(fn func(key string, value []byte) error) error {
	return s.bucket.ForEach(func(k, v []byte) error {
		return fn(string(k), v)
	})
}

// ScanPrefix iterates over items with a given prefix.
func (s *BaseStore) ScanPrefix(prefix string, fn func(key string, value []byte) error) error {
	iter := s.bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()
	for iter.Next() {
		if err := fn(string(iter.Key()), iter.Value()); err != nil {
			return err
		}
	}
	return iter.Error()
}

// DeleteByPrefix deletes all items with a given prefix.
func (s *BaseStore) DeleteByPrefix(prefix string) error {
	iter := s.bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()
	var keys [][]byte
	for iter.Next() {
		k := iter.Key()
		copied := make([]byte, len(k))
		copy(copied, k)
		keys = append(keys, copied)
	}
	if err := iter.Error(); err != nil {
		return err
	}
	for _, k := range keys {
		if err := s.bucket.Delete(k); err != nil {
			return err
		}
	}
	return nil
}

// ScanRange iterates over items within a key range.
func (s *BaseStore) ScanRange(start, end string, fn func(key string, value []byte) error) error {
	iter := s.bucket.ScanRange([]byte(start), []byte(end))
	defer iter.Close()
	for iter.Next() {
		if err := fn(string(iter.Key()), iter.Value()); err != nil {
			return err
		}
	}
	return iter.Error()
}

// ListResult represents the result of a list operation.
type ListResult[T any] struct {
	// Items contains the list of retrieved items.
	Items []*T
	// NextMarker is the marker for the next page of results.
	NextMarker string
	// IsTruncated indicates whether there are more results to retrieve.
	IsTruncated bool
}

// ListOptions contains options for listing items from the store.
type ListOptions struct {
	// Prefix filters items by key prefix.
	Prefix string
	// Marker specifies the starting point for the next page of results.
	Marker string
	// MaxItems specifies the maximum number of items to return.
	MaxItems int
}

// FilterFunc is a function type that filters items during listing.
type FilterFunc[T any] func(*T) bool

// ListProtoResult represents the result of a protobuf list operation.
type ListProtoResult[T proto.Message] struct {
	Items       []T
	NextMarker  string
	IsTruncated bool
}

// ListProto retrieves a list of protobuf items from the store.
func ListProto[T proto.Message](store *BaseStore, opts ListOptions, newFunc func() T, filter func(T) bool) (*ListProtoResult[T], error) {
	if opts.MaxItems <= 0 {
		opts.MaxItems = DefaultMaxItems
	}
	if opts.MaxItems > AbsoluteMaxItems {
		opts.MaxItems = AbsoluteMaxItems
	}

	var items []T
	count := 0
	started := opts.Marker == ""
	var lastKey string
	hasMore := false

	err := store.ForEach(func(key string, value []byte) error {
		if opts.Prefix != "" && !strings.HasPrefix(key, opts.Prefix) {
			return nil
		}

		if !started {
			if key == opts.Marker {
				started = true
				return nil
			} else if key > opts.Marker {
				started = true
			}
			if !started {
				return nil
			}
		}

		if count >= opts.MaxItems {
			hasMore = true
			return nil
		}

		item := newFunc()
		if err := proto.Unmarshal(value, item); err != nil {
			return err
		}

		if filter == nil || filter(item) {
			items = append(items, item)
			lastKey = key
			count++
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError(store.service, "list_proto", err)
	}

	nextMarker := ""
	if hasMore {
		nextMarker = lastKey
	}
	return &ListProtoResult[T]{
		Items:       items,
		NextMarker:  nextMarker,
		IsTruncated: hasMore,
	}, nil
}

// List retrieves a list of items from the store with the given options and optional filter.
func List[T any](store *BaseStore, opts ListOptions, filter FilterFunc[T]) (*ListResult[T], error) {
	if opts.MaxItems <= 0 {
		opts.MaxItems = DefaultMaxItems
	}
	if opts.MaxItems > AbsoluteMaxItems {
		opts.MaxItems = AbsoluteMaxItems
	}

	var items []*T
	count := 0
	started := opts.Marker == ""
	var lastKey string
	hasMore := false

	err := store.ForEach(func(key string, value []byte) error {
		if !strings.HasPrefix(key, opts.Prefix) {
			return nil
		}

		if !started {
			if key == opts.Marker {
				started = true
				return nil
			} else if key > opts.Marker {
				started = true
			}
			if !started {
				return nil
			}
		}

		if count >= opts.MaxItems {
			hasMore = true
			return nil
		}

		var item T
		if err := json.Unmarshal(value, &item); err != nil {
			return err
		}

		if filter == nil || filter(&item) {
			items = append(items, &item)
			lastKey = key
			count++
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError(store.service, "list", err)
	}

	nextMarker := ""
	if hasMore {
		nextMarker = lastKey
	}
	return &ListResult[T]{
		Items:       items,
		NextMarker:  nextMarker,
		IsTruncated: hasMore,
	}, nil
}

// GetOrCreateStore retrieves a cached store or creates a new one using the provided function.
// This is a thread-safe pattern using sync.Map.LoadOrStore.
func GetOrCreateStore[T any](stores *sync.Map, region string, createFn func() T) T {
	if cached, ok := stores.Load(region); ok {
		if typed, ok := cached.(T); ok {
			return typed
		}
	}
	store := createFn()
	actual, _ := stores.LoadOrStore(region, store)
	if typed, ok := actual.(T); ok {
		return typed
	}
	return store
}

// GetOrCreateStoreE retrieves a cached store or creates a new one with error handling.
// This is a thread-safe pattern using sync.Map.LoadOrStore.
func GetOrCreateStoreE[T any](stores *sync.Map, region string, createFn func() (T, error)) (T, error) {
	if cached, ok := stores.Load(region); ok {
		if typed, ok := cached.(T); ok {
			return typed, nil
		}
	}
	store, err := createFn()
	if err != nil {
		var zero T
		return zero, err
	}
	actual, _ := stores.LoadOrStore(region, store)
	if typed, ok := actual.(T); ok {
		return typed, nil
	}
	return store, nil
}
