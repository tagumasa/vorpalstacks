// Package iam provides AWS IAM store functionality for vorpalstacks.
package iam

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const userBucketName = "iam_users"

// UserStore manages IAM users.
type UserStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
}

// NewUserStore creates a new UserStore.
func NewUserStore(store storage.BasicStorage, accountId string) *UserStore {
	return &UserStore{
		BaseStore:  common.NewBaseStore(store.Bucket(userBucketName), "iam"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves a user by name.
func (s *UserStore) Get(userName string) (*User, error) {
	var user User
	if err := s.BaseStore.Get(userName, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves a user by ID.
func (s *UserStore) GetByID(userID string) (*User, error) {
	var found *User
	err := s.ForEach(func(key string, value []byte) error {
		var user User
		if err := json.Unmarshal(value, &user); err != nil {
			return err
		}
		if user.ID == userID && found == nil {
			found = &user
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_user_by_id", err)
	}
	if found == nil {
		return nil, NewStoreError("get_user_by_id", ErrUserNotFound)
	}
	return found, nil
}

// List returns a list of users with optional filtering by path prefix.
func (s *UserStore) List(pathPrefix string, marker string, maxItems int) (*UserListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var users []*User
	count := 0
	started := marker == ""
	hasMore := false

	err := s.ForEach(func(key string, value []byte) error {
		var user User
		if err := json.Unmarshal(value, &user); err != nil {
			return err
		}

		if !started {
			if user.UserName == marker {
				started = true
			}
			return nil
		}

		if pathPrefix != "" && !strings.HasPrefix(user.Path, pathPrefix) {
			return nil
		}

		if count < maxItems {
			users = append(users, &user)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_users", err)
	}

	result := &UserListResult{
		Users:       users,
		IsTruncated: hasMore,
	}
	if len(users) > 0 {
		result.Marker = users[len(users)-1].UserName
	}

	return result, nil
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

// Delete removes a user by name.
func (s *UserStore) Delete(userName string) error {
	return s.BaseStore.Delete(userName)
}

// Exists checks whether a user exists.
func (s *UserStore) Exists(userName string) bool {
	return s.BaseStore.Exists(userName)
}

// Create creates a new user.
func (s *UserStore) Create(userName, path, accountId string, tags []Tag) (*User, error) {
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
	user, err := s.Get(userName)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	user.PasswordLastUsed = &now

	return s.Put(user)
}

// Count returns the total number of users.
func (s *UserStore) Count() int {
	return s.BaseStore.Count()
}
