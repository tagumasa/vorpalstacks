package appsync

import (
	"context"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/services/aws/common/request"
)

// CreateResolver creates a new resolver for a GraphQL API type and field.
// POST /v1/apis/{apiId}/types/{typeName}/resolvers
func (s *AppSyncService) CreateResolver(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
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

// resolverToMap converts a Resolver struct to a response map with correct wire format.
func resolverToMap(r *appsyncstore.Resolver) map[string]interface{} {
	m := map[string]interface{}{
		"typeName":  r.TypeName,
		"fieldName": r.FieldName,
	}

	if r.ResolverArn != "" {
		m["resolverArn"] = r.ResolverArn
	}
	if r.Kind != "" {
		m["kind"] = r.Kind
	}
	if r.DataSourceName != "" {
		m["dataSourceName"] = r.DataSourceName
	}
	if r.RequestMappingTemplate != "" {
		m["requestMappingTemplate"] = r.RequestMappingTemplate
	}
	if r.ResponseMappingTemplate != "" {
		m["responseMappingTemplate"] = r.ResponseMappingTemplate
	}
	if r.PipelineConfig != nil {
		m["pipelineConfig"] = pipelineConfigToMap(r.PipelineConfig)
	}
	if r.Runtime != nil {
		m["runtime"] = appSyncRuntimeToMap(r.Runtime)
	}
	if r.Code != "" {
		m["code"] = r.Code
	}
	if r.CachingConfig != nil {
		m["cachingConfig"] = cachingConfigToMap(r.CachingConfig)
	}
	if r.MaxBatchSize > 0 {
		m["maxBatchSize"] = r.MaxBatchSize
	}
	if r.MetricsConfig != "" {
		m["metricsConfig"] = r.MetricsConfig
	}
	if r.SyncConfig != nil {
		m["syncConfig"] = syncConfigToMap(r.SyncConfig)
	}

	return m
}

// parseAppSyncRuntime parses an AppSyncRuntime from request parameters.
func parseAppSyncRuntime(params map[string]interface{}) *appsyncstore.AppSyncRuntime {
	raw := request.GetMapParam(params, "runtime")
	if raw == nil {
		return nil
	}
	return &appsyncstore.AppSyncRuntime{
		Name:           request.GetStringParam(raw, "name"),
		RuntimeVersion: request.GetStringParam(raw, "runtimeVersion"),
	}
}

// parsePipelineConfig parses a PipelineConfig from request parameters.
func parsePipelineConfig(params map[string]interface{}) *appsyncstore.PipelineConfig {
	raw := request.GetMapParam(params, "pipelineConfig")
	if raw == nil {
		return nil
	}
	functions := request.GetStringList(raw, "functions")
	if len(functions) == 0 {
		return nil
	}
	return &appsyncstore.PipelineConfig{
		Functions: functions,
	}
}

// parseCachingConfig parses a CachingConfig from request parameters.
func parseCachingConfig(params map[string]interface{}) *appsyncstore.CachingConfig {
	raw := request.GetMapParam(params, "cachingConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.CachingConfig{
		CachingKeys: request.GetStringList(raw, "cachingKeys"),
		Ttl:         request.GetInt64Param(raw, "ttl"),
	}
}

// parseSyncConfig parses a SyncConfig from request parameters.
func parseSyncConfig(params map[string]interface{}) *appsyncstore.SyncConfig {
	raw := request.GetMapParam(params, "syncConfig")
	if raw == nil {
		return nil
	}
	cfg := &appsyncstore.SyncConfig{
		ConflictDetection: request.GetStringParam(raw, "conflictDetection"),
		ConflictHandler:   request.GetStringParam(raw, "conflictHandler"),
	}
	if lambdaRaw := request.GetMapParam(raw, "lambdaConflictHandlerConfig"); lambdaRaw != nil {
		cfg.LambdaConflictHandlerConfig = &appsyncstore.LambdaConflictHandlerConfig{
			LambdaConflictHandlerArn: request.GetStringParam(lambdaRaw, "lambdaConflictHandlerArn"),
		}
	}
	return cfg
}

// appSyncRuntimeToMap converts an AppSyncRuntime to a wire-format map.
func appSyncRuntimeToMap(r *appsyncstore.AppSyncRuntime) map[string]interface{} {
	return map[string]interface{}{
		"name":           r.Name,
		"runtimeVersion": r.RuntimeVersion,
	}
}

// pipelineConfigToMap converts a PipelineConfig to a wire-format map.
func pipelineConfigToMap(c *appsyncstore.PipelineConfig) map[string]interface{} {
	functions := make([]interface{}, 0, len(c.Functions))
	for _, f := range c.Functions {
		functions = append(functions, f)
	}
	return map[string]interface{}{
		"functions": functions,
	}
}

// cachingConfigToMap converts a CachingConfig to a wire-format map.
func cachingConfigToMap(c *appsyncstore.CachingConfig) map[string]interface{} {
	m := map[string]interface{}{
		"ttl": c.Ttl,
	}
	if len(c.CachingKeys) > 0 {
		keys := make([]interface{}, 0, len(c.CachingKeys))
		for _, k := range c.CachingKeys {
			keys = append(keys, k)
		}
		m["cachingKeys"] = keys
	}
	return m
}

// syncConfigToMap converts a SyncConfig to a wire-format map.
func syncConfigToMap(c *appsyncstore.SyncConfig) map[string]interface{} {
	m := map[string]interface{}{}
	if c.ConflictDetection != "" {
		m["conflictDetection"] = c.ConflictDetection
	}
	if c.ConflictHandler != "" {
		m["conflictHandler"] = c.ConflictHandler
	}
	if c.LambdaConflictHandlerConfig != nil {
		m["lambdaConflictHandlerConfig"] = map[string]interface{}{
			"lambdaConflictHandlerArn": c.LambdaConflictHandlerConfig.LambdaConflictHandlerArn,
		}
	}
	return m
}
