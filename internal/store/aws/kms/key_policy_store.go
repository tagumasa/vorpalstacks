package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
)

// DefaultKeyPolicyName is the name of the default key policy.
const DefaultKeyPolicyName = "default"

func keyPolicyBucketName(region string) string {
	return "kms_key_policies-" + region
}

// DefaultKeyPolicy is the default key policy document.
var DefaultKeyPolicy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "Enable IAM User Permissions",
			"Effect": "Allow",
			"Principal": {
				"AWS": "*"
			},
			"Action": "kms:*",
			"Resource": "*"
		}
	]
}`

// KeyPolicyStore manages KMS key policies.
type KeyPolicyStore struct {
	bucket storage.Bucket
}

// NewKeyPolicyStore creates a new KeyPolicyStore.
func NewKeyPolicyStore(store storage.BasicStorage, region string) *KeyPolicyStore {
	return &KeyPolicyStore{
		bucket: store.Bucket(keyPolicyBucketName(region)),
	}
}

// Create creates a new key policy.
func (s *KeyPolicyStore) Create(keyID, policyName, policyDocument string) (*KeyPolicy, error) {
	if policyName == "" {
		policyName = DefaultKeyPolicyName
	}

	key := keyPolicyKey(keyID, policyName)

	policy := &KeyPolicy{
		KeyID:          keyID,
		PolicyName:     policyName,
		PolicyDocument: policyDocument,
		CreateDate:     time.Now(),
	}

	if err := s.save(key, policy); err != nil {
		return nil, err
	}

	return policy, nil
}

// Get retrieves a key policy.
func (s *KeyPolicyStore) Get(keyID, policyName string) (*KeyPolicy, error) {
	if policyName == "" {
		policyName = DefaultKeyPolicyName
	}

	key := keyPolicyKey(keyID, policyName)
	data, err := s.bucket.Get([]byte(key))
	if err != nil {
		return nil, NewStoreError("get_key_policy", err)
	}
	if data == nil {
		return nil, NewStoreError("get_key_policy", ErrKeyPolicyNotFound)
	}

	var policy KeyPolicy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, NewStoreError("get_key_policy", err)
	}
	return &policy, nil
}

func (s *KeyPolicyStore) save(key string, policy *KeyPolicy) error {
	data, err := json.Marshal(policy)
	if err != nil {
		return NewStoreError("save_key_policy", err)
	}
	return s.bucket.Put([]byte(key), data)
}

// Put stores a key policy (creates or updates).
func (s *KeyPolicyStore) Put(keyID, policyName, policyDocument string) error {
	if policyName == "" {
		policyName = DefaultKeyPolicyName
	}

	key := keyPolicyKey(keyID, policyName)

	existing, err := s.Get(keyID, policyName)
	if err != nil {
		_, err = s.Create(keyID, policyName, policyDocument)
		return err
	}

	existing.PolicyDocument = policyDocument
	return s.save(key, existing)
}

// Delete removes a key policy.
func (s *KeyPolicyStore) Delete(keyID, policyName string) error {
	if policyName == "" {
		policyName = DefaultKeyPolicyName
	}

	key := keyPolicyKey(keyID, policyName)
	if !s.bucket.Has([]byte(key)) {
		return NewStoreError("delete_key_policy", ErrKeyPolicyNotFound)
	}
	return s.bucket.Delete([]byte(key))
}

// List returns the policy names for a key.
func (s *KeyPolicyStore) List(keyID string) ([]string, error) {
	var policyNames []string
	prefix := keyID + "/"

	err := s.bucket.ForEach(func(k, v []byte) error {
		keyStr := string(k)
		if len(keyStr) > len(prefix) && keyStr[:len(prefix)] == prefix {
			policyName := keyStr[len(prefix):]
			policyNames = append(policyNames, policyName)
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_key_policies", err)
	}

	if len(policyNames) == 0 {
		policyNames = []string{DefaultKeyPolicyName}
	}

	return policyNames, nil
}

// DeleteAllForKey removes all policies for a key.
func (s *KeyPolicyStore) DeleteAllForKey(keyID string) error {
	policies, err := s.List(keyID)
	if err != nil {
		return err
	}

	for _, name := range policies {
		if err := s.Delete(keyID, name); err != nil {
			return err
		}
	}

	return nil
}

// GetDefault retrieves the default policy for a key.
func (s *KeyPolicyStore) GetDefault(keyID string) (*KeyPolicy, error) {
	return s.Get(keyID, DefaultKeyPolicyName)
}

// PutDefault sets the default policy for a key.
func (s *KeyPolicyStore) PutDefault(keyID, policyDocument string) error {
	return s.Put(keyID, DefaultKeyPolicyName, policyDocument)
}

func keyPolicyKey(keyID, policyName string) string {
	return keyID + "/" + policyName
}
