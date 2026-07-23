package appsync

import (
	"context"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// CreateResolver creates a new resolver for a GraphQL API type and field.
// POST /v1/apis/{apiId}/types/{typeName}/resolvers
func (s *AppSyncService) CreateResolver(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	typeName := request.GetStringParam(req.Parameters, "typeName")
	fieldName := request.GetStringParam(req.Parameters, "fieldName")

	if apiId == "" || typeName == "" || fieldName == "" {
		return nil, NewBadRequestException("apiId, typeName, and fieldName are required")
	}

	r := &appsyncstore.Resolver{
		ApiId:                   apiId,
		TypeName:                typeName,
		FieldName:               fieldName,
		Kind:                    request.GetStringParam(req.Parameters, "kind"),
		DataSourceName:          request.GetStringParam(req.Parameters, "dataSourceName"),
		RequestMappingTemplate:  request.GetStringParam(req.Parameters, "requestMappingTemplate"),
		ResponseMappingTemplate: request.GetStringParam(req.Parameters, "responseMappingTemplate"),
		Runtime:                 parseAppSyncRuntime(req.Parameters),
		Code:                    request.GetStringParam(req.Parameters, "code"),
		CachingConfig:           parseCachingConfig(req.Parameters),
		MaxBatchSize:            int32(request.GetIntParam(req.Parameters, "maxBatchSize")),
		MetricsConfig:           request.GetStringParam(req.Parameters, "metricsConfig"),
		PipelineConfig:          parsePipelineConfig(req.Parameters),
		SyncConfig:              parseSyncConfig(req.Parameters),
	}

	if err := validateCachingConfig(r.CachingConfig); err != nil {
		return nil, err
	}
	if err := validateAppSyncRuntime(r.Runtime); err != nil {
		return nil, err
	}

	created, err := store.CreateResolver(r)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"resolver": resolverToMap(created),
	}, nil
}

// GetResolver retrieves a resolver by API ID, type name, and field name.
// GET /v1/apis/{apiId}/types/{typeName}/resolvers/{fieldName}
func (s *AppSyncService) GetResolver(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	typeName := request.GetStringParam(req.Parameters, "typeName")
	fieldName := request.GetStringParam(req.Parameters, "fieldName")

	if apiId == "" || typeName == "" || fieldName == "" {
		return nil, NewBadRequestException("apiId, typeName, and fieldName are required")
	}

	r, err := store.GetResolver(apiId, typeName, fieldName)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"resolver": resolverToMap(r),
	}, nil
}

// UpdateResolver updates an existing resolver.
// POST /v1/apis/{apiId}/types/{typeName}/resolvers/{fieldName}
func (s *AppSyncService) UpdateResolver(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	typeName := request.GetStringParam(req.Parameters, "typeName")
	fieldName := request.GetStringParam(req.Parameters, "fieldName")

	if apiId == "" || typeName == "" || fieldName == "" {
		return nil, NewBadRequestException("apiId, typeName, and fieldName are required")
	}

	r := &appsyncstore.Resolver{
		ApiId:                   apiId,
		TypeName:                typeName,
		FieldName:               fieldName,
		Kind:                    request.GetStringParam(req.Parameters, "kind"),
		DataSourceName:          request.GetStringParam(req.Parameters, "dataSourceName"),
		RequestMappingTemplate:  request.GetStringParam(req.Parameters, "requestMappingTemplate"),
		ResponseMappingTemplate: request.GetStringParam(req.Parameters, "responseMappingTemplate"),
		Runtime:                 parseAppSyncRuntime(req.Parameters),
		Code:                    request.GetStringParam(req.Parameters, "code"),
		CachingConfig:           parseCachingConfig(req.Parameters),
		MaxBatchSize:            int32(request.GetIntParam(req.Parameters, "maxBatchSize")),
		MetricsConfig:           request.GetStringParam(req.Parameters, "metricsConfig"),
		PipelineConfig:          parsePipelineConfig(req.Parameters),
		SyncConfig:              parseSyncConfig(req.Parameters),
	}

	if err := validateCachingConfig(r.CachingConfig); err != nil {
		return nil, err
	}
	if err := validateAppSyncRuntime(r.Runtime); err != nil {
		return nil, err
	}

	updated, err := store.UpdateResolver(r)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"resolver": resolverToMap(updated),
	}, nil
}

// DeleteResolver removes a resolver.
// DELETE /v1/apis/{apiId}/types/{typeName}/resolvers/{fieldName}
func (s *AppSyncService) DeleteResolver(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	typeName := request.GetStringParam(req.Parameters, "typeName")
	fieldName := request.GetStringParam(req.Parameters, "fieldName")

	if apiId == "" || typeName == "" || fieldName == "" {
		return nil, NewBadRequestException("apiId, typeName, and fieldName are required")
	}

	if err := store.DeleteResolver(apiId, typeName, fieldName); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// ListResolvers returns a paginated list of resolvers for a GraphQL API type.
// GET /v1/apis/{apiId}/types/{typeName}/resolvers
func (s *AppSyncService) ListResolvers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	typeName := request.GetStringParam(req.Parameters, "typeName")
	if apiId == "" || typeName == "" {
		return nil, NewBadRequestException("apiId and typeName are required")
	}

	opts := parsePaginationOptions(req)
	resolvers, nextToken, err := store.ListResolvers(apiId, typeName, opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(resolvers))
	for _, r := range resolvers {
		items = append(items, resolverToMap(r))
	}

	response := map[string]interface{}{
		"resolvers": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

// ListResolversByFunction returns resolvers that reference a given function.
// GET /v1/apis/{apiId}/functions/{functionId}/resolvers
func (s *AppSyncService) ListResolversByFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	functionId := request.GetStringParam(req.Parameters, "functionId")
	if apiId == "" || functionId == "" {
		return nil, NewBadRequestException("apiId and functionId are required")
	}

	opts := parsePaginationOptions(req)
	resolvers, nextToken, err := store.ListResolversByFunction(apiId, functionId, opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(resolvers))
	for _, r := range resolvers {
		items = append(items, resolverToMap(r))
	}

	response := map[string]interface{}{
		"resolvers": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}
