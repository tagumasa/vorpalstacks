package appsync

import (
	"context"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// CreateFunction creates a new AppSync function (a reusable resolver unit).
// POST /v1/apis/{apiId}/functions
func (s *AppSyncService) CreateFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}
	if err := validateGraphqlApiExists(store, apiId); err != nil {
		return nil, err
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if err := validateResourceName(name); err != nil {
		return nil, err
	}

	dataSourceName := request.GetStringParam(req.Parameters, "dataSourceName")
	if dataSourceName == "" {
		return nil, NewBadRequestException("dataSourceName is required")
	}

	description := request.GetStringParam(req.Parameters, "description")
	if err := validateDescription(description); err != nil {
		return nil, err
	}

	syncCfg, err := parseSyncConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	f := &appsyncstore.FunctionConfiguration{
		ApiId:                   apiId,
		Name:                    name,
		DataSourceName:          dataSourceName,
		Description:             description,
		FunctionVersion:         request.GetStringParam(req.Parameters, "functionVersion"),
		RequestMappingTemplate:  request.GetStringParam(req.Parameters, "requestMappingTemplate"),
		ResponseMappingTemplate: request.GetStringParam(req.Parameters, "responseMappingTemplate"),
		Runtime:                 parseAppSyncRuntime(req.Parameters),
		Code:                    request.GetStringParam(req.Parameters, "code"),
		MaxBatchSize:            int32(request.GetIntParam(req.Parameters, "maxBatchSize")),
		SyncConfig:              syncCfg,
	}

	if err := validateAppSyncRuntime(f.Runtime); err != nil {
		return nil, err
	}
	if _, ok := req.Parameters["maxBatchSize"]; ok {
		if err := validateMaxBatchSize(f.MaxBatchSize); err != nil {
			return nil, err
		}
	}
	if f.Code != "" {
		if err := validateCode(f.Code); err != nil {
			return nil, err
		}
	}
	if f.RequestMappingTemplate != "" {
		if err := validateMappingTemplate(f.RequestMappingTemplate); err != nil {
			return nil, err
		}
	}
	if f.ResponseMappingTemplate != "" {
		if err := validateMappingTemplate(f.ResponseMappingTemplate); err != nil {
			return nil, err
		}
	}

	created, err := store.CreateFunction(f)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"functionConfiguration": functionToMap(created),
	}, nil
}

// GetFunction retrieves an AppSync function by API ID and function ID.
// GET /v1/apis/{apiId}/functions/{functionId}
func (s *AppSyncService) GetFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	functionId := request.GetStringParam(req.Parameters, "functionId")
	if apiId == "" || functionId == "" {
		return nil, NewBadRequestException("apiId and functionId are required")
	}

	f, err := store.GetFunction(apiId, functionId)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"functionConfiguration": functionToMap(f),
	}, nil
}

// UpdateFunction updates an existing AppSync function.
// POST /v1/apis/{apiId}/functions/{functionId}
func (s *AppSyncService) UpdateFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	functionId := request.GetStringParam(req.Parameters, "functionId")
	if apiId == "" || functionId == "" {
		return nil, NewBadRequestException("apiId and functionId are required")
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if err := validateResourceName(name); err != nil {
		return nil, err
	}

	dataSourceName := request.GetStringParam(req.Parameters, "dataSourceName")
	if dataSourceName == "" {
		return nil, NewBadRequestException("dataSourceName is required")
	}

	description := request.GetStringParam(req.Parameters, "description")
	if err := validateDescription(description); err != nil {
		return nil, err
	}

	syncCfg, err := parseSyncConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	f := &appsyncstore.FunctionConfiguration{
		ApiId:                   apiId,
		FunctionId:              functionId,
		Name:                    name,
		DataSourceName:          dataSourceName,
		Description:             description,
		FunctionVersion:         request.GetStringParam(req.Parameters, "functionVersion"),
		RequestMappingTemplate:  request.GetStringParam(req.Parameters, "requestMappingTemplate"),
		ResponseMappingTemplate: request.GetStringParam(req.Parameters, "responseMappingTemplate"),
		Runtime:                 parseAppSyncRuntime(req.Parameters),
		Code:                    request.GetStringParam(req.Parameters, "code"),
		MaxBatchSize:            int32(request.GetIntParam(req.Parameters, "maxBatchSize")),
		SyncConfig:              syncCfg,
	}

	if err := validateAppSyncRuntime(f.Runtime); err != nil {
		return nil, err
	}
	if _, ok := req.Parameters["maxBatchSize"]; ok {
		if err := validateMaxBatchSize(f.MaxBatchSize); err != nil {
			return nil, err
		}
	}
	if f.Code != "" {
		if err := validateCode(f.Code); err != nil {
			return nil, err
		}
	}
	if f.RequestMappingTemplate != "" {
		if err := validateMappingTemplate(f.RequestMappingTemplate); err != nil {
			return nil, err
		}
	}
	if f.ResponseMappingTemplate != "" {
		if err := validateMappingTemplate(f.ResponseMappingTemplate); err != nil {
			return nil, err
		}
	}

	updated, err := store.UpdateFunction(f)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"functionConfiguration": functionToMap(updated),
	}, nil
}

// DeleteFunction removes an AppSync function.
// DELETE /v1/apis/{apiId}/functions/{functionId}
func (s *AppSyncService) DeleteFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	functionId := request.GetStringParam(req.Parameters, "functionId")
	if apiId == "" || functionId == "" {
		return nil, NewBadRequestException("apiId and functionId are required")
	}

	if err := store.DeleteFunction(apiId, functionId); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// ListFunctions returns a paginated list of AppSync functions for a GraphQL API.
// GET /v1/apis/{apiId}/functions
func (s *AppSyncService) ListFunctions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	opts, err := parsePaginationOptions(req)
	if err != nil {
		return nil, err
	}
	functions, nextToken, err := store.ListFunctions(apiId, opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(functions))
	for _, f := range functions {
		items = append(items, functionToMap(f))
	}

	response := map[string]interface{}{
		"functions": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}
