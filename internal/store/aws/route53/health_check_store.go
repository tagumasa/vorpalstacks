package route53

// Package route53 provides Route 53 data store implementations for vorpalstacks.

import (
	"fmt"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const healthCheckBucketName = "route53_health_checks"

// HealthCheckStore manages Route 53 health checks.
type HealthCheckStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
}

// NewHealthCheckStore creates a new HealthCheckStore.
func NewHealthCheckStore(store storage.BasicStorage, accountId string) *HealthCheckStore {
	return &HealthCheckStore{
		BaseStore:  common.NewBaseStore(store.Bucket(healthCheckBucketName), "route53"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves a health check by its ID. A deletion tombstone behaves as
// a missing record.
func (s *HealthCheckStore) Get(id string) (*HealthCheck, error) {
	var healthCheck HealthCheck
	if err := s.BaseStore.Get(id, &healthCheck); err != nil {
		return nil, NewStoreError("get_health_check", err)
	}
	if healthCheck.Deleted {
		return nil, NewStoreError("get_health_check", common.ErrNotFound)
	}
	return &healthCheck, nil
}

// GetByCallerReference retrieves a health check by its CallerReference,
// preferring live records over deletion tombstones. Used for the
// CreateHealthCheck retry semantics: the same CallerReference with the same
// settings returns the existing health check, and one matching a recently
// deleted health check fails with HealthCheckAlreadyExists.
func (s *HealthCheckStore) GetByCallerReference(callerRef string) (*HealthCheck, error) {
	if hc, err := common.FindFirst[HealthCheck](s.BaseStore, func(hc *HealthCheck) bool {
		return hc.CallerReference == callerRef && !hc.Deleted
	}); err == nil {
		return hc, nil
	}
	return common.FindFirst[HealthCheck](s.BaseStore, func(hc *HealthCheck) bool {
		return hc.CallerReference == callerRef && hc.Deleted
	})
}

// List returns health checks with pagination. Deletion tombstones are
// excluded.
func (s *HealthCheckStore) List(marker string, maxItems int) (*HealthCheckListResult, error) {
	result, err := common.List[HealthCheck](s.BaseStore, common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	}, func(hc *HealthCheck) bool { return !hc.Deleted })
	if err != nil {
		return nil, NewStoreError("list_health_checks", err)
	}

	return &HealthCheckListResult{
		HealthChecks: result.Items,
		IsTruncated:  result.IsTruncated,
		Marker:       result.NextMarker,
	}, nil
}

// Create creates a new health check.
func (s *HealthCheckStore) Create(healthCheck *HealthCheck) error {
	if s.Exists(healthCheck.ID) {
		return NewStoreError("create_health_check", common.ErrAlreadyExists)
	}
	healthCheck.CreatedAt = time.Now()
	healthCheck.HealthCheckVersion = "1"
	if err := s.BaseStore.Put(healthCheck.ID, healthCheck); err != nil {
		return NewStoreError("create_health_check", err)
	}
	return nil
}

// Update updates an existing health check.
func (s *HealthCheckStore) Update(healthCheck *HealthCheck) error {
	existing, err := s.Get(healthCheck.ID)
	if err != nil {
		return NewStoreError("update_health_check", err)
	}
	version := 1
	if _, err := fmt.Sscanf(existing.HealthCheckVersion, "%d", &version); err == nil {
		version++
	}
	healthCheck.HealthCheckVersion = fmt.Sprintf("%d", version)
	healthCheck.CreatedAt = existing.CreatedAt
	if err := s.BaseStore.Put(healthCheck.ID, healthCheck); err != nil {
		return NewStoreError("update_health_check", err)
	}
	return nil
}

// Delete retains the health check as a deletion tombstone instead of
// removing the record: the CallerReference must stay resolvable for a while
// so a CreateHealthCheck retry with the same reference can fail with
// HealthCheckAlreadyExists, matching the documented retry semantics.
func (s *HealthCheckStore) Delete(id string) error {
	healthCheck, err := s.Get(id)
	if err != nil {
		return NewStoreError("delete_health_check", common.ErrNotFound)
	}
	healthCheck.Deleted = true
	healthCheck.DeletedAt = time.Now()
	if err := s.BaseStore.Put(healthCheck.ID, healthCheck); err != nil {
		return NewStoreError("delete_health_check", err)
	}
	return nil
}

// Exists checks whether a live health check exists. Deletion tombstones do
// not count.
func (s *HealthCheckStore) Exists(id string) bool {
	hc, err := s.Get(id)
	return err == nil && hc != nil
}
