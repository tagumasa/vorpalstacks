package iam

import (
	"encoding/json"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const accountSettingsBucketName = "iam_account_settings"

// AccountSettings represents the account-level settings for IAM.
type AccountSettings struct {
	OutboundWebIdentityFederationEnabled bool   `json:"outbound_web_identity_federation_enabled"`
	IssuerIdentifier                     string `json:"issuer_identifier"`
	GlobalEndpointTokenVersion           string `json:"global_endpoint_token_version"`
}

// AccountSettingsStore provides storage operations for IAM account settings.
type AccountSettingsStore struct {
	*common.BaseStore
}

// NewAccountSettingsStore creates a new AccountSettingsStore instance.
func NewAccountSettingsStore(store storage.BasicStorage) *AccountSettingsStore {
	return &AccountSettingsStore{
		BaseStore: common.NewBaseStore(store.Bucket(accountSettingsBucketName), "iam"),
	}
}

// Get retrieves the current account settings, returning defaults if none are stored.
func (s *AccountSettingsStore) Get() (*AccountSettings, error) {
	data, err := s.Bucket().Get([]byte("settings"))
	if err != nil {
		return nil, NewStoreError("get_account_settings", err)
	}
	if data == nil {
		return &AccountSettings{
			OutboundWebIdentityFederationEnabled: false,
			IssuerIdentifier:                     "",
			GlobalEndpointTokenVersion:           "v1Token",
		}, nil
	}
	var settings AccountSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, NewStoreError("get_account_settings", err)
	}
	return &settings, nil
}

// Put stores the account settings.
func (s *AccountSettingsStore) Put(settings *AccountSettings) error {
	return s.BaseStore.Put("settings", settings)
}
