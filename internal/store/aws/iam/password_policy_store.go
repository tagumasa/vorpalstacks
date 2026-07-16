package iam

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const passwordPolicyBucketName = "iam_password_policy"

// PasswordPolicyStore manages account password policies.
type PasswordPolicyStore struct {
	*common.BaseStore
}

// NewPasswordPolicyStore creates a new PasswordPolicyStore.
func NewPasswordPolicyStore(store storage.BasicStorage) *PasswordPolicyStore {
	return &PasswordPolicyStore{
		BaseStore: common.NewBaseStore(store.Bucket(passwordPolicyBucketName), "iam"),
	}
}

// Get retrieves the account password policy.
func (s *PasswordPolicyStore) Get() (*AccountPasswordPolicy, error) {
	var policy AccountPasswordPolicy
	if err := s.BaseStore.Get("default", &policy); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_password_policy", ErrPasswordPolicyNotFound)
		}
		return nil, NewStoreError("get_password_policy", err)
	}
	return &policy, nil
}

// Put stores the account password policy.
func (s *PasswordPolicyStore) Put(policy *AccountPasswordPolicy) error {
	return s.BaseStore.Put("default", policy)
}

// Delete removes the account password policy.
func (s *PasswordPolicyStore) Delete() error {
	return s.BaseStore.Delete("default")
}

// Exists checks whether a password policy exists.
func (s *PasswordPolicyStore) Exists() bool {
	return s.BaseStore.Exists("default")
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
