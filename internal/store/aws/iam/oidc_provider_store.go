package iam

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const oidcProviderBucketName = "iam_oidc_providers"

// OpenIDConnectProviderStore provides storage operations for IAM OpenID Connect providers.
type OpenIDConnectProviderStore struct {
	bucket     storage.Bucket
	arnBuilder *ARNBuilder
	kl         common.KeyLocker
}

// NewOpenIDConnectProviderStore creates a new OpenIDConnectProviderStore instance.
func NewOpenIDConnectProviderStore(store storage.BasicStorage, accountID string) *OpenIDConnectProviderStore {
	return &OpenIDConnectProviderStore{
		bucket:     store.Bucket(oidcProviderBucketName),
		arnBuilder: NewARNBuilder(accountID),
	}
}

// Get retrieves an OpenID Connect provider by its ARN.
func (s *OpenIDConnectProviderStore) Get(arn string) (*OpenIDConnectProvider, error) {
	data, err := s.bucket.Get([]byte(arn))
	if err != nil {
		return nil, NewStoreError("get_oidc_provider", err)
	}
	if data == nil {
		return nil, NewStoreError("get_oidc_provider", ErrOpenIDConnectProviderNotFound)
	}
	var provider OpenIDConnectProvider
	if err := json.Unmarshal(data, &provider); err != nil {
		return nil, NewStoreError("get_oidc_provider", err)
	}
	return &provider, nil
}

// Put stores an OpenID Connect provider, keyed by its ARN.
func (s *OpenIDConnectProviderStore) Put(provider *OpenIDConnectProvider) error {
	if provider.Tags == nil {
		provider.Tags = []Tag{}
	}
	data, err := json.Marshal(provider)
	if err != nil {
		return NewStoreError("put_oidc_provider", err)
	}
	if err := s.bucket.Put([]byte(provider.Arn), data); err != nil {
		return NewStoreError("put_oidc_provider", err)
	}
	return nil
}

// Delete removes an OpenID Connect provider by its ARN.
func (s *OpenIDConnectProviderStore) Delete(arn string) error {
	if err := s.bucket.Delete([]byte(arn)); err != nil {
		return NewStoreError("delete_oidc_provider", err)
	}
	return nil
}

// Exists reports whether an OpenID Connect provider exists for the given ARN.
func (s *OpenIDConnectProviderStore) Exists(arn string) bool {
	return s.bucket.Has([]byte(arn))
}

// Create creates a new OpenID Connect provider with the specified URL, thumbprints, client IDs, and tags.
func (s *OpenIDConnectProviderStore) Create(url string, thumbprintList, clientIdList []string, tags []Tag) (*OpenIDConnectProvider, error) {
	arn := s.arnBuilder.OpenIDConnectProviderARN(url)
	if s.Exists(arn) {
		return nil, ErrOpenIDConnectProviderAlreadyExists
	}
	id, err := GenerateOpenIDConnectProviderID()
	if err != nil {
		return nil, err
	}
	provider := &OpenIDConnectProvider{
		Arn:            arn,
		AccountId:      id,
		URL:            url,
		ThumbprintList: thumbprintList,
		ClientIDList:   clientIdList,
		CreateDate:     time.Now().UTC(),
		Tags:           tags,
	}
	if err := s.Put(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

// Update modifies the thumbprint list and client ID list of an existing OpenID Connect provider.
func (s *OpenIDConnectProviderStore) Update(arn string, thumbprintList, clientIdList []string) error {
	return s.kl.WithLock(arn, func() error {
		provider, err := s.Get(arn)
		if err != nil {
			return err
		}
		if thumbprintList != nil {
			provider.ThumbprintList = thumbprintList
		}
		if clientIdList != nil {
			provider.ClientIDList = clientIdList
		}
		now := time.Now().UTC()
		provider.LastModifiedDate = &now
		return s.Put(provider)
	})
}

// List returns all OpenID Connect providers.
func (s *OpenIDConnectProviderStore) List() (*OpenIDConnectProviderListResult, error) {
	var providers []*OpenIDConnectProvider
	err := s.bucket.ForEach(func(k, v []byte) error {
		var provider OpenIDConnectProvider
		if err := json.Unmarshal(v, &provider); err != nil {
			return err
		}
		providers = append(providers, &provider)
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_oidc_providers", err)
	}
	return &OpenIDConnectProviderListResult{OpenIDConnectProviderList: providers}, nil
}

// GetByArn retrieves an OpenID Connect provider by its ARN.
func (s *OpenIDConnectProviderStore) GetByArn(arn string) (*OpenIDConnectProvider, error) {
	return s.Get(arn)
}

// Count returns the total number of OpenID Connect providers.
func (s *OpenIDConnectProviderStore) Count() int {
	return s.bucket.Count()
}

// ListByPrefix returns OpenID Connect providers whose ARNs match the given prefix.
func (s *OpenIDConnectProviderStore) ListByPrefix(prefix string) ([]*OpenIDConnectProvider, error) {
	var providers []*OpenIDConnectProvider
	err := s.bucket.ForEach(func(k, v []byte) error {
		var provider OpenIDConnectProvider
		if err := json.Unmarshal(v, &provider); err != nil {
			return err
		}
		if prefix != "" && !strings.HasPrefix(provider.Arn, prefix) {
			return nil
		}
		providers = append(providers, &provider)
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_oidc_providers_by_prefix", err)
	}
	return providers, nil
}
