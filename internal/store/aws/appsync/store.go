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

	if s.apisStore.Exists(api.Name) {
		return nil, ErrApiAlreadyExists
	}

	api.ApiId = s.GenerateId()
	api.Arn = s.BuildApiARN(api.ApiId)
	api.Created = time.Now().UTC()
	if api.Dns == nil {
		eventsPort := config.GetInt("ports.appsync_events")
		api.Dns = map[string]string{
			"HTTP":     fmt.Sprintf("http://127.0.0.1:%d/event", eventsPort),
			"REALTIME": fmt.Sprintf("ws://127.0.0.1:%d/event/realtime", eventsPort),
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

	// If the name changed, remove the old entry before saving under the new key.
	if oldName != existing.Name {
		if err := s.apisStore.Delete(oldName); err != nil {
			logs.Warn("failed to delete stale Event API name during rename",
				logs.String("apiId", existing.ApiId), logs.String("oldName", oldName), logs.Err(err))
		}
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
	if err := s.channelsStore.DeleteByPrefix(apiId + "/"); err != nil {
		logs.Warn("failed to delete channel namespaces during Event API deletion",
			logs.String("apiId", apiId), logs.Err(err))
	}

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

	mergedTags := mergeFn(ns.Tags)
	ns.Tags = mergedTags
	ns.LastModified = time.Now().UTC()

	if err := s.channelsStore.Put(key, ns); err != nil {
		return err
	}

	if len(mergedTags) > 0 {
		if err := s.TagStore.Tag(ns.ChannelNamespaceArn, mergedTags); err != nil {
			logs.Warn("failed to update tags for channel namespace",
				logs.String("arn", ns.ChannelNamespaceArn), logs.Err(err))
		}
	}

	return nil
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
