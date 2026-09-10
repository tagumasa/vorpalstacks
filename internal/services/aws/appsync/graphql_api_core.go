package appsync

import (
	"fmt"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// createGraphqlApiInput is the transport-agnostic input for creating a
// GraphQL API from either the admin console or the AWS API layer.
type createGraphqlApiInput struct {
	Name                              string
	AuthenticationType                string
	AdditionalAuthenticationProviders []appsyncstore.AdditionalAuthenticationProvider
	ApiType                           string
	EnhancedMetricsConfig             *appsyncstore.EnhancedMetricsConfig
	IntrospectionConfig               string
	LambdaAuthorizerConfig            *appsyncstore.LambdaAuthorizerConfig
	LogConfig                         *appsyncstore.LogConfig
	MergedApiExecutionRoleArn         string
	OpenIDConnectConfig               *appsyncstore.OpenIDConnectConfig
	OwnerContact                      string
	QueryDepthLimit                   int32
	HasQueryDepthLimit                bool
	ResolverCountLimit                int32
	HasResolverCountLimit             bool
	Tags                              map[string]string
	UserPoolConfig                    *appsyncstore.UserPoolConfig
	Visibility                        string
	WafWebAclArn                      string
	XrayEnabled                       bool
}

// createGraphqlApiCore creates a GraphQL API via the store layer. The
// returned tags are the tag-store view read back after creation (nil when
// empty).
func (s *AppSyncService) createGraphqlApiCore(store *appsyncstore.AppSyncStore, in createGraphqlApiInput) (*appsyncstore.GraphqlApi, map[string]string, error) {
	if in.Name == "" {
		return nil, nil, NewBadRequestException("name is required")
	}
	if err := validateApiName(in.Name); err != nil {
		return nil, nil, err
	}
	if in.AuthenticationType == "" {
		return nil, nil, NewBadRequestException("authenticationType is required")
	}
	if !validateAuthenticationType(in.AuthenticationType) {
		return nil, nil, NewBadRequestException(fmt.Sprintf("Invalid authenticationType: %s", in.AuthenticationType))
	}
	if in.ApiType != "" && !validateApiType(in.ApiType) {
		return nil, nil, NewBadRequestException(fmt.Sprintf("Invalid apiType: %s", in.ApiType))
	}
	if in.Visibility != "" && !validateVisibility(in.Visibility) {
		return nil, nil, NewBadRequestException(fmt.Sprintf("Invalid visibility: %s", in.Visibility))
	}
	if in.IntrospectionConfig != "" && !validateIntrospectionConfig(in.IntrospectionConfig) {
		return nil, nil, NewBadRequestException(fmt.Sprintf("Invalid introspectionConfig: %s", in.IntrospectionConfig))
	}
	if in.LogConfig != nil && in.LogConfig.FieldLogLevel != "" && !validateFieldLogLevel(in.LogConfig.FieldLogLevel) {
		return nil, nil, NewBadRequestException(fmt.Sprintf("Invalid logConfig.fieldLogLevel: %s", in.LogConfig.FieldLogLevel))
	}
	if in.HasQueryDepthLimit {
		if err := validateQueryDepthLimit(in.QueryDepthLimit); err != nil {
			return nil, nil, err
		}
	}
	if in.HasResolverCountLimit {
		if err := validateResolverCountLimit(in.ResolverCountLimit); err != nil {
			return nil, nil, err
		}
	}

	api := &appsyncstore.GraphqlApi{
		Name:                              in.Name,
		AuthenticationType:                in.AuthenticationType,
		AdditionalAuthenticationProviders: in.AdditionalAuthenticationProviders,
		ApiType:                           in.ApiType,
		EnhancedMetricsConfig:             in.EnhancedMetricsConfig,
		IntrospectionConfig:               in.IntrospectionConfig,
		LambdaAuthorizerConfig:            in.LambdaAuthorizerConfig,
		LogConfig:                         in.LogConfig,
		MergedApiExecutionRoleArn:         in.MergedApiExecutionRoleArn,
		OpenIDConnectConfig:               in.OpenIDConnectConfig,
		OwnerContact:                      in.OwnerContact,
		QueryDepthLimit:                   in.QueryDepthLimit,
		QueryDepthLimitSet:                in.HasQueryDepthLimit,
		ResolverCountLimit:                in.ResolverCountLimit,
		ResolverCountLimitSet:             in.HasResolverCountLimit,
		UserPoolConfig:                    in.UserPoolConfig,
		Visibility:                        in.Visibility,
		WafWebAclArn:                      in.WafWebAclArn,
		XrayEnabled:                       in.XrayEnabled,
	}

	created, err := store.CreateGraphqlApi(api)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	if err := applyCreateTags(store, created.Arn, in.Tags); err != nil {
		return nil, nil, err
	}

	return created, listTagsIfAny(store, created.Arn), nil
}

// deleteGraphqlApiCore deletes a GraphQL API by ID.
func (s *AppSyncService) deleteGraphqlApiCore(store *appsyncstore.AppSyncStore, apiID string) error {
	if apiID == "" {
		return NewBadRequestException("apiId is required")
	}
	if _, err := store.GetGraphqlApiById(apiID); err != nil {
		return mapStoreErrorE(err)
	}
	if err := store.DeleteGraphqlApiById(apiID); err != nil {
		return mapStoreErrorE(err)
	}
	s.schemaCache.Delete(apiID)
	return nil
}

// getGraphqlApiCore fetches a GraphQL API (v1) by ID together with its tags.
func (s *AppSyncService) getGraphqlApiCore(store *appsyncstore.AppSyncStore, apiId string) (*appsyncstore.GraphqlApi, map[string]string, error) {
	if apiId == "" {
		return nil, nil, NewBadRequestException("apiId is required")
	}

	api, err := store.GetGraphqlApiById(apiId)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	return api, listTagsIfAny(store, api.Arn), nil
}

// updateGraphqlApiInput carries the parsed UpdateGraphqlApi (v1) request
// payload. The Has* flags distinguish explicitly supplied members from
// omitted ones.
type updateGraphqlApiInput struct {
	ApiId                             string
	Name                              string
	AuthenticationType                string
	AdditionalAuthenticationProviders []appsyncstore.AdditionalAuthenticationProvider
	EnhancedMetricsConfig             *appsyncstore.EnhancedMetricsConfig
	IntrospectionConfig               string
	LambdaAuthorizerConfig            *appsyncstore.LambdaAuthorizerConfig
	LogConfig                         *appsyncstore.LogConfig
	MergedApiExecutionRoleArn         string
	OpenIDConnectConfig               *appsyncstore.OpenIDConnectConfig
	OwnerContact                      string
	QueryDepthLimit                   int32
	HasQueryDepthLimit                bool
	ResolverCountLimit                int32
	HasResolverCountLimit             bool
	UserPoolConfig                    *appsyncstore.UserPoolConfig
	WafWebAclArn                      string
	HasWafWebAclArn                   bool
	XrayEnabled                       bool
	HasXrayEnabled                    bool
}

// updateGraphqlApiCore applies an update to an existing GraphQL API (v1),
// preserving members that were not provided in the request.
func (s *AppSyncService) updateGraphqlApiCore(store *appsyncstore.AppSyncStore, in updateGraphqlApiInput) (*appsyncstore.GraphqlApi, map[string]string, error) {
	if in.ApiId == "" {
		return nil, nil, NewBadRequestException("apiId is required")
	}

	// Per Smithy model, name and authenticationType are required for
	// UpdateGraphqlApiRequest. The client must resend current values.
	if in.Name == "" {
		return nil, nil, NewBadRequestException("name is required")
	}
	if err := validateApiName(in.Name); err != nil {
		return nil, nil, err
	}
	if in.AuthenticationType == "" {
		return nil, nil, NewBadRequestException("authenticationType is required")
	}
	if !validateAuthenticationType(in.AuthenticationType) {
		return nil, nil, NewBadRequestException(fmt.Sprintf("Invalid authenticationType: %s", in.AuthenticationType))
	}

	if in.IntrospectionConfig != "" && !validateIntrospectionConfig(in.IntrospectionConfig) {
		return nil, nil, NewBadRequestException(fmt.Sprintf("Invalid introspectionConfig: %s", in.IntrospectionConfig))
	}

	if in.LogConfig != nil && in.LogConfig.FieldLogLevel != "" && !validateFieldLogLevel(in.LogConfig.FieldLogLevel) {
		return nil, nil, NewBadRequestException(fmt.Sprintf("Invalid logConfig.fieldLogLevel: %s", in.LogConfig.FieldLogLevel))
	}

	// Fetch existing to preserve fields that were not provided in the request.
	// Without this, WafWebAclArn and XrayEnabled would be overwritten with
	// Go zero values on every update call that omits them.
	existing, err := store.GetGraphqlApiById(in.ApiId)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	wafWebAclArn := existing.WafWebAclArn
	if in.HasWafWebAclArn {
		wafWebAclArn = in.WafWebAclArn
	}

	xrayEnabled := existing.XrayEnabled
	if in.HasXrayEnabled {
		xrayEnabled = in.XrayEnabled
	}

	if in.HasQueryDepthLimit {
		if err := validateQueryDepthLimit(in.QueryDepthLimit); err != nil {
			return nil, nil, err
		}
	}
	if in.HasResolverCountLimit {
		if err := validateResolverCountLimit(in.ResolverCountLimit); err != nil {
			return nil, nil, err
		}
	}

	api := &appsyncstore.GraphqlApi{
		Name:                              in.Name,
		AuthenticationType:                in.AuthenticationType,
		AdditionalAuthenticationProviders: in.AdditionalAuthenticationProviders,
		EnhancedMetricsConfig:             in.EnhancedMetricsConfig,
		IntrospectionConfig:               in.IntrospectionConfig,
		LambdaAuthorizerConfig:            in.LambdaAuthorizerConfig,
		LogConfig:                         in.LogConfig,
		MergedApiExecutionRoleArn:         in.MergedApiExecutionRoleArn,
		OpenIDConnectConfig:               in.OpenIDConnectConfig,
		OwnerContact:                      in.OwnerContact,
		QueryDepthLimit:                   in.QueryDepthLimit,
		QueryDepthLimitSet:                in.HasQueryDepthLimit,
		ResolverCountLimit:                in.ResolverCountLimit,
		ResolverCountLimitSet:             in.HasResolverCountLimit,
		UserPoolConfig:                    in.UserPoolConfig,
		WafWebAclArn:                      wafWebAclArn,
		XrayEnabled:                       xrayEnabled,
	}

	updated, err := store.UpdateGraphqlApiById(in.ApiId, api)
	if err != nil {
		return nil, nil, mapStoreErrorE(err)
	}

	return updated, listTagsIfAny(store, updated.Arn), nil
}

// graphqlApiWithTags pairs a listed GraphQL API with its tag-store view.
type graphqlApiWithTags struct {
	Api  *appsyncstore.GraphqlApi
	Tags map[string]string
}

// listGraphqlApisCore lists GraphQL APIs with pagination, enriching each
// entry with its tag-store view. ownerFilter carries the Ownership enum of
// the list request: every API in a regional store belongs to that account,
// so CURRENT_ACCOUNT is the unfiltered view, while OTHER_ACCOUNTS — APIs
// shared into the account — has no substrate and lists nothing.
func (s *AppSyncService) listGraphqlApisCore(store *appsyncstore.AppSyncStore, maxResults int, nextToken string, apiTypeFilter, ownerFilter string) ([]graphqlApiWithTags, string, error) {
	switch ownerFilter {
	case "", "CURRENT_ACCOUNT":
	case "OTHER_ACCOUNTS":
		return []graphqlApiWithTags{}, "", nil
	default:
		return nil, "", NewBadRequestException(fmt.Sprintf("Invalid owner: %s", ownerFilter))
	}

	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}
	apis, nextToken, err := store.ListGraphqlApis(opts, apiTypeFilter)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}
	entries := make([]graphqlApiWithTags, 0, len(apis))
	for _, api := range apis {
		entries = append(entries, graphqlApiWithTags{Api: api, Tags: listTagsIfAny(store, api.Arn)})
	}
	return entries, nextToken, nil
}
