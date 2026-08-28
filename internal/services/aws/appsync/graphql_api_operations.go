package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateGraphqlApi creates a new GraphQL API (v1).
// POST /v1/apis
func (s *AppSyncService) CreateGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	additionalAuthProviders, err := parseAdditionalAuthProviders(req.Parameters)
	if err != nil {
		return nil, err
	}
	lambdaAuthConfig, err := parseLambdaAuthorizerConfig(req.Parameters)
	if err != nil {
		return nil, err
	}
	tagMap, err := parseTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	queryDepthLimit := int32(request.GetIntParam(req.Parameters, "queryDepthLimit"))
	_, hasQueryDepthLimit := req.Parameters["queryDepthLimit"]
	resolverCountLimit := int32(request.GetIntParam(req.Parameters, "resolverCountLimit"))
	_, hasResolverCountLimit := req.Parameters["resolverCountLimit"]

	enhancedMetrics, err := parseEnhancedMetricsConfig(req.Parameters)
	if err != nil {
		return nil, err
	}
	oidcCfg, err := parseOpenIDConnectConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	created, tags, err := s.createGraphqlApiCore(store, createGraphqlApiInput{
		Name:                              request.GetStringParam(req.Parameters, "name"),
		AuthenticationType:                request.GetStringParam(req.Parameters, "authenticationType"),
		AdditionalAuthenticationProviders: additionalAuthProviders,
		ApiType:                           request.GetStringParam(req.Parameters, "apiType"),
		EnhancedMetricsConfig:             enhancedMetrics,
		IntrospectionConfig:               request.GetStringParam(req.Parameters, "introspectionConfig"),
		LambdaAuthorizerConfig:            lambdaAuthConfig,
		LogConfig:                         parseLogConfig(req.Parameters),
		MergedApiExecutionRoleArn:         request.GetStringParam(req.Parameters, "mergedApiExecutionRoleArn"),
		OpenIDConnectConfig:               oidcCfg,
		OwnerContact:                      request.GetStringParam(req.Parameters, "ownerContact"),
		QueryDepthLimit:                   queryDepthLimit,
		HasQueryDepthLimit:                hasQueryDepthLimit,
		ResolverCountLimit:                resolverCountLimit,
		HasResolverCountLimit:             hasResolverCountLimit,
		Tags:                              tagMap,
		UserPoolConfig:                    parseUserPoolConfig(req.Parameters),
		Visibility:                        request.GetStringParam(req.Parameters, "visibility"),
		WafWebAclArn:                      request.GetStringParam(req.Parameters, "wafWebAclArn"),
		XrayEnabled:                       request.GetBoolParam(req.Parameters, "xrayEnabled"),
	})
	if err != nil {
		return nil, err
	}

	result := graphqlApiToMap(created)
	if tags != nil {
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

	api, tags, err := s.getGraphqlApiCore(store, request.GetStringParam(req.Parameters, "apiId"))
	if err != nil {
		return nil, err
	}

	result := graphqlApiToMap(api)
	if tags != nil {
		result["tags"] = tags
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

	additionalAuthProviders, err := parseAdditionalAuthProviders(req.Parameters)
	if err != nil {
		return nil, err
	}
	lambdaAuthConfig, err := parseLambdaAuthorizerConfig(req.Parameters)
	if err != nil {
		return nil, err
	}
	enhancedMetrics, err := parseEnhancedMetricsConfig(req.Parameters)
	if err != nil {
		return nil, err
	}
	oidcCfg, err := parseOpenIDConnectConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	queryDepthLimit := int32(request.GetIntParam(req.Parameters, "queryDepthLimit"))
	_, hasQueryDepthLimit := req.Parameters["queryDepthLimit"]
	resolverCountLimit := int32(request.GetIntParam(req.Parameters, "resolverCountLimit"))
	_, hasResolverCountLimit := req.Parameters["resolverCountLimit"]

	updated, tags, err := s.updateGraphqlApiCore(store, updateGraphqlApiInput{
		ApiId:                             request.GetStringParam(req.Parameters, "apiId"),
		Name:                              request.GetStringParam(req.Parameters, "name"),
		AuthenticationType:                request.GetStringParam(req.Parameters, "authenticationType"),
		AdditionalAuthenticationProviders: additionalAuthProviders,
		EnhancedMetricsConfig:             enhancedMetrics,
		IntrospectionConfig:               request.GetStringParam(req.Parameters, "introspectionConfig"),
		LambdaAuthorizerConfig:            lambdaAuthConfig,
		LogConfig:                         parseLogConfig(req.Parameters),
		MergedApiExecutionRoleArn:         request.GetStringParam(req.Parameters, "mergedApiExecutionRoleArn"),
		OpenIDConnectConfig:               oidcCfg,
		OwnerContact:                      request.GetStringParam(req.Parameters, "ownerContact"),
		QueryDepthLimit:                   queryDepthLimit,
		HasQueryDepthLimit:                hasQueryDepthLimit,
		ResolverCountLimit:                resolverCountLimit,
		HasResolverCountLimit:             hasResolverCountLimit,
		UserPoolConfig:                    parseUserPoolConfig(req.Parameters),
		WafWebAclArn:                      request.GetStringParam(req.Parameters, "wafWebAclArn"),
		HasWafWebAclArn:                   request.HasParam(req.Parameters, "wafWebAclArn"),
		XrayEnabled:                       request.GetBoolParam(req.Parameters, "xrayEnabled"),
		HasXrayEnabled:                    request.HasParam(req.Parameters, "xrayEnabled"),
	})
	if err != nil {
		return nil, err
	}

	result := graphqlApiToMap(updated)
	if tags != nil {
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
	if err := s.deleteGraphqlApiCore(store, apiId); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListGraphqlApis returns a paginated list of GraphQL APIs (v1).
// GET /v1/apis?apiType=GRAPHQL&maxResults=25&nextToken=...
func (s *AppSyncService) ListGraphqlApis(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	opts, err := parsePaginationOptions(req)
	if err != nil {
		return nil, err
	}
	apiTypeFilter := request.GetStringParam(req.Parameters, "apiType")

	entries, nextToken, err := s.listGraphqlApisCore(store, int(opts.MaxItems), opts.Marker, apiTypeFilter)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(entries))
	for _, entry := range entries {
		item := graphqlApiToMap(entry.Api)
		if entry.Tags != nil {
			item["tags"] = entry.Tags
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
