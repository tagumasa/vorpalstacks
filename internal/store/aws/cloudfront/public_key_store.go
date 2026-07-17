package cloudfront

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const publicKeyBucketName = "cloudfront_public_keys"

// PublicKeyStore provides storage operations for CloudFront public keys.
type PublicKeyStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	mu         sync.Mutex
}

// NewPublicKeyStore creates a new PublicKeyStore instance.
func NewPublicKeyStore(store storage.BasicStorage, accountId string) *PublicKeyStore {
	return &PublicKeyStore{
		BaseStore:  common.NewBaseStore(store.Bucket(publicKeyBucketName), "cloudfront"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves a public key by its ID.
func (s *PublicKeyStore) Get(id string) (*PublicKey, error) {
	var pk PublicKey
	if err := s.BaseStore.Get(id, &pk); err != nil {
		return nil, NewStoreError("get_public_key", err)
	}
	return &pk, nil
}

// GetByName retrieves a public key by its name.
func (s *PublicKeyStore) GetByName(name string) (*PublicKey, error) {
	return common.FindFirst[PublicKey](s.BaseStore, func(pk *PublicKey) bool {
		return pk.PublicKeyConfig != nil && pk.PublicKeyConfig.Name == name
	})
}

// Create creates a new public key.
func (s *PublicKeyStore) Create(config *PublicKeyConfig) (*PublicKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := generateDistributionID()
	if err != nil {
		return nil, NewStoreError("create_public_key", err)
	}
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("create_public_key", err)
	}

	now := time.Now()
	pk := &PublicKey{
		ID:              id,
		ARN:             s.arnBuilder.BuildPublicKeyARN(id),
		ETag:            etag,
		CreatedTime:     now,
		PublicKeyConfig: config,
	}

	if err := s.BaseStore.Put(id, pk); err != nil {
		return nil, NewStoreError("create_public_key", err)
	}
	return pk, nil
}

// Update updates an existing public key.
func (s *PublicKeyStore) Update(id string, config *PublicKeyConfig) (*PublicKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pk, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	pk.PublicKeyConfig = config
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("update_public_key", err)
	}
	pk.ETag = etag

	if err := s.BaseStore.Put(id, pk); err != nil {
		return nil, NewStoreError("update_public_key", err)
	}
	return pk, nil
}

// Delete removes a public key by its ID.
func (s *PublicKeyStore) Delete(id string) error {
	if _, err := s.Get(id); err != nil {
		return err
	}
	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_public_key", err)
	}
	return nil
}

// List returns a paginated list of public keys.
func (s *PublicKeyStore) List(marker string, maxItems int) (*PublicKeyListResult, error) {
	result, err := common.List[PublicKey](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, nil)
	if err != nil {
		return nil, NewStoreError("list_public_keys", err)
	}
	return &PublicKeyListResult{
		PublicKeys:  result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}
