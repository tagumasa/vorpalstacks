package appsync

import (
	"fmt"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
	storecommon "vorpalstacks/internal/store/aws/common"
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

// createGraphqlApiCore creates a GraphQL API via the store layer.
func (s *AppSyncService) createGraphqlApiCore(store *appsyncstore.AppSyncStore, in createGraphqlApiInput) (*appsyncstore.GraphqlApi, error) {
	if in.Name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if err := validateApiName(in.Name); err != nil {
		return nil, err
	}
	if in.AuthenticationType == "" {
		return nil, NewBadRequestException("authenticationType is required")
	}
	if !validateAuthenticationType(in.AuthenticationType) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid authenticationType: %s", in.AuthenticationType))
	}
	if in.ApiType != "" && !validateApiType(in.ApiType) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid apiType: %s", in.ApiType))
	}
	if in.Visibility != "" && !validateVisibility(in.Visibility) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid visibility: %s", in.Visibility))
	}
	if in.IntrospectionConfig != "" && !validateIntrospectionConfig(in.IntrospectionConfig) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid introspectionConfig: %s", in.IntrospectionConfig))
	}
	if in.LogConfig != nil && in.LogConfig.FieldLogLevel != "" && !validateFieldLogLevel(in.LogConfig.FieldLogLevel) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid logConfig.fieldLogLevel: %s", in.LogConfig.FieldLogLevel))
	}
	if in.HasQueryDepthLimit {
		if err := validateQueryDepthLimit(in.QueryDepthLimit); err != nil {
			return nil, err
		}
	}
	if in.HasResolverCountLimit {
		if err := validateResolverCountLimit(in.ResolverCountLimit); err != nil {
			return nil, err
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
		ResolverCountLimit:                in.ResolverCountLimit,
		Tags:                              in.Tags,
		UserPoolConfig:                    in.UserPoolConfig,
		Visibility:                        in.Visibility,
		WafWebAclArn:                      in.WafWebAclArn,
		XrayEnabled:                       in.XrayEnabled,
	}

	created, err := store.CreateGraphqlApi(api)
	if err != nil {
		return nil, mapStoreErrorE(err)
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

	return created, nil
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

// listApisCore lists APIs with pagination.
func (s *AppSyncService) listApisCore(store *appsyncstore.AppSyncStore, maxResults int, nextToken string) ([]*appsyncstore.Api, string, error) {
	if maxResults < 0 {
		maxResults = 0
	}
	if maxResults == 0 {
		maxResults = 25
	}
	if maxResults > 25 {
		return nil, "", NewBadRequestException("maxResults must be between 1 and 25")
	}
	return store.ListApis(storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   nextToken,
	})
}

// listGraphqlApisCore lists GraphQL APIs with pagination.
func (s *AppSyncService) listGraphqlApisCore(store *appsyncstore.AppSyncStore, maxResults int, nextToken string, apiTypeFilter string) ([]*appsyncstore.GraphqlApi, string, error) {
	if maxResults < 0 {
		maxResults = 0
	}
	if maxResults == 0 {
		maxResults = 25
	}
	if maxResults > 25 {
		return nil, "", NewBadRequestException("maxResults must be between 1 and 25")
	}
	return store.ListGraphqlApis(storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   nextToken,
	}, apiTypeFilter)
}
