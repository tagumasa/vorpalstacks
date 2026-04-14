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
		return nil, err
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

	return map[string]interface{}{
		"graphqlApi": graphqlApiToMap(created),
	}, nil
}

// GetGraphqlApi retrieves a GraphQL API (v1) by its ID.
// GET /v1/apis/{apiId}
func (s *AppSyncService) GetGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	api, err := store.GetGraphqlApiById(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"graphqlApi": graphqlApiToMap(api),
	}, nil
}

// UpdateGraphqlApi updates an existing GraphQL API (v1).
// POST /v1/apis/{apiId}
func (s *AppSyncService) UpdateGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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
		XrayEnabled:                       request.GetBoolParam(req.Parameters, "xrayEnabled"),
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
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
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
		return nil, err
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

// graphqlApiToMap converts a GraphqlApi struct to a response map with correct wire format.
// Timestamps are serialised as epoch seconds (float64) per REST-JSON 1.0 protocol.
func graphqlApiToMap(api *appsyncstore.GraphqlApi) map[string]interface{} {
	m := map[string]interface{}{
		"apiId":              api.ApiId,
		"name":               api.Name,
		"arn":                api.Arn,
		"authenticationType": api.AuthenticationType,
		"xrayEnabled":        api.XrayEnabled,
	}

	if api.ApiType != "" {
		m["apiType"] = api.ApiType
	}
	if len(api.AdditionalAuthenticationProviders) > 0 {
		providers := make([]interface{}, 0, len(api.AdditionalAuthenticationProviders))
		for _, p := range api.AdditionalAuthenticationProviders {
			providers = append(providers, additionalAuthProviderToMap(&p))
		}
		m["additionalAuthenticationProviders"] = providers
	}
	if len(api.Dns) > 0 {
		m["dns"] = api.Dns
	}
	if api.EnhancedMetricsConfig != nil {
		m["enhancedMetricsConfig"] = enhancedMetricsConfigToMap(api.EnhancedMetricsConfig)
	}
	if api.IntrospectionConfig != "" {
		m["introspectionConfig"] = api.IntrospectionConfig
	}
	if api.LambdaAuthorizerConfig != nil {
		m["lambdaAuthorizerConfig"] = lambdaAuthorizerConfigToMap(api.LambdaAuthorizerConfig)
	}
	if api.LogConfig != nil {
		m["logConfig"] = logConfigToMap(api.LogConfig)
	}
	if api.MergedApiExecutionRoleArn != "" {
		m["mergedApiExecutionRoleArn"] = api.MergedApiExecutionRoleArn
	}
	if api.OpenIDConnectConfig != nil {
		m["openIDConnectConfig"] = openIDConnectConfigToMap(api.OpenIDConnectConfig)
	}
	if api.Owner != "" {
		m["owner"] = api.Owner
	}
	if api.OwnerContact != "" {
		m["ownerContact"] = api.OwnerContact
	}
	if api.QueryDepthLimit > 0 {
		m["queryDepthLimit"] = api.QueryDepthLimit
	}
	if api.ResolverCountLimit > 0 {
		m["resolverCountLimit"] = api.ResolverCountLimit
	}
	if len(api.Tags) > 0 {
		m["tags"] = api.Tags
	}
	if len(api.Uris) > 0 {
		m["uris"] = api.Uris
	}
	if api.UserPoolConfig != nil {
		m["userPoolConfig"] = userPoolConfigToMap(api.UserPoolConfig)
	}
	if api.Visibility != "" {
		m["visibility"] = api.Visibility
	}
	if api.WafWebAclArn != "" {
		m["wafWebAclArn"] = api.WafWebAclArn
	}

	return m
}

// parseAdditionalAuthProviders parses additional authentication providers from request parameters.
func parseAdditionalAuthProviders(params map[string]interface{}) []appsyncstore.AdditionalAuthenticationProvider {
	raw := request.GetArrayParam(params, "additionalAuthenticationProviders")
	var providers []appsyncstore.AdditionalAuthenticationProvider
	for _, p := range raw {
		if pMap, ok := p.(map[string]interface{}); ok {
			ap := appsyncstore.AdditionalAuthenticationProvider{
				AuthenticationType: request.GetStringParam(pMap, "authenticationType"),
			}
			if lambdaRaw := request.GetMapParam(pMap, "lambdaAuthorizerConfig"); lambdaRaw != nil {
				ap.LambdaAuthorizerConfig = parseLambdaAuthorizerConfigFromMap(lambdaRaw)
			}
			if oidcRaw := request.GetMapParam(pMap, "openIDConnectConfig"); oidcRaw != nil {
				ap.OpenIDConnectConfig = parseOpenIDConnectConfigFromMap(oidcRaw)
			}
			if cognitoRaw := request.GetMapParam(pMap, "userPoolConfig"); cognitoRaw != nil {
				ap.UserPoolConfig = &appsyncstore.CognitoUserPoolConfig{
					AwsRegion:        request.GetStringParam(cognitoRaw, "awsRegion"),
					UserPoolId:       request.GetStringParam(cognitoRaw, "userPoolId"),
					AppIdClientRegex: request.GetStringParam(cognitoRaw, "appIdClientRegex"),
				}
			}
			providers = append(providers, ap)
		}
	}
	return providers
}

// parseEnhancedMetricsConfig parses enhanced metrics configuration from request parameters.
func parseEnhancedMetricsConfig(params map[string]interface{}) *appsyncstore.EnhancedMetricsConfig {
	raw := request.GetMapParam(params, "enhancedMetricsConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.EnhancedMetricsConfig{
		DataSourceLevelMetricsBehavior: request.GetStringParam(raw, "dataSourceLevelMetricsBehavior"),
		OperationLevelMetricsConfig:    request.GetStringParam(raw, "operationLevelMetricsConfig"),
		ResolverLevelMetricsBehavior:   request.GetStringParam(raw, "resolverLevelMetricsBehavior"),
	}
}

// parseLambdaAuthorizerConfig parses a Lambda authorizer config from request parameters.
func parseLambdaAuthorizerConfig(params map[string]interface{}) *appsyncstore.LambdaAuthorizerConfig {
	raw := request.GetMapParam(params, "lambdaAuthorizerConfig")
	if raw == nil {
		return nil
	}
	return parseLambdaAuthorizerConfigFromMap(raw)
}

// parseLambdaAuthorizerConfigFromMap parses a Lambda authorizer config from a map.
func parseLambdaAuthorizerConfigFromMap(raw map[string]interface{}) *appsyncstore.LambdaAuthorizerConfig {
	return &appsyncstore.LambdaAuthorizerConfig{
		AuthorizerUri:                request.GetStringParam(raw, "authorizerUri"),
		AuthorizerResultTtlInSeconds: int32(request.GetIntParam(raw, "authorizerResultTtlInSeconds")),
		IdentityValidationExpression: request.GetStringParam(raw, "identityValidationExpression"),
	}
}

// parseLogConfig parses a LogConfig from request parameters.
func parseLogConfig(params map[string]interface{}) *appsyncstore.LogConfig {
	raw := request.GetMapParam(params, "logConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.LogConfig{
		CloudWatchLogsRoleArn: request.GetStringParam(raw, "cloudWatchLogsRoleArn"),
		FieldLogLevel:         request.GetStringParam(raw, "fieldLogLevel"),
		ExcludeVerboseContent: request.GetBoolParam(raw, "excludeVerboseContent"),
	}
}

// parseOpenIDConnectConfig parses an OpenID Connect config from request parameters.
func parseOpenIDConnectConfig(params map[string]interface{}) *appsyncstore.OpenIDConnectConfig {
	raw := request.GetMapParam(params, "openIDConnectConfig")
	if raw == nil {
		return nil
	}
	return parseOpenIDConnectConfigFromMap(raw)
}

// parseOpenIDConnectConfigFromMap parses an OpenID Connect config from a map.
func parseOpenIDConnectConfigFromMap(raw map[string]interface{}) *appsyncstore.OpenIDConnectConfig {
	return &appsyncstore.OpenIDConnectConfig{
		Issuer:   request.GetStringParam(raw, "issuer"),
		AuthTTL:  request.GetInt64Param(raw, "authTTL"),
		ClientId: request.GetStringParam(raw, "clientId"),
		IatTTL:   request.GetInt64Param(raw, "iatTTL"),
	}
}

// parseUserPoolConfig parses a UserPoolConfig from request parameters.
func parseUserPoolConfig(params map[string]interface{}) *appsyncstore.UserPoolConfig {
	raw := request.GetMapParam(params, "userPoolConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.UserPoolConfig{
		AwsRegion:        request.GetStringParam(raw, "awsRegion"),
		DefaultAction:    request.GetStringParam(raw, "defaultAction"),
		UserPoolId:       request.GetStringParam(raw, "userPoolId"),
		AppIdClientRegex: request.GetStringParam(raw, "appIdClientRegex"),
	}
}

// additionalAuthProviderToMap converts an AdditionalAuthenticationProvider to a wire-format map.
func additionalAuthProviderToMap(ap *appsyncstore.AdditionalAuthenticationProvider) map[string]interface{} {
	m := map[string]interface{}{
		"authenticationType": ap.AuthenticationType,
	}
	if ap.LambdaAuthorizerConfig != nil {
		m["lambdaAuthorizerConfig"] = lambdaAuthorizerConfigToMap(ap.LambdaAuthorizerConfig)
	}
	if ap.OpenIDConnectConfig != nil {
		m["openIDConnectConfig"] = openIDConnectConfigToMap(ap.OpenIDConnectConfig)
	}
	if ap.UserPoolConfig != nil {
		m["userPoolConfig"] = map[string]interface{}{
			"awsRegion":        ap.UserPoolConfig.AwsRegion,
			"userPoolId":       ap.UserPoolConfig.UserPoolId,
			"appIdClientRegex": ap.UserPoolConfig.AppIdClientRegex,
		}
	}
	return m
}

// enhancedMetricsConfigToMap converts an EnhancedMetricsConfig to a wire-format map.
func enhancedMetricsConfigToMap(c *appsyncstore.EnhancedMetricsConfig) map[string]interface{} {
	return map[string]interface{}{
		"dataSourceLevelMetricsBehavior": c.DataSourceLevelMetricsBehavior,
		"operationLevelMetricsConfig":    c.OperationLevelMetricsConfig,
		"resolverLevelMetricsBehavior":   c.ResolverLevelMetricsBehavior,
	}
}

// lambdaAuthorizerConfigToMap converts a LambdaAuthorizerConfig to a wire-format map.
func lambdaAuthorizerConfigToMap(c *appsyncstore.LambdaAuthorizerConfig) map[string]interface{} {
	m := map[string]interface{}{
		"authorizerUri": c.AuthorizerUri,
	}
	if c.AuthorizerResultTtlInSeconds > 0 {
		m["authorizerResultTtlInSeconds"] = c.AuthorizerResultTtlInSeconds
	}
	if c.IdentityValidationExpression != "" {
		m["identityValidationExpression"] = c.IdentityValidationExpression
	}
	return m
}

// logConfigToMap converts a LogConfig to a wire-format map.
func logConfigToMap(c *appsyncstore.LogConfig) map[string]interface{} {
	m := map[string]interface{}{
		"cloudWatchLogsRoleArn": c.CloudWatchLogsRoleArn,
		"fieldLogLevel":         c.FieldLogLevel,
	}
	if c.ExcludeVerboseContent {
		m["excludeVerboseContent"] = c.ExcludeVerboseContent
	}
	return m
}

// openIDConnectConfigToMap converts an OpenIDConnectConfig to a wire-format map.
func openIDConnectConfigToMap(c *appsyncstore.OpenIDConnectConfig) map[string]interface{} {
	m := map[string]interface{}{
		"issuer": c.Issuer,
	}
	if c.AuthTTL > 0 {
		m["authTTL"] = c.AuthTTL
	}
	if c.ClientId != "" {
		m["clientId"] = c.ClientId
	}
	if c.IatTTL > 0 {
		m["iatTTL"] = c.IatTTL
	}
	return m
}

// userPoolConfigToMap converts a UserPoolConfig to a wire-format map.
func userPoolConfigToMap(c *appsyncstore.UserPoolConfig) map[string]interface{} {
	m := map[string]interface{}{
		"awsRegion":     c.AwsRegion,
		"defaultAction": c.DefaultAction,
		"userPoolId":    c.UserPoolId,
	}
	if c.AppIdClientRegex != "" {
		m["appIdClientRegex"] = c.AppIdClientRegex
	}
	return m
}
