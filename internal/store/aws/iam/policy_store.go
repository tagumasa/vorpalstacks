package iam

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
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
	return getByKey[Policy](s.BaseStore, policyArn, "get_policy", ErrPolicyNotFound)
}

// List retrieves policies with optional filtering by scope, path prefix, and attachment status.
// Scope follows the ListPolicies contract: "AWS" lists only AWS managed
// policies (AccountId "aws"), "Local" only customer managed policies, and
// "All" (or any other value) everything.
func (s *PolicyStore) List(scope, pathPrefix string, onlyAttached bool, marker string, maxItems int) (*PolicyListResult, error) {
	awsManagedOnly, customerManagedOnly := scope == "AWS", scope == "Local"
	filter := func(p *Policy) bool {
		if awsManagedOnly && p.AccountId != "aws" {
			return false
		}
		if customerManagedOnly && p.AccountId == "aws" {
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

	var policy *Policy
	err := s.kl.WithLock(arn, func() error {
		if s.Exists(arn) {
			return NewStoreError("create_policy", ErrPolicyAlreadyExists)
		}

		policyID, err := GeneratePolicyID()
		if err != nil {
			return NewStoreError("generate_policy_id", err)
		}

		policy = &Policy{
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
			return err
		}

		version := &PolicyVersion{
			VersionId:        "v1",
			PolicyArn:        arn,
			IsDefaultVersion: true,
			CreateDate:       time.Now().UTC(),
			Document:         document,
		}
		return s.PutVersion(version)
	})
	if err != nil {
		return nil, err
	}
	return policy, nil
}

// adjustUsageCount applies a signed adjustment to one of the policy's usage
// counters inside the policy lock scope: read-modify-write under WithLock so
// concurrent attach/detach or boundary changes cannot lose an update. A
// decrement never drops below zero (the counters are usage tallies, not
// reference counts with phantom ownership).
func (s *PolicyStore) adjustUsageCount(policyArn string, counter func(*Policy) *int, delta int) error {
	return s.kl.WithLock(policyArn, func() error {
		policy, err := s.Get(policyArn)
		if err != nil {
			return err
		}
		value := counter(policy)
		if delta > 0 || *value > 0 {
			*value += delta
		}
		return s.Put(policy)
	})
}

// IncrementAttachmentCount increments the attachment count for a policy.
func (s *PolicyStore) IncrementAttachmentCount(policyArn string) error {
	return s.adjustUsageCount(policyArn, func(p *Policy) *int { return &p.AttachmentCount }, 1)
}

// DecrementAttachmentCount decrements the attachment count for a policy.
func (s *PolicyStore) DecrementAttachmentCount(policyArn string) error {
	return s.adjustUsageCount(policyArn, func(p *Policy) *int { return &p.AttachmentCount }, -1)
}

// IncrementPermissionsBoundaryUsageCount increments the permissions boundary
// usage count for a policy. Called when a user or role's permissions boundary
// is set to this policy.
func (s *PolicyStore) IncrementPermissionsBoundaryUsageCount(policyArn string) error {
	return s.adjustUsageCount(policyArn, func(p *Policy) *int { return &p.PermissionsBoundaryUsageCount }, 1)
}

// DecrementPermissionsBoundaryUsageCount decrements the permissions boundary
// usage count for a policy. Called when a user or role's permissions boundary
// is removed or changed away from this policy.
func (s *PolicyStore) DecrementPermissionsBoundaryUsageCount(policyArn string) error {
	return s.adjustUsageCount(policyArn, func(p *Policy) *int { return &p.PermissionsBoundaryUsageCount }, -1)
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
	return getByKey[PolicyVersion](s.versionStore, policyArn+":"+versionId, "get_policy_version", ErrPolicyNotFound)
}

// DeleteVersion removes a specific version of a policy.
func (s *PolicyStore) DeleteVersion(policyArn, versionId string) error {
	key := policyArn + ":" + versionId
	return s.versionStore.Delete(key)
}

// ListVersions retrieves all versions of a policy, newest first. AWS does
// not contractually guarantee an order; its documented example lists the
// newest version first (v3, v2, v1), and the numeric comparison is what
// keeps the sequence sensible: bucket key order is lexicographic, which
// would list v10 before v6 once cumulative numbering passes nine
// versions, breaking any numeric expectation callers build on the
// v-identifiers and the marker pagination that follows them.
func (s *PolicyStore) ListVersions(policyArn string, marker string, maxItems int) (*PolicyVersionListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var all []*PolicyVersion
	prefix := policyArn + ":"

	err := s.versionStore.ForEach(func(k string, v []byte) error {
		if !strings.HasPrefix(k, prefix) {
			return nil
		}

		var version PolicyVersion
		if err := json.Unmarshal(v, &version); err != nil {
			return err
		}
		all = append(all, &version)
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_policy_versions", err)
	}

	slices.SortFunc(all, func(a, b *PolicyVersion) int {
		return cmp.Compare(extractVersionNumber(b.VersionId), extractVersionNumber(a.VersionId))
	})

	var versions []*PolicyVersion
	started := marker == ""
	hasMore := false
	for _, version := range all {
		if !started {
			if version.VersionId == marker {
				started = true
			}
			continue
		}

		if len(versions) < maxItems {
			versions = append(versions, version)
		} else {
			hasMore = true
			break
		}
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
		return s.setDefaultVersionUnlocked(policyArn, versionId)
	})
}

// setDefaultVersionUnlocked swaps the default version marker from the
// current default to versionId.  The caller MUST already hold the
// policyArn lock (used by CreateVersion which operates inside the same
// lock scope to avoid a non-reentrant deadlock).
func (s *PolicyStore) setDefaultVersionUnlocked(policyArn, versionId string) error {
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

// CreateVersion atomically creates a new policy version inside the
// policy lock scope, enforcing the MaxPolicyVersions quota and
// optionally swapping the default marker.  This prevents the race
// condition where concurrent requests could both observe a version
// count below the limit and both succeed in creating a new version.
func (s *PolicyStore) CreateVersion(policyArn, document string, setAsDefault bool, maxVersions int) (*PolicyVersion, error) {
	var version *PolicyVersion
	err := s.kl.WithLock(policyArn, func() error {
		count, err := s.CountVersions(policyArn)
		if err != nil {
			return err
		}
		if count >= maxVersions {
			return NewStoreError("create_policy_version", ErrPolicyVersionLimitExceeded)
		}

		maxVer, err := s.GetMaxVersion(policyArn)
		if err != nil {
			return err
		}
		versionId := fmt.Sprintf("v%d", maxVer+1)

		version = &PolicyVersion{
			VersionId:        versionId,
			PolicyArn:        policyArn,
			IsDefaultVersion: false,
			Document:         document,
		}

		if err := s.PutVersion(version); err != nil {
			return err
		}

		if setAsDefault {
			if err := s.setDefaultVersionUnlocked(policyArn, versionId); err != nil {
				return err
			}
			// Reflect the persisted default marker on the local copy.
			version.IsDefaultVersion = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return version, nil
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

// awsManagedPolicyARNPrefix identifies seeded AWS managed policies. Their
// permissions cannot be changed or deleted, so mutation operations reject
// ARNs carrying this prefix.
const awsManagedPolicyARNPrefix = "arn:aws:iam::aws:policy/"

// IsAWSManagedPolicyARN reports whether the ARN refers to an AWS managed
// policy, which cannot be modified or deleted.
func IsAWSManagedPolicyARN(arn string) bool {
	return strings.HasPrefix(arn, awsManagedPolicyARNPrefix)
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
