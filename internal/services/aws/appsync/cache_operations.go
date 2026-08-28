package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateApiCache creates a cache for a GraphQL API.
func (s *AppSyncService) CreateApiCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	cache, err := s.createApiCacheCore(store, createApiCacheInput{
		ApiId:                    request.GetStringParam(req.Parameters, "apiId"),
		Type:                     request.GetStringParam(req.Parameters, "type"),
		Ttl:                      request.GetInt64Param(req.Parameters, "ttl"),
		ApiCachingBehavior:       request.GetStringParam(req.Parameters, "apiCachingBehavior"),
		AtRestEncryptionEnabled:  request.GetBoolParam(req.Parameters, "atRestEncryptionEnabled"),
		TransitEncryptionEnabled: request.GetBoolParam(req.Parameters, "transitEncryptionEnabled"),
		HealthMetricsConfig:      request.GetStringParam(req.Parameters, "healthMetricsConfig"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"apiCache": apiCacheToMap(cache)}, nil
}

// GetApiCache retrieves the cache configuration for a GraphQL API.
func (s *AppSyncService) GetApiCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	cache, err := s.getApiCacheCore(store, request.GetStringParam(req.Parameters, "apiId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"apiCache": apiCacheToMap(cache)}, nil
}

// UpdateApiCache updates the cache configuration for a GraphQL API.
func (s *AppSyncService) UpdateApiCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	cache, err := s.updateApiCacheCore(store, updateApiCacheInput{
		ApiId:                 request.GetStringParam(req.Parameters, "apiId"),
		Type:                  request.GetStringParam(req.Parameters, "type"),
		HasType:               request.HasParam(req.Parameters, "type"),
		Ttl:                   request.GetInt64Param(req.Parameters, "ttl"),
		HasTtl:                request.HasParam(req.Parameters, "ttl"),
		ApiCachingBehavior:    request.GetStringParam(req.Parameters, "apiCachingBehavior"),
		HasApiCachingBehavior: request.HasParam(req.Parameters, "apiCachingBehavior"),
		HealthMetricsConfig:   request.GetStringParam(req.Parameters, "healthMetricsConfig"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"apiCache": apiCacheToMap(cache)}, nil
}

// DeleteApiCache deletes the cache for a GraphQL API.
func (s *AppSyncService) DeleteApiCache(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	if err := s.deleteApiCacheCore(store, request.GetStringParam(req.Parameters, "apiId")); err != nil {
		return nil, err
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

	if err := s.flushApiCacheCore(store, request.GetStringParam(req.Parameters, "apiId")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}
