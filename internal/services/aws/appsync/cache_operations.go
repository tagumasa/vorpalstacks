package appsync

import (
	"context"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// CreateApiCache creates a cache for a GraphQL API.
func (s *AppSyncService) CreateApiCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if _, err := store.GetGraphqlApiById(apiId); err != nil {
		return mapStoreError(err)
	}

	cacheType := request.GetStringParam(req.Parameters, "type")
	if cacheType == "" {
		return nil, NewBadRequestException("type is required")
	}
	ttl := request.GetInt64Param(req.Parameters, "ttl")
	apiCachingBehavior := request.GetStringParam(req.Parameters, "apiCachingBehavior")
	if apiCachingBehavior == "" {
		return nil, NewBadRequestException("apiCachingBehavior is required")
	}

	cache := &appsyncstore.ApiCache{
		Type:                     cacheType,
		Ttl:                      ttl,
		ApiCachingBehavior:       apiCachingBehavior,
		AtRestEncryptionEnabled:  request.GetBoolParam(req.Parameters, "atRestEncryptionEnabled"),
		TransitEncryptionEnabled: request.GetBoolParam(req.Parameters, "transitEncryptionEnabled"),
		HealthMetricsConfig:      request.GetStringParam(req.Parameters, "healthMetricsConfig"),
		Status:                   "CREATING",
	}

	if err := store.CreateApiCache(apiId, cache); err != nil {
		return mapStoreError(err)
	}

	// Simulate async cache creation: transition CREATING → AVAILABLE.
	go func() {
		defer func() { resilience.RecoverPanic("appsync cache creation async") }()
		time.Sleep(2 * time.Second)
		cache.Status = "AVAILABLE"
		if err := store.UpdateApiCache(apiId, cache); err != nil {
			logs.Warn("failed to persist cache AVAILABLE status",
				logs.String("apiId", apiId),
				logs.Err(err))
		}
	}()

	return map[string]interface{}{"apiCache": apiCacheToMap(cache)}, nil
}

// GetApiCache retrieves the cache configuration for a GraphQL API.
func (s *AppSyncService) GetApiCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	cache, err := store.GetApiCache(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"apiCache": apiCacheToMap(cache)}, nil
}

// UpdateApiCache updates the cache configuration for a GraphQL API.
func (s *AppSyncService) UpdateApiCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	cache, err := store.GetApiCache(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	cacheType := request.GetStringParam(req.Parameters, "type")
	ttl := request.GetInt64Param(req.Parameters, "ttl")
	apiCachingBehavior := request.GetStringParam(req.Parameters, "apiCachingBehavior")

	if cacheType != "" {
		cache.Type = cacheType
	}
	// Use HasParam to distinguish "ttl not provided" from "ttl explicitly set to 0".
	if request.HasParam(req.Parameters, "ttl") {
		cache.Ttl = ttl
	}
	if apiCachingBehavior != "" {
		cache.ApiCachingBehavior = apiCachingBehavior
	}
	healthMetrics := request.GetStringParam(req.Parameters, "healthMetricsConfig")
	if healthMetrics != "" {
		cache.HealthMetricsConfig = healthMetrics
	}

	if err := store.UpdateApiCache(apiId, cache); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"apiCache": apiCacheToMap(cache)}, nil
}

// DeleteApiCache deletes the cache for a GraphQL API.
func (s *AppSyncService) DeleteApiCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if err := store.DeleteApiCache(apiId); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// FlushApiCache flushes the cache for a GraphQL API.
// Deletes all cached resolver results and invalidates the schema parse cache.
func (s *AppSyncService) FlushApiCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if _, err := store.GetApiCache(apiId); err != nil {
		return mapStoreError(err)
	}

	// Flush all cached resolver results for this API.
	if err := store.FlushResolverCache(apiId); err != nil {
		return mapStoreError(err)
	}

	// Invalidate the in-memory schema parse cache so subsequent requests
	// re-fetch fresh schema and resolver definitions.
	s.schemaCache.Delete(apiId)

	return map[string]interface{}{}, nil
}
