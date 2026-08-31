package timestreamquery

import (
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// ---------------------------------------------------------------------------
// DescribeAccountSettings / UpdateAccountSettings — DTO inputs and Core
// functions. Enum validation and store access live here.
// ---------------------------------------------------------------------------

// UpdateAccountSettingsInput is the DTO input for the UpdateAccountSettings
// operation. The enum members stay in raw wire form; the Core owns their
// validation. MaxQueryTCU is nil when the request omits the member.
type UpdateAccountSettingsInput struct {
	MaxQueryTCU         *int64
	QueryPricingModel   string
	HasComputeMode      bool
	ComputeMode         string
	ProvisionedCapacity *tsstore.ProvisionedCapacitySettings
}

// describeAccountSettingsCore returns the stored account settings.
func (s *TimestreamQueryService) describeAccountSettingsCore(stores *tsQueryStores) (*tsstore.AccountSettings, error) {
	return stores.accountSettingsStore.GetAccountSettings()
}

// updateAccountSettingsCore validates the request and updates the account
// settings.
func (s *TimestreamQueryService) updateAccountSettingsCore(stores *tsQueryStores, input UpdateAccountSettingsInput) (*tsstore.AccountSettings, error) {
	// Validate QueryPricingModel enum (Smithy: BYTES_SCANNED, COMPUTE_UNITS).
	if input.QueryPricingModel != "" && !validQueryPricingModels[input.QueryPricingModel] {
		return nil, ErrValidationException
	}

	queryComputeType := ""
	if input.HasComputeMode {
		// Validate ComputeMode enum (Smithy: ON_DEMAND, PROVISIONED).
		if input.ComputeMode != "" && !validComputeModes[input.ComputeMode] {
			return nil, ErrValidationException
		}
		queryComputeType = input.ComputeMode
	}

	return stores.accountSettingsStore.UpdateAccountSettings(input.MaxQueryTCU, input.QueryPricingModel, queryComputeType, "", input.ProvisionedCapacity)
}
