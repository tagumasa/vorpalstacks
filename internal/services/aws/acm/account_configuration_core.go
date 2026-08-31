package acm

import (
	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	acmstore "vorpalstacks/internal/store/aws/acm"
)

// defaultDaysBeforeExpiry mirrors the AWS default: accounts receive expiry
// events starting 45 days before certificate expiration.
const defaultDaysBeforeExpiry = 45

// AccountConfigurationResult is the transport-agnostic account configuration.
type AccountConfigurationResult struct {
	DaysBeforeExpiry int
}

// getAccountConfigurationCore is the single read path for the ACM account
// configuration, shared by the HTTP API handler.
func (s *ACMService) getAccountConfigurationCore(reqCtx *request.RequestContext) (*AccountConfigurationResult, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := stores.certificates.GetAccountConfiguration(reqCtx.GetAccountID(), reqCtx.GetRegion())
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
func (s *ACMService) putAccountConfigurationCore(reqCtx *request.RequestContext, in PutAccountConfigurationInput) error {
	// IdempotencyToken is @required per Smithy model.
	// Validate both presence and format (@length(1-32) + @pattern(^\w+$)).
	if in.IdempotencyToken == "" {
		return awserrors.NewValidationException("IdempotencyToken is required")
	}
	if err := validateIdempotencyToken(in.IdempotencyToken); err != nil {
		return err
	}

	daysBeforeExpiry := defaultDaysBeforeExpiry
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

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	return stores.certificates.PutAccountConfiguration(reqCtx.GetAccountID(), reqCtx.GetRegion(), config)
}
