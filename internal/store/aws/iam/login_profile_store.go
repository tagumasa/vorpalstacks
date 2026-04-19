package iam

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const loginProfileBucketName = "iam_login_profiles"

const bcryptCost = 12

// LoginProfileStore manages IAM login profile data in persistent storage.
type LoginProfileStore struct {
	*common.BaseStore
	kl common.KeyLocker
}

// NewLoginProfileStore creates a new store for IAM login profiles.
func NewLoginProfileStore(store storage.BasicStorage) *LoginProfileStore {
	return &LoginProfileStore{
		BaseStore: common.NewBaseStore(store.Bucket(loginProfileBucketName), "iam"),
	}
}

// Get retrieves a login profile by username.
func (s *LoginProfileStore) Get(userName string) (*LoginProfile, error) {
	var profile LoginProfile
	if err := s.BaseStore.Get(userName, &profile); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_login_profile", ErrLoginProfileNotFound)
		}
		return nil, NewStoreError("get_login_profile", err)
	}
	return &profile, nil
}

// Put stores a login profile.
func (s *LoginProfileStore) Put(profile *LoginProfile) error {
	if profile.CreateDate.IsZero() {
		profile.CreateDate = time.Now().UTC()
	}
	return s.BaseStore.Put(profile.UserName, profile)
}

// Delete removes a login profile by username.
func (s *LoginProfileStore) Delete(userName string) error {
	return s.BaseStore.Delete(userName)
}

// Exists checks whether a login profile exists for a user.
func (s *LoginProfileStore) Exists(userName string) bool {
	return s.BaseStore.Exists(userName)
}

// Create creates a new login profile for a user.
func (s *LoginProfileStore) Create(userName, password string, passwordResetRequired bool) (*LoginProfile, error) {
	if s.Exists(userName) {
		return nil, NewStoreError("create_login_profile", ErrLoginProfileExists)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, NewStoreError("create_login_profile", err)
	}

	profile := &LoginProfile{
		UserName:              userName,
		PasswordHash:          string(passwordHash),
		PasswordResetRequired: passwordResetRequired,
		CreateDate:            time.Now().UTC(),
	}

	if err := s.Put(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// UpdatePassword changes the password for a login profile.
func (s *LoginProfileStore) UpdatePassword(userName, password string) error {
	return s.kl.WithLock(userName, func() error {
		profile, err := s.Get(userName)
		if err != nil {
			return err
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
		if err != nil {
			return NewStoreError("update_password", err)
		}

		profile.PasswordHash = string(passwordHash)
		profile.PasswordResetRequired = false

		return s.Put(profile)
	})
}

// UpdatePasswordResetRequired changes the password reset requirement flag.
func (s *LoginProfileStore) UpdatePasswordResetRequired(userName string, required bool) error {
	return s.kl.WithLock(userName, func() error {
		profile, err := s.Get(userName)
		if err != nil {
			return err
		}

		profile.PasswordResetRequired = required
		return s.Put(profile)
	})
}

// VerifyPassword validates a password against the stored hash.
func (s *LoginProfileStore) VerifyPassword(userName, password string) (bool, error) {
	profile, err := s.Get(userName)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(profile.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, NewStoreError("verify_password", err)
	}
	return true, nil
}

// Count returns the total number of login profiles.
func (s *LoginProfileStore) Count() int {
	return s.BaseStore.Count()
}
