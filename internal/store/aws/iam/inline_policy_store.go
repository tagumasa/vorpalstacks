package iam

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const inlinePolicyBucketName = "iam_inline_policies"

// InlinePolicyStore manages IAM inline policies attached to principals.
type InlinePolicyStore struct {
	pk principalKeyed[InlinePolicy]
}

// NewInlinePolicyStore creates a new store for IAM inline policies.
func NewInlinePolicyStore(store storage.BasicStorage) *InlinePolicyStore {
	bs := common.NewBaseStore(store.Bucket(inlinePolicyBucketName), "iam")
	return &InlinePolicyStore{
		pk: newPrincipalKeyed[InlinePolicy](bs,
			func(p *InlinePolicy) string { return p.PolicyName },
			func(p *InlinePolicy, newName string) { p.PrincipalName = newName },
		),
	}
}

func (s *InlinePolicyStore) policyKey(principalType, principalName, policyName string) string {
	return s.pk.buildPrefix(principalType, principalName) + policyName
}

// Put stores an inline policy for the specified principal.
func (s *InlinePolicyStore) Put(principalType, principalName, policyName, document string) error {
	lockKey := s.policyKey(principalType, principalName, policyName)
	return s.pk.kl.WithLock(lockKey, func() error {
		policy := &InlinePolicy{
			PrincipalType:  principalType,
			PrincipalName:  principalName,
			PolicyName:     policyName,
			PolicyDocument: document,
		}
		return s.pk.BaseStore.Put(s.policyKey(principalType, principalName, policyName), policy)
	})
}

// Get retrieves an inline policy by principal type, name, and policy name.
func (s *InlinePolicyStore) Get(principalType, principalName, policyName string) (*InlinePolicy, error) {
	var policy InlinePolicy
	if err := s.pk.BaseStore.Get(s.policyKey(principalType, principalName, policyName), &policy); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_inline_policy", ErrPolicyNotFound)
		}
		return nil, NewStoreError("get_inline_policy", err)
	}
	return &policy, nil
}

// Delete removes an inline policy.
func (s *InlinePolicyStore) Delete(principalType, principalName, policyName string) error {
	lockKey := s.policyKey(principalType, principalName, policyName)
	return s.pk.kl.WithLock(lockKey, func() error {
		return s.pk.BaseStore.Delete(s.policyKey(principalType, principalName, policyName))
	})
}

// Exists checks whether an inline policy exists.
func (s *InlinePolicyStore) Exists(principalType, principalName, policyName string) bool {
	return s.pk.BaseStore.Exists(s.policyKey(principalType, principalName, policyName))
}

// List returns all inline policy names for a principal.
func (s *InlinePolicyStore) List(principalType, principalName string) ([]string, error) {
	names, err := s.pk.listByPrincipal(principalType, principalName)
	if err != nil {
		return nil, NewStoreError("list_inline_policies", err)
	}
	return names, nil
}

// DeleteAllForPrincipal removes all inline policies for a principal.
func (s *InlinePolicyStore) DeleteAllForPrincipal(principalType, principalName string) error {
	return s.pk.deleteAllForPrincipal(principalType, principalName, "delete_all_inline_policies")
}

// Count returns the number of inline policies for a principal.
func (s *InlinePolicyStore) Count(principalType, principalName string) int {
	policyNames, _ := s.List(principalType, principalName)
	return len(policyNames)
}

// MigratePrincipal re-keys all inline policies when a principal is renamed.
func (s *InlinePolicyStore) MigratePrincipal(oldName, newName, principalType string) error {
	return s.pk.migratePrincipal(oldName, newName, principalType, "migrate_principal")
}
