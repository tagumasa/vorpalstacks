// Package cloudfront provides CloudFront storage functionality for vorpalstacks.
package cloudfront

import (
	"encoding/json"
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
	var found *OriginAccessControl
	err := s.ForEach(func(key string, value []byte) error {
		var oac OriginAccessControl
		if err := json.Unmarshal(value, &oac); err != nil {
			return err
		}
		if oac.Name == name {
			found = &oac
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_origin_access_control_by_name", err)
	}
	if found == nil {
		return nil, NewStoreError("get_origin_access_control_by_name", ErrNotFound)
	}
	return found, nil
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
	if maxItems <= 0 {
		maxItems = 100
	}

	var oacs []*OriginAccessControl
	count := 0
	started := marker == ""
	hasMore := false
	var lastKey string

	err := s.ForEach(func(key string, value []byte) error {
		var oac OriginAccessControl
		if err := json.Unmarshal(value, &oac); err != nil {
			return err
		}

		if !started {
			if oac.ID == marker {
				started = true
				return nil
			}
			return nil
		}

		if count < maxItems {
			oacs = append(oacs, &oac)
			lastKey = key
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_origin_access_controls", err)
	}

	var nextMarker string
	if hasMore {
		nextMarker = lastKey
	}

	return &OriginAccessControlListResult{
		OriginAccessControls: oacs,
		IsTruncated:          hasMore,
		NextMarker:           nextMarker,
	}, nil
}

// Raw returns the underlying origin access control store.
func (s *OriginAccessControlStore) Raw() *OriginAccessControlStore {
	return s
}
