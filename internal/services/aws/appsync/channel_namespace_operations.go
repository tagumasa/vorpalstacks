package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateChannelNamespace creates a new channel namespace within an Event API (v2).
// POST /v2/apis/{apiId}/channelNamespaces
func (s *AppSyncService) CreateChannelNamespace(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	publishAuthModes, err := parseAuthModes(request.GetArrayParam(req.Parameters, "publishAuthModes"))
	if err != nil {
		return nil, err
	}
	subscribeAuthModes, err := parseAuthModes(request.GetArrayParam(req.Parameters, "subscribeAuthModes"))
	if err != nil {
		return nil, err
	}

	tagMap, err := parseTags(req.Parameters)
	if err != nil {
		return nil, err
	}
	handlerCfgs, err := parseHandlerConfigs(req.Parameters)
	if err != nil {
		return nil, err
	}

	created, tags, err := s.createChannelNamespaceCore(store, createChannelNamespaceInput{
		ApiId:              request.GetStringParam(req.Parameters, "apiId"),
		Name:               request.GetStringParam(req.Parameters, "name"),
		CodeHandlers:       request.GetStringParam(req.Parameters, "codeHandlers"),
		PublishAuthModes:   publishAuthModes,
		SubscribeAuthModes: subscribeAuthModes,
		Tags:               tagMap,
		HandlerConfigs:     handlerCfgs,
	})
	if err != nil {
		return nil, err
	}

	result := channelNamespaceToMap(created)
	if tags != nil {
		result["tags"] = tags
	}

	return map[string]interface{}{
		"channelNamespace": result,
	}, nil
}

// GetChannelNamespace retrieves a channel namespace by API ID and name.
// GET /v2/apis/{apiId}/channelNamespaces/{name}
func (s *AppSyncService) GetChannelNamespace(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	ns, tags, err := s.getChannelNamespaceCore(store, request.GetStringParam(req.Parameters, "apiId"), request.GetStringParam(req.Parameters, "name"))
	if err != nil {
		return nil, err
	}

	result := channelNamespaceToMap(ns)
	if tags != nil {
		result["tags"] = tags
	}

	return map[string]interface{}{
		"channelNamespace": result,
	}, nil
}

// UpdateChannelNamespace updates an existing channel namespace.
// POST /v2/apis/{apiId}/channelNamespaces/{name}
func (s *AppSyncService) UpdateChannelNamespace(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	publishAuthModes, err := parseAuthModes(request.GetArrayParam(req.Parameters, "publishAuthModes"))
	if err != nil {
		return nil, err
	}
	subscribeAuthModes, err := parseAuthModes(request.GetArrayParam(req.Parameters, "subscribeAuthModes"))
	if err != nil {
		return nil, err
	}

	handlerCfgs, err := parseHandlerConfigs(req.Parameters)
	if err != nil {
		return nil, err
	}

	updated, tags, err := s.updateChannelNamespaceCore(store, updateChannelNamespaceInput{
		ApiId:              request.GetStringParam(req.Parameters, "apiId"),
		Name:               request.GetStringParam(req.Parameters, "name"),
		CodeHandlers:       request.GetStringParam(req.Parameters, "codeHandlers"),
		HasCodeHandlers:    request.HasParam(req.Parameters, "codeHandlers"),
		PublishAuthModes:   publishAuthModes,
		SubscribeAuthModes: subscribeAuthModes,
		HandlerConfigs:     handlerCfgs,
	})
	if err != nil {
		return nil, err
	}

	result := channelNamespaceToMap(updated)
	if tags != nil {
		result["tags"] = tags
	}

	return map[string]interface{}{
		"channelNamespace": result,
	}, nil
}

// DeleteChannelNamespace removes a channel namespace.
// DELETE /v2/apis/{apiId}/channelNamespaces/{name}
func (s *AppSyncService) DeleteChannelNamespace(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	if err := s.deleteChannelNamespaceCore(store, request.GetStringParam(req.Parameters, "apiId"), request.GetStringParam(req.Parameters, "name")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListChannelNamespaces returns a paginated list of channel namespaces for an API.
// GET /v2/apis/{apiId}/channelNamespaces
func (s *AppSyncService) ListChannelNamespaces(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")

	entries, nextToken, err := s.listChannelNamespacesCore(store, apiId, request.GetIntParam(req.Parameters, "maxResults"), request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(entries))
	for _, entry := range entries {
		item := channelNamespaceToMap(entry.Namespace)
		if entry.Tags != nil {
			item["tags"] = entry.Tags
		}
		items = append(items, item)
	}

	response := map[string]interface{}{
		"channelNamespaces": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}
