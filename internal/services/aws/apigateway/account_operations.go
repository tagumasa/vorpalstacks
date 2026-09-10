package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/store/aws/apigateway"
)

// GetAccount returns the account-level API Gateway configuration
// (GET /account).
func (s *APIGatewayService) GetAccount(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	account, err := s.getAccountCore(stores)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toAccountResponse(account), nil
}

// UpdateAccount patches the account-level configuration
// (PATCH /account).
func (s *APIGatewayService) UpdateAccount(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	account, err := s.updateAccountCore(stores, ops)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toAccountResponse(account), nil
}

// toAccountResponse renders the account configuration in the wire shape.
func (s *APIGatewayService) toAccountResponse(account *apigateway.Account) map[string]interface{} {
	response := map[string]interface{}{
		"throttleSettings": map[string]interface{}{
			"burstLimit": account.ThrottleSettings.BurstLimit,
			"rateLimit":  account.ThrottleSettings.RateLimit,
		},
	}
	if account.CloudwatchRoleArn != "" {
		response["cloudwatchRoleArn"] = account.CloudwatchRoleArn
	}
	if len(account.Features) > 0 {
		response["features"] = account.Features
	}
	return response
}
