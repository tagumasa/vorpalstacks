package iam

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
)

const accountAliasBucketName = "iam_account_alias"

// AccountAliasStore provides storage operations for the IAM account alias.
type AccountAliasStore struct {
	bucket storage.Bucket
}

// NewAccountAliasStore creates a new AccountAliasStore instance.
func NewAccountAliasStore(store storage.BasicStorage) *AccountAliasStore {
	return &AccountAliasStore{
		bucket: store.Bucket(accountAliasBucketName),
	}
}

// Get retrieves the current account alias.
func (s *AccountAliasStore) Get() (*AccountAlias, error) {
	data, err := s.bucket.Get([]byte("alias"))
	if err != nil {
		return nil, NewStoreError("get_account_alias", err)
	}
	if data == nil {
		return nil, nil
	}
	var alias AccountAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return nil, NewStoreError("get_account_alias", err)
	}
	return &alias, nil
}

// Put stores the account alias.
func (s *AccountAliasStore) Put(accountAlias string) error {
	alias := &AccountAlias{
		AccountAlias: accountAlias,
		CreateDate:   time.Now().UTC(),
	}
	data, err := json.Marshal(alias)
	if err != nil {
		return NewStoreError("put_account_alias", err)
	}
	if err := s.bucket.Put([]byte("alias"), data); err != nil {
		return NewStoreError("put_account_alias", err)
	}
	return nil
}

// Delete removes the account alias.
func (s *AccountAliasStore) Delete() error {
	if err := s.bucket.Delete([]byte("alias")); err != nil {
		return NewStoreError("delete_account_alias", err)
	}
	return nil
}

// Exists reports whether an account alias is set.
func (s *AccountAliasStore) Exists() bool {
	return s.bucket.Has([]byte("alias"))
}
