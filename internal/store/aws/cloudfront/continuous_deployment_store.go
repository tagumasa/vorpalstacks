package cloudfront

import (
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const continuousDeploymentPolicyBucketName = "cloudfront_continuous_deployment_policies"

// ContinuousDeploymentPolicyStore provides storage operations for
// CloudFront continuous deployment policies.
type ContinuousDeploymentPolicyStore struct {
	*common.BaseStore
	mu sync.Mutex
}

// NewContinuousDeploymentPolicyStore creates a new store instance with the
// specified storage.
func NewContinuousDeploymentPolicyStore(store storage.BasicStorage) *ContinuousDeploymentPolicyStore {
	return &ContinuousDeploymentPolicyStore{
		BaseStore: common.NewBaseStore(store.Bucket(continuousDeploymentPolicyBucketName), "cloudfront"),
	}
}

// Get retrieves a continuous deployment policy by its ID.
func (s *ContinuousDeploymentPolicyStore) Get(id string) (*ContinuousDeploymentPolicy, error) {
	var policy ContinuousDeploymentPolicy
	if err := s.BaseStore.Get(id, &policy); err != nil {
		return nil, NewStoreError("get_continuous_deployment_policy", err)
	}
	return &policy, nil
}

// Create stores a new continuous deployment policy with a generated ID and
// ETag.
func (s *ContinuousDeploymentPolicyStore) Create(config *ContinuousDeploymentPolicyConfig) (*ContinuousDeploymentPolicy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := generateDistributionID()
	if err != nil {
		return nil, NewStoreError("create_continuous_deployment_policy", err)
	}
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("create_continuous_deployment_policy", err)
	}

	now := time.Now()
	policy := &ContinuousDeploymentPolicy{
		ID:                               id,
		ETag:                             etag,
		LastModifiedTime:                 now,
		ContinuousDeploymentPolicyConfig: config,
	}
	if err := s.BaseStore.Put(id, policy); err != nil {
		return nil, NewStoreError("create_continuous_deployment_policy", err)
	}
	return policy, nil
}

// Update replaces the configuration of an existing policy and rotates its
// ETag and last-modified timestamp.
func (s *ContinuousDeploymentPolicyStore) Update(id string, config *ContinuousDeploymentPolicyConfig) (*ContinuousDeploymentPolicy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	policy, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	policy.ContinuousDeploymentPolicyConfig = config
	etag, err := generateETag()
	if err != nil {
		return nil, NewStoreError("update_continuous_deployment_policy", err)
	}
	policy.ETag = etag
	policy.LastModifiedTime = time.Now()

	if err := s.BaseStore.Put(id, policy); err != nil {
		return nil, NewStoreError("update_continuous_deployment_policy", err)
	}
	return policy, nil
}

// Delete removes a policy by its ID.
func (s *ContinuousDeploymentPolicyStore) Delete(id string) error {
	if _, err := s.Get(id); err != nil {
		return err
	}
	if err := s.BaseStore.Delete(id); err != nil {
		return NewStoreError("delete_continuous_deployment_policy", err)
	}
	return nil
}

// List returns a paginated list of continuous deployment policies.
func (s *ContinuousDeploymentPolicyStore) List(marker string, maxItems int) (*ContinuousDeploymentPolicyListResult, error) {
	result, err := common.List[ContinuousDeploymentPolicy](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, nil)
	if err != nil {
		return nil, NewStoreError("list_continuous_deployment_policies", err)
	}
	return &ContinuousDeploymentPolicyListResult{
		Policies:    result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}

// ContinuousDeploymentPolicyListResult carries one page of policies.
type ContinuousDeploymentPolicyListResult struct {
	Policies    []*ContinuousDeploymentPolicy
	IsTruncated bool
	NextMarker  string
}
