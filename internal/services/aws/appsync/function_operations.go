package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// parseFunctionInput extracts the shared function payload members from the
// request parameters.
func parseFunctionInput(req *request.ParsedRequest) (functionInput, error) {
	syncCfg, err := parseSyncConfig(req.Parameters)
	if err != nil {
		return functionInput{}, err
	}

	in := functionInput{
		ApiId:                   request.GetStringParam(req.Parameters, "apiId"),
		FunctionId:              request.GetStringParam(req.Parameters, "functionId"),
		Name:                    request.GetStringParam(req.Parameters, "name"),
		DataSourceName:          request.GetStringParam(req.Parameters, "dataSourceName"),
		Description:             request.GetStringParam(req.Parameters, "description"),
		FunctionVersion:         request.GetStringParam(req.Parameters, "functionVersion"),
		RequestMappingTemplate:  request.GetStringParam(req.Parameters, "requestMappingTemplate"),
		ResponseMappingTemplate: request.GetStringParam(req.Parameters, "responseMappingTemplate"),
		Runtime:                 parseAppSyncRuntime(req.Parameters),
		Code:                    request.GetStringParam(req.Parameters, "code"),
		MaxBatchSize:            int32(request.GetIntParam(req.Parameters, "maxBatchSize")),
		SyncConfig:              syncCfg,
	}
	_, in.HasMaxBatchSize = req.Parameters["maxBatchSize"]

	return in, nil
}

// CreateFunction creates a new AppSync function (a reusable resolver unit).
// POST /v1/apis/{apiId}/functions
func (s *AppSyncService) CreateFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	in, err := parseFunctionInput(req)
	if err != nil {
		return nil, err
	}

	created, err := s.createFunctionCore(store, in)
	if err != nil {
		return nil, err
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

	f, err := s.getFunctionCore(store, request.GetStringParam(req.Parameters, "apiId"), request.GetStringParam(req.Parameters, "functionId"))
	if err != nil {
		return nil, err
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

	in, err := parseFunctionInput(req)
	if err != nil {
		return nil, err
	}

	updated, err := s.updateFunctionCore(store, in)
	if err != nil {
		return nil, err
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

	if err := s.deleteFunctionCore(store, request.GetStringParam(req.Parameters, "apiId"), request.GetStringParam(req.Parameters, "functionId")); err != nil {
		return nil, err
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

	functions, nextToken, err := s.listFunctionsCore(store, apiId, request.GetIntParam(req.Parameters, "maxResults"), request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
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
