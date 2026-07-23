package appsync

import (
	"context"
	"errors"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// CreateApi creates a new Event API (v2).
// POST /v2/apis
func (s *AppSyncService) CreateApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	eventConfig, err := parseEventConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	// Check combined API count quota (GraphQL + Event APIs per region).
	graphqlCount, _ := store.CountGraphqlApis()
	eventCount, _ := store.CountApis()
	if graphqlCount+eventCount >= maxApisPerRegion {
		return nil, ErrApiLimitExceededException
	}

	api := &appsyncstore.Api{
		Name:         name,
		EventConfig:  eventConfig,
		OwnerContact: request.GetStringParam(req.Parameters, "ownerContact"),
		Tags:         parseTags(req.Parameters),
		WafWebAclArn: request.GetStringParam(req.Parameters, "wafWebAclArn"),
		XrayEnabled:  request.GetBoolParam(req.Parameters, "xrayEnabled"),
	}

	created, err := store.CreateApi(api)
	if err != nil {
		return mapStoreError(err)
	}

	if len(created.Tags) > 0 {
		tagMap := make(map[string]string, len(created.Tags))
		for k, v := range created.Tags {
			tagMap[k] = v
		}
		if err := store.TagStore.Tag(created.Arn, tagMap); err != nil {
			return nil, err
		}
	}

	result := apiToMap(created)
	if tags, err := store.TagStore.List(created.Arn); err == nil && len(tags) > 0 {
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

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	api, err := store.GetApiById(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	result := apiToMap(api)
	if tagsFromStore, err := store.TagStore.List(api.Arn); err == nil && len(tagsFromStore) > 0 {
		result["tags"] = tagsFromStore
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

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	// Per Smithy model, name is required for UpdateApiRequest.
	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	// Fetch existing to preserve fields that were not provided in the request.
	// Without this, WafWebAclArn and XrayEnabled would be overwritten with
	// Go zero values on every update call that omits them.
	existing, err := store.GetApiById(apiId)
	if err != nil {
		return mapStoreError(err)
	}

	wafWebAclArn := existing.WafWebAclArn
	if request.HasParam(req.Parameters, "wafWebAclArn") {
		wafWebAclArn = request.GetStringParam(req.Parameters, "wafWebAclArn")
	}

	xrayEnabled := existing.XrayEnabled
	if request.HasParam(req.Parameters, "xrayEnabled") {
		xrayEnabled = request.GetBoolParam(req.Parameters, "xrayEnabled")
	}

	api := &appsyncstore.Api{
		Name:         request.GetStringParam(req.Parameters, "name"),
		OwnerContact: request.GetStringParam(req.Parameters, "ownerContact"),
		WafWebAclArn: wafWebAclArn,
		XrayEnabled:  xrayEnabled,
	}

	// Per Smithy UpdateApiRequest, eventConfig is @required.
	if !request.HasParam(req.Parameters, "eventConfig") {
		return nil, NewBadRequestException("eventConfig is required")
	}
	eventConfig, err := parseEventConfig(req.Parameters)
	if err != nil {
		return nil, err
	}
	api.EventConfig = eventConfig

	updated, err := store.UpdateApiById(apiId, api)
	if err != nil {
		return mapStoreError(err)
	}

	result := apiToMap(updated)
	if tags, err := store.TagStore.List(updated.Arn); err == nil && len(tags) > 0 {
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

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if _, err := store.GetApiById(apiId); err != nil {
		return mapStoreError(err)
	}

	if err := store.DeleteApiById(apiId); err != nil {
		return mapStoreError(err)
	}

	s.eventServer.DisconnectByApiId(apiId)

	return map[string]interface{}{}, nil
}

// ListApis returns a paginated list of Event APIs (v2).
// GET /v2/apis
func (s *AppSyncService) ListApis(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	opts := parsePaginationOptions(req)
	apis, nextToken, err := store.ListApis(opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(apis))
	for _, api := range apis {
		item := apiToMap(api)
		if tags, err := store.TagStore.List(api.Arn); err == nil && len(tags) > 0 {
			item["tags"] = tags
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

// mapStoreError converts a store-level error to the corresponding AppSync service error.
// Uses errors.Is for unwrapping so that wrapped sentinel errors are correctly matched.
func mapStoreError(err error) (interface{}, error) {
	switch {
	case errors.Is(err, appsyncstore.ErrApiNotFound):
		return nil, NewNotFoundException("API")
	case errors.Is(err, appsyncstore.ErrApiAlreadyExists):
		return nil, NewConflictException("API already exists")
	case errors.Is(err, appsyncstore.ErrChannelNamespaceNotFound):
		return nil, NewNotFoundException("Channel namespace")
	case errors.Is(err, appsyncstore.ErrChannelNamespaceExists):
		return nil, NewConflictException("Channel namespace already exists")
	case errors.Is(err, appsyncstore.ErrGraphqlApiNotFound):
		return nil, NewNotFoundException("GraphQL API")
	case errors.Is(err, appsyncstore.ErrGraphqlApiAlreadyExists):
		return nil, NewConflictException("GraphQL API already exists")
	case errors.Is(err, appsyncstore.ErrDataSourceNotFound):
		return nil, NewNotFoundException("Data source")
	case errors.Is(err, appsyncstore.ErrDataSourceAlreadyExists):
		return nil, NewConflictException("Data source already exists")
	case errors.Is(err, appsyncstore.ErrResolverNotFound):
		return nil, NewNotFoundException("Resolver")
	case errors.Is(err, appsyncstore.ErrResolverAlreadyExists):
		return nil, NewConflictException("Resolver already exists")
	case errors.Is(err, appsyncstore.ErrFunctionNotFound):
		return nil, NewNotFoundException("Function")
	case errors.Is(err, appsyncstore.ErrFunctionAlreadyExists):
		return nil, NewConflictException("Function already exists")
	case errors.Is(err, appsyncstore.ErrTypeNotFound):
		return nil, NewNotFoundException("Type")
	case errors.Is(err, appsyncstore.ErrTypeAlreadyExists):
		return nil, NewConflictException("Type already exists")
	case errors.Is(err, appsyncstore.ErrApiKeyNotFound):
		return nil, NewNotFoundException("API key")
	case errors.Is(err, appsyncstore.ErrApiCacheNotFound):
		return nil, NewNotFoundException("API cache")
	case errors.Is(err, appsyncstore.ErrApiCacheAlreadyExists):
		return nil, NewConflictException("API cache already exists")
	case errors.Is(err, appsyncstore.ErrDomainNameNotFound):
		return nil, NewNotFoundException("Domain name")
	case errors.Is(err, appsyncstore.ErrDomainNameAlreadyExists):
		return nil, NewConflictException("Domain name already exists")
	case errors.Is(err, appsyncstore.ErrApiAssociationNotFound):
		return nil, NewNotFoundException("API association")
	case errors.Is(err, appsyncstore.ErrMergedApiAssociationNotFound):
		return nil, NewNotFoundException("Merged API association")
	default:
		return nil, ErrInternalFailureException
	}
}
