package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateApiKey generates a new API key for a GraphQL API.
func (s *AppSyncService) CreateApiKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")

	apiKey, err := s.createApiKeyCore(store, createApiKeyInput{
		ApiId:       apiId,
		Description: request.GetStringParam(req.Parameters, "description"),
		Expires:     request.GetInt64Param(req.Parameters, "expires"),
	})
	if err != nil {
		return nil, err
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

	keys, nextToken, err := s.listApiKeysCore(store, apiId, request.GetIntParam(req.Parameters, "maxResults"), request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
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

	apiKey, err := s.updateApiKeyCore(store, updateApiKeyInput{
		ApiId:       apiId,
		Id:          request.GetStringParam(req.Parameters, "id"),
		Description: request.GetStringParam(req.Parameters, "description"),
		Expires:     request.GetInt64Param(req.Parameters, "expires"),
	})
	if err != nil {
		return nil, err
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

	if err := s.deleteApiKeyCore(store, request.GetStringParam(req.Parameters, "apiId"), request.GetStringParam(req.Parameters, "id")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}
