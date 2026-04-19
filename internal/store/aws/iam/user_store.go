// Package iam provides AWS IAM store functionality for vorpalstacks.
package iam

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/utils/aws/types"
)

const userBucketName = "iam_users"

// UserStore manages IAM users.
type UserStore struct {
	entityStore[User]
	arnBuilder *ARNBuilder
}

// NewUserStore creates a new UserStore.
func NewUserStore(store storage.BasicStorage, accountId string) *UserStore {
	return &UserStore{
		entityStore: newEntityStore[User](store, userBucketName),
		arnBuilder:  NewARNBuilder(accountId),
	}
}

// GetByID retrieves a user by ID.
func (s *UserStore) GetByID(userID string) (*User, error) {
	return getEntityByID(s.BaseStore, userID, func(u *User) string { return u.ID }, "get_user_by_id", ErrUserNotFound)
}

// List returns a list of users with optional filtering by path prefix.
func (s *UserStore) List(pathPrefix string, marker string, maxItems int) (*UserListResult, error) {
	items, truncated, nextMarker, err := listEntitiesWithPathPrefix(s.BaseStore, pathPrefix, marker, maxItems, func(u *User) string { return u.Path })
	if err != nil {
		return nil, err
	}
	return &UserListResult{
		Users:       items,
		IsTruncated: truncated,
		Marker:      nextMarker,
	}, nil
}

// Put stores a user.
func (s *UserStore) Put(user *User) error {
	if user.CreateDate.IsZero() {
		user.CreateDate = time.Now().UTC()
	}
	if user.Path == "" {
		user.Path = "/"
	}
	return s.BaseStore.Put(user.UserName, user)
}

// Create creates a new user.
func (s *UserStore) Create(userName, path, accountId string, tags []types.Tag) (*User, error) {
	if s.Exists(userName) {
		return nil, NewStoreError("create_user", ErrUserAlreadyExists)
	}

	userID, err := GenerateUserID()
	if err != nil {
		return nil, NewStoreError("generate_user_id", err)
	}

	user := &User{
		ID:         userID,
		Path:       path,
		UserName:   userName,
		AccountId:  accountId,
		CreateDate: time.Now().UTC(),
		Tags:       tags,
	}
	user.Arn = s.arnBuilder.UserARN(path, userName)

	if err := s.Put(user); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdatePasswordLastUsed updates the password last used timestamp.
func (s *UserStore) UpdatePasswordLastUsed(userName string) error {
	return s.kl.WithLock(userName, func() error {
		user, err := s.Get(userName)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		user.PasswordLastUsed = &now

		return s.Put(user)
	})
}
