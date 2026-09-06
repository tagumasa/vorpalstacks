package acm

import (
	awserrors "vorpalstacks/internal/common/errors"
	acmstore "vorpalstacks/internal/store/aws/acm"
)

// AccountConfigurationResult is the transport-agnostic account configuration.
type AccountConfigurationResult struct {
	DaysBeforeExpiry int
}

// getAccountConfigurationCore is the single read path for the ACM account
// configuration, shared by the HTTP API handler.
func (s *ACMService) getAccountConfigurationCore(stores *acmStores) (*AccountConfigurationResult, error) {
	config, err := stores.certificates.GetAccountConfiguration(stores.accountID, stores.region)
	if err != nil {
		return nil, err
	}

	return &AccountConfigurationResult{
		DaysBeforeExpiry: config.ExpiryEvents.DaysBeforeExpiry,
	}, nil
}

// PutAccountConfigurationInput carries the wire-extracted fields of
// PutAccountConfiguration. DaysBeforeExpiryRaw is the raw wire value; numeric
// type validation happens in the Core.
type PutAccountConfigurationInput struct {
	IdempotencyToken    string
	DaysBeforeExpiryRaw interface{}
	DaysBeforeExpirySet bool
}

// putAccountConfigurationCore is the single validation + persistence path for
// PutAccountConfiguration.
func (s *ACMService) putAccountConfigurationCore(stores *acmStores, in PutAccountConfigurationInput) error {
	// IdempotencyToken is @required per Smithy model.
	// Validate both presence and format (@length(1-32) + @pattern(^\w+$)).
	if in.IdempotencyToken == "" {
		return awserrors.NewValidationException("IdempotencyToken is required")
	}
	if err := validateIdempotencyToken(in.IdempotencyToken); err != nil {
		return err
	}

	daysBeforeExpiry := acmstore.DefaultDaysBeforeExpiry
	if in.DaysBeforeExpirySet {
		switch val := in.DaysBeforeExpiryRaw.(type) {
		case float64:
			daysBeforeExpiry = int(val)
		case int:
			daysBeforeExpiry = val
		default:
			return awserrors.NewInvalidParameterException("DaysBeforeExpiry must be a numeric value")
		}
	}

	// Smithy: DaysBeforeExpiry is a PositiveInteger (@range min 1).
	if daysBeforeExpiry < 1 {
		return awserrors.NewValidationException("DaysBeforeExpiry must be a positive integer (>= 1)")
	}

	config := &acmstore.AccountConfiguration{
		ExpiryEvents: acmstore.ExpiryEventsConfiguration{
			DaysBeforeExpiry: daysBeforeExpiry,
		},
	}

	return stores.certificates.PutAccountConfiguration(stores.accountID, stores.region, config)
}
