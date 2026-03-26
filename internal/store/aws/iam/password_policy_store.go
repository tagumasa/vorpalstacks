package iam

// Package iam provides IAM (Identity and Access Management) data store implementations
// for vorpalstacks.

import (
	"encoding/json"

	"vorpalstacks/internal/core/storage"
)

const passwordPolicyBucketName = "iam_password_policy"

// PasswordPolicyStore manages account password policies.
type PasswordPolicyStore struct {
	bucket storage.Bucket
}

// NewPasswordPolicyStore creates a new PasswordPolicyStore.
func NewPasswordPolicyStore(store storage.BasicStorage) *PasswordPolicyStore {
	return &PasswordPolicyStore{
		bucket: store.Bucket(passwordPolicyBucketName),
	}
}

// Get retrieves the account password policy.
func (s *PasswordPolicyStore) Get() (*AccountPasswordPolicy, error) {
	data, err := s.bucket.Get([]byte("default"))
	if err != nil {
		return nil, NewStoreError("get_password_policy", err)
	}
	if data == nil {
		return nil, NewStoreError("get_password_policy", ErrPasswordPolicyNotFound)
	}
	var policy AccountPasswordPolicy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, NewStoreError("get_password_policy", err)
	}
	return &policy, nil
}

// Put stores the account password policy.
func (s *PasswordPolicyStore) Put(policy *AccountPasswordPolicy) error {
	data, err := json.Marshal(policy)
	if err != nil {
		return NewStoreError("put_password_policy", err)
	}

	if err := s.bucket.Put([]byte("default"), data); err != nil {
		return NewStoreError("put_password_policy", err)
	}
	return nil
}

// Delete removes the account password policy.
func (s *PasswordPolicyStore) Delete() error {
	if err := s.bucket.Delete([]byte("default")); err != nil {
		return NewStoreError("delete_password_policy", err)
	}
	return nil
}

// Exists checks whether a password policy exists.
func (s *PasswordPolicyStore) Exists() bool {
	return s.bucket.Has([]byte("default"))
}

// GetOrDefault returns the password policy or the default if not found.
func (s *PasswordPolicyStore) GetOrDefault() *AccountPasswordPolicy {
	if policy, err := s.Get(); err == nil {
		return policy
	}
	return s.DefaultPolicy()
}

// DefaultPolicy returns the default password policy configuration.
func (s *PasswordPolicyStore) DefaultPolicy() *AccountPasswordPolicy {
	return &AccountPasswordPolicy{
		MinimumPasswordLength:      8,
		RequireSymbols:             true,
		RequireNumbers:             true,
		RequireUppercaseCharacters: true,
		RequireLowercaseCharacters: true,
		AllowUsersToChangePassword: true,
		MaxPasswordAge:             90,
		PasswordReusePrevention:    24,
		HardExpiry:                 false,
		ExpirePasswords:            false,
	}
}
