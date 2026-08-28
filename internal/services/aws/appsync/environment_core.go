package appsync

import (
	"fmt"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// getGraphqlApiEnvironmentVariablesCore fetches the environment variables of
// a GraphQL API.
func (s *AppSyncService) getGraphqlApiEnvironmentVariablesCore(store *appsyncstore.AppSyncStore, apiId string) (map[string]string, error) {
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if _, err := store.GetGraphqlApiById(apiId); err != nil {
		return nil, mapStoreErrorE(err)
	}

	envVars, err := store.GetEnvironmentVariables(apiId)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return envVars.EnvironmentVariables, nil
}

// putGraphqlApiEnvironmentVariablesCore validates and persists the
// environment variables of a GraphQL API, returning the stored string map.
// The environmentVariables member is required; an explicit empty map clears
// the variables while an omitted member is rejected.
func (s *AppSyncService) putGraphqlApiEnvironmentVariablesCore(store *appsyncstore.AppSyncStore, apiId string, envVars map[string]interface{}) (map[string]string, error) {
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if envVars == nil {
		return nil, NewBadRequestException("environmentVariables is required")
	}

	if _, err := store.GetGraphqlApiById(apiId); err != nil {
		return nil, mapStoreErrorE(err)
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
		return nil, mapStoreErrorE(err)
	}

	return stringMap, nil
}
