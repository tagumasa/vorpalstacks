package iam

// Package iam provides IAM (Identity and Access Management) data store implementations
// for vorpalstacks.

import (
	"encoding/json"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const inlinePolicyBucketName = "iam_inline_policies"

// InlinePolicyStore manages IAM inline policy data in persistent storage.
type InlinePolicyStore struct {
	bucket storage.Bucket
	kl     common.KeyLocker
}

// NewInlinePolicyStore creates a new store for IAM inline policies.
func NewInlinePolicyStore(store storage.BasicStorage) *InlinePolicyStore {
	return &InlinePolicyStore{
		bucket: store.Bucket(inlinePolicyBucketName),
	}
}

func (s *InlinePolicyStore) policyKey(principalType, principalName, policyName string) []byte {
	return []byte(principalType + ":" + principalName + ":" + policyName)
}

// Put stores an inline policy for a principal.
func (s *InlinePolicyStore) Put(principalType, principalName, policyName, document string) error {
	lockKey := principalType + ":" + principalName + ":" + policyName
	return s.kl.WithLock(lockKey, func() error {
		policy := &InlinePolicy{
			PrincipalType:  principalType,
			PrincipalName:  principalName,
			PolicyName:     policyName,
			PolicyDocument: document,
		}

		data, err := json.Marshal(policy)
		if err != nil {
			return NewStoreError("put_inline_policy", err)
		}

		if err := s.bucket.Put(s.policyKey(principalType, principalName, policyName), data); err != nil {
			return NewStoreError("put_inline_policy", err)
		}
		return nil
	})
}

// Get retrieves an inline policy for a principal.
func (s *InlinePolicyStore) Get(principalType, principalName, policyName string) (*InlinePolicy, error) {
	data, err := s.bucket.Get(s.policyKey(principalType, principalName, policyName))
	if err != nil {
		return nil, NewStoreError("get_inline_policy", err)
	}
	if data == nil {
		return nil, NewStoreError("get_inline_policy", ErrPolicyNotFound)
	}
	var policy InlinePolicy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, NewStoreError("get_inline_policy", err)
	}
	return &policy, nil
}

// Delete removes an inline policy from a principal.
func (s *InlinePolicyStore) Delete(principalType, principalName, policyName string) error {
	lockKey := principalType + ":" + principalName + ":" + policyName
	return s.kl.WithLock(lockKey, func() error {
		if err := s.bucket.Delete(s.policyKey(principalType, principalName, policyName)); err != nil {
			return NewStoreError("delete_inline_policy", err)
		}
		return nil
	})
}

// Exists checks whether an inline policy exists for a principal.
func (s *InlinePolicyStore) Exists(principalType, principalName, policyName string) bool {
	return s.bucket.Has(s.policyKey(principalType, principalName, policyName))
}

// List returns all inline policy names for a principal.
func (s *InlinePolicyStore) List(principalType, principalName string) ([]string, error) {
	var policyNames []string
	prefix := principalType + ":" + principalName + ":"

	err := s.bucket.ForEach(func(k, v []byte) error {
		key := string(k)
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			policyNames = append(policyNames, key[len(prefix):])
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_inline_policies", err)
	}
	return policyNames, nil
}

// DeleteAllForPrincipal removes all inline policies for a principal.
func (s *InlinePolicyStore) DeleteAllForPrincipal(principalType, principalName string) error {
	lockKey := "removeall:" + principalType + ":" + principalName
	return s.kl.WithLock(lockKey, func() error {
		prefix := principalType + ":" + principalName + ":"
		var keysToDelete [][]byte

		err := s.bucket.ForEach(func(k, v []byte) error {
			key := string(k)
			if len(key) > len(prefix) && key[:len(prefix)] == prefix {
				keysToDelete = append(keysToDelete, k)
			}
			return nil
		})

		if err != nil {
			return NewStoreError("delete_all_inline_policies", err)
		}

		for _, key := range keysToDelete {
			if err := s.bucket.Delete(key); err != nil {
				return NewStoreError("delete_all_inline_policies", err)
			}
		}
		return nil
	})
}

// Count returns the number of inline policies for a principal.
func (s *InlinePolicyStore) Count(principalType, principalName string) int {
	policyNames, _ := s.List(principalType, principalName)
	return len(policyNames)
}

// MigratePrincipal moves all inline policies from one principal name to another.
func (s *InlinePolicyStore) MigratePrincipal(oldName, newName, principalType string) error {
	lockKey := "migrate:" + principalType + ":" + oldName
	return s.kl.WithLock(lockKey, func() error {
		prefix := principalType + ":" + oldName + ":"
		var policies []InlinePolicy

		err := s.bucket.ForEach(func(k, v []byte) error {
			key := string(k)
			if len(key) > len(prefix) && key[:len(prefix)] == prefix {
				var policy InlinePolicy
				if err := json.Unmarshal(v, &policy); err != nil {
					return err
				}
				policies = append(policies, policy)
			}
			return nil
		})

		if err != nil {
			return NewStoreError("migrate_principal", err)
		}

		for _, policy := range policies {
			if err := s.bucket.Delete(s.policyKey(principalType, oldName, policy.PolicyName)); err != nil {
				return NewStoreError("migrate_principal", err)
			}
			policy.PrincipalName = newName
			data, err := json.Marshal(policy)
			if err != nil {
				return NewStoreError("migrate_principal", err)
			}
			if err := s.bucket.Put(s.policyKey(principalType, newName, policy.PolicyName), data); err != nil {
				return NewStoreError("migrate_principal", err)
			}
		}

		return nil
	})
}
