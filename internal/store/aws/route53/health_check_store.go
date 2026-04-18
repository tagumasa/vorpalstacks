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

// Get retrieves a health check by its ID.
func (s *HealthCheckStore) Get(id string) (*HealthCheck, error) {
	var healthCheck HealthCheck
	if err := s.BaseStore.Get(id, &healthCheck); err != nil {
		return nil, NewStoreError("get_health_check", err)
	}
	return &healthCheck, nil
}

// List returns health checks with pagination.
func (s *HealthCheckStore) List(marker string, maxItems int) (*HealthCheckListResult, error) {
	result, err := common.List[HealthCheck](s.BaseStore, common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	}, nil)
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

// Delete removes a health check.
func (s *HealthCheckStore) Delete(id string) error {
	if !s.Exists(id) {
		return NewStoreError("delete_health_check", common.ErrNotFound)
	}
	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_health_check", err)
	}
	return nil
}

// Exists checks whether a health check exists.
func (s *HealthCheckStore) Exists(id string) bool {
	return s.BaseStore.Exists(id)
}
