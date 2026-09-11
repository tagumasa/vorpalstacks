package iam

import (
	"encoding/json"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const accountPropertiesBucketName = "iam_account_properties"

// AccountPropertiesStore provides storage operations for IAM account
// properties: the account-wide key-value configuration pairs in
// Namespace/PropertyName format (RoleManager/Enabled and future
// namespaces), served by the Get/PutAccountProperties operations.
type AccountPropertiesStore struct {
	*common.BaseStore
}

// NewAccountPropertiesStore creates a new AccountPropertiesStore instance.
func NewAccountPropertiesStore(store storage.BasicStorage) *AccountPropertiesStore {
	return &AccountPropertiesStore{
		BaseStore: common.NewBaseStore(store.Bucket(accountPropertiesBucketName), "iam"),
	}
}

// Get retrieves the stored account properties overlaid on the documented
// defaults, returning the defaults alone if nothing has been stored.
func (s *AccountPropertiesStore) Get() (map[string]string, error) {
	properties := map[string]string{
		// RoleManager/Enabled is the one documented account property; a
		// fresh account has the Role Manager feature off.
		"RoleManager/Enabled": "false",
	}

	data, err := s.GetRaw("properties")
	if err != nil {
		return nil, NewStoreError("get_account_properties", err)
	}
	if data == nil {
		return properties, nil
	}
	var stored map[string]string
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, NewStoreError("get_account_properties", err)
	}
	for key, value := range stored {
		properties[key] = value
	}
	return properties, nil
}

// Put merges the supplied property pairs into the stored set — the
// documented set semantics of PutAccountProperties, which sets the named
// pairs and leaves every other property untouched.
func (s *AccountPropertiesStore) Put(properties map[string]string) error {
	current := map[string]string{}
	data, err := s.GetRaw("properties")
	if err != nil {
		return NewStoreError("put_account_properties", err)
	}
	if data != nil {
		if err := json.Unmarshal(data, &current); err != nil {
			return NewStoreError("put_account_properties", err)
		}
	}
	for key, value := range properties {
		current[key] = value
	}
	encoded, err := json.Marshal(current)
	if err != nil {
		return NewStoreError("put_account_properties", err)
	}
	return s.PutRaw("properties", encoded)
}
