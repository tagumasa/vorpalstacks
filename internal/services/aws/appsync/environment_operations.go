package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// GetGraphqlApiEnvironmentVariables retrieves the environment variables for a GraphQL API.
func (s *AppSyncService) GetGraphqlApiEnvironmentVariables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	envVars, err := s.getGraphqlApiEnvironmentVariablesCore(store, request.GetStringParam(req.Parameters, "apiId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"environmentVariables": envVars,
	}, nil
}

// PutGraphqlApiEnvironmentVariables sets the environment variables for a GraphQL API.
func (s *AppSyncService) PutGraphqlApiEnvironmentVariables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	stringMap, err := s.putGraphqlApiEnvironmentVariablesCore(store, request.GetStringParam(req.Parameters, "apiId"), request.GetMapParam(req.Parameters, "environmentVariables"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"environmentVariables": stringMap,
	}, nil
}
