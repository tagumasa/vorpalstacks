package acm

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	acmstore "vorpalstacks/internal/store/aws/acm"
)

// GetAccountConfiguration retrieves the account configuration for ACM.
func (s *ACMService) GetAccountConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	config, err := store.certificates.GetAccountConfiguration(reqCtx.GetAccountID(), reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ExpiryEvents": map[string]interface{}{
			"DaysBeforeExpiry": config.ExpiryEvents.DaysBeforeExpiry,
		},
	}, nil
}

// PutAccountConfiguration updates the account configuration for ACM.
func (s *ACMService) PutAccountConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// IdempotencyToken is @required per Smithy model.
	// Validate both presence and format (@length(1-32) + @pattern(^\w+$)).
	token := request.GetStringParam(req.Parameters, "IdempotencyToken")
	if token == "" {
		return nil, awserrors.NewValidationException("IdempotencyToken is required")
	}
	if err := validateIdempotencyToken(token); err != nil {
		return nil, err
	}

	daysBeforeExpiry := 45
	if raw, ok := req.Parameters["ExpiryEvents"]; ok {
		if m, ok := raw.(map[string]interface{}); ok {
			if v, ok := m["DaysBeforeExpiry"]; ok {
				switch val := v.(type) {
				case float64:
					daysBeforeExpiry = int(val)
				case int:
					daysBeforeExpiry = val
				default:
					return nil, awserrors.NewInvalidParameterException("DaysBeforeExpiry must be a numeric value")
				}
			}
		}
	}

	// Smithy: DaysBeforeExpiry is a PositiveInteger (@range min 1).
	if daysBeforeExpiry < 1 {
		return nil, awserrors.NewValidationException("DaysBeforeExpiry must be a positive integer (>= 1)")
	}

	config := &acmstore.AccountConfiguration{
		ExpiryEvents: acmstore.ExpiryEventsConfiguration{
			DaysBeforeExpiry: daysBeforeExpiry,
		},
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.certificates.PutAccountConfiguration(reqCtx.GetAccountID(), reqCtx.GetRegion(), config); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
