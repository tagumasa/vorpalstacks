package route53

// Package route53 provides Route 53 data store implementations for vorpalstacks.

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const changeBucketName = "route53_changes"

// ChangeStore manages Route 53 change records.
type ChangeStore struct {
	*common.BaseStore
}

// NewChangeStore creates a new ChangeStore.
func NewChangeStore(store storage.BasicStorage) *ChangeStore {
	return &ChangeStore{
		BaseStore: common.NewBaseStore(store.Bucket(changeBucketName), "route53"),
	}
}

// Get retrieves a change by its ID.
func (s *ChangeStore) Get(id string) (*ChangeInfo, error) {
	var change ChangeInfo
	if err := s.BaseStore.Get(id, &change); err != nil {
		return nil, NewStoreError("get_change", err)
	}
	return &change, nil
}

// Create creates a new change record.
func (s *ChangeStore) Create(change *ChangeInfo) error {
	if change.SubmittedAt.IsZero() {
		change.SubmittedAt = time.Now()
	}
	if change.Status == "" {
		change.Status = "INSYNC"
	}
	if err := s.BaseStore.Put(change.ID, change); err != nil {
		return NewStoreError("create_change", err)
	}
	return nil
}

// UpdateStatus updates the status of a change.
func (s *ChangeStore) UpdateStatus(id, status string) error {
	change, err := s.Get(id)
	if err != nil {
		return NewStoreError("update_change_status", err)
	}
	change.Status = status
	if err := s.BaseStore.Put(id, change); err != nil {
		return NewStoreError("update_change_status", err)
	}
	return nil
}

// UpdateStatusAndGet updates the status and returns the change.
func (s *ChangeStore) UpdateStatusAndGet(id, status string) (*ChangeInfo, error) {
	change, err := s.Get(id)
	if err != nil {
		return nil, NewStoreError("update_change_status", err)
	}
	change.Status = status
	if err := s.BaseStore.Put(id, change); err != nil {
		return nil, NewStoreError("update_change_status", err)
	}
	return change, nil
}

// Exists checks whether a change exists.
func (s *ChangeStore) Exists(id string) bool {
	return s.BaseStore.Exists(id)
}
