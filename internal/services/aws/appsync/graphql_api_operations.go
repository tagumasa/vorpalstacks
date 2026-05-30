package appsync

import (
	"context"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

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

	api := &appsyncstore.GraphqlApi{
		Name:                              name,
		AuthenticationType:                authType,
		AdditionalAuthenticationProviders: parseAdditionalAuthProviders(req.Parameters),
		ApiType:                           request.GetStringParam(req.Parameters, "apiType"),
		EnhancedMetricsConfig:             parseEnhancedMetricsConfig(req.Parameters),
		IntrospectionConfig:               request.GetStringParam(req.Parameters, "introspectionConfig"),
		LambdaAuthorizerConfig:            parseLambdaAuthorizerConfig(req.Parameters),
		LogConfig:                         parseLogConfig(req.Parameters),
		MergedApiExecutionRoleArn:         request.GetStringParam(req.Parameters, "mergedApiExecutionRoleArn"),
		OpenIDConnectConfig:               parseOpenIDConnectConfig(req.Parameters),
		OwnerContact:                      request.GetStringParam(req.Parameters, "ownerContact"),
		QueryDepthLimit:                   int32(request.GetIntParam(req.Parameters, "queryDepthLimit")),
		ResolverCountLimit:                int32(request.GetIntParam(req.Parameters, "resolverCountLimit")),
		Tags:                              parseTags(req.Parameters),
		UserPoolConfig:                    parseUserPoolConfig(req.Parameters),
		Visibility:                        request.GetStringParam(req.Parameters, "visibility"),
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

	return map[string]interface{}{
		"graphqlApi": graphqlApiToMap(created),
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

	// In AWS, name and authenticationType are optional for updates;
	// only non-empty values are applied via the store's merge logic.
	name := request.GetStringParam(req.Parameters, "name")
	authType := request.GetStringParam(req.Parameters, "authenticationType")

	api := &appsyncstore.GraphqlApi{
		Name:                              name,
		AuthenticationType:                authType,
		AdditionalAuthenticationProviders: parseAdditionalAuthProviders(req.Parameters),
		EnhancedMetricsConfig:             parseEnhancedMetricsConfig(req.Parameters),
		IntrospectionConfig:               request.GetStringParam(req.Parameters, "introspectionConfig"),
		LambdaAuthorizerConfig:            parseLambdaAuthorizerConfig(req.Parameters),
		LogConfig:                         parseLogConfig(req.Parameters),
		MergedApiExecutionRoleArn:         request.GetStringParam(req.Parameters, "mergedApiExecutionRoleArn"),
		OpenIDConnectConfig:               parseOpenIDConnectConfig(req.Parameters),
		OwnerContact:                      request.GetStringParam(req.Parameters, "ownerContact"),
		QueryDepthLimit:                   int32(request.GetIntParam(req.Parameters, "queryDepthLimit")),
		ResolverCountLimit:                int32(request.GetIntParam(req.Parameters, "resolverCountLimit")),
		UserPoolConfig:                    parseUserPoolConfig(req.Parameters),
		WafWebAclArn:                      request.GetStringParam(req.Parameters, "wafWebAclArn"),
	}
	if request.HasParam(req.Parameters, "xrayEnabled") {
		api.XrayEnabled = request.GetBoolParam(req.Parameters, "xrayEnabled")
	}

	updated, err := store.UpdateGraphqlApiById(apiId, api)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"graphqlApi": graphqlApiToMap(updated),
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

	api, err := store.GetGraphqlApiById(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	if err := store.DeleteGraphqlApiById(apiId); err != nil {
		return mapStoreError(err)
	}

	store.TagStore.Delete(api.Arn)

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
		items = append(items, graphqlApiToMap(api))
	}

	response := map[string]interface{}{
		"graphqlApis": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}
