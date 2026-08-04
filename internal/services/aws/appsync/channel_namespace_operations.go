package appsync

import (
	"context"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// CreateChannelNamespace creates a new channel namespace within an Event API (v2).
// POST /v2/apis/{apiId}/channelNamespaces
func (s *AppSyncService) CreateChannelNamespace(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}
	if err := validateEventApiExists(store, apiId); err != nil {
		return nil, err
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if err := validateNamespace(name); err != nil {
		return nil, err
	}

	publishAuthModes, err := parseAuthModes(request.GetArrayParam(req.Parameters, "publishAuthModes"))
	if err != nil {
		return nil, err
	}
	subscribeAuthModes, err := parseAuthModes(request.GetArrayParam(req.Parameters, "subscribeAuthModes"))
	if err != nil {
		return nil, err
	}

	codeHandlers := request.GetStringParam(req.Parameters, "codeHandlers")
	if codeHandlers != "" {
		if err := validateCode(codeHandlers); err != nil {
			return nil, err
		}
	}

	tagMap, err := parseTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	ns := &appsyncstore.ChannelNamespace{
		ApiId:              apiId,
		Name:               name,
		CodeHandlers:       codeHandlers,
		HandlerConfigs:     parseHandlerConfigs(req.Parameters),
		PublishAuthModes:   publishAuthModes,
		SubscribeAuthModes: subscribeAuthModes,
		Tags:               tagMap,
	}

	created, err := store.CreateChannelNamespace(ns)
	if err != nil {
		return mapStoreError(err)
	}

	if len(created.Tags) > 0 {
		tagMap := make(map[string]string, len(created.Tags))
		for k, v := range created.Tags {
			tagMap[k] = v
		}
		if err := store.TagStore.Tag(created.ChannelNamespaceArn, tagMap); err != nil {
			return nil, err
		}
	}

	result := channelNamespaceToMap(created)
	if tags, err := store.TagStore.List(created.ChannelNamespaceArn); err == nil && len(tags) > 0 {
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

	apiId := request.GetStringParam(req.Parameters, "apiId")
	name := request.GetStringParam(req.Parameters, "name")
	if apiId == "" || name == "" {
		return nil, NewBadRequestException("apiId and name are required")
	}

	ns, err := store.GetChannelNamespace(apiId, name)
	if err != nil {
		return mapStoreError(err)
	}

	result := channelNamespaceToMap(ns)
	if tagsFromStore, err := store.TagStore.List(ns.ChannelNamespaceArn); err == nil && len(tagsFromStore) > 0 {
		result["tags"] = tagsFromStore
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

	apiId := request.GetStringParam(req.Parameters, "apiId")
	name := request.GetStringParam(req.Parameters, "name")
	if apiId == "" || name == "" {
		return nil, NewBadRequestException("apiId and name are required")
	}
	if err := validateNamespace(name); err != nil {
		return nil, err
	}

	publishAuthModes, err := parseAuthModes(request.GetArrayParam(req.Parameters, "publishAuthModes"))
	if err != nil {
		return nil, err
	}
	subscribeAuthModes, err := parseAuthModes(request.GetArrayParam(req.Parameters, "subscribeAuthModes"))
	if err != nil {
		return nil, err
	}

	codeHandlers := request.GetStringParam(req.Parameters, "codeHandlers")
	if codeHandlers != "" {
		if err := validateCode(codeHandlers); err != nil {
			return nil, err
		}
	}

	ns := &appsyncstore.ChannelNamespace{
		ApiId:              apiId,
		Name:               name,
		CodeHandlers:       codeHandlers,
		HandlerConfigs:     parseHandlerConfigs(req.Parameters),
		PublishAuthModes:   publishAuthModes,
		SubscribeAuthModes: subscribeAuthModes,
	}

	updated, err := store.UpdateChannelNamespace(ns)
	if err != nil {
		return mapStoreError(err)
	}

	result := channelNamespaceToMap(updated)
	if tags, err := store.TagStore.List(updated.ChannelNamespaceArn); err == nil && len(tags) > 0 {
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

	apiId := request.GetStringParam(req.Parameters, "apiId")
	name := request.GetStringParam(req.Parameters, "name")
	if apiId == "" || name == "" {
		return nil, NewBadRequestException("apiId and name are required")
	}

	if _, err := store.GetChannelNamespace(apiId, name); err != nil {
		return mapStoreError(err)
	}

	if err := store.DeleteChannelNamespace(apiId, name); err != nil {
		return mapStoreError(err)
	}

	// Clean up active subscriptions on the deleted namespace to prevent
	// stale data delivery.
	s.eventServer.RemoveSubscriptionsByNamespace(name)

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
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	opts, err := parsePaginationOptions(req)
	if err != nil {
		return nil, err
	}
	namespaces, nextToken, err := store.ListChannelNamespaces(apiId, opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(namespaces))
	for _, ns := range namespaces {
		item := channelNamespaceToMap(ns)
		if tags, err := store.TagStore.List(ns.ChannelNamespaceArn); err == nil && len(tags) > 0 {
			item["tags"] = tags
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
