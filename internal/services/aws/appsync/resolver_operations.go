package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// parseResolverInput extracts the shared resolver payload members from the
// request parameters.
func parseResolverInput(req *request.ParsedRequest) (resolverInput, error) {
	syncCfg, err := parseSyncConfig(req.Parameters)
	if err != nil {
		return resolverInput{}, err
	}

	in := resolverInput{
		ApiId:                   request.GetStringParam(req.Parameters, "apiId"),
		TypeName:                request.GetStringParam(req.Parameters, "typeName"),
		FieldName:               request.GetStringParam(req.Parameters, "fieldName"),
		Kind:                    request.GetStringParam(req.Parameters, "kind"),
		DataSourceName:          request.GetStringParam(req.Parameters, "dataSourceName"),
		RequestMappingTemplate:  request.GetStringParam(req.Parameters, "requestMappingTemplate"),
		ResponseMappingTemplate: request.GetStringParam(req.Parameters, "responseMappingTemplate"),
		Code:                    request.GetStringParam(req.Parameters, "code"),
		MetricsConfig:           request.GetStringParam(req.Parameters, "metricsConfig"),
		Runtime:                 parseAppSyncRuntime(req.Parameters),
		CachingConfig:           parseCachingConfig(req.Parameters),
		PipelineConfig:          parsePipelineConfig(req.Parameters),
		SyncConfig:              syncCfg,
		MaxBatchSize:            int32(request.GetIntParam(req.Parameters, "maxBatchSize")),
	}
	_, in.HasMaxBatchSize = req.Parameters["maxBatchSize"]

	return in, nil
}

// CreateResolver creates a new resolver for a GraphQL API type and field.
// POST /v1/apis/{apiId}/types/{typeName}/resolvers
func (s *AppSyncService) CreateResolver(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	in, err := parseResolverInput(req)
	if err != nil {
		return nil, err
	}

	created, err := s.createResolverCore(store, in)
	if err != nil {
		return nil, err
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

	r, err := s.getResolverCore(store,
		request.GetStringParam(req.Parameters, "apiId"),
		request.GetStringParam(req.Parameters, "typeName"),
		request.GetStringParam(req.Parameters, "fieldName"))
	if err != nil {
		return nil, err
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

	in, err := parseResolverInput(req)
	if err != nil {
		return nil, err
	}

	updated, err := s.updateResolverCore(store, in)
	if err != nil {
		return nil, err
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

	if err := s.deleteResolverCore(store,
		request.GetStringParam(req.Parameters, "apiId"),
		request.GetStringParam(req.Parameters, "typeName"),
		request.GetStringParam(req.Parameters, "fieldName")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListResolvers lists the resolvers attached to one type of a GraphQL API.
// GET /v1/apis/{apiId}/types/{typeName}/resolvers
func (s *AppSyncService) ListResolvers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	resolvers, nextToken, err := s.listResolversCore(store,
		request.GetStringParam(req.Parameters, "apiId"),
		request.GetStringParam(req.Parameters, "typeName"),
		request.GetIntParam(req.Parameters, "maxResults"),
		request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
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

// ListResolversByFunction lists the resolvers that reference a given AppSync
// function.
// GET /v1/apis/{apiId}/functions/{functionId}/resolvers
func (s *AppSyncService) ListResolversByFunction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	resolvers, nextToken, err := s.listResolversByFunctionCore(store,
		request.GetStringParam(req.Parameters, "apiId"),
		request.GetStringParam(req.Parameters, "functionId"),
		request.GetIntParam(req.Parameters, "maxResults"),
		request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
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
