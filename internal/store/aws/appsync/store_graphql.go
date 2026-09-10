package appsync

import (
	"encoding/json"
	"strings"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
)

// --- GraphQL API (v1) ---

// graphqlApiIdIndexKey returns the index key for looking up a GraphQL API name by its UUID.
func graphqlApiIdIndexKey(apiId string) string {
	return "#id:" + apiId
}

// putGraphqlApiIdIndex writes the apiId→name index entry. Must be called within createMu.
func (s *AppSyncStore) putGraphqlApiIdIndex(apiId, name string) error {
	return s.graphqlApisStore.Put(graphqlApiIdIndexKey(apiId), map[string]string{"name": name})
}

// getGraphqlApiNameByIndex retrieves the GraphQL API name from the apiId
// index. A missing index entry means the API does not exist: every fresh
// write records the index entry alongside the record. Any other read
// failure is a storage error and propagates as such rather than
// masquerading as absence.
func (s *AppSyncStore) getGraphqlApiNameByIndex(apiId string) (string, error) {
	var m map[string]string
	if err := s.graphqlApisStore.Get(graphqlApiIdIndexKey(apiId), &m); err != nil {
		if !common.IsNotFound(err) {
			return "", err
		}
		return "", ErrGraphqlApiNotFound
	}
	if name, ok := m["name"]; ok {
		return name, nil
	}
	return "", ErrGraphqlApiNotFound
}

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

	graphqlCount, _ := s.CountGraphqlApis()
	eventCount, _ := s.CountApis()
	if graphqlCount+eventCount >= MaxApisPerRegion {
		return nil, ErrApiLimitExceeded
	}

	if s.graphqlApisStore.Exists(api.Name) {
		return nil, ErrGraphqlApiAlreadyExists
	}

	api.ApiId = s.GenerateId()
	api.Arn = s.BuildGraphQLApiARN(api.ApiId)
	if api.ApiType == "" {
		api.ApiType = "GRAPHQL"
	}
	if api.Uris == nil {
		api.Uris = graphqlEndpointURLs(api.ApiId)
	}
	if api.Dns == nil {
		api.Dns = graphqlEndpointURLs(api.ApiId)
	}

	if err := s.graphqlApisStore.Put(api.Name, api); err != nil {
		return nil, err
	}
	if err := s.putGraphqlApiIdIndex(api.ApiId, api.Name); err != nil {
		// The index is the only ID-lookup path, so a failed index write
		// must fail the create: the record is rolled back best-effort
		// rather than left unreachable by ID.
		if delErr := s.graphqlApisStore.Delete(api.Name); delErr != nil {
			logs.Error("failed to roll back GraphQL API record after index write failure",
				logs.String("name", api.Name), logs.Err(delErr))
		}
		return nil, err
	}
	return api, nil
}

// GetGraphqlApi retrieves a GraphQL API by name.
func (s *AppSyncStore) GetGraphqlApi(name string) (*GraphqlApi, error) {
	var api GraphqlApi
	if err := s.graphqlApisStore.Get(name, &api); err != nil {
		if !common.IsNotFound(err) {
			return nil, err
		}
		return nil, ErrGraphqlApiNotFound
	}
	return &api, nil
}

// GetGraphqlApiById retrieves a GraphQL API by its UUID via the apiId→name index.
func (s *AppSyncStore) GetGraphqlApiById(apiId string) (*GraphqlApi, error) {
	name, err := s.getGraphqlApiNameByIndex(apiId)
	if err != nil {
		return nil, err
	}
	return s.GetGraphqlApi(name)
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
	if update.QueryDepthLimitSet {
		existing.QueryDepthLimit = update.QueryDepthLimit
	}
	if update.ResolverCountLimitSet {
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
		if s.graphqlApisStore.Exists(existing.Name) {
			return nil, ErrGraphqlApiAlreadyExists
		}
	}

	if err := s.graphqlApisStore.Put(existing.Name, existing); err != nil {
		return nil, err
	}

	// Delete old key after successful Put to prevent data loss on rename.
	if oldName != existing.Name {
		_ = s.graphqlApisStore.Delete(oldName)
	}

	if err := s.putGraphqlApiIdIndex(existing.ApiId, existing.Name); err != nil {
		logs.Warn("failed to update graphqlApiId index during rename",
			logs.String("apiId", existing.ApiId), logs.Err(err))
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

	// Child tag rows are keyed by ARN, so they are swept before the bulk
	// prefix deletes remove the records they derive from.
	s.sweepChildTags(s.dataSourcesStore, apiId, func(segments []string) string {
		return s.BuildDataSourceARN(apiId, segments[0])
	})
	s.sweepChildTags(s.resolversStore, apiId, func(segments []string) string {
		if len(segments) < 2 {
			return ""
		}
		return s.BuildResolverARN(apiId, segments[0], segments[1])
	})
	s.sweepChildTags(s.functionsStore, apiId, func(segments []string) string {
		return s.BuildFunctionARN(apiId, segments[0])
	})
	s.sweepChildTags(s.typesStore, apiId, func(segments []string) string {
		return s.BuildTypeARN(apiId, segments[0])
	})
	s.sweepChildTags(s.apiKeysStore, apiId, func(segments []string) string {
		return s.BuildApiKeyARN(apiId, segments[0])
	})
	_ = s.TagStore.Delete(s.BuildApiCacheARN(apiId))

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
		{"resolverCache", s.resolverCacheStore.DeleteByPrefix},
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
	// Use prefix scan for MergedApiId matches (O(log n)); full scan for
	// SourceApiId matches (no secondary index, but API deletion is rare).
	type assocCleanup struct {
		key           string
		associationId string
	}
	var assocsToDelete []assocCleanup
	seen := make(map[string]bool)

	// Prefix scan: associations where this API is the merged API.
	_ = s.mergedApiAssociationsStore.ScanPrefix(apiId+"/", func(key string, value []byte) error {
		var assoc SourceApiAssociation
		if json.Unmarshal(value, &assoc) == nil && !seen[assoc.AssociationId] {
			seen[assoc.AssociationId] = true
			assocsToDelete = append(assocsToDelete, assocCleanup{key: key, associationId: assoc.AssociationId})
		}
		return nil
	})

	// Full scan: associations where this API is the source API.
	_ = s.mergedApiAssociationsStore.ForEach(func(key string, value []byte) error {
		var assoc SourceApiAssociation
		if json.Unmarshal(value, &assoc) == nil && assoc.SourceApiId == apiId && !seen[assoc.AssociationId] {
			seen[assoc.AssociationId] = true
			assocsToDelete = append(assocsToDelete, assocCleanup{key: key, associationId: assoc.AssociationId})
		}
		return nil
	})

	for _, entry := range assocsToDelete {
		if err := s.mergedApiAssociationsStore.Delete(entry.key); err != nil {
			logs.Warn("failed to delete merged API association during API deletion",
				logs.String("apiId", apiId), logs.String("assocKey", entry.key), logs.Err(err))
		}
		_ = s.mergedApiAssocIndexStore.Delete(entry.associationId)
	}

	// Remove domain name associations referencing this API.
	var domainAssocKeys []string
	_ = s.apiAssociationsStore.ForEach(func(key string, value []byte) error {
		var assoc ApiAssociation
		if json.Unmarshal(value, &assoc) == nil && assoc.ApiId == apiId {
			domainAssocKeys = append(domainAssocKeys, key)
		}
		return nil
	})
	for _, k := range domainAssocKeys {
		if err := s.apiAssociationsStore.Delete(k); err != nil {
			logs.Warn("failed to delete domain association during API deletion",
				logs.String("apiId", apiId), logs.String("domainName", k), logs.Err(err))
		}
	}

	_ = s.TagStore.Delete(existing.Arn)

	_ = s.graphqlApisStore.Delete(graphqlApiIdIndexKey(apiId))

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

// CountGraphqlApis returns the total number of GraphQL APIs in the store.
func (s *AppSyncStore) CountGraphqlApis() (int, error) {
	count := 0
	err := s.graphqlApisStore.ScanPrefix("", func(key string, value []byte) error {
		if !strings.HasPrefix(key, "#id:") {
			count++
		}
		return nil
	})
	return count, err
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
	if ds.NeptuneConfig != nil {
		existing.NeptuneConfig = ds.NeptuneConfig
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

// DeleteDataSource removes a data source by API ID and name. Tag rows are
// keyed by the data source ARN and removed with the record.
func (s *AppSyncStore) DeleteDataSource(apiId, name string) error {
	key := apiId + "/" + name
	if !s.dataSourcesStore.Exists(key) {
		return ErrDataSourceNotFound
	}
	if err := s.dataSourcesStore.Delete(key); err != nil {
		return err
	}
	_ = s.TagStore.Delete(s.BuildDataSourceARN(apiId, name))
	return nil
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
	if r.MaxBatchSizeSet {
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
	if err := s.resolversStore.Delete(key); err != nil {
		return err
	}
	_ = s.TagStore.Delete(s.BuildResolverARN(apiId, typeName, fieldName))
	return nil
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
// Uses common.List with a filter function for opaque-token pagination
// consistent with all other List operations.
func (s *AppSyncStore) ListResolversByFunction(apiId, functionId string, opts common.ListOptions) ([]*Resolver, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   apiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[Resolver](s.resolversStore, prefixOpts, func(r *Resolver) bool {
		if r.PipelineConfig == nil {
			return false
		}
		for _, fnId := range r.PipelineConfig.Functions {
			if fnId == functionId {
				return true
			}
		}
		return false
	})
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// GetAllResolversForApi returns all resolvers for an API without pagination limits.
func (s *AppSyncStore) GetAllResolversForApi(apiId string) ([]*Resolver, error) {
	return common.ListMatching[Resolver](s.resolversStore, apiId+"/", nil)
}
