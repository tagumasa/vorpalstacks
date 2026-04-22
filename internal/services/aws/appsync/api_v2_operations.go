package appsync

import (
	"context"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateApi creates a new Event API (v2).
// POST /v2/apis
func (s *AppSyncService) CreateApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	eventConfig, err := parseEventConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	api := &appsyncstore.Api{
		Name:         name,
		EventConfig:  eventConfig,
		OwnerContact: request.GetStringParam(req.Parameters, "ownerContact"),
		Tags:         parseTags(req.Parameters),
		WafWebAclArn: request.GetStringParam(req.Parameters, "wafWebAclArn"),
		XrayEnabled:  request.GetBoolParam(req.Parameters, "xrayEnabled"),
	}

	created, err := store.CreateApi(api)
	if err != nil {
		return mapStoreError(err)
	}

	if len(api.Tags) > 0 {
		tagMap := make(map[string]string, len(api.Tags))
		for k, v := range api.Tags {
			tagMap[k] = v
		}
		if err := store.TagStore.Tag(created.Arn, tagMap); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"api": apiToMap(created),
	}, nil
}

// GetApi retrieves an Event API (v2) by its ID.
// GET /v2/apis/{apiId}
func (s *AppSyncService) GetApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	api, err := store.GetApiById(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	result := apiToMap(api)
	if tagsFromStore, err := store.TagStore.List(api.Arn); err == nil && len(tagsFromStore) > 0 {
		result["tags"] = tagsFromStore
	}

	return map[string]interface{}{
		"api": result,
	}, nil
}

// UpdateApi updates an existing Event API (v2).
// POST /v2/apis/{apiId}
func (s *AppSyncService) UpdateApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	api := &appsyncstore.Api{
		Name:         request.GetStringParam(req.Parameters, "name"),
		OwnerContact: request.GetStringParam(req.Parameters, "ownerContact"),
		WafWebAclArn: request.GetStringParam(req.Parameters, "wafWebAclArn"),
		XrayEnabled:  request.GetBoolParam(req.Parameters, "xrayEnabled"),
	}

	if eventConfig, err := parseEventConfig(req.Parameters); err == nil {
		api.EventConfig = eventConfig
	} else if request.HasParam(req.Parameters, "eventConfig") {
		return nil, err
	}

	updated, err := store.UpdateApiById(apiId, api)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"api": apiToMap(updated),
	}, nil
}

// DeleteApi deletes an Event API (v2).
// DELETE /v2/apis/{apiId}
func (s *AppSyncService) DeleteApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	api, err := store.GetApiById(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	if err := store.DeleteApiById(apiId); err != nil {
		return mapStoreError(err)
	}

	store.TagStore.Delete(api.Arn)

	s.eventServer.DisconnectByApiId(apiId)

	return map[string]interface{}{}, nil
}

// ListApis returns a paginated list of Event APIs (v2).
// GET /v2/apis
func (s *AppSyncService) ListApis(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parsePaginationOptions(req)
	apis, nextToken, err := store.ListApis(opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(apis))
	for _, api := range apis {
		items = append(items, apiToMap(api))
	}

	response := map[string]interface{}{
		"apis": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

// mapStoreError converts a store-level error to the corresponding AppSync service error.
func mapStoreError(err error) (interface{}, error) {
	switch err {
	case appsyncstore.ErrApiNotFound:
		return nil, NewNotFoundException("API")
	case appsyncstore.ErrApiAlreadyExists:
		return nil, NewConflictException("API already exists")
	case appsyncstore.ErrChannelNamespaceNotFound:
		return nil, NewNotFoundException("Channel namespace")
	case appsyncstore.ErrChannelNamespaceExists:
		return nil, NewConflictException("Channel namespace already exists")
	case appsyncstore.ErrGraphqlApiNotFound:
		return nil, NewNotFoundException("GraphQL API")
	case appsyncstore.ErrGraphqlApiAlreadyExists:
		return nil, NewConflictException("GraphQL API already exists")
	case appsyncstore.ErrDataSourceNotFound:
		return nil, NewNotFoundException("Data source")
	case appsyncstore.ErrDataSourceAlreadyExists:
		return nil, NewConflictException("Data source already exists")
	case appsyncstore.ErrResolverNotFound:
		return nil, NewNotFoundException("Resolver")
	case appsyncstore.ErrResolverAlreadyExists:
		return nil, NewConflictException("Resolver already exists")
	case appsyncstore.ErrFunctionNotFound:
		return nil, NewNotFoundException("Function")
	case appsyncstore.ErrFunctionAlreadyExists:
		return nil, NewConflictException("Function already exists")
	case appsyncstore.ErrTypeNotFound:
		return nil, NewNotFoundException("Type")
	case appsyncstore.ErrTypeAlreadyExists:
		return nil, NewConflictException("Type already exists")
	case appsyncstore.ErrApiKeyNotFound:
		return nil, NewNotFoundException("API key")
	case appsyncstore.ErrApiCacheNotFound:
		return nil, NewNotFoundException("API cache")
	case appsyncstore.ErrApiCacheAlreadyExists:
		return nil, NewConflictException("API cache already exists")
	case appsyncstore.ErrDomainNameNotFound:
		return nil, NewNotFoundException("Domain name")
	case appsyncstore.ErrDomainNameAlreadyExists:
		return nil, NewConflictException("Domain name already exists")
	case appsyncstore.ErrApiAssociationNotFound:
		return nil, NewNotFoundException("API association")
	case appsyncstore.ErrMergedApiAssociationNotFound:
		return nil, NewNotFoundException("Merged API association")
	default:
		return nil, ErrInternalFailureException
	}
}

// apiToMap converts an Api struct to a response map with correct wire format.
// Timestamps are serialised as epoch seconds (float64) per REST-JSON 1.0 protocol.
func apiToMap(api *appsyncstore.Api) map[string]interface{} {
	m := map[string]interface{}{
		"apiId":       api.ApiId,
		"name":        api.Name,
		"apiArn":      api.Arn,
		"created":     timeutils.FormatEpochSeconds(api.Created),
		"xrayEnabled": api.XrayEnabled,
	}

	if len(api.Tags) > 0 {
		m["tags"] = api.Tags
	}
	if api.OwnerContact != "" {
		m["ownerContact"] = api.OwnerContact
	}
	if api.WafWebAclArn != "" {
		m["wafWebAclArn"] = api.WafWebAclArn
	}
	if len(api.Dns) > 0 {
		m["dns"] = api.Dns
	}
	if api.EventConfig != nil {
		m["eventConfig"] = eventConfigToMap(api.EventConfig)
	}

	return m
}

// eventConfigToMap converts an EventConfig to a wire-format map.
func eventConfigToMap(ec *appsyncstore.EventConfig) map[string]interface{} {
	m := map[string]interface{}{}

	if len(ec.AuthProviders) > 0 {
		providers := make([]interface{}, 0, len(ec.AuthProviders))
		for _, ap := range ec.AuthProviders {
			providers = append(providers, authProviderToMap(&ap))
		}
		m["authProviders"] = providers
	}
	if len(ec.ConnectionAuthModes) > 0 {
		m["connectionAuthModes"] = authModesToMap(ec.ConnectionAuthModes)
	}
	if len(ec.DefaultPublishAuthModes) > 0 {
		m["defaultPublishAuthModes"] = authModesToMap(ec.DefaultPublishAuthModes)
	}
	if len(ec.DefaultSubscribeAuthModes) > 0 {
		m["defaultSubscribeAuthModes"] = authModesToMap(ec.DefaultSubscribeAuthModes)
	}
	if ec.LogConfig != nil {
		m["logConfig"] = map[string]interface{}{
			"cloudWatchLogsRoleArn": ec.LogConfig.CloudWatchLogsRoleArn,
			"logLevel":              ec.LogConfig.LogLevel,
		}
	}

	return m
}

// authProviderToMap converts an AuthProvider to a wire-format map.
func authProviderToMap(ap *appsyncstore.AuthProvider) map[string]interface{} {
	m := map[string]interface{}{
		"authType": ap.AuthType,
	}
	if ap.CognitoConfig != nil {
		m["cognitoConfig"] = map[string]interface{}{
			"awsRegion":        ap.CognitoConfig.AwsRegion,
			"userPoolId":       ap.CognitoConfig.UserPoolId,
			"appIdClientRegex": ap.CognitoConfig.AppIdClientRegex,
		}
	}
	if ap.LambdaAuthorizerConfig != nil {
		lc := map[string]interface{}{
			"authorizerUri": ap.LambdaAuthorizerConfig.AuthorizerUri,
		}
		if ap.LambdaAuthorizerConfig.AuthorizerResultTtlInSeconds > 0 {
			lc["authorizerResultTtlInSeconds"] = ap.LambdaAuthorizerConfig.AuthorizerResultTtlInSeconds
		}
		if ap.LambdaAuthorizerConfig.IdentityValidationExpression != "" {
			lc["identityValidationExpression"] = ap.LambdaAuthorizerConfig.IdentityValidationExpression
		}
		m["lambdaAuthorizerConfig"] = lc
	}
	if ap.OpenIDConnectConfig != nil {
		oidc := map[string]interface{}{
			"issuer": ap.OpenIDConnectConfig.Issuer,
		}
		if ap.OpenIDConnectConfig.AuthTTL > 0 {
			oidc["authTTL"] = ap.OpenIDConnectConfig.AuthTTL
		}
		if ap.OpenIDConnectConfig.ClientId != "" {
			oidc["clientId"] = ap.OpenIDConnectConfig.ClientId
		}
		if ap.OpenIDConnectConfig.IatTTL > 0 {
			oidc["iatTTL"] = ap.OpenIDConnectConfig.IatTTL
		}
		m["openIDConnectConfig"] = oidc
	}
	return m
}

// authModesToMap converts a slice of AuthMode to a wire-format slice.
func authModesToMap(modes []appsyncstore.AuthMode) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(modes))
	for _, mode := range modes {
		result = append(result, map[string]interface{}{
			"authType": mode.AuthType,
		})
	}
	return result
}
