package appsync

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
)

// --- GraphQL API (v1) ---

// CreateGraphqlApi persists a new GraphQL API (v1).
// Generates apiId, ARN, timestamps, and default URIs.
func (s *AppSyncStore) CreateGraphqlApi(api *GraphqlApi) (*GraphqlApi, error) {
	if api.Name == "" {
		return nil, common.NewStoreError("appsync", "create_graphql_api", common.ErrInvalidInput)
	}
	if api.AuthenticationType == "" {
		return nil, common.NewStoreError("appsync", "create_graphql_api", common.ErrInvalidInput)
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	if s.graphqlApisStore.Exists(api.Name) {
		return nil, ErrGraphqlApiAlreadyExists
	}

	api.ApiId = s.GenerateId()
	api.Arn = s.BuildGraphQLApiARN(api.ApiId)
	if api.ApiType == "" {
		api.ApiType = "GRAPHQL"
	}
	if api.Uris == nil {
		baseURL := config.BaseURL()
		wsBase := strings.Replace(baseURL, "http://", "ws://", 1)
		api.Uris = map[string]string{
			"GRAPHQL":  fmt.Sprintf("%s/v1/apis/%s/graphql", baseURL, api.ApiId),
			"REALTIME": fmt.Sprintf("%s/v1/apis/%s/realtime", wsBase, api.ApiId),
		}
	}
	if api.Dns == nil {
		baseURL := config.BaseURL()
		wsBase := strings.Replace(baseURL, "http://", "ws://", 1)
		api.Dns = map[string]string{
			"GRAPHQL":  fmt.Sprintf("%s/v1/apis/%s/graphql", baseURL, api.ApiId),
			"REALTIME": fmt.Sprintf("%s/v1/apis/%s/realtime", wsBase, api.ApiId),
		}
	}

	if err := s.graphqlApisStore.Put(api.Name, api); err != nil {
		return nil, err
	}
	return api, nil
}

// GetGraphqlApi retrieves a GraphQL API by name.
func (s *AppSyncStore) GetGraphqlApi(name string) (*GraphqlApi, error) {
	var api GraphqlApi
	if err := s.graphqlApisStore.Get(name, &api); err != nil {
		return nil, ErrGraphqlApiNotFound
	}
	return &api, nil
}

// GetGraphqlApiById retrieves a GraphQL API by its UUID.
// Scans all GraphQL APIs to find the one matching the given ID.
func (s *AppSyncStore) GetGraphqlApiById(apiId string) (*GraphqlApi, error) {
	apis, _, err := s.ListGraphqlApis(common.ListOptions{MaxItems: 10000}, "")
	if err != nil {
		return nil, err
	}
	for _, api := range apis {
		if api.ApiId == apiId {
			return api, nil
		}
	}
	return nil, ErrGraphqlApiNotFound
}

// UpdateGraphqlApiById updates a GraphQL API identified by apiId.
// Merges non-zero fields from the update.
func (s *AppSyncStore) UpdateGraphqlApiById(apiId string, update *GraphqlApi) (*GraphqlApi, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	existing, err := s.GetGraphqlApiById(apiId)
	if err != nil {
		return nil, err
	}

	oldName := existing.Name

	if update.Name != "" {
		existing.Name = update.Name
	}
	if update.AuthenticationType != "" {
		existing.AuthenticationType = update.AuthenticationType
	}
	if update.AdditionalAuthenticationProviders != nil {
		existing.AdditionalAuthenticationProviders = update.AdditionalAuthenticationProviders
	}
	if update.EnhancedMetricsConfig != nil {
		existing.EnhancedMetricsConfig = update.EnhancedMetricsConfig
	}
	if update.IntrospectionConfig != "" {
		existing.IntrospectionConfig = update.IntrospectionConfig
	}
	if update.LambdaAuthorizerConfig != nil {
		existing.LambdaAuthorizerConfig = update.LambdaAuthorizerConfig
	}
	if update.LogConfig != nil {
		existing.LogConfig = update.LogConfig
	}
	if update.MergedApiExecutionRoleArn != "" {
		existing.MergedApiExecutionRoleArn = update.MergedApiExecutionRoleArn
	}
	if update.OpenIDConnectConfig != nil {
		existing.OpenIDConnectConfig = update.OpenIDConnectConfig
	}
	if update.OwnerContact != "" {
		existing.OwnerContact = update.OwnerContact
	}
	if update.QueryDepthLimit > 0 {
		existing.QueryDepthLimit = update.QueryDepthLimit
	}
	if update.ResolverCountLimit > 0 {
		existing.ResolverCountLimit = update.ResolverCountLimit
	}
	if update.UserPoolConfig != nil {
		existing.UserPoolConfig = update.UserPoolConfig
	}
	if update.Visibility != "" {
		existing.Visibility = update.Visibility
	}
	existing.WafWebAclArn = update.WafWebAclArn
	existing.XrayEnabled = update.XrayEnabled

	if oldName != existing.Name {
		_ = s.graphqlApisStore.Delete(oldName)
	}

	if err := s.graphqlApisStore.Put(existing.Name, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteGraphqlApiById removes a GraphQL API by its UUID.
// Cascades deletion to all child resources: data sources, resolvers, functions,
// types, schema status, environment variables, API keys, API cache, and any
// merged API associations that reference this API.
func (s *AppSyncStore) DeleteGraphqlApiById(apiId string) error {
	existing, err := s.GetGraphqlApiById(apiId)
	if err != nil {
		return err
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	prefix := apiId + "/"

	// Remove all prefix-scoped child resources.
	for _, op := range []struct {
		name string
		fn   func(string) error
	}{
		{"dataSources", s.dataSourcesStore.DeleteByPrefix},
		{"resolvers", s.resolversStore.DeleteByPrefix},
		{"functions", s.functionsStore.DeleteByPrefix},
		{"types", s.typesStore.DeleteByPrefix},
		{"apiKeys", s.apiKeysStore.DeleteByPrefix},
	} {
		if err := op.fn(prefix); err != nil {
			logs.Warn("failed to delete child resources during API deletion",
				logs.String("apiId", apiId), logs.String("resource", op.name), logs.Err(err))
		}
	}

	// Remove exact-key resources.
	for _, op := range []struct {
		name string
		fn   func(string) error
	}{
		{"schemaStatuses", s.schemaStatusesStore.Delete},
		{"envVariables", s.envVariablesStore.Delete},
		{"apiCaches", s.apiCachesStore.Delete},
	} {
		if err := op.fn(apiId); err != nil {
			logs.Warn("failed to delete resource during API deletion",
				logs.String("apiId", apiId), logs.String("resource", op.name), logs.Err(err))
		}
	}

	// Remove merged API associations referencing this API as source or merged.
	var assocKeys []string
	_ = s.mergedApiAssociationsStore.ForEach(func(key string, value []byte) error {
		var assoc SourceApiAssociation
		if json.Unmarshal(value, &assoc) == nil {
			if assoc.SourceApiId == apiId || assoc.MergedApiId == apiId {
				assocKeys = append(assocKeys, key)
			}
		}
		return nil
	})
	for _, k := range assocKeys {
		if err := s.mergedApiAssociationsStore.Delete(k); err != nil {
			logs.Warn("failed to delete merged API association during API deletion",
				logs.String("apiId", apiId), logs.String("assocKey", k), logs.Err(err))
		}
	}

	return s.graphqlApisStore.Delete(existing.Name)
}

// ListGraphqlApis returns a paginated list of GraphQL APIs.
// When apiType is non-empty, only APIs matching that type are returned
// and the filter is applied before pagination.
func (s *AppSyncStore) ListGraphqlApis(opts common.ListOptions, apiType string) ([]*GraphqlApi, string, error) {
	var filter common.FilterFunc[GraphqlApi]
	if apiType != "" {
		filter = func(api *GraphqlApi) bool {
			return api.ApiType == apiType
		}
	}
	result, err := common.List[GraphqlApi](s.graphqlApisStore, opts, filter)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// --- DataSource ---

// CreateDataSource persists a new data source scoped to a GraphQL API.
func (s *AppSyncStore) CreateDataSource(ds *DataSource) (*DataSource, error) {
	if ds.ApiId == "" || ds.Name == "" || ds.Type == "" {
		return nil, common.NewStoreError("appsync", "create_data_source", common.ErrInvalidInput)
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := ds.ApiId + "/" + ds.Name
	if s.dataSourcesStore.Exists(key) {
		return nil, ErrDataSourceAlreadyExists
	}

	ds.DataSourceArn = s.BuildDataSourceARN(ds.ApiId, ds.Name)

	if err := s.dataSourcesStore.Put(key, ds); err != nil {
		return nil, err
	}
	return ds, nil
}

// GetDataSource retrieves a data source by API ID and name.
func (s *AppSyncStore) GetDataSource(apiId, name string) (*DataSource, error) {
	key := apiId + "/" + name
	var ds DataSource
	if err := s.dataSourcesStore.Get(key, &ds); err != nil {
		return nil, ErrDataSourceNotFound
	}
	return &ds, nil
}

// UpdateDataSource merges non-zero fields from the input into the existing data source.
func (s *AppSyncStore) UpdateDataSource(ds *DataSource) (*DataSource, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := ds.ApiId + "/" + ds.Name
	existing, err := s.GetDataSource(ds.ApiId, ds.Name)
	if err != nil {
		return nil, err
	}

	if ds.Type != "" {
		existing.Type = ds.Type
	}
	if ds.Description != "" {
		existing.Description = ds.Description
	}
	if ds.ServiceRoleArn != "" {
		existing.ServiceRoleArn = ds.ServiceRoleArn
	}
	if ds.DynamodbConfig != nil {
		existing.DynamodbConfig = ds.DynamodbConfig
	}
	if ds.ElasticsearchConfig != nil {
		existing.ElasticsearchConfig = ds.ElasticsearchConfig
	}
	if ds.EventBridgeConfig != nil {
		existing.EventBridgeConfig = ds.EventBridgeConfig
	}
	if ds.HttpConfig != nil {
		existing.HttpConfig = ds.HttpConfig
	}
	if ds.LambdaConfig != nil {
		existing.LambdaConfig = ds.LambdaConfig
	}
	if ds.MetricsConfig != "" {
		existing.MetricsConfig = ds.MetricsConfig
	}
	if ds.OpenSearchServiceConfig != nil {
		existing.OpenSearchServiceConfig = ds.OpenSearchServiceConfig
	}
	if ds.RelationalDatabaseConfig != nil {
		existing.RelationalDatabaseConfig = ds.RelationalDatabaseConfig
	}

	if err := s.dataSourcesStore.Put(key, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteDataSource removes a data source by API ID and name.
func (s *AppSyncStore) DeleteDataSource(apiId, name string) error {
	key := apiId + "/" + name
	if !s.dataSourcesStore.Exists(key) {
		return ErrDataSourceNotFound
	}
	return s.dataSourcesStore.Delete(key)
}

// ListDataSources returns a paginated list of data sources for a given GraphQL API.
func (s *AppSyncStore) ListDataSources(apiId string, opts common.ListOptions) ([]*DataSource, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   apiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[DataSource](s.dataSourcesStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// --- Resolver ---

// CreateResolver persists a new resolver scoped to a GraphQL API type and field.
func (s *AppSyncStore) CreateResolver(r *Resolver) (*Resolver, error) {
	if r.ApiId == "" || r.TypeName == "" || r.FieldName == "" {
		return nil, common.NewStoreError("appsync", "create_resolver", common.ErrInvalidInput)
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := r.ApiId + "/" + r.TypeName + "/" + r.FieldName
	if s.resolversStore.Exists(key) {
		return nil, ErrResolverAlreadyExists
	}

	r.ResolverArn = s.BuildResolverARN(r.ApiId, r.TypeName, r.FieldName)

	if err := s.resolversStore.Put(key, r); err != nil {
		return nil, err
	}
	return r, nil
}

// GetResolver retrieves a resolver by API ID, type name, and field name.
func (s *AppSyncStore) GetResolver(apiId, typeName, fieldName string) (*Resolver, error) {
	key := apiId + "/" + typeName + "/" + fieldName
	var r Resolver
	if err := s.resolversStore.Get(key, &r); err != nil {
		return nil, ErrResolverNotFound
	}
	return &r, nil
}

// UpdateResolver merges non-zero fields from the input into the existing resolver.
func (s *AppSyncStore) UpdateResolver(r *Resolver) (*Resolver, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := r.ApiId + "/" + r.TypeName + "/" + r.FieldName
	existing, err := s.GetResolver(r.ApiId, r.TypeName, r.FieldName)
	if err != nil {
		return nil, err
	}

	if r.Kind != "" {
		existing.Kind = r.Kind
	}
	if r.DataSourceName != "" {
		existing.DataSourceName = r.DataSourceName
	}
	if r.RequestMappingTemplate != "" {
		existing.RequestMappingTemplate = r.RequestMappingTemplate
	}
	if r.ResponseMappingTemplate != "" {
		existing.ResponseMappingTemplate = r.ResponseMappingTemplate
	}
	if r.PipelineConfig != nil {
		existing.PipelineConfig = r.PipelineConfig
	}
	if r.Runtime != nil {
		existing.Runtime = r.Runtime
	}
	if r.Code != "" {
		existing.Code = r.Code
	}
	if r.CachingConfig != nil {
		existing.CachingConfig = r.CachingConfig
	}
	if r.MaxBatchSize > 0 {
		existing.MaxBatchSize = r.MaxBatchSize
	}
	if r.MetricsConfig != "" {
		existing.MetricsConfig = r.MetricsConfig
	}
	if r.SyncConfig != nil {
		existing.SyncConfig = r.SyncConfig
	}

	if err := s.resolversStore.Put(key, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteResolver removes a resolver by API ID, type name, and field name.
func (s *AppSyncStore) DeleteResolver(apiId, typeName, fieldName string) error {
	key := apiId + "/" + typeName + "/" + fieldName
	if !s.resolversStore.Exists(key) {
		return ErrResolverNotFound
	}
	return s.resolversStore.Delete(key)
}

// ListResolvers returns a paginated list of resolvers for a given GraphQL API type.
func (s *AppSyncStore) ListResolvers(apiId, typeName string, opts common.ListOptions) ([]*Resolver, string, error) {
	prefix := apiId + "/"
	if typeName != "" {
		prefix += typeName + "/"
	}
	prefixOpts := common.ListOptions{
		Prefix:   prefix,
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[Resolver](s.resolversStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// ListResolversByFunction returns resolvers that reference a given function ID.
// Performs a full scan filtered by functionId, then applies offset-based pagination
// using the Marker as a decimal skip count.
func (s *AppSyncStore) ListResolversByFunction(apiId, functionId string, opts common.ListOptions) ([]*Resolver, string, error) {
	allResolvers, _, err := s.ListResolvers(apiId, "", common.ListOptions{MaxItems: 10000, Prefix: apiId + "/"})
	if err != nil {
		return nil, "", err
	}

	var filtered []*Resolver
	for _, r := range allResolvers {
		if r.PipelineConfig != nil {
			for _, fnId := range r.PipelineConfig.Functions {
				if fnId == functionId {
					filtered = append(filtered, r)
					break
				}
			}
		}
	}

	offset := 0
	if opts.Marker != "" {
		if n, err := strconv.Atoi(opts.Marker); err == nil {
			offset = n
		}
	}

	if offset >= len(filtered) {
		return []*Resolver{}, "", nil
	}

	remaining := filtered[offset:]
	if len(remaining) > opts.MaxItems {
		nextToken := fmt.Sprintf("%d", offset+opts.MaxItems)
		return remaining[:opts.MaxItems], nextToken, nil
	}
	return remaining, "", nil
}
