package iam

import (
	"encoding/json"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const attachedPolicyBucketName = "iam_attached_policies"

// AttachedPolicyStore manages IAM policy attachments in persistent storage.
type AttachedPolicyStore struct {
	pk principalKeyed[AttachedPolicyRef]
}

// NewAttachedPolicyStore creates a new store for IAM policy attachments.
func NewAttachedPolicyStore(store storage.BasicStorage) *AttachedPolicyStore {
	bs := common.NewBaseStore(store.Bucket(attachedPolicyBucketName), "iam")
	return &AttachedPolicyStore{
		pk: newPrincipalKeyed[AttachedPolicyRef](bs,
			func(r *AttachedPolicyRef) string { return r.PolicyArn },
			func(r *AttachedPolicyRef, newName string) { r.PrincipalName = newName },
		),
	}
}

func (s *AttachedPolicyStore) attachmentKey(principalType, principalName, policyArn string) string {
	return s.pk.buildPrefix(principalType, principalName) + policyArn
}

// Attach attaches a policy to a principal (user, group, or role).
func (s *AttachedPolicyStore) Attach(principalType, principalName, policyArn string) error {
	lockKey := s.attachmentKey(principalType, principalName, policyArn)
	return s.pk.kl.WithLock(lockKey, func() error {
		ref := &AttachedPolicyRef{
			PrincipalType: principalType,
			PrincipalName: principalName,
			PolicyArn:     policyArn,
		}
		return s.pk.BaseStore.Put(s.attachmentKey(principalType, principalName, policyArn), ref)
	})
}

// Detach removes a policy attachment from a principal.
func (s *AttachedPolicyStore) Detach(principalType, principalName, policyArn string) error {
	lockKey := s.attachmentKey(principalType, principalName, policyArn)
	return s.pk.kl.WithLock(lockKey, func() error {
		return s.pk.BaseStore.Delete(s.attachmentKey(principalType, principalName, policyArn))
	})
}

// IsAttached checks whether a policy is attached to a principal.
func (s *AttachedPolicyStore) IsAttached(principalType, principalName, policyArn string) bool {
	return s.pk.BaseStore.Exists(s.attachmentKey(principalType, principalName, policyArn))
}

// ListAttachedPolicies returns all policies attached to a principal.
func (s *AttachedPolicyStore) ListAttachedPolicies(principalType, principalName string) ([]string, error) {
	arns, err := s.pk.listByPrincipal(principalType, principalName)
	if err != nil {
		return nil, NewStoreError("list_attached_policies", err)
	}
	return arns, nil
}

// ListPrincipalsForPolicy returns all principals that have a policy attached.
func (s *AttachedPolicyStore) ListPrincipalsForPolicy(policyArn string) ([]AttachedPolicyRef, error) {
	var refs []AttachedPolicyRef
	suffix := ":" + policyArn

	err := s.pk.ForEach(func(k string, v []byte) error {
		if len(k) > len(suffix) && k[len(k)-len(suffix):] == suffix {
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
	return s.pk.deleteAllForPrincipal(principalType, principalName, "detach_all_policies")
}

// DetachAllForPolicy removes a policy from all principals.
func (s *AttachedPolicyStore) DetachAllForPolicy(policyArn string) error {
	lockKey := "removeall:policy:" + policyArn
	return s.pk.kl.WithLock(lockKey, func() error {
		suffix := ":" + policyArn
		var keysToDelete []string

		err := s.pk.ForEach(func(k string, v []byte) error {
			if len(k) > len(suffix) && k[len(k)-len(suffix):] == suffix {
				keysToDelete = append(keysToDelete, k)
			}
			return nil
		})

		if err != nil {
			return NewStoreError("detach_all_principals", err)
		}

		for _, key := range keysToDelete {
			if err := s.pk.BaseStore.Delete(key); err != nil {
				return NewStoreError("detach_all_principals", err)
			}
		}
		return nil
	})
}

// CountAttachedPolicies returns the number of policies attached to a principal.
func (s *AttachedPolicyStore) CountAttachedPolicies(principalType, principalName string) int {
	arns, _ := s.ListAttachedPolicies(principalType, principalName)
	return len(arns)
}

// MigratePrincipal moves all policy attachments from one principal name to another.
func (s *AttachedPolicyStore) MigratePrincipal(oldName, newName, principalType string) error {
	return s.pk.migratePrincipal(oldName, newName, principalType, "migrate_principal")
}
