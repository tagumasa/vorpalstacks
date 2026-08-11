package appsync

import (
	"context"
	"fmt"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// GetGraphqlApiEnvironmentVariables retrieves the environment variables for a GraphQL API.
func (s *AppSyncService) GetGraphqlApiEnvironmentVariables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if _, err := store.GetGraphqlApiById(apiId); err != nil {
		return mapStoreError(err)
	}

	envVars, err := store.GetEnvironmentVariables(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"environmentVariables": envVars.EnvironmentVariables,
	}, nil
}

// PutGraphqlApiEnvironmentVariables sets the environment variables for a GraphQL API.
func (s *AppSyncService) PutGraphqlApiEnvironmentVariables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if _, err := store.GetGraphqlApiById(apiId); err != nil {
		return mapStoreError(err)
	}

	envVars := request.GetMapParam(req.Parameters, "environmentVariables")
	if envVars == nil {
		envVars = make(map[string]interface{})
	}

	stringMap := make(map[string]string)
	for k, v := range envVars {
		s, ok := v.(string)
		if !ok {
			return nil, NewBadRequestException(fmt.Sprintf("environment variable value for %q must be a string", k))
		}
		if err := validateEnvVarKey(k); err != nil {
			return nil, err
		}
		if err := validateEnvVarValue(s); err != nil {
			return nil, err
		}
		stringMap[k] = s
	}

	if err := validateEnvironmentVariableMapSize(stringMap); err != nil {
		return nil, err
	}

	toSave := &appsyncstore.EnvironmentVariables{
		EnvironmentVariables: stringMap,
	}
	if err := store.SaveEnvironmentVariables(apiId, toSave); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"environmentVariables": stringMap,
	}, nil
}
