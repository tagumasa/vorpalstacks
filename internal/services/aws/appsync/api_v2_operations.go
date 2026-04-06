package appsync

import (
	"context"
	"fmt"
	"strings"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/utils/aws/arn"
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

	return map[string]interface{}{
		"api": apiToMap(api),
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

	if err := store.DeleteApiById(apiId); err != nil {
		return mapStoreError(err)
	}

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

// resolveResourceFromArn determines the resource type and ID from an ARN.
// Supports Event API (apis/), GraphQL API (apis/), ChannelNamespace, and DataSource ARNs.
// Since both v1 and v2 use apis/{apiId} format, the caller must check both stores.
func resolveResourceFromArn(resourceArn string) (resourceType string, region string, apiId string, name string, err error) {
	_, _, region, _, resource := arn.SplitARN(resourceArn)
	if resource == "" {
		return "", "", "", "", fmt.Errorf("invalid ARN: %s", resourceArn)
	}

	// Channel namespace ARN: .../apis/{apiId}/channelNamespaces/{name}
	if strings.Contains(resource, "/channelNamespaces/") {
		parts := strings.SplitN(resource, "/", 4)
		if len(parts) < 4 {
			return "", "", "", "", fmt.Errorf("cannot parse channel namespace ARN: %s", resourceArn)
		}
		return "channelNamespace", region, parts[1], parts[3], nil
	}

	// DataSource ARN: .../apis/{apiId}/datasources/{name}
	if strings.Contains(resource, "/datasources/") {
		parts := strings.SplitN(resource, "/", 4)
		if len(parts) < 4 {
			return "", "", "", "", fmt.Errorf("cannot parse data source ARN: %s", resourceArn)
		}
		return "dataSource", region, parts[1], parts[3], nil
	}

	// Event API and GraphQL API ARN: .../apis/{apiId}
	// Both v1 and v2 use the same format; try v2 first, then v1.
	if strings.HasPrefix(resource, "apis/") {
		id := strings.TrimPrefix(resource, "apis/")
		if idx := strings.Index(id, "/"); idx >= 0 {
			id = id[:idx]
		}
		if id == "" {
			return "", "", "", "", fmt.Errorf("cannot extract ID from ARN: %s", resourceArn)
		}
		return "api", region, id, "", nil
	}

	return "", "", "", "", fmt.Errorf("unsupported ARN resource type: %s", resource)
}

// tagResource applies tag operations to either an API or a channel namespace.
func (s *AppSyncService) tagResource(store *appsyncstore.AppSyncStore, resourceArn string, newTags map[string]string) error {
	resourceType, _, apiId, name, err := resolveResourceFromArn(resourceArn)
	if err != nil {
		return NewNotFoundException(resourceArn)
	}

	switch resourceType {
	case "api":
		if err := s.tagApi(store, apiId, newTags); err == nil {
			return nil
		}
		return s.tagGraphqlApi(store, apiId, newTags)
	case "channelNamespace":
		return s.tagChannelNamespace(store, apiId, name, newTags)
	case "dataSource":
		return s.tagDataSource(store, apiId, name, newTags)
	default:
		return NewNotFoundException(resourceArn)
	}
}

// untagResource removes tag keys from a resource identified by ARN.
func (s *AppSyncService) untagResource(store *appsyncstore.AppSyncStore, resourceArn string, keysToRemove []string) error {
	resourceType, _, apiId, name, err := resolveResourceFromArn(resourceArn)
	if err != nil {
		return NewNotFoundException(resourceArn)
	}

	switch resourceType {
	case "api":
		if err := s.untagApi(store, apiId, keysToRemove); err == nil {
			return nil
		}
		return s.untagGraphqlApi(store, apiId, keysToRemove)
	case "channelNamespace":
		return s.untagChannelNamespace(store, apiId, name, keysToRemove)
	case "dataSource":
		return s.untagDataSource(store, apiId, name, keysToRemove)
	default:
		return NewNotFoundException(resourceArn)
	}
}

// listResourceTags returns the tags for a resource identified by ARN.
func (s *AppSyncService) listResourceTags(store *appsyncstore.AppSyncStore, resourceArn string) (map[string]string, error) {
	resourceType, _, apiId, name, err := resolveResourceFromArn(resourceArn)
	if err != nil {
		return nil, NewNotFoundException(resourceArn)
	}

	switch resourceType {
	case "api":
		tags, err := s.listApiTags(store, apiId)
		if err == nil {
			return tags, nil
		}
		return s.listGraphqlApiTags(store, apiId)
	case "channelNamespace":
		return s.listChannelNamespaceTags(store, apiId, name)
	case "dataSource":
		return s.listDataSourceTags(store, apiId, name)
	default:
		return nil, NewNotFoundException(resourceArn)
	}
}

// tagApi merges new tags into an existing Event API's tag set.
func (s *AppSyncService) tagApi(store *appsyncstore.AppSyncStore, apiId string, newTags map[string]string) error {
	api, err := store.GetApiById(apiId)
	if err != nil {
		if err == appsyncstore.ErrApiNotFound {
			return NewNotFoundException(fmt.Sprintf("API with ID %s", apiId))
		}
		return ErrInternalFailureException
	}

	if api.Tags == nil {
		api.Tags = make(map[string]string)
	}
	for k, v := range newTags {
		api.Tags[k] = v
	}

	_, err = store.UpdateApiById(apiId, api)
	return err
}

// untagApi removes specified keys from an Event API's tag set.
func (s *AppSyncService) untagApi(store *appsyncstore.AppSyncStore, apiId string, keysToRemove []string) error {
	api, err := store.GetApiById(apiId)
	if err != nil {
		if err == appsyncstore.ErrApiNotFound {
			return NewNotFoundException(fmt.Sprintf("API with ID %s", apiId))
		}
		return ErrInternalFailureException
	}

	if api.Tags == nil {
		return nil
	}
	for _, key := range keysToRemove {
		delete(api.Tags, key)
	}

	_, err = store.UpdateApiById(apiId, api)
	return err
}

// listApiTags returns the tag map for an Event API.
func (s *AppSyncService) listApiTags(store *appsyncstore.AppSyncStore, apiId string) (map[string]string, error) {
	api, err := store.GetApiById(apiId)
	if err != nil {
		if err == appsyncstore.ErrApiNotFound {
			return nil, NewNotFoundException(fmt.Sprintf("API with ID %s", apiId))
		}
		return nil, ErrInternalFailureException
	}

	if api.Tags == nil {
		return map[string]string{}, nil
	}
	return api.Tags, nil
}

// tagChannelNamespace merges new tags into an existing channel namespace.
func (s *AppSyncService) tagChannelNamespace(store *appsyncstore.AppSyncStore, apiId, name string, newTags map[string]string) error {
	return store.UpdateChannelNamespaceTags(apiId, name, func(tags map[string]string) map[string]string {
		if tags == nil {
			tags = make(map[string]string)
		}
		for k, v := range newTags {
			tags[k] = v
		}
		return tags
	})
}

// untagChannelNamespace removes specified keys from a channel namespace's tags.
func (s *AppSyncService) untagChannelNamespace(store *appsyncstore.AppSyncStore, apiId, name string, keysToRemove []string) error {
	return store.UpdateChannelNamespaceTags(apiId, name, func(tags map[string]string) map[string]string {
		if tags == nil {
			return nil
		}
		for _, key := range keysToRemove {
			delete(tags, key)
		}
		return tags
	})
}

// listChannelNamespaceTags returns the tag map for a channel namespace.
func (s *AppSyncService) listChannelNamespaceTags(store *appsyncstore.AppSyncStore, apiId, name string) (map[string]string, error) {
	ns, err := store.GetChannelNamespace(apiId, name)
	if err != nil {
		if err == appsyncstore.ErrChannelNamespaceNotFound {
			return nil, NewNotFoundException(fmt.Sprintf("Channel namespace %s", name))
		}
		return nil, ErrInternalFailureException
	}
	if ns.Tags == nil {
		return map[string]string{}, nil
	}
	return ns.Tags, nil
}

// tagGraphqlApi merges new tags into an existing GraphQL API's tag set.
func (s *AppSyncService) tagGraphqlApi(store *appsyncstore.AppSyncStore, apiId string, newTags map[string]string) error {
	api, err := store.GetGraphqlApiById(apiId)
	if err != nil {
		if err == appsyncstore.ErrGraphqlApiNotFound {
			return NewNotFoundException(fmt.Sprintf("GraphQL API with ID %s", apiId))
		}
		return ErrInternalFailureException
	}

	if api.Tags == nil {
		api.Tags = make(map[string]string)
	}
	for k, v := range newTags {
		api.Tags[k] = v
	}

	_, err = store.UpdateGraphqlApiById(apiId, api)
	return err
}

// untagGraphqlApi removes specified keys from a GraphQL API's tag set.
func (s *AppSyncService) untagGraphqlApi(store *appsyncstore.AppSyncStore, apiId string, keysToRemove []string) error {
	api, err := store.GetGraphqlApiById(apiId)
	if err != nil {
		if err == appsyncstore.ErrGraphqlApiNotFound {
			return NewNotFoundException(fmt.Sprintf("GraphQL API with ID %s", apiId))
		}
		return ErrInternalFailureException
	}

	if api.Tags == nil {
		return nil
	}
	for _, key := range keysToRemove {
		delete(api.Tags, key)
	}

	_, err = store.UpdateGraphqlApiById(apiId, api)
	return err
}

// listGraphqlApiTags returns the tag map for a GraphQL API.
func (s *AppSyncService) listGraphqlApiTags(store *appsyncstore.AppSyncStore, apiId string) (map[string]string, error) {
	api, err := store.GetGraphqlApiById(apiId)
	if err != nil {
		if err == appsyncstore.ErrGraphqlApiNotFound {
			return nil, NewNotFoundException(fmt.Sprintf("GraphQL API with ID %s", apiId))
		}
		return nil, ErrInternalFailureException
	}

	if api.Tags == nil {
		return map[string]string{}, nil
	}
	return api.Tags, nil
}

// tagDataSource merges new tags into an existing data source's tag set.
func (s *AppSyncService) tagDataSource(store *appsyncstore.AppSyncStore, apiId, name string, newTags map[string]string) error {
	ds, err := store.GetDataSource(apiId, name)
	if err != nil {
		if err == appsyncstore.ErrDataSourceNotFound {
			return NewNotFoundException(fmt.Sprintf("Data source %s", name))
		}
		return ErrInternalFailureException
	}

	if ds.Tags == nil {
		ds.Tags = make(map[string]string)
	}
	for k, v := range newTags {
		ds.Tags[k] = v
	}

	_, err = store.UpdateDataSource(ds)
	return err
}

// untagDataSource removes specified keys from a data source's tag set.
func (s *AppSyncService) untagDataSource(store *appsyncstore.AppSyncStore, apiId, name string, keysToRemove []string) error {
	ds, err := store.GetDataSource(apiId, name)
	if err != nil {
		if err == appsyncstore.ErrDataSourceNotFound {
			return NewNotFoundException(fmt.Sprintf("Data source %s", name))
		}
		return ErrInternalFailureException
	}

	if ds.Tags == nil {
		return nil
	}
	for _, key := range keysToRemove {
		delete(ds.Tags, key)
	}

	_, err = store.UpdateDataSource(ds)
	return err
}

// listDataSourceTags returns the tag map for a data source.
func (s *AppSyncService) listDataSourceTags(store *appsyncstore.AppSyncStore, apiId, name string) (map[string]string, error) {
	ds, err := store.GetDataSource(apiId, name)
	if err != nil {
		if err == appsyncstore.ErrDataSourceNotFound {
			return nil, NewNotFoundException(fmt.Sprintf("Data source %s", name))
		}
		return nil, ErrInternalFailureException
	}
	if ds.Tags == nil {
		return map[string]string{}, nil
	}
	return ds.Tags, nil
}
