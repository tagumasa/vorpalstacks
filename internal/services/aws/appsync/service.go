package appsync

import (
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/eventbus"
	awscommon "vorpalstacks/internal/services/aws/common"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
)

// AppSyncService provides AWS AppSync service operations for vorpalstacks.
// Implements control-plane CRUD for Event APIs (v2), GraphQL APIs (v1),
// channel namespaces, data sources, resolvers, GraphQL execution, and tag operations.
type AppSyncService struct {
	accountID     string
	stores        sync.Map
	lambdaInvoker awscommon.LambdaInvoker
	bus           eventbus.Bus
}

// NewAppSyncService creates a new AppSync service instance scoped to the given account.
func NewAppSyncService(accountID string) *AppSyncService {
	return &AppSyncService{
		accountID: accountID,
	}
}

// SetLambdaInvoker injects a Lambda invoker for Lambda DataSource calls
// within GraphQL resolvers. Follows the same pattern as Step Functions,
// EventBridge, SNS, and other services that invoke Lambda.
func (s *AppSyncService) SetLambdaInvoker(invoker awscommon.LambdaInvoker) {
	s.lambdaInvoker = invoker
}

// SetEventBus injects the global event bus for WebSocket pub/sub fan-out
// and cross-service event delivery.
func (s *AppSyncService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
}

// GetStoreForRegion returns the cached AppSync store for the given region,
// creating one if not already cached. Used by both HTTP handlers and the
// admin console to ensure a single store instance per region.
func (s *AppSyncService) GetStoreForRegion(region string, rs storage.BasicStorage) *appsyncstore.AppSyncStore {
	key := region
	if cached, ok := s.stores.Load(key); ok {
		return cached.(*appsyncstore.AppSyncStore)
	}
	store := appsyncstore.NewAppSyncStore(rs, s.accountID, region)
	actual, _ := s.stores.LoadOrStore(key, store)
	return actual.(*appsyncstore.AppSyncStore)
}

// store resolves the AppSync store for the current request context.
// Uses sync.Map for per-region caching of store instances.
func (s *AppSyncService) store(reqCtx *request.RequestContext) (*appsyncstore.AppSyncStore, error) {
	region := reqCtx.GetRegion()
	key := region
	if cached, ok := s.stores.Load(key); ok {
		return cached.(*appsyncstore.AppSyncStore), nil
	}

	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}

	store := appsyncstore.NewAppSyncStore(storage, s.accountID, region)
	actual, _ := s.stores.LoadOrStore(key, store)
	return actual.(*appsyncstore.AppSyncStore), nil
}

// RegisterHandlers registers all AppSync control-plane operation handlers with the dispatcher.
// Phase 1: Event API (v2) operations + tag operations.
// Phase 2: GraphQL API (v1) core — data sources, resolvers, functions, types, schema.
func (s *AppSyncService) RegisterHandlers(d *dispatcher.Dispatcher) {
	// Event API (v2) operations
	d.RegisterHandlerForService("appsync", "CreateApi", s.CreateApi)
	d.RegisterHandlerForService("appsync", "GetApi", s.GetApi)
	d.RegisterHandlerForService("appsync", "UpdateApi", s.UpdateApi)
	d.RegisterHandlerForService("appsync", "DeleteApi", s.DeleteApi)
	d.RegisterHandlerForService("appsync", "ListApis", s.ListApis)

	// Channel namespace operations
	d.RegisterHandlerForService("appsync", "CreateChannelNamespace", s.CreateChannelNamespace)
	d.RegisterHandlerForService("appsync", "GetChannelNamespace", s.GetChannelNamespace)
	d.RegisterHandlerForService("appsync", "UpdateChannelNamespace", s.UpdateChannelNamespace)
	d.RegisterHandlerForService("appsync", "DeleteChannelNamespace", s.DeleteChannelNamespace)
	d.RegisterHandlerForService("appsync", "ListChannelNamespaces", s.ListChannelNamespaces)

	// GraphQL API (v1) operations
	d.RegisterHandlerForService("appsync", "CreateGraphqlApi", s.CreateGraphqlApi)
	d.RegisterHandlerForService("appsync", "GetGraphqlApi", s.GetGraphqlApi)
	d.RegisterHandlerForService("appsync", "UpdateGraphqlApi", s.UpdateGraphqlApi)
	d.RegisterHandlerForService("appsync", "DeleteGraphqlApi", s.DeleteGraphqlApi)
	d.RegisterHandlerForService("appsync", "ListGraphqlApis", s.ListGraphqlApis)

	// Data source operations
	d.RegisterHandlerForService("appsync", "CreateDataSource", s.CreateDataSource)
	d.RegisterHandlerForService("appsync", "GetDataSource", s.GetDataSource)
	d.RegisterHandlerForService("appsync", "UpdateDataSource", s.UpdateDataSource)
	d.RegisterHandlerForService("appsync", "DeleteDataSource", s.DeleteDataSource)
	d.RegisterHandlerForService("appsync", "ListDataSources", s.ListDataSources)

	// Resolver operations
	d.RegisterHandlerForService("appsync", "CreateResolver", s.CreateResolver)
	d.RegisterHandlerForService("appsync", "GetResolver", s.GetResolver)
	d.RegisterHandlerForService("appsync", "UpdateResolver", s.UpdateResolver)
	d.RegisterHandlerForService("appsync", "DeleteResolver", s.DeleteResolver)
	d.RegisterHandlerForService("appsync", "ListResolvers", s.ListResolvers)
	d.RegisterHandlerForService("appsync", "ListResolversByFunction", s.ListResolversByFunction)

	// Function (AppSync) operations
	d.RegisterHandlerForService("appsync", "CreateFunction", s.CreateFunction)
	d.RegisterHandlerForService("appsync", "GetFunction", s.GetFunction)
	d.RegisterHandlerForService("appsync", "UpdateFunction", s.UpdateFunction)
	d.RegisterHandlerForService("appsync", "DeleteFunction", s.DeleteFunction)
	d.RegisterHandlerForService("appsync", "ListFunctions", s.ListFunctions)

	// Type operations
	d.RegisterHandlerForService("appsync", "CreateType", s.CreateType)
	d.RegisterHandlerForService("appsync", "GetType", s.GetType)
	d.RegisterHandlerForService("appsync", "UpdateType", s.UpdateType)
	d.RegisterHandlerForService("appsync", "DeleteType", s.DeleteType)
	d.RegisterHandlerForService("appsync", "ListTypes", s.ListTypes)

	// Schema operations
	d.RegisterHandlerForService("appsync", "StartSchemaCreation", s.StartSchemaCreation)
	d.RegisterHandlerForService("appsync", "GetSchemaCreationStatus", s.GetSchemaCreationStatus)
	d.RegisterHandlerForService("appsync", "GetIntrospectionSchema", s.GetIntrospectionSchema)

	// Data source introspection (501)
	d.RegisterHandlerForService("appsync", "StartDataSourceIntrospection", s.StartDataSourceIntrospection)
	d.RegisterHandlerForService("appsync", "GetDataSourceIntrospection", s.GetDataSourceIntrospection)

	// Merged API — ListTypesByAssociation (501)
	d.RegisterHandlerForService("appsync", "ListTypesByAssociation", s.ListTypesByAssociation)

	// Code evaluation (501)
	d.RegisterHandlerForService("appsync", "EvaluateCode", s.EvaluateCode)
	d.RegisterHandlerForService("appsync", "EvaluateMappingTemplate", s.EvaluateMappingTemplate)

	// GraphQL execution endpoint — dispatched via sentinel operation name
	// from appsync_parser.go for POST /v1/apis/{apiId}/graphql
	d.RegisterHandlerForService("appsync", "GraphQLExecution", s.HandleGraphQLExecution)

	// Environment variable operations
	d.RegisterHandlerForService("appsync", "GetGraphqlApiEnvironmentVariables", s.GetGraphqlApiEnvironmentVariables)
	d.RegisterHandlerForService("appsync", "PutGraphqlApiEnvironmentVariables", s.PutGraphqlApiEnvironmentVariables)

	// API key operations
	d.RegisterHandlerForService("appsync", "CreateApiKey", s.CreateApiKey)
	d.RegisterHandlerForService("appsync", "ListApiKeys", s.ListApiKeys)
	d.RegisterHandlerForService("appsync", "UpdateApiKey", s.UpdateApiKey)
	d.RegisterHandlerForService("appsync", "DeleteApiKey", s.DeleteApiKey)

	// Cache operations
	d.RegisterHandlerForService("appsync", "CreateApiCache", s.CreateApiCache)
	d.RegisterHandlerForService("appsync", "GetApiCache", s.GetApiCache)
	d.RegisterHandlerForService("appsync", "UpdateApiCache", s.UpdateApiCache)
	d.RegisterHandlerForService("appsync", "DeleteApiCache", s.DeleteApiCache)
	d.RegisterHandlerForService("appsync", "FlushApiCache", s.FlushApiCache)

	// Domain name and association operations
	d.RegisterHandlerForService("appsync", "CreateDomainName", s.CreateDomainName)
	d.RegisterHandlerForService("appsync", "ListDomainNames", s.ListDomainNames)
	d.RegisterHandlerForService("appsync", "GetDomainName", s.GetDomainName)
	d.RegisterHandlerForService("appsync", "UpdateDomainName", s.UpdateDomainName)
	d.RegisterHandlerForService("appsync", "DeleteDomainName", s.DeleteDomainName)
	d.RegisterHandlerForService("appsync", "AssociateApi", s.AssociateApi)
	d.RegisterHandlerForService("appsync", "DisassociateApi", s.DisassociateApi)
	d.RegisterHandlerForService("appsync", "GetApiAssociation", s.GetApiAssociation)

	// Merged API association operations
	d.RegisterHandlerForService("appsync", "AssociateSourceGraphqlApi", s.AssociateSourceGraphqlApi)
	d.RegisterHandlerForService("appsync", "GetSourceApiAssociation", s.GetSourceApiAssociation)
	d.RegisterHandlerForService("appsync", "UpdateSourceApiAssociation", s.UpdateSourceApiAssociation)
	d.RegisterHandlerForService("appsync", "DisassociateSourceGraphqlApi", s.DisassociateSourceGraphqlApi)
	d.RegisterHandlerForService("appsync", "StartSchemaMerge", s.StartSchemaMerge)
	d.RegisterHandlerForService("appsync", "AssociateMergedGraphqlApi", s.AssociateMergedGraphqlApi)
	d.RegisterHandlerForService("appsync", "DisassociateMergedGraphqlApi", s.DisassociateMergedGraphqlApi)
	d.RegisterHandlerForService("appsync", "ListSourceApiAssociations", s.ListSourceApiAssociations)

	// Tag operations
	d.RegisterHandlerForService("appsync", "TagResource", s.TagResource)
	d.RegisterHandlerForService("appsync", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("appsync", "ListTagsForResource", s.ListTagsForResource)
}
