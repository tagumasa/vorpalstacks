package iam

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

const samlProviderBucketName = "iam_saml_providers"

// SAMLProviderStore provides storage operations for IAM SAML providers.
type SAMLProviderStore struct {
	entityStore[SAMLProvider]
	arnBuilder *ARNBuilder
}

// NewSAMLProviderStore creates a new SAMLProviderStore instance.
func NewSAMLProviderStore(store storage.BasicStorage, accountID string) *SAMLProviderStore {
	return &SAMLProviderStore{
		entityStore: newEntityStore[SAMLProvider](store, samlProviderBucketName),
		arnBuilder:  NewARNBuilder(accountID),
	}
}

// Get retrieves a SAML provider by its ARN.
func (s *SAMLProviderStore) Get(arn string) (*SAMLProvider, error) {
	var provider SAMLProvider
	if err := s.BaseStore.Get(arn, &provider); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_saml_provider", ErrSAMLProviderNotFound)
		}
		return nil, NewStoreError("get_saml_provider", err)
	}
	return &provider, nil
}

// Put stores a SAML provider, keyed by its ARN.
func (s *SAMLProviderStore) Put(provider *SAMLProvider) error {
	if provider.Tags == nil {
		provider.Tags = []types.Tag{}
	}
	return s.BaseStore.Put(provider.Arn, provider)
}

// Create creates a new SAML provider with the given name, metadata document, and optional validity period.
func (s *SAMLProviderStore) Create(name, metadataDocument string, validUntil *time.Time, tags []types.Tag) (*SAMLProvider, error) {
	arn := s.arnBuilder.SAMLProviderARN(name)
	if s.Exists(arn) {
		return nil, NewStoreError("create_saml_provider", ErrSAMLProviderAlreadyExists)
	}
	provider := &SAMLProvider{
		Arn:                  arn,
		AccountId:            s.arnBuilder.AccountID(),
		SAMLMetadataDocument: metadataDocument,
		ValidUntil:           validUntil,
		CreateDate:           time.Now().UTC(),
		Tags:                 tags,
	}
	if err := s.Put(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

// Update modifies the metadata document and validity period of an existing SAML provider.
func (s *SAMLProviderStore) Update(arn, metadataDocument string, validUntil *time.Time) error {
	return s.kl.WithLock(arn, func() error {
		provider, err := s.Get(arn)
		if err != nil {
			return err
		}
		if metadataDocument != "" {
			provider.SAMLMetadataDocument = metadataDocument
		}
		if validUntil != nil {
			provider.ValidUntil = validUntil
		}
		return s.Put(provider)
	})
}

// List returns all SAML providers.
func (s *SAMLProviderStore) List() (*SAMLProviderListResult, error) {
	items, err := common.ListAll[SAMLProvider](s.BaseStore)
	if err != nil {
		return nil, NewStoreError("list_saml_providers", err)
	}
	return &SAMLProviderListResult{SAMLProviders: items}, nil
}

// GetByArn retrieves a SAML provider by its ARN.
func (s *SAMLProviderStore) GetByArn(arn string) (*SAMLProvider, error) {
	return s.Get(arn)
}

// ListByPrefix returns SAML providers whose ARNs match the given prefix.
func (s *SAMLProviderStore) ListByPrefix(prefix string) ([]*SAMLProvider, error) {
	return listEntitiesByPrefix(s.BaseStore, prefix, func(p *SAMLProvider) string { return p.Arn })
}
