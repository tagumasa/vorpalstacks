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
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	dataSourceName := request.GetStringParam(req.Parameters, "dataSourceName")
	if dataSourceName == "" {
		return nil, NewBadRequestException("dataSourceName is required")
	}

	f := &appsyncstore.FunctionConfiguration{
		ApiId:                   apiId,
		Name:                    name,
		DataSourceName:          dataSourceName,
		Description:             request.GetStringParam(req.Parameters, "description"),
		FunctionVersion:         request.GetStringParam(req.Parameters, "functionVersion"),
		RequestMappingTemplate:  request.GetStringParam(req.Parameters, "requestMappingTemplate"),
		ResponseMappingTemplate: request.GetStringParam(req.Parameters, "responseMappingTemplate"),
		Runtime:                 parseAppSyncRuntime(req.Parameters),
		Code:                    request.GetStringParam(req.Parameters, "code"),
		MaxBatchSize:            int32(request.GetIntParam(req.Parameters, "maxBatchSize")),
		SyncConfig:              parseSyncConfig(req.Parameters),
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
		return nil, err
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
		return nil, err
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

	dataSourceName := request.GetStringParam(req.Parameters, "dataSourceName")
	if dataSourceName == "" {
		return nil, NewBadRequestException("dataSourceName is required")
	}

	f := &appsyncstore.FunctionConfiguration{
		ApiId:                   apiId,
		FunctionId:              functionId,
		Name:                    name,
		DataSourceName:          dataSourceName,
		Description:             request.GetStringParam(req.Parameters, "description"),
		FunctionVersion:         request.GetStringParam(req.Parameters, "functionVersion"),
		RequestMappingTemplate:  request.GetStringParam(req.Parameters, "requestMappingTemplate"),
		ResponseMappingTemplate: request.GetStringParam(req.Parameters, "responseMappingTemplate"),
		Runtime:                 parseAppSyncRuntime(req.Parameters),
		Code:                    request.GetStringParam(req.Parameters, "code"),
		MaxBatchSize:            int32(request.GetIntParam(req.Parameters, "maxBatchSize")),
		SyncConfig:              parseSyncConfig(req.Parameters),
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
		return nil, err
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
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	opts := parsePaginationOptions(req)
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

// functionToMap converts a FunctionConfiguration struct to a response map with correct wire format.
func functionToMap(f *appsyncstore.FunctionConfiguration) map[string]interface{} {
	m := map[string]interface{}{
		"functionId":     f.FunctionId,
		"name":           f.Name,
		"dataSourceName": f.DataSourceName,
	}

	if f.FunctionArn != "" {
		m["functionArn"] = f.FunctionArn
	}
	if f.FunctionVersion != "" {
		m["functionVersion"] = f.FunctionVersion
	}
	if f.Description != "" {
		m["description"] = f.Description
	}
	if f.RequestMappingTemplate != "" {
		m["requestMappingTemplate"] = f.RequestMappingTemplate
	}
	if f.ResponseMappingTemplate != "" {
		m["responseMappingTemplate"] = f.ResponseMappingTemplate
	}
	if f.Runtime != nil {
		m["runtime"] = appSyncRuntimeToMap(f.Runtime)
	}
	if f.Code != "" {
		m["code"] = f.Code
	}
	if f.MaxBatchSize > 0 {
		m["maxBatchSize"] = f.MaxBatchSize
	}
	if f.SyncConfig != nil {
		m["syncConfig"] = syncConfigToMap(f.SyncConfig)
	}

	return m
}
