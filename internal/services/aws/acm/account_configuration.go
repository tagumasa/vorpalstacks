package acm

import (
	"context"

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
	daysBeforeExpiry := 45
	if raw, ok := req.Parameters["ExpiryEvents"]; ok {
		if m, ok := raw.(map[string]interface{}); ok {
			if v, ok := m["DaysBeforeExpiry"]; ok {
				switch val := v.(type) {
				case float64:
					daysBeforeExpiry = int(val)
				case int:
					daysBeforeExpiry = val
				}
			}
		}
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
