package cognitoidentityprovider

import (
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// Core functions for the provisioned-limit family. Provisioned limits are
// the adjustable per-category API request-rate quotas; they are regional
// account-level resources, so the Core persists them through the regional
// store instead of process memory.

// limitClassAPICategory is the sole member of the model's LimitClass enum:
// every provisioned limit identifies a per-API-category rate quota.
const limitClassAPICategory = "API_CATEGORY"

// ProvisionedLimitInput carries the LimitDefinition wire members shared by
// GetProvisionedLimit and UpdateProvisionedLimit. The attributes that
// identify a user-pool rate limit are the single Category key.
type ProvisionedLimitInput struct {
	Region     string
	LimitClass string
	Category   string
	// RequestedValue and RequestedValueSet carry the update request's
	// required RequestedLimitValue; get requests leave Set false.
	RequestedValue    int
	RequestedValueSet bool
}

// ProvisionedLimitResult is the LimitType wire shape: the limit definition
// echoed back with the provisioned and free values.
type ProvisionedLimitResult struct {
	LimitClass            string
	Category              string
	ProvisionedLimitValue int
	FreeLimitValue        int
}

// resolveProvisionedCategory validates the limit definition: the class must
// be the model enum's sole member, and the category must be an adjustable
// category — only adjustable quota categories support provisioning, so any
// other category identifies no provisionable limit. Returns the category's
// free quota.
func resolveProvisionedCategory(in ProvisionedLimitInput) (int, error) {
	if in.LimitClass != limitClassAPICategory || in.Category == "" {
		return 0, ErrInvalidParameter
	}
	free, ok := cognitostore.AdjustableAPICategoryFreeLimits[in.Category]
	if !ok {
		return 0, ErrResourceNotFound
	}
	return free, nil
}

// getProvisionedLimitCore returns the provisioned and free values for an
// adjustable category. A category with no stored limit answers its free
// quota as both values, matching the documented initial state.
func (s *CognitoService) getProvisionedLimitCore(in ProvisionedLimitInput) (*ProvisionedLimitResult, error) {
	free, err := resolveProvisionedCategory(in)
	if err != nil {
		return nil, err
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, err
	}

	provisioned := free
	if stored, err := store.GetProvisionedLimit(in.Category); err == nil {
		provisioned = stored.Value
	}

	return &ProvisionedLimitResult{
		LimitClass:            limitClassAPICategory,
		Category:              in.Category,
		ProvisionedLimitValue: provisioned,
		FreeLimitValue:        free,
	}, nil
}

// updateProvisionedLimitCore sets the provisioned rate for an adjustable
// category. The requested value must be at least the category's free quota
// — the documented lower bound of the adjustable range, inclusive, since a
// provisioned rate can be reduced back to the default. The range's upper
// bound is the account-level maximum in Service Quotas, which has no local
// counterpart, so no local ceiling is enforced.
func (s *CognitoService) updateProvisionedLimitCore(in ProvisionedLimitInput) (*ProvisionedLimitResult, error) {
	free, err := resolveProvisionedCategory(in)
	if err != nil {
		return nil, err
	}
	if !in.RequestedValueSet || in.RequestedValue < free {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(in.Region)
	if err != nil {
		return nil, err
	}

	limit := &cognitostore.ProvisionedLimit{Category: in.Category, Value: in.RequestedValue}
	if err := store.SaveProvisionedLimit(limit); err != nil {
		return nil, ErrInternalError
	}

	return &ProvisionedLimitResult{
		LimitClass:            limitClassAPICategory,
		Category:              in.Category,
		ProvisionedLimitValue: in.RequestedValue,
		FreeLimitValue:        free,
	}, nil
}
