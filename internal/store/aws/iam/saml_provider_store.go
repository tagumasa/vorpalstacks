package iam

import (
	"time"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
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

// Create creates a new SAML provider with the given name, metadata document,
// and optional validity period.  When addPrivateKey is non-empty, a
// SAMLPrivateKey entry with a generated KeyId is stored.
func (s *SAMLProviderStore) Create(name, metadataDocument string, validUntil *time.Time, assertionEncryptionMode string, addPrivateKey string, tags []types.Tag) (*SAMLProvider, error) {
	arn := s.arnBuilder.SAMLProviderARN(name)
	var provider *SAMLProvider
	err := s.kl.WithLock(arn, func() error {
		if s.Exists(arn) {
			return NewStoreError("create_saml_provider", ErrSAMLProviderAlreadyExists)
		}
		now := time.Now().UTC()
		providerUUID, err := generatePrivateKeyID()
		if err != nil {
			return NewStoreError("create_saml_provider", err)
		}
		provider = &SAMLProvider{
			Arn:                     arn,
			AccountId:               s.arnBuilder.AccountID(),
			UUID:                    providerUUID,
			SAMLMetadataDocument:    metadataDocument,
			ValidUntil:              validUntil,
			CreateDate:              now,
			AssertionEncryptionMode: assertionEncryptionMode,
			Tags:                    tags,
		}
		if addPrivateKey != "" {
			keyId, err := generatePrivateKeyID()
			if err != nil {
				return NewStoreError("create_saml_provider", err)
			}
			provider.PrivateKeys = []SAMLPrivateKey{
				{KeyId: keyId, Key: addPrivateKey, AddedAt: now},
			}
		}
		return s.Put(provider)
	})
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// Update modifies the metadata document, validity period, and encryption
// settings of an existing SAML provider.  When addPrivateKey is non-empty,
// a new SAMLPrivateKey entry is appended.  When removePrivateKey is
// non-empty, the entry whose KeyId matches is removed; if no match is
// found, ErrSAMLPrivateKeyNotFound is returned.
func (s *SAMLProviderStore) Update(arn, metadataDocument string, validUntil *time.Time, assertionEncryptionMode string, addPrivateKey string, removePrivateKey string) error {
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
		if assertionEncryptionMode != "" {
			provider.AssertionEncryptionMode = assertionEncryptionMode
		}
		if addPrivateKey != "" {
			keyId, err := generatePrivateKeyID()
			if err != nil {
				return NewStoreError("update_saml_provider", err)
			}
			provider.PrivateKeys = append(provider.PrivateKeys, SAMLPrivateKey{
				KeyId:   keyId,
				Key:     addPrivateKey,
				AddedAt: time.Now().UTC(),
			})
		}
		if removePrivateKey != "" {
			found := false
			for i, pk := range provider.PrivateKeys {
				if pk.KeyId == removePrivateKey {
					provider.PrivateKeys = append(provider.PrivateKeys[:i], provider.PrivateKeys[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				return NewStoreError("update_saml_provider", ErrSAMLPrivateKeyNotFound)
			}
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
