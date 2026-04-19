package iam

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

const oidcProviderBucketName = "iam_oidc_providers"

// OpenIDConnectProviderStore provides storage operations for IAM OpenID Connect providers.
type OpenIDConnectProviderStore struct {
	entityStore[OpenIDConnectProvider]
	arnBuilder *ARNBuilder
}

// NewOpenIDConnectProviderStore creates a new OpenIDConnectProviderStore instance.
func NewOpenIDConnectProviderStore(store storage.BasicStorage, accountID string) *OpenIDConnectProviderStore {
	return &OpenIDConnectProviderStore{
		entityStore: newEntityStore[OpenIDConnectProvider](store, oidcProviderBucketName),
		arnBuilder:  NewARNBuilder(accountID),
	}
}

// Get retrieves an OpenID Connect provider by its ARN.
func (s *OpenIDConnectProviderStore) Get(arn string) (*OpenIDConnectProvider, error) {
	var provider OpenIDConnectProvider
	if err := s.BaseStore.Get(arn, &provider); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_oidc_provider", ErrOpenIDConnectProviderNotFound)
		}
		return nil, NewStoreError("get_oidc_provider", err)
	}
	return &provider, nil
}

// Put stores an OpenID Connect provider, keyed by its ARN.
func (s *OpenIDConnectProviderStore) Put(provider *OpenIDConnectProvider) error {
	if provider.Tags == nil {
		provider.Tags = []types.Tag{}
	}
	return s.BaseStore.Put(provider.Arn, provider)
}

// Create creates a new OpenID Connect provider with the specified URL, thumbprints, client IDs, and tags.
func (s *OpenIDConnectProviderStore) Create(url string, thumbprintList, clientIdList []string, tags []types.Tag) (*OpenIDConnectProvider, error) {
	arn := s.arnBuilder.OpenIDConnectProviderARN(url)
	if s.Exists(arn) {
		return nil, NewStoreError("create_oidc_provider", ErrOpenIDConnectProviderAlreadyExists)
	}
	provider := &OpenIDConnectProvider{
		Arn:            arn,
		AccountId:      s.arnBuilder.AccountID(),
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
	items, err := common.ListAll[OpenIDConnectProvider](s.BaseStore)
	if err != nil {
		return nil, NewStoreError("list_oidc_providers", err)
	}
	return &OpenIDConnectProviderListResult{OpenIDConnectProviderList: items}, nil
}

// GetByArn retrieves an OpenID Connect provider by its ARN.
func (s *OpenIDConnectProviderStore) GetByArn(arn string) (*OpenIDConnectProvider, error) {
	return s.Get(arn)
}

// ListByPrefix returns OpenID Connect providers whose ARNs match the given prefix.
func (s *OpenIDConnectProviderStore) ListByPrefix(prefix string) ([]*OpenIDConnectProvider, error) {
	return listEntitiesByPrefix(s.BaseStore, prefix, func(p *OpenIDConnectProvider) string { return p.Arn })
}
