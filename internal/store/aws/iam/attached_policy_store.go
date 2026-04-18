package iam

// Package iam provides IAM (Identity and Access Management) data store implementations
// for vorpalstacks.

import (
	"encoding/json"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const attachedPolicyBucketName = "iam_attached_policies"

// AttachedPolicyStore manages IAM policy attachments in persistent storage.
type AttachedPolicyStore struct {
	bucket storage.Bucket
	kl     common.KeyLocker
}

// NewAttachedPolicyStore creates a new store for IAM policy attachments.
func NewAttachedPolicyStore(store storage.BasicStorage) *AttachedPolicyStore {
	return &AttachedPolicyStore{
		bucket: store.Bucket(attachedPolicyBucketName),
	}
}

func (s *AttachedPolicyStore) attachmentKey(principalType, principalName, policyArn string) []byte {
	return []byte(principalType + ":" + principalName + ":" + policyArn)
}

// Attach attaches a policy to a principal (user, group, or role).
func (s *AttachedPolicyStore) Attach(principalType, principalName, policyArn string) error {
	lockKey := principalType + ":" + principalName + ":" + policyArn
	return s.kl.WithLock(lockKey, func() error {
		ref := &AttachedPolicyRef{
			PrincipalType: principalType,
			PrincipalName: principalName,
			PolicyArn:     policyArn,
		}

		data, err := json.Marshal(ref)
		if err != nil {
			return NewStoreError("attach_policy", err)
		}

		if err := s.bucket.Put(s.attachmentKey(principalType, principalName, policyArn), data); err != nil {
			return NewStoreError("attach_policy", err)
		}
		return nil
	})
}

// Detach removes a policy attachment from a principal.
func (s *AttachedPolicyStore) Detach(principalType, principalName, policyArn string) error {
	lockKey := principalType + ":" + principalName + ":" + policyArn
	return s.kl.WithLock(lockKey, func() error {
		if err := s.bucket.Delete(s.attachmentKey(principalType, principalName, policyArn)); err != nil {
			return NewStoreError("detach_policy", err)
		}
		return nil
	})
}

// IsAttached checks whether a policy is attached to a principal.
func (s *AttachedPolicyStore) IsAttached(principalType, principalName, policyArn string) bool {
	return s.bucket.Has(s.attachmentKey(principalType, principalName, policyArn))
}

// ListAttachedPolicies returns all policies attached to a principal.
func (s *AttachedPolicyStore) ListAttachedPolicies(principalType, principalName string) ([]string, error) {
	var policyArns []string
	prefix := principalType + ":" + principalName + ":"

	err := s.bucket.ForEach(func(k, v []byte) error {
		key := string(k)
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			policyArns = append(policyArns, key[len(prefix):])
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_attached_policies", err)
	}
	return policyArns, nil
}

// ListPrincipalsForPolicy returns all principals that have a policy attached.
func (s *AttachedPolicyStore) ListPrincipalsForPolicy(policyArn string) ([]AttachedPolicyRef, error) {
	var refs []AttachedPolicyRef
	suffix := ":" + policyArn

	err := s.bucket.ForEach(func(k, v []byte) error {
		key := string(k)
		if len(key) > len(suffix) && key[len(key)-len(suffix):] == suffix {
			var ref AttachedPolicyRef
			if err := json.Unmarshal(v, &ref); err != nil {
				return err
			}
			refs = append(refs, ref)
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_principals_for_policy", err)
	}
	return refs, nil
}

// DetachAllForPrincipal removes all policy attachments for a principal.
func (s *AttachedPolicyStore) DetachAllForPrincipal(principalType, principalName string) error {
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
			return NewStoreError("detach_all_policies", err)
		}

		for _, key := range keysToDelete {
			if err := s.bucket.Delete(key); err != nil {
				return NewStoreError("detach_all_policies", err)
			}
		}
		return nil
	})
}

// DetachAllForPolicy removes a policy from all principals.
func (s *AttachedPolicyStore) DetachAllForPolicy(policyArn string) error {
	lockKey := "removeall:policy:" + policyArn
	return s.kl.WithLock(lockKey, func() error {
		suffix := ":" + policyArn
		var keysToDelete [][]byte

		err := s.bucket.ForEach(func(k, v []byte) error {
			key := string(k)
			if len(key) > len(suffix) && key[len(key)-len(suffix):] == suffix {
				keysToDelete = append(keysToDelete, k)
			}
			return nil
		})

		if err != nil {
			return NewStoreError("detach_all_principals", err)
		}

		for _, key := range keysToDelete {
			if err := s.bucket.Delete(key); err != nil {
				return NewStoreError("detach_all_principals", err)
			}
		}
		return nil
	})
}

// CountAttachedPolicies returns the number of policies attached to a principal.
func (s *AttachedPolicyStore) CountAttachedPolicies(principalType, principalName string) int {
	policyArns, _ := s.ListAttachedPolicies(principalType, principalName)
	return len(policyArns)
}

// MigratePrincipal moves all policy attachments from one principal name to another.
func (s *AttachedPolicyStore) MigratePrincipal(oldName, newName, principalType string) error {
	lockKey := "migrate:" + principalType + ":" + oldName
	return s.kl.WithLock(lockKey, func() error {
		prefix := principalType + ":" + oldName + ":"
		var refs []AttachedPolicyRef

		err := s.bucket.ForEach(func(k, v []byte) error {
			key := string(k)
			if len(key) > len(prefix) && key[:len(prefix)] == prefix {
				var ref AttachedPolicyRef
				if err := json.Unmarshal(v, &ref); err != nil {
					return err
				}
				refs = append(refs, ref)
			}
			return nil
		})

		if err != nil {
			return NewStoreError("migrate_principal", err)
		}

		for _, ref := range refs {
			if err := s.bucket.Delete(s.attachmentKey(principalType, oldName, ref.PolicyArn)); err != nil {
				return NewStoreError("migrate_principal", err)
			}
			ref.PrincipalName = newName
			data, err := json.Marshal(ref)
			if err != nil {
				return NewStoreError("migrate_principal", err)
			}
			if err := s.bucket.Put(s.attachmentKey(principalType, newName, ref.PolicyArn), data); err != nil {
				return NewStoreError("migrate_principal", err)
			}
		}

		return nil
	})
}
