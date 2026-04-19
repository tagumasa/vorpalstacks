package waf

import (
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

type wafResourceAccessor[T any] struct {
	getIDFn        func(*T) string
	getARNFn       func(*T) string
	setARNFn       func(*T, string)
	getLockTokenFn func(*T) string
	setLockTokenFn func(*T, string)
	setModifiedFn  func(*T)
}

// ResourceStore provides generic CRUD operations for WAF resources.
// It encapsulates common patterns shared by WebACL, RuleGroup, IPSet, and RegexPatternSet stores.
type ResourceStore[T any] struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	mu         sync.Mutex
	accessor   wafResourceAccessor[T]
}

// NewResourceStore creates a new generic WAF resource store backed by the given storage bucket.
func NewResourceStore[T any](store storage.BasicStorage, bucketName string, arnBuilder *ARNBuilder, accessor wafResourceAccessor[T]) *ResourceStore[T] {
	return &ResourceStore[T]{
		BaseStore:  common.NewBaseStore(store.Bucket(bucketName), "waf"),
		arnBuilder: arnBuilder,
		accessor:   accessor,
	}
}

// Get retrieves a WAF resource by its ID.
func (s *ResourceStore[T]) Get(id string) (*T, error) {
	var resource T
	if err := s.BaseStore.Get(id, &resource); err != nil {
		return nil, NewStoreError("get_resource", err)
	}
	return &resource, nil
}

// GetByARN finds a WAF resource matching the given ARN.
func (s *ResourceStore[T]) GetByARN(arn string) (*T, error) {
	return common.FindFirst[T](s.BaseStore, func(r *T) bool {
		return s.accessor.getARNFn(r) == arn
	})
}

// Delete removes a WAF resource by ID, verifying the lock token if provided.
func (s *ResourceStore[T]) Delete(id, lockToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	resource, err := s.Get(id)
	if err != nil {
		return err
	}

	if lockToken != "" && s.accessor.getLockTokenFn(resource) != lockToken {
		return NewStoreError("delete_resource", ErrLockTokenMismatch)
	}

	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_resource", err)
	}
	return nil
}

// UpdateWithLockToken applies the given update function to a WAF resource, verifying the lock token if provided.
// On success, the modified timestamp and lock token are refreshed automatically.
func (s *ResourceStore[T]) UpdateWithLockToken(id, lockToken string, updateFn func(*T) error, opName string) (*T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	resource, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if lockToken != "" && s.accessor.getLockTokenFn(resource) != lockToken {
		return nil, NewStoreError(opName, ErrLockTokenMismatch)
	}

	if err := updateFn(resource); err != nil {
		return nil, err
	}

	s.accessor.setModifiedFn(resource)
	s.accessor.setLockTokenFn(resource, GenerateLockToken())

	if err := s.BaseStore.Put(id, resource); err != nil {
		return nil, NewStoreError(opName, err)
	}
	return resource, nil
}

// Put persists a WAF resource under the given ID.
func (s *ResourceStore[T]) Put(id string, resource *T, opName string) error {
	if err := s.BaseStore.Put(id, resource); err != nil {
		return NewStoreError(opName, err)
	}
	return nil
}

// SetTimestamps refreshes the modified timestamp and generates a new lock token for a WAF resource.
func SetTimestamps[T any](accessor *wafResourceAccessor[T], resource *T) {
	accessor.setModifiedFn(resource)
	accessor.setLockTokenFn(resource, GenerateLockToken())
}
