package cognitoidentityprovider

import (
	"time"
)

// AdjustableAPICategoryFreeLimits maps each adjustable Amazon Cognito user
// pools API category to its default (free) request-rate quota in requests
// per second. Source: Amazon Cognito developer guide, "Quotas in Amazon
// Cognito" — the per-category request-rate quota table, which also marks
// these six categories as the adjustable ones. Only adjustable categories
// support provisioned limits, so this table is both the registry of
// provisionable categories and the FreeLimitValue answer for each.
var AdjustableAPICategoryFreeLimits = map[string]int{
	"UserAuthentication": 120,
	"UserCreation":       50,
	"UserFederation":     25,
	"UserRead":           120,
	"UserToken":          120,
	"UserResourceRead":   50,
}

// ProvisionedLimit is the provisioned request rate for one adjustable API
// category. Provisioned limits are account-level resources scoped to one
// Region; the regional store provides that scoping, so two regional stores
// over the same storage never share a limit. A category absent from the
// store has no provisioned increase and answers its free quota.
type ProvisionedLimit struct {
	Category         string
	Value            int
	LastModifiedDate time.Time
}

// GetProvisionedLimit retrieves the provisioned limit for an API category.
// It returns a not-found store error when the category has no provisioned
// limit; callers fall back to the category's free quota.
func (s *CognitoStore) GetProvisionedLimit(category string) (*ProvisionedLimit, error) {
	var limit ProvisionedLimit
	if err := s.BaseStore.Get(provisionedLimitKey(category), &limit); err != nil {
		return nil, err
	}
	return &limit, nil
}

// SaveProvisionedLimit stores the provisioned limit for an API category.
func (s *CognitoStore) SaveProvisionedLimit(limit *ProvisionedLimit) error {
	limit.LastModifiedDate = time.Now().UTC()
	return s.BaseStore.Put(provisionedLimitKey(limit.Category), limit)
}
