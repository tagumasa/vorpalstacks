package cloudfront

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const responseHeadersPolicyBucketName = "cloudfront_response_headers_policies"

// ResponseHeadersPolicyStore manages CloudFront response headers policies.
type ResponseHeadersPolicyStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
	mu         sync.Mutex
}

// NewResponseHeadersPolicyStore creates a new response headers policy store.
func NewResponseHeadersPolicyStore(store storage.BasicStorage, accountId string) *ResponseHeadersPolicyStore {
	return &ResponseHeadersPolicyStore{
		BaseStore:  common.NewBaseStore(store.Bucket(responseHeadersPolicyBucketName), "cloudfront"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves a response headers policy by its ID.
func (s *ResponseHeadersPolicyStore) Get(id string) (*ResponseHeadersPolicy, error) {
	var policy ResponseHeadersPolicy
	if err := s.BaseStore.Get(id, &policy); err != nil {
		return nil, NewStoreError("get_response_headers_policy", err)
	}
	return &policy, nil
}

// GetByName retrieves a response headers policy by its name.
func (s *ResponseHeadersPolicyStore) GetByName(name string) (*ResponseHeadersPolicy, error) {
	return common.FindFirst[ResponseHeadersPolicy](s.BaseStore, func(p *ResponseHeadersPolicy) bool { return p.Name == name })
}

// Create creates a new response headers policy.
func (s *ResponseHeadersPolicyStore) Create(config *ResponseHeadersPolicyConfig) (*ResponseHeadersPolicy, error) {
	id, err := generateDistributionID()
	if err != nil {
		return nil, NewStoreError("create_response_headers_policy", err)
	}
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("create_response_headers_policy", err)
	}

	now := time.Now()
	policy := &ResponseHeadersPolicy{
		ID:                        id,
		ARN:                       s.arnBuilder.BuildResponseHeadersPolicyARN(id),
		ETag:                      etag,
		Name:                      config.Name,
		Comment:                   config.Comment,
		CorsConfig:                config.CorsConfig,
		CustomHeadersConfig:       config.CustomHeadersConfig,
		RemoveHeadersConfig:       config.RemoveHeadersConfig,
		SecurityHeadersConfig:     config.SecurityHeadersConfig,
		ServerTimingHeadersConfig: config.ServerTimingHeadersConfig,
		CreatedAt:                 now,
		LastModifiedAt:            now,
	}

	if err := s.BaseStore.Put(id, policy); err != nil {
		return nil, NewStoreError("create_response_headers_policy", err)
	}

	return policy, nil
}

// Update updates an existing response headers policy.
func (s *ResponseHeadersPolicyStore) Update(id string, config *ResponseHeadersPolicyConfig) (*ResponseHeadersPolicy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	policy, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	policy.Name = config.Name
	policy.Comment = config.Comment
	policy.CorsConfig = config.CorsConfig
	policy.CustomHeadersConfig = config.CustomHeadersConfig
	policy.RemoveHeadersConfig = config.RemoveHeadersConfig
	policy.SecurityHeadersConfig = config.SecurityHeadersConfig
	policy.ServerTimingHeadersConfig = config.ServerTimingHeadersConfig
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("update_response_headers_policy", err)
	}
	policy.ETag = etag
	policy.LastModifiedAt = time.Now()

	if err := s.BaseStore.Put(id, policy); err != nil {
		return nil, NewStoreError("update_response_headers_policy", err)
	}
	return policy, nil
}

// Delete removes a response headers policy by its ID.
func (s *ResponseHeadersPolicyStore) Delete(id string) error {
	if _, err := s.Get(id); err != nil {
		return err
	}
	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_response_headers_policy", err)
	}
	return nil
}

// List returns a paginated list of response headers policies.
func (s *ResponseHeadersPolicyStore) List(marker string, maxItems int) (*ResponseHeadersPolicyListResult, error) {
	result, err := common.List[ResponseHeadersPolicy](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, nil)
	if err != nil {
		return nil, NewStoreError("list_response_headers_policies", err)
	}

	return &ResponseHeadersPolicyListResult{
		ResponseHeadersPolicies: result.Items,
		IsTruncated:             result.IsTruncated,
		NextMarker:              result.NextMarker,
	}, nil
}

// Raw returns the underlying response headers policy store.
func (s *ResponseHeadersPolicyStore) Raw() *ResponseHeadersPolicyStore {
	return s
}
