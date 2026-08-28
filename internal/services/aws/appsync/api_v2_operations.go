package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateApi creates a new Event API (v2).
// POST /v2/apis
func (s *AppSyncService) CreateApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	eventConfig, err := parseEventConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	tagMap, err := parseTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	created, tags, err := s.createApiCore(store, createApiInput{
		Name:         request.GetStringParam(req.Parameters, "name"),
		OwnerContact: request.GetStringParam(req.Parameters, "ownerContact"),
		WafWebAclArn: request.GetStringParam(req.Parameters, "wafWebAclArn"),
		XrayEnabled:  request.GetBoolParam(req.Parameters, "xrayEnabled"),
	}, eventConfig, tagMap)
	if err != nil {
		return nil, err
	}

	result := apiToMap(created)
	if tags != nil {
		result["tags"] = tags
	}

	return map[string]interface{}{
		"api": result,
	}, nil
}

// GetApi retrieves an Event API (v2) by its ID.
// GET /v2/apis/{apiId}
func (s *AppSyncService) GetApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	api, tags, err := s.getApiCore(store, request.GetStringParam(req.Parameters, "apiId"))
	if err != nil {
		return nil, err
	}

	result := apiToMap(api)
	if tags != nil {
		result["tags"] = tags
	}

	return map[string]interface{}{
		"api": result,
	}, nil
}

// UpdateApi updates an existing Event API (v2).
// POST /v2/apis/{apiId}
func (s *AppSyncService) UpdateApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	eventConfig, err := parseEventConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	updated, tags, err := s.updateApiCore(store, updateApiInput{
		ApiId:           request.GetStringParam(req.Parameters, "apiId"),
		Name:            request.GetStringParam(req.Parameters, "name"),
		OwnerContact:    request.GetStringParam(req.Parameters, "ownerContact"),
		WafWebAclArn:    request.GetStringParam(req.Parameters, "wafWebAclArn"),
		HasWafWebAclArn: request.HasParam(req.Parameters, "wafWebAclArn"),
		XrayEnabled:     request.GetBoolParam(req.Parameters, "xrayEnabled"),
		HasXrayEnabled:  request.HasParam(req.Parameters, "xrayEnabled"),
	}, eventConfig)
	if err != nil {
		return nil, err
	}

	result := apiToMap(updated)
	if tags != nil {
		result["tags"] = tags
	}

	return map[string]interface{}{
		"api": result,
	}, nil
}

// DeleteApi deletes an Event API (v2).
// DELETE /v2/apis/{apiId}
func (s *AppSyncService) DeleteApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	if err := s.deleteApiCore(store, request.GetStringParam(req.Parameters, "apiId")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListApis returns a paginated list of Event APIs (v2).
// GET /v2/apis
func (s *AppSyncService) ListApis(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	opts, err := parsePaginationOptions(req)
	if err != nil {
		return nil, err
	}

	entries, nextToken, err := s.listApisCore(store, int(opts.MaxItems), opts.Marker)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(entries))
	for _, entry := range entries {
		item := apiToMap(entry.Api)
		if entry.Tags != nil {
			item["tags"] = entry.Tags
		}
		items = append(items, item)
	}

	response := map[string]interface{}{
		"apis": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}
