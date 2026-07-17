package cloudfront

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const keyGroupBucketName = "cloudfront_key_groups"

// KeyGroupStore provides storage operations for CloudFront key groups.
type KeyGroupStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	mu         sync.Mutex
}

// NewKeyGroupStore creates a new KeyGroupStore instance.
func NewKeyGroupStore(store storage.BasicStorage, accountId string) *KeyGroupStore {
	return &KeyGroupStore{
		BaseStore:  common.NewBaseStore(store.Bucket(keyGroupBucketName), "cloudfront"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves a key group by its ID.
func (s *KeyGroupStore) Get(id string) (*KeyGroup, error) {
	var kg KeyGroup
	if err := s.BaseStore.Get(id, &kg); err != nil {
		return nil, NewStoreError("get_key_group", err)
	}
	return &kg, nil
}

// GetByName retrieves a key group by its name.
func (s *KeyGroupStore) GetByName(name string) (*KeyGroup, error) {
	return common.FindFirst[KeyGroup](s.BaseStore, func(kg *KeyGroup) bool {
		return kg.KeyGroupConfig != nil && kg.KeyGroupConfig.Name == name
	})
}

// Create creates a new key group.
func (s *KeyGroupStore) Create(config *KeyGroupConfig) (*KeyGroup, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := generateDistributionID()
	if err != nil {
		return nil, NewStoreError("create_key_group", err)
	}
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("create_key_group", err)
	}

	now := time.Now()
	kg := &KeyGroup{
		ID:             id,
		ARN:            s.arnBuilder.BuildKeyGroupARN(id),
		ETag:           etag,
		LastModifiedAt: now,
		KeyGroupConfig: config,
	}

	if err := s.BaseStore.Put(id, kg); err != nil {
		return nil, NewStoreError("create_key_group", err)
	}
	return kg, nil
}

// Update updates an existing key group.
func (s *KeyGroupStore) Update(id string, config *KeyGroupConfig) (*KeyGroup, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	kg, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	kg.KeyGroupConfig = config
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("update_key_group", err)
	}
	kg.ETag = etag
	kg.LastModifiedAt = time.Now()

	if err := s.BaseStore.Put(id, kg); err != nil {
		return nil, NewStoreError("update_key_group", err)
	}
	return kg, nil
}

// Delete removes a key group by its ID.
func (s *KeyGroupStore) Delete(id string) error {
	if _, err := s.Get(id); err != nil {
		return err
	}
	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_key_group", err)
	}
	return nil
}

// List returns a paginated list of key groups.
func (s *KeyGroupStore) List(marker string, maxItems int) (*KeyGroupListResult, error) {
	result, err := common.List[KeyGroup](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, nil)
	if err != nil {
		return nil, NewStoreError("list_key_groups", err)
	}
	return &KeyGroupListResult{
		KeyGroups:   result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}
