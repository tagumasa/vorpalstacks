package appsync

import (
	"context"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/utils/timeutils"
)

// CreateChannelNamespace creates a new channel namespace within an Event API (v2).
// POST /v2/apis/{apiId}/channelNamespaces
func (s *AppSyncService) CreateChannelNamespace(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	ns := &appsyncstore.ChannelNamespace{
		ApiId:              apiId,
		Name:               name,
		CodeHandlers:       request.GetStringParam(req.Parameters, "codeHandlers"),
		HandlerConfigs:     parseHandlerConfigs(req.Parameters),
		PublishAuthModes:   parseAuthModes(request.GetArrayParam(req.Parameters, "publishAuthModes")),
		SubscribeAuthModes: parseAuthModes(request.GetArrayParam(req.Parameters, "subscribeAuthModes")),
		Tags:               parseTags(req.Parameters),
	}

	created, err := store.CreateChannelNamespace(ns)
	if err != nil {
		return mapStoreError(err)
	}

	if len(ns.Tags) > 0 {
		tagMap := make(map[string]string, len(ns.Tags))
		for k, v := range ns.Tags {
			tagMap[k] = v
		}
		if err := store.TagStore.Tag(created.ChannelNamespaceArn, tagMap); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"channelNamespace": channelNamespaceToMap(created),
	}, nil
}

// GetChannelNamespace retrieves a channel namespace by API ID and name.
// GET /v2/apis/{apiId}/channelNamespaces/{name}
func (s *AppSyncService) GetChannelNamespace(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	name := request.GetStringParam(req.Parameters, "name")
	if apiId == "" || name == "" {
		return nil, NewBadRequestException("apiId and name are required")
	}

	ns := &appsyncstore.ChannelNamespace{
		ApiId:              apiId,
		Name:               name,
		CodeHandlers:       request.GetStringParam(req.Parameters, "codeHandlers"),
		HandlerConfigs:     parseHandlerConfigs(req.Parameters),
		PublishAuthModes:   parseAuthModes(request.GetArrayParam(req.Parameters, "publishAuthModes")),
		SubscribeAuthModes: parseAuthModes(request.GetArrayParam(req.Parameters, "subscribeAuthModes")),
	}

	updated, err := store.UpdateChannelNamespace(ns)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"channelNamespace": channelNamespaceToMap(updated),
	}, nil
}

// DeleteChannelNamespace removes a channel namespace.
// DELETE /v2/apis/{apiId}/channelNamespaces/{name}
func (s *AppSyncService) DeleteChannelNamespace(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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

	if err := store.DeleteChannelNamespace(apiId, name); err != nil {
		return mapStoreError(err)
	}

	store.TagStore.Delete(ns.ChannelNamespaceArn)

	return map[string]interface{}{}, nil
}

// ListChannelNamespaces returns a paginated list of channel namespaces for an API.
// GET /v2/apis/{apiId}/channelNamespaces
func (s *AppSyncService) ListChannelNamespaces(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	opts := parsePaginationOptions(req)
	namespaces, nextToken, err := store.ListChannelNamespaces(apiId, opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(namespaces))
	for _, ns := range namespaces {
		items = append(items, channelNamespaceToMap(ns))
	}

	response := map[string]interface{}{
		"channelNamespaces": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

// channelNamespaceToMap converts a ChannelNamespace struct to a response map
// with correct wire format. Timestamps are serialised as epoch seconds.
func channelNamespaceToMap(ns *appsyncstore.ChannelNamespace) map[string]interface{} {
	m := map[string]interface{}{
		"apiId":               ns.ApiId,
		"name":                ns.Name,
		"channelNamespaceArn": ns.ChannelNamespaceArn,
		"created":             timeutils.FormatEpochSeconds(ns.Created),
		"lastModified":        timeutils.FormatEpochSeconds(ns.LastModified),
	}

	if ns.CodeHandlers != "" {
		m["codeHandlers"] = ns.CodeHandlers
	}
	if ns.HandlerConfigs != nil {
		m["handlerConfigs"] = handlerConfigsToMap(ns.HandlerConfigs)
	}
	if len(ns.PublishAuthModes) > 0 {
		m["publishAuthModes"] = authModesToMap(ns.PublishAuthModes)
	}
	if len(ns.SubscribeAuthModes) > 0 {
		m["subscribeAuthModes"] = authModesToMap(ns.SubscribeAuthModes)
	}
	if len(ns.Tags) > 0 {
		m["tags"] = ns.Tags
	}

	return m
}

// handlerConfigsToMap converts HandlerConfigs to a wire-format map.
func handlerConfigsToMap(hc *appsyncstore.HandlerConfigs) map[string]interface{} {
	m := map[string]interface{}{}
	if hc.OnPublish != nil {
		m["onPublish"] = handlerConfigToMap(hc.OnPublish)
	}
	if hc.OnSubscribe != nil {
		m["onSubscribe"] = handlerConfigToMap(hc.OnSubscribe)
	}
	return m
}

// handlerConfigToMap converts a HandlerConfig to a wire-format map.
func handlerConfigToMap(hc *appsyncstore.HandlerConfig) map[string]interface{} {
	m := map[string]interface{}{
		"behavior": hc.Behavior,
	}
	if hc.Integration != nil {
		integration := map[string]interface{}{
			"dataSourceName": hc.Integration.DataSourceName,
		}
		if hc.Integration.LambdaConfig != nil {
			integration["lambdaConfig"] = map[string]interface{}{
				"invokeType": hc.Integration.LambdaConfig.InvokeType,
			}
		}
		m["integration"] = integration
	}
	return m
}
