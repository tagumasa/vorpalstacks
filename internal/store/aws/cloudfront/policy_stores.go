// Package cloudfront provides CloudFront storage functionality for vorpalstacks.
package cloudfront

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const (
	cachePolicyBucketName         = "cloudfront_cache_policies"
	originRequestPolicyBucketName = "cloudfront_origin_request_policies"
)

// CachePolicyStore provides storage operations for CloudFront cache policies.
type CachePolicyStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
}

// NewCachePolicyStore creates a new CachePolicyStore instance.
func NewCachePolicyStore(store storage.BasicStorage, accountId string) *CachePolicyStore {
	return &CachePolicyStore{
		BaseStore:  common.NewBaseStore(store.Bucket(cachePolicyBucketName), "cloudfront"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves a cache policy by its ID.
func (s *CachePolicyStore) Get(id string) (*CachePolicy, error) {
	var cachePolicy CachePolicy
	if err := s.BaseStore.Get(id, &cachePolicy); err != nil {
		return nil, NewStoreError("get_cache_policy", err)
	}
	return &cachePolicy, nil
}

// Put stores a cache policy by its ID (used for seeding managed policies).
func (s *CachePolicyStore) Put(id string, cp *CachePolicy) error {
	if err := s.BaseStore.Put(id, cp); err != nil {
		return NewStoreError("put_cache_policy", err)
	}
	return nil
}

// Create creates a new cache policy.
func (s *CachePolicyStore) Create(name, comment string, config *CachePolicyConfig) (*CachePolicy, error) {
	id, err := generateDistributionID()
	if err != nil {
		return nil, NewStoreError("create_cache_policy", err)
	}
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("create_cache_policy", err)
	}

	cachePolicy := &CachePolicy{
		ID:                id,
		Name:              name,
		CachePolicyConfig: config,
		ARN:               s.arnBuilder.BuildCachePolicyARN(id),
		ETag:              etag,
		CreatedAt:         time.Now(),
		ModifiedAt:        time.Now(),
	}

	if err := s.BaseStore.Put(id, cachePolicy); err != nil {
		return nil, NewStoreError("create_cache_policy", err)
	}

	return cachePolicy, nil
}

// Update updates an existing cache policy.
func (s *CachePolicyStore) Update(id string, config *CachePolicyConfig) (*CachePolicy, error) {
	cachePolicy, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	cachePolicy.CachePolicyConfig = config
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("update_cache_policy", err)
	}
	cachePolicy.ETag = etag
	cachePolicy.ModifiedAt = time.Now()

	if err := s.BaseStore.Put(id, cachePolicy); err != nil {
		return nil, NewStoreError("update_cache_policy", err)
	}
	return cachePolicy, nil
}

// Delete removes a cache policy by its ID.
func (s *CachePolicyStore) Delete(id string) error {
	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_cache_policy", err)
	}
	return nil
}

// List returns a paginated list of cache policies.
func (s *CachePolicyStore) List(marker string, maxItems int) (*CachePolicyListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var cachePolicies []*CachePolicy
	count := 0
	started := marker == ""
	hasMore := false
	var lastKey string

	err := s.ForEach(func(key string, value []byte) error {
		var cp CachePolicy
		if err := json.Unmarshal(value, &cp); err != nil {
			return err
		}

		if !started {
			if cp.ID == marker {
				started = true
				return nil
			}
			if cp.ID > marker {
				started = true
			}
			if !started {
				return nil
			}
		}

		if count < maxItems {
			cachePolicies = append(cachePolicies, &cp)
			lastKey = key
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_cache_policies", err)
	}

	var nextMarker string
	if hasMore {
		nextMarker = lastKey
	}

	return &CachePolicyListResult{
		CachePolicies: cachePolicies,
		IsTruncated:   hasMore,
		NextMarker:    nextMarker,
	}, nil
}

// OriginRequestPolicyStore provides storage operations for CloudFront origin request policies.
type OriginRequestPolicyStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
}

// NewOriginRequestPolicyStore creates a new OriginRequestPolicyStore instance.
func NewOriginRequestPolicyStore(store storage.BasicStorage, accountId string) *OriginRequestPolicyStore {
	return &OriginRequestPolicyStore{
		BaseStore:  common.NewBaseStore(store.Bucket(originRequestPolicyBucketName), "cloudfront"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves an origin request policy by its ID.
func (s *OriginRequestPolicyStore) Get(id string) (*OriginRequestPolicy, error) {
	var originRequestPolicy OriginRequestPolicy
	if err := s.BaseStore.Get(id, &originRequestPolicy); err != nil {
		return nil, NewStoreError("get_origin_request_policy", err)
	}
	return &originRequestPolicy, nil
}

// Put stores an origin request policy by its ID (used for seeding managed policies).
func (s *OriginRequestPolicyStore) Put(id string, orp *OriginRequestPolicy) error {
	if err := s.BaseStore.Put(id, orp); err != nil {
		return NewStoreError("put_origin_request_policy", err)
	}
	return nil
}

// Create creates a new origin request policy.
func (s *OriginRequestPolicyStore) Create(name, comment string, config *OriginRequestPolicyConfig) (*OriginRequestPolicy, error) {
	id, err := generateDistributionID()
	if err != nil {
		return nil, NewStoreError("create_origin_request_policy", err)
	}
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("create_origin_request_policy", err)
	}

	originRequestPolicy := &OriginRequestPolicy{
		ID:                        id,
		Name:                      name,
		OriginRequestPolicyConfig: config,
		ARN:                       s.arnBuilder.BuildOriginRequestPolicyARN(id),
		ETag:                      etag,
		CreatedAt:                 time.Now(),
		ModifiedAt:                time.Now(),
	}

	if err := s.BaseStore.Put(id, originRequestPolicy); err != nil {
		return nil, NewStoreError("create_origin_request_policy", err)
	}

	return originRequestPolicy, nil
}

// Update updates an existing origin request policy.
func (s *OriginRequestPolicyStore) Update(id string, config *OriginRequestPolicyConfig) (*OriginRequestPolicy, error) {
	originRequestPolicy, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	originRequestPolicy.OriginRequestPolicyConfig = config
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("update_origin_request_policy", err)
	}
	originRequestPolicy.ETag = etag
	originRequestPolicy.ModifiedAt = time.Now()

	if err := s.BaseStore.Put(id, originRequestPolicy); err != nil {
		return nil, NewStoreError("update_origin_request_policy", err)
	}
	return originRequestPolicy, nil
}

// Delete removes an origin request policy by its ID.
func (s *OriginRequestPolicyStore) Delete(id string) error {
	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_origin_request_policy", err)
	}
	return nil
}

// List returns a paginated list of origin request policies.
func (s *OriginRequestPolicyStore) List(marker string, maxItems int) (*OriginRequestPolicyListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var originRequestPolicies []*OriginRequestPolicy
	count := 0
	started := marker == ""
	hasMore := false
	var lastKey string

	err := s.ForEach(func(key string, value []byte) error {
		var orp OriginRequestPolicy
		if err := json.Unmarshal(value, &orp); err != nil {
			return err
		}

		if !started {
			if orp.ID == marker {
				started = true
				return nil
			}
			if orp.ID > marker {
				started = true
			}
			if !started {
				return nil
			}
		}

		if count < maxItems {
			originRequestPolicies = append(originRequestPolicies, &orp)
			lastKey = key
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_origin_request_policies", err)
	}

	var nextMarker string
	if hasMore {
		nextMarker = lastKey
	}

	return &OriginRequestPolicyListResult{
		OriginRequestPolicies: originRequestPolicies,
		IsTruncated:           hasMore,
		NextMarker:            nextMarker,
	}, nil
}

// Raw returns the underlying cache policy store.
func (s *CachePolicyStore) Raw() *CachePolicyStore {
	return s
}

// Raw returns the underlying origin request policy store.
func (s *OriginRequestPolicyStore) Raw() *OriginRequestPolicyStore {
	return s
}

// SeedManagedPolicies seeds AWS-managed cache and origin request policies.
func SeedManagedPolicies(cacheStore *CachePolicyStore, orpStore *OriginRequestPolicyStore) {
	now := time.Now()

	managedCachePolicies := []*CachePolicy{
		{
			ID:   "658327ea-f89d-4fab-a63d-7e88639e58f6",
			Name: "CachingOptimized",
			ARN:  cacheStore.arnBuilder.BuildCachePolicyARN("658327ea-f89d-4fab-a63d-7e88639e58f6"),
			CachePolicyConfig: &CachePolicyConfig{
				Name:       "CachingOptimized",
				Comment:    "AWS Managed CachingOptimized Policy",
				DefaultTTL: 86400,
				MaxTTL:     31536000,
				MinTTL:     0,
			},
			ETag:       "managed",
			CreatedAt:  now,
			ModifiedAt: now,
		},
		{
			ID:   "4135ea2d-6df8-44a3-9df3-4b5a84be39ad",
			Name: "CachingOptimizedForUncompressedObjects",
			ARN:  cacheStore.arnBuilder.BuildCachePolicyARN("4135ea2d-6df8-44a3-9df3-4b5a84be39ad"),
			CachePolicyConfig: &CachePolicyConfig{
				Name:       "CachingOptimizedForUncompressedObjects",
				Comment:    "AWS Managed CachingOptimizedForUncompressedObjects Policy",
				DefaultTTL: 86400,
				MaxTTL:     31536000,
				MinTTL:     0,
			},
			ETag:       "managed",
			CreatedAt:  now,
			ModifiedAt: now,
		},
		{
			ID:   "b2884449-e4de-46a7-ac36-70bcab13c0e4",
			Name: "UseOriginCacheControlHeaders",
			ARN:  cacheStore.arnBuilder.BuildCachePolicyARN("b2884449-e4de-46a7-ac36-70bcab13c0e4"),
			CachePolicyConfig: &CachePolicyConfig{
				Name:       "UseOriginCacheControlHeaders",
				Comment:    "AWS Managed UseOriginCacheControlHeaders Policy",
				DefaultTTL: 0,
				MaxTTL:     0,
				MinTTL:     0,
			},
			ETag:       "managed",
			CreatedAt:  now,
			ModifiedAt: now,
		},
		{
			ID:   "08627262-05a9-4f76-9ded-b5072e5f7663",
			Name: "DisableCaching",
			ARN:  cacheStore.arnBuilder.BuildCachePolicyARN("08627262-05a9-4f76-9ded-b5072e5f7663"),
			CachePolicyConfig: &CachePolicyConfig{
				Name:       "DisableCaching",
				Comment:    "AWS Managed DisableCaching Policy",
				DefaultTTL: 0,
				MaxTTL:     0,
				MinTTL:     0,
			},
			ETag:       "managed",
			CreatedAt:  now,
			ModifiedAt: now,
		},
		{
			ID:   "60669652-45c0-4ac8-834e-af3d0dddcd4a",
			Name: "Elemental-MediaStore",
			ARN:  cacheStore.arnBuilder.BuildCachePolicyARN("60669652-45c0-4ac8-834e-af3d0dddcd4a"),
			CachePolicyConfig: &CachePolicyConfig{
				Name:       "Elemental-MediaStore",
				Comment:    "AWS Managed Elemental-MediaStore Policy",
				DefaultTTL: 3600,
				MaxTTL:     86400,
				MinTTL:     0,
			},
			ETag:       "managed",
			CreatedAt:  now,
			ModifiedAt: now,
		},
	}

	for _, cp := range managedCachePolicies {
		if err := cacheStore.Put(cp.ID, cp); err != nil {
			logs.Error("Failed to seed managed cache policy", logs.String("id", cp.ID), logs.Err(err))
		}
	}

	managedORPPolicies := []*OriginRequestPolicy{
		{
			ID:   "88a5eaf4-2fd4-4709-b370-b4c650ea3fcf",
			Name: "CORS-S3Origin",
			ARN:  orpStore.arnBuilder.BuildOriginRequestPolicyARN("88a5eaf4-2fd4-4709-b370-b4c650ea3fcf"),
			OriginRequestPolicyConfig: &OriginRequestPolicyConfig{
				Name:    "CORS-S3Origin",
				Comment: "AWS Managed CORS-S3Origin Policy",
			},
			ETag:       "managed",
			CreatedAt:  now,
			ModifiedAt: now,
		},
		{
			ID:   "b689b0a8-53d0-40ab-baf2-68738e2966b1",
			Name: "AllViewer",
			ARN:  orpStore.arnBuilder.BuildOriginRequestPolicyARN("b689b0a8-53d0-40ab-baf2-68738e2966b1"),
			OriginRequestPolicyConfig: &OriginRequestPolicyConfig{
				Name:    "AllViewer",
				Comment: "AWS Managed AllViewer Policy",
			},
			ETag:       "managed",
			CreatedAt:  now,
			ModifiedAt: now,
		},
		{
			ID:   "216adef6-5c7f-47e4-b989-5492e0773f65",
			Name: "AllViewerExceptHostHeader",
			ARN:  orpStore.arnBuilder.BuildOriginRequestPolicyARN("216adef6-5c7f-47e4-b989-5492e0773f65"),
			OriginRequestPolicyConfig: &OriginRequestPolicyConfig{
				Name:    "AllViewerExceptHostHeader",
				Comment: "AWS Managed AllViewerExceptHostHeader Policy",
			},
			ETag:       "managed",
			CreatedAt:  now,
			ModifiedAt: now,
		},
		{
			ID:   "59781a5b-390d-41f3-afcb-af62939cc928",
			Name: "UserAgentRefererHeaders",
			ARN:  orpStore.arnBuilder.BuildOriginRequestPolicyARN("59781a5b-390d-41f3-afcb-af62939cc928"),
			OriginRequestPolicyConfig: &OriginRequestPolicyConfig{
				Name:    "UserAgentRefererHeaders",
				Comment: "AWS Managed UserAgentRefererHeaders Policy",
			},
			ETag:       "managed",
			CreatedAt:  now,
			ModifiedAt: now,
		},
		{
			ID:   "0e24e3e5-6262-4ce8-9c12-3fb4b7ffe3f2",
			Name: "None",
			ARN:  orpStore.arnBuilder.BuildOriginRequestPolicyARN("0e24e3e5-6262-4ce8-9c12-3fb4b7ffe3f2"),
			OriginRequestPolicyConfig: &OriginRequestPolicyConfig{
				Name:    "None",
				Comment: "AWS Managed None Policy",
			},
			ETag:       "managed",
			CreatedAt:  now,
			ModifiedAt: now,
		},
	}

	for _, orp := range managedORPPolicies {
		if err := orpStore.Put(orp.ID, orp); err != nil {
			logs.Error("Failed to seed managed origin request policy", logs.String("id", orp.ID), logs.Err(err))
		}
	}
}
