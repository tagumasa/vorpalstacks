package appsync

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/config"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/graphql"

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
	resolverCacheStore         *common.BaseStore
	domainNamesStore           *common.BaseStore
	apiAssociationsStore       *common.BaseStore
	mergedApiAssociationsStore *common.BaseStore
	mergedApiAssocIndexStore   *common.BaseStore
	TagStore                   *common.TagStore
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

// mergedApiAssocIndexBucketName returns the PebbleDB bucket name for the secondary index
// that maps associationId → mergedApiId, enabling O(1) reverse lookups.
func mergedApiAssocIndexBucketName(region string) string {
	return "appsync-merged-api-assoc-index-" + region
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
		resolverCacheStore:         common.NewBaseStore(store.Bucket(resolverCacheBucketName(region)), "appsync-resolver-cache"),
		domainNamesStore:           common.NewBaseStore(store.Bucket(domainNameBucketName(region)), "appsync-domain-names"),
		apiAssociationsStore:       common.NewBaseStore(store.Bucket(apiAssociationBucketName(region)), "appsync-api-associations"),
		mergedApiAssociationsStore: common.NewBaseStore(store.Bucket(mergedApiAssociationBucketName(region)), "appsync-merged-api-associations"),
		mergedApiAssocIndexStore:   common.NewBaseStore(store.Bucket(mergedApiAssocIndexBucketName(region)), "appsync-merged-api-assoc-index"),
		TagStore:                   common.NewTagStoreWithRegion(store, "appsync", region),
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

	graphqlCount, _ := s.CountGraphqlApis()
	eventCount, _ := s.CountApis()
	if graphqlCount+eventCount >= MaxApisPerRegion {
		return nil, ErrApiLimitExceeded
	}

	if s.apisStore.Exists(api.Name) {
		return nil, ErrApiAlreadyExists
	}

	api.ApiId = s.GenerateId()
	api.Arn = s.BuildApiARN(api.ApiId)
	api.Created = time.Now().UTC()
	if api.Dns == nil {
		eventsPort := config.GetInt("ports.appsync_events")
		baseURL := config.GetString("endpoints.base_url")
		host := strings.TrimPrefix(baseURL, "http://")
		host = strings.TrimPrefix(host, "https://")
		if idx := strings.Index(host, ":"); idx >= 0 {
			host = host[:idx]
		}
		if host == "" {
			host = "127.0.0.1"
		}
		api.Dns = map[string]string{
			"HTTP":     fmt.Sprintf("http://%s:%d/event", host, eventsPort),
			"REALTIME": fmt.Sprintf("ws://%s:%d/event/realtime", host, eventsPort),
		}
	}

	if err := s.apisStore.Put(api.Name, api); err != nil {
		return nil, err
	}
	if err := s.putApiIdIndex(api.ApiId, api.Name); err != nil {
		logs.Warn("failed to write apiId index", logs.String("apiId", api.ApiId), logs.Err(err))
	}
	return api, nil
}

// apiIdIndexKey returns the index key for looking up an Event API name by its UUID.
func apiIdIndexKey(apiId string) string {
	return "#id:" + apiId
}

// putApiIdIndex writes the apiId→name index entry. Must be called within createMu.
func (s *AppSyncStore) putApiIdIndex(apiId, name string) error {
	return s.apisStore.Put(apiIdIndexKey(apiId), map[string]string{"name": name})
}

// getApiNameByIndex retrieves the API name from the apiId index.
// Falls back to full scan if the index entry is missing (pre-index data).
func (s *AppSyncStore) getApiNameByIndex(apiId string) (string, error) {
	var m map[string]string
	if err := s.apisStore.Get(apiIdIndexKey(apiId), &m); err == nil {
		if name, ok := m["name"]; ok {
			return name, nil
		}
	}

	// Fallback: full scan for pre-index data.
	apis, err := common.ListMatching[Api](s.apisStore, "", func(a *Api) bool {
		return a.ApiId == apiId
	})
	if err != nil {
		return "", err
	}
	if len(apis) > 0 {
		return apis[0].Name, nil
	}
	return "", ErrApiNotFound
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
// Uses the apiId→name index for direct lookup, with full-scan fallback.
func (s *AppSyncStore) GetApiById(apiId string) (*Api, error) {
	name, err := s.getApiNameByIndex(apiId)
	if err != nil {
		return nil, err
	}
	return s.GetApi(name)
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

	if oldName != existing.Name {
		if s.apisStore.Exists(existing.Name) {
			return nil, ErrApiAlreadyExists
		}
	}

	if err := s.apisStore.Put(existing.Name, existing); err != nil {
		return nil, err
	}

	// Delete old key after successful Put to prevent data loss on rename.
	if oldName != existing.Name {
		if err := s.apisStore.Delete(oldName); err != nil {
			logs.Warn("failed to delete stale Event API name during rename",
				logs.String("apiId", existing.ApiId), logs.String("oldName", oldName), logs.Err(err))
		}
	}

	if err := s.putApiIdIndex(existing.ApiId, existing.Name); err != nil {
		logs.Warn("failed to update apiId index during rename",
			logs.String("apiId", existing.ApiId), logs.Err(err))
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
	if err := s.channelsStore.DeleteByPrefix(apiId + "/"); err != nil {
		logs.Warn("failed to delete channel namespaces during Event API deletion",
			logs.String("apiId", apiId), logs.Err(err))
	}

	_ = s.TagStore.Delete(existing.Arn)

	_ = s.apisStore.Delete(apiIdIndexKey(apiId))

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

// CountApis returns the total number of Event APIs in the store.
func (s *AppSyncStore) CountApis() (int, error) {
	count := 0
	err := s.apisStore.ScanPrefix("", func(key string, value []byte) error {
		if !strings.HasPrefix(key, "#id:") {
			count++
		}
		return nil
	})
	return count, err
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

	if ns.CodeHandlersSet {
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
	arn := s.BuildChannelNamespaceARN(apiId, name)
	_ = s.TagStore.Delete(arn)
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

// extractTypeName parses a GraphQL SDL definition and extracts the type name.
// Handles formats like "type Post { ... }", "input PostInput { ... }", "enum Status { ... }".
func extractTypeName(definition string) string {
	return graphql.ExtractTypeName(definition)
}
