package acm

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// GetAccountConfiguration retrieves the account configuration for ACM.
func (s *ACMService) GetAccountConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.getAccountConfigurationCore(stores)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ExpiryEvents": map[string]interface{}{
			"DaysBeforeExpiry": result.DaysBeforeExpiry,
		},
	}, nil
}

// PutAccountConfiguration updates the account configuration for ACM.
func (s *ACMService) PutAccountConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := PutAccountConfigurationInput{
		IdempotencyToken: request.GetStringParam(req.Parameters, "IdempotencyToken"),
	}
	if m, ok := req.Parameters["ExpiryEvents"].(map[string]interface{}); ok {
		if v, ok := m["DaysBeforeExpiry"]; ok {
			in.DaysBeforeExpiryRaw = v
			in.DaysBeforeExpirySet = true
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.putAccountConfigurationCore(stores, in); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}
