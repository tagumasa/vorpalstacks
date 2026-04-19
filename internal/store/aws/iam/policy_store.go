// Package iam provides AWS IAM store functionality for vorpalstacks.
package iam

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

const policyBucketName = "iam_policies"
const policyVersionBucketName = "iam_policy_versions"

// PolicyStore manages IAM policies.
type PolicyStore struct {
	*common.BaseStore
	versionStore *common.BaseStore
	arnBuilder   *ARNBuilder
	kl           common.KeyLocker
}

// NewPolicyStore creates a new PolicyStore.
func NewPolicyStore(store storage.BasicStorage, accountId string) *PolicyStore {
	return &PolicyStore{
		BaseStore:    common.NewBaseStore(store.Bucket(policyBucketName), "iam"),
		versionStore: common.NewBaseStore(store.Bucket(policyVersionBucketName), "iam"),
		arnBuilder:   NewARNBuilder(accountId),
	}
}

// Get retrieves a policy by its ARN.
func (s *PolicyStore) Get(policyArn string) (*Policy, error) {
	var policy Policy
	if err := s.BaseStore.Get(policyArn, &policy); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_policy", ErrPolicyNotFound)
		}
		return nil, NewStoreError("get_policy", err)
	}
	return &policy, nil
}

// GetByPathAndName retrieves a policy by its path and name.
func (s *PolicyStore) GetByPathAndName(path, policyName string) (*Policy, error) {
	found, err := common.FindFirst[Policy](s.BaseStore, func(p *Policy) bool {
		return p.Path == path && p.PolicyName == policyName
	})
	if err != nil {
		return nil, NewStoreError("get_policy_by_name", err)
	}
	return found, nil
}

// List retrieves policies with optional filtering by scope, path prefix, and attachment status.
func (s *PolicyStore) List(scope, pathPrefix string, onlyAttached bool, marker string, maxItems int) (*PolicyListResult, error) {
	var filter common.FilterFunc[Policy]
	if scope == "AWS" {
		filter = func(p *Policy) bool {
			if p.AccountId != "aws" {
				return false
			}
			if pathPrefix != "" && !strings.HasPrefix(p.Path, pathPrefix) {
				return false
			}
			if onlyAttached && p.AttachmentCount == 0 {
				return false
			}
			return true
		}
	} else {
		filter = func(p *Policy) bool {
			if pathPrefix != "" && !strings.HasPrefix(p.Path, pathPrefix) {
				return false
			}
			if onlyAttached && p.AttachmentCount == 0 {
				return false
			}
			return true
		}
	}

	result, err := common.List[Policy](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, filter)
	if err != nil {
		return nil, err
	}
	return &PolicyListResult{
		Policies:    result.Items,
		IsTruncated: result.IsTruncated,
		Marker:      result.NextMarker,
	}, nil
}

// Put stores a policy.
func (s *PolicyStore) Put(policy *Policy) error {
	if policy.CreateDate.IsZero() {
		policy.CreateDate = time.Now().UTC()
	}
	policy.UpdateDate = time.Now().UTC()
	if policy.Path == "" {
		policy.Path = "/"
	}
	if policy.DefaultVersionId == "" {
		policy.DefaultVersionId = "v1"
	}
	policy.IsAttachable = true

	return s.BaseStore.Put(policy.Arn, policy)
}

// Delete removes a policy by its ARN.
func (s *PolicyStore) Delete(policyArn string) error {
	return s.BaseStore.Delete(policyArn)
}

// Exists checks whether a policy exists.
func (s *PolicyStore) Exists(policyArn string) bool {
	return s.BaseStore.Exists(policyArn)
}

// Create creates a new policy.
func (s *PolicyStore) Create(policyName, path, accountId, document, description string, tags []types.Tag) (*Policy, error) {
	arn := s.arnBuilder.PolicyARN(path, policyName)
	if s.Exists(arn) {
		return nil, NewStoreError("create_policy", ErrPolicyAlreadyExists)
	}

	policyID, err := GeneratePolicyID()
	if err != nil {
		return nil, NewStoreError("generate_policy_id", err)
	}

	policy := &Policy{
		ID:               policyID,
		Path:             path,
		PolicyName:       policyName,
		Arn:              arn,
		AccountId:        accountId,
		CreateDate:       time.Now().UTC(),
		DefaultVersionId: "v1",
		AttachmentCount:  0,
		IsAttachable:     true,
		Description:      description,
		Tags:             tags,
	}

	if err := s.Put(policy); err != nil {
		return nil, err
	}

	version := &PolicyVersion{
		VersionId:        "v1",
		PolicyArn:        arn,
		IsDefaultVersion: true,
		CreateDate:       time.Now().UTC(),
		Document:         document,
	}
	if err := s.PutVersion(version); err != nil {
		return nil, err
	}

	return policy, nil
}

// IncrementAttachmentCount increments the attachment count for a policy.
func (s *PolicyStore) IncrementAttachmentCount(policyArn string) error {
	return s.kl.WithLock(policyArn, func() error {
		policy, err := s.Get(policyArn)
		if err != nil {
			return err
		}
		policy.AttachmentCount++
		return s.Put(policy)
	})
}

// DecrementAttachmentCount decrements the attachment count for a policy.
func (s *PolicyStore) DecrementAttachmentCount(policyArn string) error {
	return s.kl.WithLock(policyArn, func() error {
		policy, err := s.Get(policyArn)
		if err != nil {
			return err
		}
		if policy.AttachmentCount > 0 {
			policy.AttachmentCount--
		}
		return s.Put(policy)
	})
}

// Count returns the total number of policies.
func (s *PolicyStore) Count() int {
	return s.BaseStore.Count()
}

// PutVersion stores a policy version.
func (s *PolicyStore) PutVersion(version *PolicyVersion) error {
	if version.CreateDate.IsZero() {
		version.CreateDate = time.Now().UTC()
	}
	key := version.PolicyArn + ":" + version.VersionId
	return s.versionStore.Put(key, version)
}

// GetVersion retrieves a specific version of a policy.
func (s *PolicyStore) GetVersion(policyArn, versionId string) (*PolicyVersion, error) {
	key := policyArn + ":" + versionId
	var version PolicyVersion
	if err := s.versionStore.Get(key, &version); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_policy_version", ErrPolicyNotFound)
		}
		return nil, NewStoreError("get_policy_version", err)
	}
	return &version, nil
}

// DeleteVersion removes a specific version of a policy.
func (s *PolicyStore) DeleteVersion(policyArn, versionId string) error {
	key := policyArn + ":" + versionId
	return s.versionStore.Delete(key)
}

// ListVersions retrieves all versions of a policy.
func (s *PolicyStore) ListVersions(policyArn string, marker string, maxItems int) (*PolicyVersionListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var versions []*PolicyVersion
	count := 0
	started := marker == ""
	hasMore := false
	prefix := policyArn + ":"

	err := s.versionStore.ForEach(func(k string, v []byte) error {
		if !strings.HasPrefix(k, prefix) {
			return nil
		}

		var version PolicyVersion
		if err := json.Unmarshal(v, &version); err != nil {
			return err
		}

		if !started {
			if version.VersionId == marker {
				started = true
			}
			return nil
		}

		if count < maxItems {
			versions = append(versions, &version)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_policy_versions", err)
	}

	result := &PolicyVersionListResult{
		Versions:    versions,
		IsTruncated: hasMore,
	}
	if len(versions) > 0 {
		result.Marker = versions[len(versions)-1].VersionId
	}

	return result, nil
}

// GetDefaultVersion retrieves the default version of a policy.
func (s *PolicyStore) GetDefaultVersion(policyArn string) (*PolicyVersion, error) {
	policy, err := s.Get(policyArn)
	if err != nil {
		return nil, err
	}
	return s.GetVersion(policyArn, policy.DefaultVersionId)
}

// SetDefaultVersion sets the default version for a policy.
func (s *PolicyStore) SetDefaultVersion(policyArn, versionId string) error {
	return s.kl.WithLock(policyArn, func() error {
		policy, err := s.Get(policyArn)
		if err != nil {
			return err
		}

		oldDefault, err := s.GetVersion(policyArn, policy.DefaultVersionId)
		if err == nil {
			oldDefault.IsDefaultVersion = false
			if err := s.PutVersion(oldDefault); err != nil {
				return err
			}
		}

		newDefault, err := s.GetVersion(policyArn, versionId)
		if err != nil {
			return err
		}
		newDefault.IsDefaultVersion = true
		if err := s.PutVersion(newDefault); err != nil {
			return err
		}

		policy.DefaultVersionId = versionId
		return s.Put(policy)
	})
}

func extractVersionNumber(versionId string) int {
	if len(versionId) < 2 {
		return 0
	}
	if versionId[0] != 'v' {
		return 0
	}
	num, err := strconv.Atoi(versionId[1:])
	if err != nil {
		return 0
	}
	return num
}

// GetMaxVersion returns the maximum version number for a policy.
func (s *PolicyStore) GetMaxVersion(policyArn string) (int, error) {
	prefix := policyArn + ":"
	maxVersion := 0
	err := s.versionStore.ForEach(func(k string, v []byte) error {
		if strings.HasPrefix(k, prefix) {
			versionId := strings.TrimPrefix(k, prefix)
			vnum := extractVersionNumber(versionId)
			if vnum > maxVersion {
				maxVersion = vnum
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return maxVersion, nil
}

// CountVersions returns the number of versions for a policy.
func (s *PolicyStore) CountVersions(policyArn string) (int, error) {
	prefix := policyArn + ":"
	count := 0
	err := s.versionStore.ForEach(func(k string, v []byte) error {
		if strings.HasPrefix(k, prefix) {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DeleteAllVersions removes all versions of a policy.
func (s *PolicyStore) DeleteAllVersions(policyArn string) error {
	prefix := policyArn + ":"
	var keysToDelete []string

	if err := s.versionStore.ForEach(func(k string, v []byte) error {
		if strings.HasPrefix(k, prefix) {
			keysToDelete = append(keysToDelete, k)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, key := range keysToDelete {
		if err := s.versionStore.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

// AWSManagedPolicy represents an AWS managed policy definition.
type AWSManagedPolicy struct {
	ARN         string
	Path        string
	PolicyName  string
	Description string
	Document    string
}

// CreateAWSManagedPolicy creates an AWS managed policy (arn:aws:iam::aws:policy/...).
func (s *PolicyStore) CreateAWSManagedPolicy(def AWSManagedPolicy) error {
	if s.Exists(def.ARN) {
		return nil
	}

	now := time.Now().UTC()
	policy := &Policy{
		ID:               GeneratePolicyIDFromARN(def.ARN),
		Path:             def.Path,
		PolicyName:       def.PolicyName,
		Arn:              def.ARN,
		AccountId:        "aws",
		CreateDate:       now,
		UpdateDate:       now,
		DefaultVersionId: "v1",
		AttachmentCount:  0,
		IsAttachable:     true,
		Description:      def.Description,
	}

	if err := s.BaseStore.Put(def.ARN, policy); err != nil {
		return NewStoreError("create_aws_managed_policy", err)
	}

	version := &PolicyVersion{
		VersionId:        "v1",
		PolicyArn:        def.ARN,
		IsDefaultVersion: true,
		CreateDate:       now,
		Document:         def.Document,
	}
	return s.PutVersion(version)
}

// GeneratePolicyIDFromARN derives a stable policy ID from a policy ARN.
func GeneratePolicyIDFromARN(arn string) string {
	parts := strings.Split(arn, "/")
	if len(parts) > 0 {
		return "ANPA" + parts[len(parts)-1]
	}
	return "ANPA" + arn
}
