package appsync

import (
	"fmt"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// createApiCacheInput carries the parsed CreateApiCache request payload.
type createApiCacheInput struct {
	ApiId                    string
	Type                     string
	Ttl                      int64
	ApiCachingBehavior       string
	AtRestEncryptionEnabled  bool
	TransitEncryptionEnabled bool
	HealthMetricsConfig      string
}

// updateApiCacheInput carries the parsed UpdateApiCache request payload.
// The Has* flags distinguish explicitly supplied members from omitted ones;
// type, ttl and apiCachingBehavior are required on UpdateApiCache.
type updateApiCacheInput struct {
	ApiId                 string
	Type                  string
	HasType               bool
	Ttl                   int64
	HasTtl                bool
	ApiCachingBehavior    string
	HasApiCachingBehavior bool
	HealthMetricsConfig   string
}

// createApiCacheCore validates the request and provisions an API cache,
// transitioning it CREATING → AVAILABLE asynchronously.
func (s *AppSyncService) createApiCacheCore(store *appsyncstore.AppSyncStore, in createApiCacheInput) (*appsyncstore.ApiCache, error) {
	if in.ApiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if _, err := store.GetGraphqlApiById(in.ApiId); err != nil {
		return nil, mapStoreErrorE(err)
	}

	if in.Type == "" {
		return nil, NewBadRequestException("type is required")
	}
	if !validateApiCacheType(in.Type) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid cache type: %s", in.Type))
	}
	if in.ApiCachingBehavior == "" {
		return nil, NewBadRequestException("apiCachingBehavior is required")
	}
	if !validateApiCachingBehavior(in.ApiCachingBehavior) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid apiCachingBehavior: %s", in.ApiCachingBehavior))
	}
	if err := validateApiCacheTtl(in.Ttl); err != nil {
		return nil, err
	}

	cache := &appsyncstore.ApiCache{
		Type:                     in.Type,
		Ttl:                      in.Ttl,
		ApiCachingBehavior:       in.ApiCachingBehavior,
		AtRestEncryptionEnabled:  in.AtRestEncryptionEnabled,
		TransitEncryptionEnabled: in.TransitEncryptionEnabled,
		HealthMetricsConfig:      in.HealthMetricsConfig,
		Status:                   "CREATING",
	}
	hmc := cache.HealthMetricsConfig
	if hmc != "" && !validateEnabledDisabled(hmc) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid healthMetricsConfig: %s", hmc))
	}

	if err := store.CreateApiCache(in.ApiId, cache); err != nil {
		return nil, mapStoreErrorE(err)
	}

	// Simulate async cache creation: transition CREATING → AVAILABLE.
	go func() {
		defer func() { resilience.RecoverPanic("appsync cache creation async") }()
		time.Sleep(2 * time.Second)
		cache.Status = "AVAILABLE"
		if err := store.UpdateApiCache(in.ApiId, cache); err != nil {
			logs.Warn("failed to persist cache AVAILABLE status",
				logs.String("apiId", in.ApiId),
				logs.Err(err))
		}
	}()

	return cache, nil
}

// getApiCacheCore fetches the cache configuration of a GraphQL API.
func (s *AppSyncService) getApiCacheCore(store *appsyncstore.AppSyncStore, apiId string) (*appsyncstore.ApiCache, error) {
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	cache, err := store.GetApiCache(apiId)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return cache, nil
}

// updateApiCacheCore applies an update to an existing cache. type, ttl and
// apiCachingBehavior are required members of UpdateApiCache; healthMetricsConfig
// stays presence-based.
func (s *AppSyncService) updateApiCacheCore(store *appsyncstore.AppSyncStore, in updateApiCacheInput) (*appsyncstore.ApiCache, error) {
	if in.ApiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if !in.HasType {
		return nil, NewBadRequestException("type is required")
	}
	if !in.HasTtl {
		return nil, NewBadRequestException("ttl is required")
	}
	if !in.HasApiCachingBehavior {
		return nil, NewBadRequestException("apiCachingBehavior is required")
	}

	if !validateApiCacheType(in.Type) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid cache type: %s", in.Type))
	}
	if err := validateApiCacheTtl(in.Ttl); err != nil {
		return nil, err
	}
	if !validateApiCachingBehavior(in.ApiCachingBehavior) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid apiCachingBehavior: %s", in.ApiCachingBehavior))
	}
	if in.HealthMetricsConfig != "" {
		if !validateEnabledDisabled(in.HealthMetricsConfig) {
			return nil, NewBadRequestException(fmt.Sprintf("Invalid healthMetricsConfig: %s", in.HealthMetricsConfig))
		}
	}

	cache, err := store.GetApiCache(in.ApiId)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	cache.Type = in.Type
	cache.Ttl = in.Ttl
	cache.ApiCachingBehavior = in.ApiCachingBehavior
	if in.HealthMetricsConfig != "" {
		cache.HealthMetricsConfig = in.HealthMetricsConfig
	}

	if err := store.UpdateApiCache(in.ApiId, cache); err != nil {
		return nil, mapStoreErrorE(err)
	}

	return cache, nil
}

// deleteApiCacheCore removes the cache of a GraphQL API.
func (s *AppSyncService) deleteApiCacheCore(store *appsyncstore.AppSyncStore, apiId string) error {
	if apiId == "" {
		return NewBadRequestException("apiId is required")
	}

	if err := store.DeleteApiCache(apiId); err != nil {
		return mapStoreErrorE(err)
	}

	return nil
}

// flushApiCacheCore deletes all cached resolver results for the API and
// invalidates the in-memory schema parse cache so subsequent requests
// re-fetch fresh schema and resolver definitions.
func (s *AppSyncService) flushApiCacheCore(store *appsyncstore.AppSyncStore, apiId string) error {
	if apiId == "" {
		return NewBadRequestException("apiId is required")
	}

	if _, err := store.GetApiCache(apiId); err != nil {
		return mapStoreErrorE(err)
	}

	// Flush all cached resolver results for this API.
	if err := store.FlushResolverCache(apiId); err != nil {
		return mapStoreErrorE(err)
	}

	s.schemaCache.Delete(apiId)

	return nil
}
