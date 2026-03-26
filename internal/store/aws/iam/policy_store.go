// Package iam provides AWS IAM store functionality for vorpalstacks.
package iam

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
)

const policyBucketName = "iam_policies"
const policyVersionBucketName = "iam_policy_versions"

// PolicyStore manages IAM policies.
type PolicyStore struct {
	bucket        storage.Bucket
	versionBucket storage.Bucket
	arnBuilder    *ARNBuilder
	attachMutex   sync.Mutex
}

// NewPolicyStore creates a new PolicyStore.
func NewPolicyStore(store storage.BasicStorage, accountId string) *PolicyStore {
	return &PolicyStore{
		bucket:        store.Bucket(policyBucketName),
		versionBucket: store.Bucket(policyVersionBucketName),
		arnBuilder:    NewARNBuilder(accountId),
	}
}

// Get retrieves a policy by its ARN.
func (s *PolicyStore) Get(policyArn string) (*Policy, error) {
	data, err := s.bucket.Get([]byte(policyArn))
	if err != nil {
		return nil, NewStoreError("get_policy", err)
	}
	if data == nil {
		return nil, NewStoreError("get_policy", ErrPolicyNotFound)
	}
	var policy Policy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, NewStoreError("get_policy", err)
	}
	return &policy, nil
}

// GetByPathAndName retrieves a policy by its path and name.
func (s *PolicyStore) GetByPathAndName(path, policyName string) (*Policy, error) {
	var found *Policy
	err := s.bucket.ForEach(func(k, v []byte) error {
		var policy Policy
		if err := json.Unmarshal(v, &policy); err != nil {
			return err
		}
		if policy.Path == path && policy.PolicyName == policyName {
			found = &policy
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_policy_by_name", err)
	}
	if found == nil {
		return nil, NewStoreError("get_policy_by_name", ErrPolicyNotFound)
	}
	return found, nil
}

// List retrieves policies with optional filtering by scope, path prefix, and attachment status.
func (s *PolicyStore) List(scope, pathPrefix string, onlyAttached bool, marker string, maxItems int) (*PolicyListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	if scope == "AWS" {
		return s.listWithFilter(pathPrefix, onlyAttached, marker, maxItems, func(p *Policy) bool {
			return p.AccountId == "aws"
		})
	}

	return s.listWithFilter(pathPrefix, onlyAttached, marker, maxItems, nil)
}

func (s *PolicyStore) listWithFilter(pathPrefix string, onlyAttached bool, marker string, maxItems int, extraFilter func(*Policy) bool) (*PolicyListResult, error) {
	var policies []*Policy
	count := 0
	started := marker == ""
	hasMore := false

	err := s.bucket.ForEach(func(k, v []byte) error {
		var policy Policy
		if err := json.Unmarshal(v, &policy); err != nil {
			return err
		}

		if !started {
			if policy.Arn == marker {
				started = true
			}
			return nil
		}

		if pathPrefix != "" && !strings.HasPrefix(policy.Path, pathPrefix) {
			return nil
		}

		if onlyAttached && policy.AttachmentCount == 0 {
			return nil
		}

		if extraFilter != nil && !extraFilter(&policy) {
			return nil
		}

		if count < maxItems {
			policies = append(policies, &policy)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_policies", err)
	}

	result := &PolicyListResult{
		Policies:    policies,
		IsTruncated: hasMore,
	}
	if len(policies) > 0 {
		result.Marker = policies[len(policies)-1].Arn
	}

	return result, nil
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

	data, err := json.Marshal(policy)
	if err != nil {
		return NewStoreError("put_policy", err)
	}

	if err := s.bucket.Put([]byte(policy.Arn), data); err != nil {
		return NewStoreError("put_policy", err)
	}
	return nil
}

// Delete removes a policy by its ARN.
func (s *PolicyStore) Delete(policyArn string) error {
	if err := s.bucket.Delete([]byte(policyArn)); err != nil {
		return NewStoreError("delete_policy", err)
	}
	return nil
}

// Exists checks whether a policy exists.
func (s *PolicyStore) Exists(policyArn string) bool {
	return s.bucket.Has([]byte(policyArn))
}

// Create creates a new policy.
func (s *PolicyStore) Create(policyName, path, accountId, document, description string, tags []Tag) (*Policy, error) {
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
	s.attachMutex.Lock()
	defer s.attachMutex.Unlock()
	policy, err := s.Get(policyArn)
	if err != nil {
		return err
	}
	policy.AttachmentCount++
	return s.Put(policy)
}

// DecrementAttachmentCount decrements the attachment count for a policy.
func (s *PolicyStore) DecrementAttachmentCount(policyArn string) error {
	s.attachMutex.Lock()
	defer s.attachMutex.Unlock()
	policy, err := s.Get(policyArn)
	if err != nil {
		return err
	}
	if policy.AttachmentCount > 0 {
		policy.AttachmentCount--
	}
	return s.Put(policy)
}

// Count returns the total number of policies.
func (s *PolicyStore) Count() int {
	return s.bucket.Count()
}

// PutVersion stores a policy version.
func (s *PolicyStore) PutVersion(version *PolicyVersion) error {
	if version.CreateDate.IsZero() {
		version.CreateDate = time.Now().UTC()
	}

	key := []byte(version.PolicyArn + ":" + version.VersionId)
	data, err := json.Marshal(version)
	if err != nil {
		return NewStoreError("put_policy_version", err)
	}

	if err := s.versionBucket.Put(key, data); err != nil {
		return NewStoreError("put_policy_version", err)
	}
	return nil
}

// GetVersion retrieves a specific version of a policy.
func (s *PolicyStore) GetVersion(policyArn, versionId string) (*PolicyVersion, error) {
	key := []byte(policyArn + ":" + versionId)
	data, err := s.versionBucket.Get(key)
	if err != nil {
		return nil, NewStoreError("get_policy_version", err)
	}
	if data == nil {
		return nil, NewStoreError("get_policy_version", ErrPolicyNotFound)
	}
	var version PolicyVersion
	if err := json.Unmarshal(data, &version); err != nil {
		return nil, NewStoreError("get_policy_version", err)
	}
	return &version, nil
}

// DeleteVersion removes a specific version of a policy.
func (s *PolicyStore) DeleteVersion(policyArn, versionId string) error {
	key := []byte(policyArn + ":" + versionId)
	if err := s.versionBucket.Delete(key); err != nil {
		return NewStoreError("delete_policy_version", err)
	}
	return nil
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

	err := s.versionBucket.ForEach(func(k, v []byte) error {
		key := string(k)
		if !strings.HasPrefix(key, prefix) {
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
	s.attachMutex.Lock()
	defer s.attachMutex.Unlock()

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
	err := s.versionBucket.ForEach(func(k, v []byte) error {
		if strings.HasPrefix(string(k), prefix) {
			versionId := strings.TrimPrefix(string(k), prefix)
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
	err := s.versionBucket.ForEach(func(k, v []byte) error {
		if strings.HasPrefix(string(k), prefix) {
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
	var keysToDelete [][]byte

	if err := s.versionBucket.ForEach(func(k, v []byte) error {
		if strings.HasPrefix(string(k), prefix) {
			keysToDelete = append(keysToDelete, k)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, key := range keysToDelete {
		if err := s.versionBucket.Delete(key); err != nil {
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

	data, err := json.Marshal(policy)
	if err != nil {
		return NewStoreError("create_aws_managed_policy", err)
	}

	if err := s.bucket.Put([]byte(def.ARN), data); err != nil {
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
