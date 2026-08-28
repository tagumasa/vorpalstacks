package appsync

import (
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// functionInput carries the parsed function request payload shared by the
// create and update operations. HasMaxBatchSize distinguishes an explicitly
// supplied maxBatchSize from an omitted one.
type functionInput struct {
	ApiId                   string
	FunctionId              string
	Name                    string
	DataSourceName          string
	Description             string
	FunctionVersion         string
	RequestMappingTemplate  string
	ResponseMappingTemplate string
	Runtime                 *appsyncstore.AppSyncRuntime
	Code                    string
	MaxBatchSize            int32
	HasMaxBatchSize         bool
	SyncConfig              *appsyncstore.SyncConfig
}

// createFunctionCore validates the request and persists a new AppSync
// function (a reusable resolver unit).
func (s *AppSyncService) createFunctionCore(store *appsyncstore.AppSyncStore, in functionInput) (*appsyncstore.FunctionConfiguration, error) {
	if in.ApiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}
	if err := validateGraphqlApiExists(store, in.ApiId); err != nil {
		return nil, err
	}

	if in.Name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if err := validateResourceName(in.Name); err != nil {
		return nil, err
	}

	if in.DataSourceName == "" {
		return nil, NewBadRequestException("dataSourceName is required")
	}

	if err := validateDescription(in.Description); err != nil {
		return nil, err
	}

	f := &appsyncstore.FunctionConfiguration{
		ApiId:                   in.ApiId,
		Name:                    in.Name,
		DataSourceName:          in.DataSourceName,
		Description:             in.Description,
		FunctionVersion:         in.FunctionVersion,
		RequestMappingTemplate:  in.RequestMappingTemplate,
		ResponseMappingTemplate: in.ResponseMappingTemplate,
		Runtime:                 in.Runtime,
		Code:                    in.Code,
		MaxBatchSize:            in.MaxBatchSize,
		SyncConfig:              in.SyncConfig,
	}

	if err := validateFunctionPayload(f, in.HasMaxBatchSize); err != nil {
		return nil, err
	}

	created, err := store.CreateFunction(f)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return created, nil
}

// getFunctionCore fetches an AppSync function by API ID and function ID.
func (s *AppSyncService) getFunctionCore(store *appsyncstore.AppSyncStore, apiId, functionId string) (*appsyncstore.FunctionConfiguration, error) {
	if apiId == "" || functionId == "" {
		return nil, NewBadRequestException("apiId and functionId are required")
	}

	f, err := store.GetFunction(apiId, functionId)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return f, nil
}

// updateFunctionCore validates the request and applies an update to an
// existing AppSync function.
func (s *AppSyncService) updateFunctionCore(store *appsyncstore.AppSyncStore, in functionInput) (*appsyncstore.FunctionConfiguration, error) {
	if in.ApiId == "" || in.FunctionId == "" {
		return nil, NewBadRequestException("apiId and functionId are required")
	}

	if in.Name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if err := validateResourceName(in.Name); err != nil {
		return nil, err
	}

	if in.DataSourceName == "" {
		return nil, NewBadRequestException("dataSourceName is required")
	}

	if err := validateDescription(in.Description); err != nil {
		return nil, err
	}

	f := &appsyncstore.FunctionConfiguration{
		ApiId:                   in.ApiId,
		FunctionId:              in.FunctionId,
		Name:                    in.Name,
		DataSourceName:          in.DataSourceName,
		Description:             in.Description,
		FunctionVersion:         in.FunctionVersion,
		RequestMappingTemplate:  in.RequestMappingTemplate,
		ResponseMappingTemplate: in.ResponseMappingTemplate,
		Runtime:                 in.Runtime,
		Code:                    in.Code,
		MaxBatchSize:            in.MaxBatchSize,
		SyncConfig:              in.SyncConfig,
	}

	if err := validateFunctionPayload(f, in.HasMaxBatchSize); err != nil {
		return nil, err
	}

	updated, err := store.UpdateFunction(f)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return updated, nil
}

// deleteFunctionCore removes an AppSync function.
func (s *AppSyncService) deleteFunctionCore(store *appsyncstore.AppSyncStore, apiId, functionId string) error {
	if apiId == "" || functionId == "" {
		return NewBadRequestException("apiId and functionId are required")
	}

	if err := store.DeleteFunction(apiId, functionId); err != nil {
		return mapStoreErrorE(err)
	}

	return nil
}

// listFunctionsCore lists the AppSync functions of a GraphQL API.
func (s *AppSyncService) listFunctionsCore(store *appsyncstore.AppSyncStore, apiId string, maxResults int, nextToken string) ([]*appsyncstore.FunctionConfiguration, string, error) {
	if apiId == "" {
		return nil, "", NewBadRequestException("apiId is required")
	}

	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	functions, nextToken, err := store.ListFunctions(apiId, opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	return functions, nextToken, nil
}

// validateFunctionPayload applies the shared runtime/code/template/batch-size
// checks to a built function configuration in the order the operations
// historically applied them.
func validateFunctionPayload(f *appsyncstore.FunctionConfiguration, hasMaxBatchSize bool) error {
	if err := validateAppSyncRuntime(f.Runtime); err != nil {
		return err
	}
	if hasMaxBatchSize {
		if err := validateMaxBatchSize(f.MaxBatchSize); err != nil {
			return err
		}
	}
	if f.Code != "" {
		if err := validateCode(f.Code); err != nil {
			return err
		}
	}
	if f.RequestMappingTemplate != "" {
		if err := validateMappingTemplate(f.RequestMappingTemplate); err != nil {
			return err
		}
	}
	if f.ResponseMappingTemplate != "" {
		if err := validateMappingTemplate(f.ResponseMappingTemplate); err != nil {
			return err
		}
	}
	return nil
}
