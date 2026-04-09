package appsync

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/arn"

	"github.com/google/uuid"
)

// AppSyncStore provides persistent storage for AppSync resources.
// Uses separate buckets for each resource collection (apis, channel-namespaces, etc.).
type AppSyncStore struct {
	apisStore                  *common.BaseStore
	channelsStore              *common.BaseStore
	graphqlApisStore           *common.BaseStore
	dataSourcesStore           *common.BaseStore
	resolversStore             *common.BaseStore
	functionsStore             *common.BaseStore
	typesStore                 *common.BaseStore
	schemaStatusesStore        *common.BaseStore
	envVariablesStore          *common.BaseStore
	apiKeysStore               *common.BaseStore
	apiCachesStore             *common.BaseStore
	domainNamesStore           *common.BaseStore
	apiAssociationsStore       *common.BaseStore
	mergedApiAssociationsStore *common.BaseStore
	arnBuilder                 *arn.ARNBuilder
	accountId                  string
	region                     string
	createMu                   sync.Mutex
}

// apiBucketName returns the PebbleDB bucket name for Event API resources.
func apiBucketName(region string) string {
	return "appsync-apis-" + region
}

// channelNamespaceBucketName returns the PebbleDB bucket name for channel namespace resources.
func channelNamespaceBucketName(region string) string {
	return "appsync-channel-namespaces-" + region
}

// graphqlApiBucketName returns the PebbleDB bucket name for GraphQL API resources.
func graphqlApiBucketName(region string) string {
	return "appsync-graphql-apis-" + region
}

// dataSourceBucketName returns the PebbleDB bucket name for data source resources.
func dataSourceBucketName(region string) string {
	return "appsync-datasources-" + region
}

// resolverBucketName returns the PebbleDB bucket name for resolver resources.
func resolverBucketName(region string) string {
	return "appsync-resolvers-" + region
}

// functionBucketName returns the PebbleDB bucket name for AppSync function resources.
func functionBucketName(region string) string {
	return "appsync-functions-" + region
}

// typeBucketName returns the PebbleDB bucket name for GraphQL type definition resources.
func typeBucketName(region string) string {
	return "appsync-types-" + region
}

// schemaStatusBucketName returns the PebbleDB bucket name for schema creation status resources.
func schemaStatusBucketName(region string) string {
	return "appsync-schema-statuses-" + region
}

// envVariablesBucketName returns the PebbleDB bucket name for environment variable resources.
func envVariablesBucketName(region string) string {
	return "appsync-env-variables-" + region
}

// apiKeyBucketName returns the PebbleDB bucket name for API key resources.
func apiKeyBucketName(region string) string {
	return "appsync-api-keys-" + region
}

// apiCacheBucketName returns the PebbleDB bucket name for API cache resources.
func apiCacheBucketName(region string) string {
	return "appsync-api-caches-" + region
}

// domainNameBucketName returns the PebbleDB bucket name for domain name resources.
func domainNameBucketName(region string) string {
	return "appsync-domain-names-" + region
}

// apiAssociationBucketName returns the PebbleDB bucket name for API association resources.
func apiAssociationBucketName(region string) string {
	return "appsync-api-associations-" + region
}

// mergedApiAssociationBucketName returns the PebbleDB bucket name for merged API association resources.
func mergedApiAssociationBucketName(region string) string {
	return "appsync-merged-api-associations-" + region
}

// NewAppSyncStore creates a new store backed by the given storage and scoped to the specified account and region.
func NewAppSyncStore(store storage.BasicStorage, accountId, region string) *AppSyncStore {
	b := arn.NewARNBuilder(accountId, region)
	return &AppSyncStore{
		apisStore:                  common.NewBaseStore(store.Bucket(apiBucketName(region)), "appsync"),
		channelsStore:              common.NewBaseStore(store.Bucket(channelNamespaceBucketName(region)), "appsync-channel-namespaces"),
		graphqlApisStore:           common.NewBaseStore(store.Bucket(graphqlApiBucketName(region)), "appsync-graphql-apis"),
		dataSourcesStore:           common.NewBaseStore(store.Bucket(dataSourceBucketName(region)), "appsync-datasources"),
		resolversStore:             common.NewBaseStore(store.Bucket(resolverBucketName(region)), "appsync-resolvers"),
		functionsStore:             common.NewBaseStore(store.Bucket(functionBucketName(region)), "appsync-functions"),
		typesStore:                 common.NewBaseStore(store.Bucket(typeBucketName(region)), "appsync-types"),
		schemaStatusesStore:        common.NewBaseStore(store.Bucket(schemaStatusBucketName(region)), "appsync-schema-statuses"),
		envVariablesStore:          common.NewBaseStore(store.Bucket(envVariablesBucketName(region)), "appsync-env-variables"),
		apiKeysStore:               common.NewBaseStore(store.Bucket(apiKeyBucketName(region)), "appsync-api-keys"),
		apiCachesStore:             common.NewBaseStore(store.Bucket(apiCacheBucketName(region)), "appsync-api-caches"),
		domainNamesStore:           common.NewBaseStore(store.Bucket(domainNameBucketName(region)), "appsync-domain-names"),
		apiAssociationsStore:       common.NewBaseStore(store.Bucket(apiAssociationBucketName(region)), "appsync-api-associations"),
		mergedApiAssociationsStore: common.NewBaseStore(store.Bucket(mergedApiAssociationBucketName(region)), "appsync-merged-api-associations"),
		arnBuilder:                 b,
		accountId:                  accountId,
		region:                     region,
	}
}

// GetAccountID returns the AWS account ID this store is scoped to.
func (s *AppSyncStore) GetAccountID() string { return s.accountId }

// GetRegion returns the AWS region this store is scoped to.
func (s *AppSyncStore) GetRegion() string { return s.region }

// GenerateId generates a new UUID suitable for use as a resource identifier.
func (s *AppSyncStore) GenerateId() string {
	return uuid.New().String()
}

// BuildApiARN constructs an ARN for an Event API (v2).
func (s *AppSyncStore) BuildApiARN(apiId string) string {
	return s.arnBuilder.AppSync().Api(apiId)
}

// BuildChannelNamespaceARN constructs an ARN for a channel namespace.
func (s *AppSyncStore) BuildChannelNamespaceARN(apiId, name string) string {
	return s.arnBuilder.AppSync().ChannelNamespace(apiId, name)
}

// BuildGraphQLApiARN constructs an ARN for a GraphQL API (v1).
func (s *AppSyncStore) BuildGraphQLApiARN(apiId string) string {
	return s.arnBuilder.AppSync().GraphQLApi(apiId)
}

// BuildDataSourceARN constructs an ARN for a data source.
func (s *AppSyncStore) BuildDataSourceARN(apiId, name string) string {
	return s.arnBuilder.AppSync().DataSource(apiId, name)
}

// BuildResolverARN constructs an ARN for a resolver.
func (s *AppSyncStore) BuildResolverARN(apiId, typeName, fieldName string) string {
	return s.arnBuilder.AppSync().Resolver(apiId, typeName, fieldName)
}

// BuildFunctionARN constructs an ARN for an AppSync function.
func (s *AppSyncStore) BuildFunctionARN(apiId, functionId string) string {
	return s.arnBuilder.AppSync().Function(apiId, functionId)
}

// BuildTypeARN constructs an ARN for a type definition.
func (s *AppSyncStore) BuildTypeARN(apiId, typeName string) string {
	return s.arnBuilder.AppSync().Type(apiId, typeName)
}

// BuildDomainNameARN constructs an ARN for a custom domain name.
func (s *AppSyncStore) BuildDomainNameARN(name string) string {
	return s.arnBuilder.AppSync().DomainName(name)
}

// --- Event API (v2) ---

// CreateApi persists a new Event API. Generates apiId, ARN, timestamps, and default DNS endpoints.
func (s *AppSyncStore) CreateApi(api *Api) (*Api, error) {
	if api.Name == "" {
		return nil, common.NewStoreError("appsync", "create_api", common.ErrInvalidInput)
	}
	if api.EventConfig == nil {
		return nil, common.NewStoreError("appsync", "create_api", common.ErrInvalidInput)
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	if s.apisStore.Exists(api.Name) {
		return nil, ErrApiAlreadyExists
	}

	api.ApiId = s.GenerateId()
	api.Arn = s.BuildApiARN(api.ApiId)
	api.Created = time.Now().UTC()
	if api.Dns == nil {
		api.Dns = map[string]string{
			"HTTP":     "http://127.0.0.1:8086/event",
			"REALTIME": "ws://127.0.0.1:8086/event/realtime",
		}
	}

	if err := s.apisStore.Put(api.Name, api); err != nil {
		return nil, err
	}
	return api, nil
}

// GetApi retrieves an Event API by name.
func (s *AppSyncStore) GetApi(name string) (*Api, error) {
	var api Api
	if err := s.apisStore.Get(name, &api); err != nil {
		return nil, ErrApiNotFound
	}
	return &api, nil
}

// GetApiById retrieves an Event API by its UUID.
// Scans all APIs to find the one matching the given ID.
func (s *AppSyncStore) GetApiById(apiId string) (*Api, error) {
	apis, _, err := s.ListApis(common.ListOptions{MaxItems: 10000})
	if err != nil {
		return nil, err
	}
	for _, api := range apis {
		if api.ApiId == apiId {
			return api, nil
		}
	}
	return nil, ErrApiNotFound
}

// UpdateApiById updates an Event API identified by apiId.
// Merges non-zero fields from the update; always copies Tags if present.
func (s *AppSyncStore) UpdateApiById(apiId string, update *Api) (*Api, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	existing, err := s.GetApiById(apiId)
	if err != nil {
		return nil, err
	}

	oldName := existing.Name

	if update.Name != "" {
		existing.Name = update.Name
	}
	if update.EventConfig != nil {
		existing.EventConfig = update.EventConfig
	}
	if update.OwnerContact != "" {
		existing.OwnerContact = update.OwnerContact
	}
	existing.WafWebAclArn = update.WafWebAclArn
	existing.XrayEnabled = update.XrayEnabled
	if update.Tags != nil {
		existing.Tags = update.Tags
	}

	// If the name changed, remove the old entry before saving under the new key.
	if oldName != existing.Name {
		s.apisStore.Delete(oldName)
	}

	if err := s.apisStore.Put(existing.Name, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteApiById removes an Event API by its UUID.
// Cascades deletion to all channel namespaces belonging to the API.
func (s *AppSyncStore) DeleteApiById(apiId string) error {
	existing, err := s.GetApiById(apiId)
	if err != nil {
		return err
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	// Remove all channel namespaces for this Event API.
	_ = s.channelsStore.DeleteByPrefix(apiId + "/")

	return s.apisStore.Delete(existing.Name)
}

// ListApis returns a paginated list of Event APIs.
func (s *AppSyncStore) ListApis(opts common.ListOptions) ([]*Api, string, error) {
	result, err := common.List[Api](s.apisStore, opts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// ListAssociationsBySourceApi lists merged API associations for a given source API.
// Filters all associations by sourceApiId since the store key is composite (mergedApiId/associationId).
func (s *AppSyncStore) ListAssociationsBySourceApi(sourceApiId string, opts common.ListOptions) ([]*SourceApiAssociation, string, error) {
	result, err := common.List[SourceApiAssociation](s.mergedApiAssociationsStore, opts, func(a *SourceApiAssociation) bool {
		return a.SourceApiId == sourceApiId
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

// CreateChannelNamespace persists a new channel namespace scoped to an Event API.
func (s *AppSyncStore) CreateChannelNamespace(ns *ChannelNamespace) (*ChannelNamespace, error) {
	if ns.Name == "" || ns.ApiId == "" {
		return nil, common.NewStoreError(s.channelsStore.Service(), "create_channel_namespace", common.ErrInvalidInput)
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := ns.ApiId + "/" + ns.Name
	if s.channelsStore.Exists(key) {
		return nil, ErrChannelNamespaceExists
	}

	ns.ChannelNamespaceArn = s.BuildChannelNamespaceARN(ns.ApiId, ns.Name)
	now := time.Now().UTC()
	ns.Created = now
	ns.LastModified = now

	if err := s.channelsStore.Put(key, ns); err != nil {
		return nil, err
	}
	return ns, nil
}

// GetChannelNamespace retrieves a channel namespace by API ID and name.
func (s *AppSyncStore) GetChannelNamespace(apiId, name string) (*ChannelNamespace, error) {
	key := apiId + "/" + name
	var ns ChannelNamespace
	if err := s.channelsStore.Get(key, &ns); err != nil {
		return nil, ErrChannelNamespaceNotFound
	}
	return &ns, nil
}

// UpdateChannelNamespace merges non-zero fields from the input into the existing channel namespace.
func (s *AppSyncStore) UpdateChannelNamespace(ns *ChannelNamespace) (*ChannelNamespace, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := ns.ApiId + "/" + ns.Name
	existing, err := s.GetChannelNamespace(ns.ApiId, ns.Name)
	if err != nil {
		return nil, err
	}

	if ns.CodeHandlers != "" {
		existing.CodeHandlers = ns.CodeHandlers
	}
	if ns.HandlerConfigs != nil {
		existing.HandlerConfigs = ns.HandlerConfigs
	}
	if ns.PublishAuthModes != nil {
		existing.PublishAuthModes = ns.PublishAuthModes
	}
	if ns.SubscribeAuthModes != nil {
		existing.SubscribeAuthModes = ns.SubscribeAuthModes
	}
	existing.LastModified = time.Now().UTC()

	if err := s.channelsStore.Put(key, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteChannelNamespace removes a channel namespace by API ID and name.
func (s *AppSyncStore) DeleteChannelNamespace(apiId, name string) error {
	key := apiId + "/" + name
	if !s.channelsStore.Exists(key) {
		return ErrChannelNamespaceNotFound
	}
	return s.channelsStore.Delete(key)
}

// ListChannelNamespaces returns a paginated list of channel namespaces for a given API.
func (s *AppSyncStore) ListChannelNamespaces(apiId string, opts common.ListOptions) ([]*ChannelNamespace, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   apiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[ChannelNamespace](s.channelsStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// UpdateChannelNamespaceTags atomically updates the tags on a channel namespace
// using the provided merge function.
func (s *AppSyncStore) UpdateChannelNamespaceTags(apiId, name string, mergeFn func(map[string]string) map[string]string) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := apiId + "/" + name
	var ns ChannelNamespace
	if err := s.channelsStore.Get(key, &ns); err != nil {
		return ErrChannelNamespaceNotFound
	}

	ns.Tags = mergeFn(ns.Tags)
	ns.LastModified = time.Now().UTC()

	return s.channelsStore.Put(key, ns)
}

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
		api.Uris = map[string]string{
			"GRAPHQL":  "http://127.0.0.1:8080/v1/apis/" + api.ApiId + "/graphql",
			"REALTIME": "ws://127.0.0.1:8080/v1/apis/" + api.ApiId + "/realtime",
		}
	}
	if api.Dns == nil {
		api.Dns = map[string]string{
			"GRAPHQL":  "http://127.0.0.1:8080/v1/apis/" + api.ApiId + "/graphql",
			"REALTIME": "ws://127.0.0.1:8080/v1/apis/" + api.ApiId + "/realtime",
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
	if update.Tags != nil {
		existing.Tags = update.Tags
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
		s.graphqlApisStore.Delete(oldName)
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
	_ = s.dataSourcesStore.DeleteByPrefix(prefix)
	_ = s.resolversStore.DeleteByPrefix(prefix)
	_ = s.functionsStore.DeleteByPrefix(prefix)
	_ = s.typesStore.DeleteByPrefix(prefix)
	_ = s.apiKeysStore.DeleteByPrefix(prefix)

	// Remove exact-key resources.
	_ = s.schemaStatusesStore.Delete(apiId)
	_ = s.envVariablesStore.Delete(apiId)
	_ = s.apiCachesStore.Delete(apiId)

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
		_ = s.mergedApiAssociationsStore.Delete(k)
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

// --- Function (AppSync) ---

// CreateFunction persists a new AppSync function (a reusable resolver unit).
func (s *AppSyncStore) CreateFunction(f *FunctionConfiguration) (*FunctionConfiguration, error) {
	if f.ApiId == "" || f.Name == "" || f.DataSourceName == "" {
		return nil, common.NewStoreError("appsync", "create_function", common.ErrInvalidInput)
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	if f.FunctionId == "" {
		f.FunctionId = s.GenerateId()
	}

	key := f.ApiId + "/" + f.FunctionId
	if s.functionsStore.Exists(key) {
		return nil, ErrFunctionAlreadyExists
	}

	f.FunctionArn = s.BuildFunctionARN(f.ApiId, f.FunctionId)
	if f.FunctionVersion == "" {
		f.FunctionVersion = "2018-05-29"
	}

	if err := s.functionsStore.Put(key, f); err != nil {
		return nil, err
	}
	return f, nil
}

// GetFunction retrieves an AppSync function by API ID and function ID.
func (s *AppSyncStore) GetFunction(apiId, functionId string) (*FunctionConfiguration, error) {
	key := apiId + "/" + functionId
	var f FunctionConfiguration
	if err := s.functionsStore.Get(key, &f); err != nil {
		return nil, ErrFunctionNotFound
	}
	return &f, nil
}

// UpdateFunction merges non-zero fields from the input into the existing function.
func (s *AppSyncStore) UpdateFunction(f *FunctionConfiguration) (*FunctionConfiguration, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := f.ApiId + "/" + f.FunctionId
	existing, err := s.GetFunction(f.ApiId, f.FunctionId)
	if err != nil {
		return nil, err
	}

	if f.Name != "" {
		existing.Name = f.Name
	}
	if f.DataSourceName != "" {
		existing.DataSourceName = f.DataSourceName
	}
	if f.Description != "" {
		existing.Description = f.Description
	}
	if f.FunctionVersion != "" {
		existing.FunctionVersion = f.FunctionVersion
	}
	if f.RequestMappingTemplate != "" {
		existing.RequestMappingTemplate = f.RequestMappingTemplate
	}
	if f.ResponseMappingTemplate != "" {
		existing.ResponseMappingTemplate = f.ResponseMappingTemplate
	}
	if f.Runtime != nil {
		existing.Runtime = f.Runtime
	}
	if f.Code != "" {
		existing.Code = f.Code
	}
	if f.MaxBatchSize > 0 {
		existing.MaxBatchSize = f.MaxBatchSize
	}
	if f.SyncConfig != nil {
		existing.SyncConfig = f.SyncConfig
	}

	if err := s.functionsStore.Put(key, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteFunction removes an AppSync function by API ID and function ID.
func (s *AppSyncStore) DeleteFunction(apiId, functionId string) error {
	key := apiId + "/" + functionId
	if !s.functionsStore.Exists(key) {
		return ErrFunctionNotFound
	}
	return s.functionsStore.Delete(key)
}

// ListFunctions returns a paginated list of AppSync functions for a given GraphQL API.
func (s *AppSyncStore) ListFunctions(apiId string, opts common.ListOptions) ([]*FunctionConfiguration, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   apiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[FunctionConfiguration](s.functionsStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// --- Type ---

// CreateType persists a new type definition scoped to a GraphQL API.
// Extracts the type name from the SDL definition if the Name field is not set.
func (s *AppSyncStore) CreateType(t *Type) (*Type, error) {
	if t.ApiId == "" || t.Definition == "" || t.Format == "" {
		return nil, common.NewStoreError("appsync", "create_type", common.ErrInvalidInput)
	}

	if t.Name == "" {
		t.Name = extractTypeName(t.Definition)
	}
	if t.Name == "" {
		return nil, common.NewStoreError("appsync", "create_type", common.ErrInvalidInput)
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := t.ApiId + "/" + t.Name
	if s.typesStore.Exists(key) {
		return nil, ErrTypeAlreadyExists
	}

	t.Arn = s.BuildTypeARN(t.ApiId, t.Name)

	if err := s.typesStore.Put(key, t); err != nil {
		return nil, err
	}
	return t, nil
}

// extractTypeName parses a GraphQL SDL definition and extracts the type name.
// Handles formats like "type Post { ... }", "input PostInput { ... }", "enum Status { ... }".
func extractTypeName(definition string) string {
	definition = strings.TrimSpace(definition)
	if definition == "" {
		return ""
	}

	parts := strings.Fields(definition)
	if len(parts) < 2 {
		return ""
	}

	for _, keyword := range []string{"type", "input", "interface", "enum", "union", "scalar"} {
		if parts[0] == keyword {
			name := parts[1]
			if idx := strings.Index(name, "{"); idx >= 0 {
				name = strings.TrimSpace(name[:idx])
			}
			return name
		}
	}

	return ""
}

// GetType retrieves a type definition by API ID and type name.
func (s *AppSyncStore) GetType(apiId, typeName string) (*Type, error) {
	key := apiId + "/" + typeName
	var t Type
	if err := s.typesStore.Get(key, &t); err != nil {
		return nil, ErrTypeNotFound
	}
	return &t, nil
}

// UpdateType merges non-zero fields from the input into the existing type definition.
func (s *AppSyncStore) UpdateType(t *Type) (*Type, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := t.ApiId + "/" + t.Name
	existing, err := s.GetType(t.ApiId, t.Name)
	if err != nil {
		return nil, err
	}

	if t.Definition != "" {
		existing.Definition = t.Definition
	}
	if t.Format != "" {
		existing.Format = t.Format
	}
	if t.Description != "" {
		existing.Description = t.Description
	}

	if err := s.typesStore.Put(key, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteType removes a type definition by API ID and type name.
func (s *AppSyncStore) DeleteType(apiId, typeName string) error {
	key := apiId + "/" + typeName
	if !s.typesStore.Exists(key) {
		return ErrTypeNotFound
	}
	return s.typesStore.Delete(key)
}

// ListTypes returns a paginated list of type definitions for a given GraphQL API.
func (s *AppSyncStore) ListTypes(apiId string, opts common.ListOptions) ([]*Type, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   apiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[Type](s.typesStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// --- Schema Creation ---

// SaveSchemaCreationStatus persists the status of a schema creation operation.
func (s *AppSyncStore) SaveSchemaCreationStatus(apiId string, status *SchemaCreationStatus) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	return s.schemaStatusesStore.Put(apiId, status)
}

// GetSchemaCreationStatus retrieves the status of a schema creation operation.
func (s *AppSyncStore) GetSchemaCreationStatus(apiId string) (*SchemaCreationStatus, error) {
	var status SchemaCreationStatus
	if err := s.schemaStatusesStore.Get(apiId, &status); err != nil {
		return nil, ErrSchemaCreationNotFound
	}
	return &status, nil
}

// --- Environment Variables ---

// SaveEnvironmentVariables persists environment variables for a GraphQL API.
func (s *AppSyncStore) SaveEnvironmentVariables(apiId string, envVars *EnvironmentVariables) error {
	return s.envVariablesStore.Put(apiId, envVars)
}

// GetEnvironmentVariables retrieves environment variables for a GraphQL API.
func (s *AppSyncStore) GetEnvironmentVariables(apiId string) (*EnvironmentVariables, error) {
	var envVars EnvironmentVariables
	if err := s.envVariablesStore.Get(apiId, &envVars); err != nil {
		return &EnvironmentVariables{EnvironmentVariables: map[string]string{}}, nil
	}
	return &envVars, nil
}

// --- API Keys ---
// API keys are stored with composite keys "apiId/keyId" to ensure uniqueness
// across GraphQL APIs and to support scoped listing by apiId prefix.

// CreateApiKey persists a new API key.
func (s *AppSyncStore) CreateApiKey(apiId string, apiKey *ApiKey) error {
	return s.apiKeysStore.Put(apiId+"/"+apiKey.Id, apiKey)
}

// GetApiKey retrieves an API key by ID (requires apiId for composite key lookup).
func (s *AppSyncStore) GetApiKey(apiId, id string) (*ApiKey, error) {
	var apiKey ApiKey
	if err := s.apiKeysStore.Get(apiId+"/"+id, &apiKey); err != nil {
		return nil, ErrApiKeyNotFound
	}
	return &apiKey, nil
}

// UpdateApiKey updates an existing API key.
func (s *AppSyncStore) UpdateApiKey(apiId string, apiKey *ApiKey) error {
	return s.apiKeysStore.Put(apiId+"/"+apiKey.Id, apiKey)
}

// DeleteApiKey removes an API key by ID.
func (s *AppSyncStore) DeleteApiKey(apiId, id string) error {
	key := apiId + "/" + id
	if !s.apiKeysStore.Exists(key) {
		return ErrApiKeyNotFound
	}
	return s.apiKeysStore.Delete(key)
}

// ListApiKeys lists API keys for a given API, optionally filtered by expiry.
func (s *AppSyncStore) ListApiKeys(apiId string, opts common.ListOptions) ([]*ApiKey, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   apiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[ApiKey](s.apiKeysStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// --- API Cache ---

// CreateApiCache persists a new API cache configuration.
func (s *AppSyncStore) CreateApiCache(apiId string, cache *ApiCache) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	if s.apiCachesStore.Exists(apiId) {
		return ErrApiCacheAlreadyExists
	}

	return s.apiCachesStore.Put(apiId, cache)
}

// GetApiCache retrieves the API cache configuration for a given API.
func (s *AppSyncStore) GetApiCache(apiId string) (*ApiCache, error) {
	var cache ApiCache
	if err := s.apiCachesStore.Get(apiId, &cache); err != nil {
		return nil, ErrApiCacheNotFound
	}
	return &cache, nil
}

// UpdateApiCache updates the API cache configuration.
func (s *AppSyncStore) UpdateApiCache(apiId string, cache *ApiCache) error {
	return s.apiCachesStore.Put(apiId, cache)
}

// DeleteApiCache removes the API cache configuration.
// Validates existence first, as BaseStore.Delete silently succeeds on missing keys.
func (s *AppSyncStore) DeleteApiCache(apiId string) error {
	if !s.apiCachesStore.Exists(apiId) {
		return ErrApiCacheNotFound
	}
	return s.apiCachesStore.Delete(apiId)
}

// --- Domain Names ---

// CreateDomainName persists a new domain name configuration.
func (s *AppSyncStore) CreateDomainName(domainName *DomainNameConfig) error {
	return s.domainNamesStore.Put(domainName.DomainName, domainName)
}

// GetDomainName retrieves a domain name configuration.
func (s *AppSyncStore) GetDomainName(domainName string) (*DomainNameConfig, error) {
	var config DomainNameConfig
	if err := s.domainNamesStore.Get(domainName, &config); err != nil {
		return nil, ErrDomainNameNotFound
	}
	return &config, nil
}

// UpdateDomainName updates a domain name configuration.
func (s *AppSyncStore) UpdateDomainName(domainName *DomainNameConfig) error {
	return s.domainNamesStore.Put(domainName.DomainName, domainName)
}

// DeleteDomainName removes a domain name configuration.
// Validates existence first, as BaseStore.Delete silently succeeds on missing keys.
func (s *AppSyncStore) DeleteDomainName(domainName string) error {
	if !s.domainNamesStore.Exists(domainName) {
		return ErrDomainNameNotFound
	}
	return s.domainNamesStore.Delete(domainName)
}

// ListDomainNames lists all domain names.
func (s *AppSyncStore) ListDomainNames(opts common.ListOptions) ([]*DomainNameConfig, string, error) {
	result, err := common.List[DomainNameConfig](s.domainNamesStore, opts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// --- API Associations ---

// AssociateApi persists an API association with a domain name.
func (s *AppSyncStore) AssociateApi(association *ApiAssociation) error {
	return s.apiAssociationsStore.Put(association.DomainName, association)
}

// GetApiAssociation retrieves the API association for a domain name.
func (s *AppSyncStore) GetApiAssociation(domainName string) (*ApiAssociation, error) {
	var assoc ApiAssociation
	if err := s.apiAssociationsStore.Get(domainName, &assoc); err != nil {
		return nil, ErrApiAssociationNotFound
	}
	return &assoc, nil
}

// DisassociateApi removes the API association for a domain name.
// Validates existence first, as BaseStore.Delete silently succeeds on missing keys.
func (s *AppSyncStore) DisassociateApi(domainName string) error {
	if !s.apiAssociationsStore.Exists(domainName) {
		return ErrApiAssociationNotFound
	}
	return s.apiAssociationsStore.Delete(domainName)
}

// --- Merged API Associations ---
// Associations are stored with composite keys "mergedApiId/associationId" to
// support efficient prefix-based listing by merged API. For source-side lookups
// (where only associationId is known), GetMergedApiAssociationById performs a
// full scan with filter.

// CreateMergedApiAssociation persists a merged API association.
func (s *AppSyncStore) CreateMergedApiAssociation(assoc *SourceApiAssociation) error {
	key := assoc.MergedApiId + "/" + assoc.AssociationId
	return s.mergedApiAssociationsStore.Put(key, assoc)
}

// GetMergedApiAssociation retrieves a merged API association by merged API ID and association ID.
func (s *AppSyncStore) GetMergedApiAssociation(mergedApiId, associationId string) (*SourceApiAssociation, error) {
	key := mergedApiId + "/" + associationId
	var assoc SourceApiAssociation
	if err := s.mergedApiAssociationsStore.Get(key, &assoc); err != nil {
		return nil, ErrMergedApiAssociationNotFound
	}
	return &assoc, nil
}

// GetMergedApiAssociationById retrieves a merged API association by its UUID alone.
// This is required for source-side operations (e.g. DisassociateMergedGraphqlApi)
// where the mergedApiId is not available from the request path.
func (s *AppSyncStore) GetMergedApiAssociationById(associationId string) (*SourceApiAssociation, error) {
	var found *SourceApiAssociation
	err := s.mergedApiAssociationsStore.ForEach(func(key string, value []byte) error {
		var assoc SourceApiAssociation
		if err := json.Unmarshal(value, &assoc); err != nil {
			return err
		}
		if assoc.AssociationId == associationId {
			found = &assoc
		}
		return nil
	})
	if err != nil || found == nil {
		return nil, ErrMergedApiAssociationNotFound
	}
	return found, nil
}

// UpdateMergedApiAssociation updates an existing merged API association.
func (s *AppSyncStore) UpdateMergedApiAssociation(assoc *SourceApiAssociation) error {
	key := assoc.MergedApiId + "/" + assoc.AssociationId
	return s.mergedApiAssociationsStore.Put(key, assoc)
}

// DeleteMergedApiAssociation removes a merged API association by merged API ID and association ID.
func (s *AppSyncStore) DeleteMergedApiAssociation(mergedApiId, associationId string) error {
	key := mergedApiId + "/" + associationId
	return s.mergedApiAssociationsStore.Delete(key)
}

// ListMergedApiAssociations lists merged API associations for a given merged API.
func (s *AppSyncStore) ListMergedApiAssociations(mergedApiId string, opts common.ListOptions) ([]*SourceApiAssociation, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   mergedApiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[SourceApiAssociation](s.mergedApiAssociationsStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}
