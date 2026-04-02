// Package cloudfront provides CloudFront storage functionality for vorpalstacks.
package cloudfront

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const distributionBucketName = "cloudfront_distributions"

// DistributionStore provides storage operations for CloudFront distributions.
type DistributionStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
}

// NewDistributionStore creates a new DistributionStore instance with the specified storage and account ID.
func NewDistributionStore(store storage.BasicStorage, accountId string) *DistributionStore {
	return &DistributionStore{
		BaseStore:  common.NewBaseStore(store.Bucket(distributionBucketName), "cloudfront"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves a CloudFront distribution by its ID from the store.
// Returns the distribution or an error if not found.
func (s *DistributionStore) Get(id string) (*Distribution, error) {
	var distribution Distribution
	if err := s.BaseStore.Get(id, &distribution); err != nil {
		return nil, NewStoreError("get_distribution", err)
	}
	return &distribution, nil
}

// GetByCallerReference retrieves a CloudFront distribution by its caller reference.
// Returns the distribution or ErrDistributionNotFound if not found.
func (s *DistributionStore) GetByCallerReference(callerRef string) (*Distribution, error) {
	normalizedRef := normalizeCallerReference(callerRef)
	var found *Distribution
	err := s.ForEach(func(key string, value []byte) error {
		var dist Distribution
		if err := json.Unmarshal(value, &dist); err != nil {
			return err
		}
		if dist.CallerReference == normalizedRef {
			found = &dist
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_by_caller_reference", err)
	}
	if found == nil {
		return nil, ErrDistributionNotFound
	}
	return found, nil
}

// Create creates a new CloudFront distribution in the store.
// If a distribution with the same CallerReference already exists, it returns the existing distribution (idempotent).
// Returns the created distribution or an error if creation fails.
func (s *DistributionStore) Create(callerReference string, config *DistributionConfig) (*Distribution, error) {
	normalizedRef := normalizeCallerReference(callerReference)

	if existing, err := s.GetByCallerReference(callerReference); err == nil {
		return existing, nil
	}

	id, err := generateDistributionID()
	if err != nil {
		return nil, NewStoreError("generate_distribution_id", err)
	}
	domainName := generateDomainName(id)
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("generate_etag", err)
	}

	now := time.Now()
	distribution := &Distribution{
		ID:                 id,
		CallerReference:    normalizedRef,
		DistributionConfig: config,
		DomainName:         domainName,
		ARN:                s.arnBuilder.BuildDistributionARN(id),
		ETag:               etag,
		Status:             "Deployed",
		CreatedAt:          now,
		LastModifiedAt:     now,
		Enabled:            config.Enabled,
	}

	if err := s.BaseStore.Put(id, distribution); err != nil {
		return nil, NewStoreError("create_distribution", err)
	}

	return distribution, nil
}

// Update updates an existing CloudFront distribution in the store.
// Returns the updated distribution or an error if the distribution does not exist.
func (s *DistributionStore) Update(id string, config *DistributionConfig) (*Distribution, error) {
	distribution, err := s.Get(id)
	if err != nil {
		return nil, NewStoreError("update_distribution", err)
	}
	return distribution, nil
}

// UpdateWithLastModified updates an existing CloudFront distribution with last modified time.
// Returns updated distribution or an error if distribution does not exist.
func (s *DistributionStore) UpdateWithLastModified(id string, distribution *Distribution) error {
	distribution.LastModifiedAt = time.Now()
	etag, err := generateETag()
	if err != nil {
		return NewStoreError("update_distribution", err)
	}
	distribution.ETag = etag

	if err := s.BaseStore.Put(id, distribution); err != nil {
		return NewStoreError("update_distribution", err)
	}
	return nil
}

// Delete deletes a CloudFront distribution by its ID from the store.
func (s *DistributionStore) Delete(id string) error {
	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_distribution", err)
	}
	return nil
}

// List returns a list of CloudFront distributions from the store with pagination support.
func (s *DistributionStore) List(marker string, maxItems int) (*DistributionListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var distributions []*Distribution
	count := 0
	started := marker == ""
	hasMore := false
	var lastKey string

	err := s.ForEach(func(key string, value []byte) error {
		var dist Distribution
		if err := json.Unmarshal(value, &dist); err != nil {
			return err
		}

		if !started {
			if dist.ID == marker || dist.ID > marker {
				started = true
			}
			if !started {
				return nil
			}
			if dist.ID == marker {
				return nil
			}
		}

		if count < maxItems {
			distributions = append(distributions, &dist)
			lastKey = key
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_distributions", err)
	}

	var nextMarker string
	if hasMore {
		nextMarker = lastKey
	}

	return &DistributionListResult{
		Distributions: distributions,
		IsTruncated:   hasMore,
		NextMarker:    nextMarker,
	}, nil
}

// Raw returns the underlying distribution store.
func (s *DistributionStore) Raw() *DistributionStore {
	return s
}
