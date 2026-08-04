package appsync

import (
	"fmt"
	"net/http"
	"sync"

	"vorpalstacks/internal/common/auth"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// AppSyncService provides AWS AppSync service operations for vorpalstacks.
// Implements control-plane CRUD for Event APIs (v2), GraphQL APIs (v1),
// channel namespaces, data sources, resolvers, GraphQL execution, and tag operations.
type AppSyncService struct {
	accountID      string
	stores         sync.Map
	bus            eventbus.Bus
	schemaCache    sync.Map
	schemaWg       sync.WaitGroup
	eventServer    *EventServer
	storageManager *storage.RegionStorageManager
}

// NewAppSyncService creates a new AppSync service instance scoped to the given account.
func NewAppSyncService(accountID string) *AppSyncService {
	return &AppSyncService{
		accountID:   accountID,
		eventServer: NewEventServer(),
	}
}

// SetEventBus injects the global event bus for WebSocket pub/sub fan-out
// and cross-service event delivery.
func (s *AppSyncService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
	s.eventServer.SetEventBus(bus)
	s.eventServer.SetStoreLookup(s.lookupStoreByApiId)
}

// SetSigVerifier injects the SigV4 verifier for AWS_IAM auth mode
// enforcement on the Event API data-plane (P-5).
func (s *AppSyncService) SetSigVerifier(v *auth.SignatureV4Verifier) {
	s.eventServer.SetSigVerifier(v)
}

// lookupStoreByApiId searches all cached regional stores for one that
// contains the given API ID. Used by EventServer for auth enforcement.
func (s *AppSyncService) lookupStoreByApiId(apiId string) (*appsyncstore.AppSyncStore, error) {
	var found *appsyncstore.AppSyncStore
	s.stores.Range(func(_, v interface{}) bool {
		store := v.(*appsyncstore.AppSyncStore)
		if _, err := store.GetApiById(apiId); err == nil {
			found = store
			return false
		}
		if _, err := store.GetGraphqlApiById(apiId); err == nil {
			found = store
			return false
		}
		return true
	})
	if found == nil {
		return nil, fmt.Errorf("api %s not found in any region", apiId)
	}
	return found, nil
}

// SetStorageManager injects the storage manager for admin console store access.
func (s *AppSyncService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// EventServerHandler returns an http.Handler for the AppSync events server
// (WebSocket + HTTP publish), or nil if not initialised.
func (s *AppSyncService) EventServerHandler() http.Handler {
	return s.eventServer
}

// ShutdownEventServer closes all active WebSocket connections and waits
// for in-flight schema creation goroutines to finish.
func (s *AppSyncService) ShutdownEventServer() {
	s.eventServer.Shutdown()
	s.schemaWg.Wait()
}

// GetStoreForRegion returns the cached AppSync store for the given region,
// creating one if not already cached. Used by both HTTP handlers and the
// admin console to ensure a single store instance per region.
func (s *AppSyncService) GetStoreForRegion(region string) (*appsyncstore.AppSyncStore, error) {
	if cached, ok := s.stores.Load(region); ok {
		return cached.(*appsyncstore.AppSyncStore), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("appsync storage manager not initialised")
	}
	rs, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	store := appsyncstore.NewAppSyncStore(rs, s.accountID, region)
	actual, _ := s.stores.LoadOrStore(region, store)
	return actual.(*appsyncstore.AppSyncStore), nil
}

func (s *AppSyncService) store(reqCtx *request.RequestContext) (*appsyncstore.AppSyncStore, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (*appsyncstore.AppSyncStore, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		return appsyncstore.NewAppSyncStore(storage, s.accountID, reqCtx.GetRegion()), nil
	})
}

// RegisterHandlers registers all AppSync control-plane operation handlers with the dispatcher.
// Event API (v2) operations and tag operations.
// GraphQL API (v1) core — data sources, resolvers, functions, types, schema.
func (s *AppSyncService) RegisterHandlers(d handler.Registrar) {
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
	// from graphql_handler.go for POST /v1/apis/{apiId}/graphql
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
