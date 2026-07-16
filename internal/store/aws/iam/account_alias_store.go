package iam

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const accountAliasBucketName = "iam_account_alias"

// AccountAliasStore provides storage operations for the IAM account alias.
type AccountAliasStore struct {
	*common.BaseStore
}

// NewAccountAliasStore creates a new AccountAliasStore instance.
func NewAccountAliasStore(store storage.BasicStorage) *AccountAliasStore {
	return &AccountAliasStore{
		BaseStore: common.NewBaseStore(store.Bucket(accountAliasBucketName), "iam"),
	}
}

// Get retrieves the current account alias.
func (s *AccountAliasStore) Get() (*AccountAlias, error) {
	data, err := s.Bucket().Get([]byte("alias"))
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
	return s.BaseStore.Put("alias", alias)
}

// Delete removes the account alias.
func (s *AccountAliasStore) Delete() error {
	return s.BaseStore.Delete("alias")
}

// Exists reports whether an account alias is set.
func (s *AccountAliasStore) Exists() bool {
	return s.BaseStore.Exists("alias")
}
