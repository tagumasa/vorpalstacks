package appsync

import (
	"context"
	"time"

	"github.com/google/uuid"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/services/aws/common/request"
)

// CreateApiKey generates a new API key for a GraphQL API.
func (s *AppSyncService) CreateApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	description := request.GetStringParam(req.Parameters, "description")
	expires := request.GetInt64Param(req.Parameters, "expires")
	if expires == 0 {
		expires = time.Now().Add(365 * 24 * time.Hour).Unix()
	}

	apiKey := &appsyncstore.ApiKey{
		Id:          uuid.New().String(),
		Description: description,
		Expires:     expires,
		// AWS returns Deletes equal to Expires at creation time.
		Deletes: expires,
	}

	if err := store.CreateApiKey(apiId, apiKey); err != nil {
		return mapStoreError(err)
	}

	result := apiKeyToMap(apiKey)
	result["apiId"] = apiId
	return map[string]interface{}{"apiKey": result}, nil
}

// ListApiKeys lists API keys for a GraphQL API.
func (s *AppSyncService) ListApiKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	opts := parsePaginationOptions(req)
	keys, nextToken, err := store.ListApiKeys(apiId, opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]map[string]interface{}, 0, len(keys))
	for _, k := range keys {
		m := apiKeyToMap(k)
		m["apiId"] = apiId
		items = append(items, m)
	}

	response := map[string]interface{}{"apiKeys": items}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

// UpdateApiKey updates an existing API key.
func (s *AppSyncService) UpdateApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	id := request.GetStringParam(req.Parameters, "id")
	if id == "" {
		return nil, NewBadRequestException("id is required")
	}

	apiKey, err := store.GetApiKey(apiId, id)
	if err != nil {
		return mapStoreError(err)
	}

	description := request.GetStringParam(req.Parameters, "description")
	expires := request.GetInt64Param(req.Parameters, "expires")

	if description != "" {
		apiKey.Description = description
	}
	if expires != 0 {
		apiKey.Expires = expires
		// Deleting an API key happens after the expiry TTL; sync Deletes with Expires.
		apiKey.Deletes = expires
	}

	if err := store.UpdateApiKey(apiId, apiKey); err != nil {
		return mapStoreError(err)
	}

	result := apiKeyToMap(apiKey)
	result["apiId"] = apiId
	return map[string]interface{}{"apiKey": result}, nil
}

// DeleteApiKey deletes an API key.
func (s *AppSyncService) DeleteApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	id := request.GetStringParam(req.Parameters, "id")
	if id == "" {
		return nil, NewBadRequestException("id is required")
	}

	if err := store.DeleteApiKey(apiId, id); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// apiKeyToMap converts an API key to a serialisable map.
// The deletes field is only emitted when non-zero, matching AWS behaviour.
func apiKeyToMap(k *appsyncstore.ApiKey) map[string]interface{} {
	result := map[string]interface{}{
		"id": k.Id,
	}
	if k.Description != "" {
		result["description"] = k.Description
	}
	if k.Expires != 0 {
		result["expires"] = k.Expires
	}
	if k.Deletes != 0 {
		result["deletes"] = k.Deletes
	}
	return result
}
