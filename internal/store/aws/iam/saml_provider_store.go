package iam

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const samlProviderBucketName = "iam_saml_providers"

// SAMLProviderStore provides storage operations for IAM SAML providers.
type SAMLProviderStore struct {
	bucket     storage.Bucket
	arnBuilder *ARNBuilder
	kl         common.KeyLocker
}

// NewSAMLProviderStore creates a new SAMLProviderStore instance.
func NewSAMLProviderStore(store storage.BasicStorage, accountID string) *SAMLProviderStore {
	return &SAMLProviderStore{
		bucket:     store.Bucket(samlProviderBucketName),
		arnBuilder: NewARNBuilder(accountID),
	}
}

// Get retrieves a SAML provider by its ARN.
func (s *SAMLProviderStore) Get(arn string) (*SAMLProvider, error) {
	data, err := s.bucket.Get([]byte(arn))
	if err != nil {
		return nil, NewStoreError("get_saml_provider", err)
	}
	if data == nil {
		return nil, NewStoreError("get_saml_provider", ErrSAMLProviderNotFound)
	}
	var provider SAMLProvider
	if err := json.Unmarshal(data, &provider); err != nil {
		return nil, NewStoreError("get_saml_provider", err)
	}
	return &provider, nil
}

// Put stores a SAML provider, keyed by its ARN.
func (s *SAMLProviderStore) Put(provider *SAMLProvider) error {
	if provider.Tags == nil {
		provider.Tags = []Tag{}
	}
	data, err := json.Marshal(provider)
	if err != nil {
		return NewStoreError("put_saml_provider", err)
	}
	if err := s.bucket.Put([]byte(provider.Arn), data); err != nil {
		return NewStoreError("put_saml_provider", err)
	}
	return nil
}

// Delete removes a SAML provider by its ARN.
func (s *SAMLProviderStore) Delete(arn string) error {
	if err := s.bucket.Delete([]byte(arn)); err != nil {
		return NewStoreError("delete_saml_provider", err)
	}
	return nil
}

// Exists reports whether a SAML provider exists for the given ARN.
func (s *SAMLProviderStore) Exists(arn string) bool {
	return s.bucket.Has([]byte(arn))
}

// Create creates a new SAML provider with the given name, metadata document, and optional validity period.
func (s *SAMLProviderStore) Create(name, metadataDocument string, validUntil *time.Time, tags []Tag) (*SAMLProvider, error) {
	arn := s.arnBuilder.SAMLProviderARN(name)
	if s.Exists(arn) {
		return nil, ErrSAMLProviderAlreadyExists
	}
	id, err := GenerateSAMLProviderID()
	if err != nil {
		return nil, err
	}
	provider := &SAMLProvider{
		Arn:                  arn,
		AccountId:            id,
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
	var providers []*SAMLProvider
	err := s.bucket.ForEach(func(k, v []byte) error {
		var provider SAMLProvider
		if err := json.Unmarshal(v, &provider); err != nil {
			return err
		}
		providers = append(providers, &provider)
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_saml_providers", err)
	}
	return &SAMLProviderListResult{SAMLProviders: providers}, nil
}

// GetByArn retrieves a SAML provider by its ARN.
func (s *SAMLProviderStore) GetByArn(arn string) (*SAMLProvider, error) {
	return s.Get(arn)
}

// Count returns the total number of SAML providers.
func (s *SAMLProviderStore) Count() int {
	return s.bucket.Count()
}

// ListByPrefix returns SAML providers whose ARNs match the given prefix.
func (s *SAMLProviderStore) ListByPrefix(prefix string) ([]*SAMLProvider, error) {
	var providers []*SAMLProvider
	err := s.bucket.ForEach(func(k, v []byte) error {
		var provider SAMLProvider
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
		return nil, NewStoreError("list_saml_providers_by_prefix", err)
	}
	return providers, nil
}
