package apigateway

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// ThrottleSettings records the account-level API request limits.
type ThrottleSettings struct {
	BurstLimit int64   `json:"burst_limit,omitempty"`
	RateLimit  float64 `json:"rate_limit,omitempty"`
}

// Account is the account-level API Gateway configuration singleton.
type Account struct {
	CloudwatchRoleArn string            `json:"cloudwatch_role_arn,omitempty"`
	Features          []string          `json:"features,omitempty"`
	ThrottleSettings  *ThrottleSettings `json:"throttle_settings,omitempty"`
}

// AccountStore persists the account-level configuration: one record per
// platform account inside each region's store bundle.
type AccountStore struct {
	*common.BaseStore
	accountId string
}

// NewAccountStore creates an AccountStore over the given storage. The
// record key is the account ID alone — the storage instance is already
// region-scoped, so the account configuration is stored per account per
// region, matching AWS, where account-level settings apply per account
// per region.
func NewAccountStore(store storage.BasicStorage, accountId, _ string) *AccountStore {
	bucket := store.Bucket("apigateway-account")
	return &AccountStore{
		BaseStore: common.NewBaseStore(bucket, "apigateway-account"),
		accountId: accountId,
	}
}

// Get returns the stored account configuration, defaulting to an empty one
// when nothing has been stored yet.
func (s *AccountStore) Get() (*Account, error) {
	var account Account
	if err := s.BaseStore.Get(s.accountId, &account); err != nil {
		if common.IsNotFound(err) {
			return &Account{}, nil
		}
		return nil, err
	}
	return &account, nil
}

// Update stores the account configuration.
func (s *AccountStore) Update(account *Account) error {
	return s.BaseStore.Put(s.accountId, account)
}
