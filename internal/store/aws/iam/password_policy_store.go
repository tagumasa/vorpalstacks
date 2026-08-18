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
	// The default password policy applied when no custom policy exists:
	// minimum length 8, at least three of the four character types, and
	// passwords never expire (per the IAM User Guide "Default password
	// policy": minimum length of 8, three of four character types, never
	// expire). Requiring all four types individually would be stricter
	// than the documented rule, so the mixed-class requirement is carried
	// by MinimumCharacterTypes instead of the per-class booleans.
	return &AccountPasswordPolicy{
		MinimumPasswordLength:      8,
		MinimumCharacterTypes:      3,
		AllowUsersToChangePassword: false,
		MaxPasswordAge:             0,
		PasswordReusePrevention:    0,
		HardExpiry:                 false,
		ExpirePasswords:            false,
	}
}

// ParameterDefaults returns the policy whose field values are the documented
// per-parameter defaults of UpdateAccountPasswordPolicy: any parameter the
// request omits reverts to these values rather than being merged with the
// previous policy.
func (s *PasswordPolicyStore) ParameterDefaults() *AccountPasswordPolicy {
	return &AccountPasswordPolicy{
		MinimumPasswordLength: 6,
	}
}
