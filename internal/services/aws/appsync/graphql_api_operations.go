package appsync

import (
	"context"
	"fmt"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// maxApisPerRegion is the AWS service quota for AppSync APIs per region.
// This limit applies to the combined count of GraphQL APIs and Event APIs.
const maxApisPerRegion = 25

// CreateGraphqlApi creates a new GraphQL API (v1).
// POST /v1/apis
func (s *AppSyncService) CreateGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	authType := request.GetStringParam(req.Parameters, "authenticationType")
	if authType == "" {
		return nil, NewBadRequestException("authenticationType is required")
	}

	if !validateAuthenticationType(authType) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid authenticationType: %s", authType))
	}

	apiType := request.GetStringParam(req.Parameters, "apiType")
	if apiType != "" && !validateApiType(apiType) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid apiType: %s", apiType))
	}

	visibility := request.GetStringParam(req.Parameters, "visibility")
	if visibility != "" && !validateVisibility(visibility) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid visibility: %s", visibility))
	}

	introspectionConfig := request.GetStringParam(req.Parameters, "introspectionConfig")
	if introspectionConfig != "" && !validateIntrospectionConfig(introspectionConfig) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid introspectionConfig: %s", introspectionConfig))
	}

	logCfg := parseLogConfig(req.Parameters)
	if logCfg != nil && logCfg.FieldLogLevel != "" && !validateFieldLogLevel(logCfg.FieldLogLevel) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid logConfig.fieldLogLevel: %s", logCfg.FieldLogLevel))
	}

	// Check combined API count quota (GraphQL + Event APIs per region).
	graphqlCount, _ := store.CountGraphqlApis()
	eventCount, _ := store.CountApis()
	if graphqlCount+eventCount >= maxApisPerRegion {
		return nil, ErrApiLimitExceededException
	}

	api := &appsyncstore.GraphqlApi{
		Name:                              name,
		AuthenticationType:                authType,
		AdditionalAuthenticationProviders: parseAdditionalAuthProviders(req.Parameters),
		ApiType:                           apiType,
		EnhancedMetricsConfig:             parseEnhancedMetricsConfig(req.Parameters),
		IntrospectionConfig:               introspectionConfig,
		LambdaAuthorizerConfig:            parseLambdaAuthorizerConfig(req.Parameters),
		LogConfig:                         logCfg,
		MergedApiExecutionRoleArn:         request.GetStringParam(req.Parameters, "mergedApiExecutionRoleArn"),
		OpenIDConnectConfig:               parseOpenIDConnectConfig(req.Parameters),
		OwnerContact:                      request.GetStringParam(req.Parameters, "ownerContact"),
		QueryDepthLimit:                   int32(request.GetIntParam(req.Parameters, "queryDepthLimit")),
		ResolverCountLimit:                int32(request.GetIntParam(req.Parameters, "resolverCountLimit")),
		Tags:                              parseTags(req.Parameters),
		UserPoolConfig:                    parseUserPoolConfig(req.Parameters),
		Visibility:                        visibility,
		WafWebAclArn:                      request.GetStringParam(req.Parameters, "wafWebAclArn"),
		XrayEnabled:                       request.GetBoolParam(req.Parameters, "xrayEnabled"),
	}

	created, err := store.CreateGraphqlApi(api)
	if err != nil {
		return mapStoreError(err)
	}

	if len(created.Tags) > 0 {
		tagMap := make(map[string]string, len(created.Tags))
		for k, v := range created.Tags {
			tagMap[k] = v
		}
		if err := store.TagStore.Tag(created.Arn, tagMap); err != nil {
			return nil, err
		}
	}

	result := graphqlApiToMap(created)
	if tags, err := store.TagStore.List(created.Arn); err == nil && len(tags) > 0 {
		result["tags"] = tags
	}

	return map[string]interface{}{
		"graphqlApi": result,
	}, nil
}

// GetGraphqlApi retrieves a GraphQL API (v1) by its ID.
// GET /v1/apis/{apiId}
func (s *AppSyncService) GetGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	api, err := store.GetGraphqlApiById(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	result := graphqlApiToMap(api)
	if tagsFromStore, err := store.TagStore.List(api.Arn); err == nil && len(tagsFromStore) > 0 {
		result["tags"] = tagsFromStore
	}

	return map[string]interface{}{
		"graphqlApi": result,
	}, nil
}

// UpdateGraphqlApi updates an existing GraphQL API (v1).
// POST /v1/apis/{apiId}
func (s *AppSyncService) UpdateGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	// Per Smithy model, name and authenticationType are required for
	// UpdateGraphqlApiRequest. The client must resend current values.
	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}
	authType := request.GetStringParam(req.Parameters, "authenticationType")
	if authType == "" {
		return nil, NewBadRequestException("authenticationType is required")
	}
	if !validateAuthenticationType(authType) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid authenticationType: %s", authType))
	}

	introspectionConfig := request.GetStringParam(req.Parameters, "introspectionConfig")
	if introspectionConfig != "" && !validateIntrospectionConfig(introspectionConfig) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid introspectionConfig: %s", introspectionConfig))
	}

	logCfg := parseLogConfig(req.Parameters)
	if logCfg != nil && logCfg.FieldLogLevel != "" && !validateFieldLogLevel(logCfg.FieldLogLevel) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid logConfig.fieldLogLevel: %s", logCfg.FieldLogLevel))
	}

	// Fetch existing to preserve fields that were not provided in the request.
	// Without this, WafWebAclArn and XrayEnabled would be overwritten with
	// Go zero values on every update call that omits them.
	existing, err := store.GetGraphqlApiById(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	wafWebAclArn := existing.WafWebAclArn
	if request.HasParam(req.Parameters, "wafWebAclArn") {
		wafWebAclArn = request.GetStringParam(req.Parameters, "wafWebAclArn")
	}

	xrayEnabled := existing.XrayEnabled
	if request.HasParam(req.Parameters, "xrayEnabled") {
		xrayEnabled = request.GetBoolParam(req.Parameters, "xrayEnabled")
	}

	api := &appsyncstore.GraphqlApi{
		Name:                              name,
		AuthenticationType:                authType,
		AdditionalAuthenticationProviders: parseAdditionalAuthProviders(req.Parameters),
		EnhancedMetricsConfig:             parseEnhancedMetricsConfig(req.Parameters),
		IntrospectionConfig:               introspectionConfig,
		LambdaAuthorizerConfig:            parseLambdaAuthorizerConfig(req.Parameters),
		LogConfig:                         logCfg,
		MergedApiExecutionRoleArn:         request.GetStringParam(req.Parameters, "mergedApiExecutionRoleArn"),
		OpenIDConnectConfig:               parseOpenIDConnectConfig(req.Parameters),
		OwnerContact:                      request.GetStringParam(req.Parameters, "ownerContact"),
		QueryDepthLimit:                   int32(request.GetIntParam(req.Parameters, "queryDepthLimit")),
		ResolverCountLimit:                int32(request.GetIntParam(req.Parameters, "resolverCountLimit")),
		UserPoolConfig:                    parseUserPoolConfig(req.Parameters),
		WafWebAclArn:                      wafWebAclArn,
		XrayEnabled:                       xrayEnabled,
	}

	updated, err := store.UpdateGraphqlApiById(apiId, api)
	if err != nil {
		return mapStoreError(err)
	}

	result := graphqlApiToMap(updated)
	if tags, err := store.TagStore.List(updated.Arn); err == nil && len(tags) > 0 {
		result["tags"] = tags
	}

	return map[string]interface{}{
		"graphqlApi": result,
	}, nil
}

// DeleteGraphqlApi deletes a GraphQL API (v1).
// DELETE /v1/apis/{apiId}
func (s *AppSyncService) DeleteGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	if err := store.DeleteGraphqlApiById(apiId); err != nil {
		return mapStoreError(err)
	}

	s.schemaCache.Delete(apiId)

	return map[string]interface{}{}, nil
}

// ListGraphqlApis returns a paginated list of GraphQL APIs (v1).
// GET /v1/apis?apiType=GRAPHQL&maxResults=25&nextToken=...
func (s *AppSyncService) ListGraphqlApis(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	opts := parsePaginationOptions(req)
	apiTypeFilter := request.GetStringParam(req.Parameters, "apiType")
	apis, nextToken, err := store.ListGraphqlApis(opts, apiTypeFilter)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(apis))
	for _, api := range apis {
		item := graphqlApiToMap(api)
		if tags, err := store.TagStore.List(api.Arn); err == nil && len(tags) > 0 {
			item["tags"] = tags
		}
		items = append(items, item)
	}

	response := map[string]interface{}{
		"graphqlApis": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}
