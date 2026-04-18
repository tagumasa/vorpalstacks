// Package cloudfront provides CloudFront storage functionality for vorpalstacks.
package cloudfront

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const originAccessControlBucketName = "cloudfront_origin_access_controls"

// OriginAccessControlStore provides storage operations for CloudFront origin access controls.
type OriginAccessControlStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	mu         sync.Mutex
}

// NewOriginAccessControlStore creates a new OriginAccessControlStore instance.
func NewOriginAccessControlStore(store storage.BasicStorage, accountId string) *OriginAccessControlStore {
	return &OriginAccessControlStore{
		BaseStore:  common.NewBaseStore(store.Bucket(originAccessControlBucketName), "cloudfront"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves an origin access control by its ID.
func (s *OriginAccessControlStore) Get(id string) (*OriginAccessControl, error) {
	var oac OriginAccessControl
	if err := s.BaseStore.Get(id, &oac); err != nil {
		return nil, NewStoreError("get_origin_access_control", err)
	}
	return &oac, nil
}

// GetByName retrieves an origin access control by its name.
func (s *OriginAccessControlStore) GetByName(name string) (*OriginAccessControl, error) {
	return common.FindFirst[OriginAccessControl](s.BaseStore, func(o *OriginAccessControl) bool { return o.Name == name })
}

// Create creates a new origin access control.
func (s *OriginAccessControlStore) Create(config *OriginAccessControlConfig) (*OriginAccessControl, error) {
	id, err := generateDistributionID()
	if err != nil {
		return nil, NewStoreError("create_origin_access_control", err)
	}
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("create_origin_access_control", err)
	}

	now := time.Now()
	oac := &OriginAccessControl{
		ID:                            id,
		ARN:                           s.arnBuilder.BuildOriginAccessControlARN(id),
		ETag:                          etag,
		Name:                          config.Name,
		Description:                   config.Description,
		OriginAccessControlOriginType: config.OriginAccessControlOriginType,
		SigningBehavior:               config.SigningBehavior,
		SigningProtocol:               config.SigningProtocol,
		CreatedAt:                     now,
		LastModifiedAt:                now,
	}

	if err := s.BaseStore.Put(id, oac); err != nil {
		return nil, NewStoreError("create_origin_access_control", err)
	}

	return oac, nil
}

// Update updates an existing origin access control.
func (s *OriginAccessControlStore) Update(id string, config *OriginAccessControlConfig) (*OriginAccessControl, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	oac, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	oac.Name = config.Name
	oac.Description = config.Description
	oac.OriginAccessControlOriginType = config.OriginAccessControlOriginType
	oac.SigningBehavior = config.SigningBehavior
	oac.SigningProtocol = config.SigningProtocol
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("update_origin_access_control", err)
	}
	oac.ETag = etag
	oac.LastModifiedAt = time.Now()

	if err := s.BaseStore.Put(id, oac); err != nil {
		return nil, NewStoreError("update_origin_access_control", err)
	}
	return oac, nil
}

// Delete removes an origin access control by its ID.
func (s *OriginAccessControlStore) Delete(id string) error {
	if _, err := s.Get(id); err != nil {
		return err
	}
	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_origin_access_control", err)
	}
	return nil
}

// List returns a paginated list of origin access controls.
func (s *OriginAccessControlStore) List(marker string, maxItems int) (*OriginAccessControlListResult, error) {
	result, err := common.List[OriginAccessControl](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, nil)
	if err != nil {
		return nil, NewStoreError("list_origin_access_controls", err)
	}

	return &OriginAccessControlListResult{
		OriginAccessControls: result.Items,
		IsTruncated:          result.IsTruncated,
		NextMarker:           result.NextMarker,
	}, nil
}

// Raw returns the underlying origin access control store.
func (s *OriginAccessControlStore) Raw() *OriginAccessControlStore {
	return s
}
