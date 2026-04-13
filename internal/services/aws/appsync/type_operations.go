package appsync

import (
	"context"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// CreateType creates a new type definition for a GraphQL API.
// POST /v1/apis/{apiId}/types
func (s *AppSyncService) CreateType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	definition := request.GetStringParam(req.Parameters, "definition")
	if definition == "" {
		return nil, NewBadRequestException("definition is required")
	}

	format := request.GetStringParam(req.Parameters, "format")
	if format == "" {
		return nil, NewBadRequestException("format is required")
	}

	t := &appsyncstore.Type{
		ApiId:       apiId,
		Definition:  definition,
		Format:      format,
		Description: request.GetStringParam(req.Parameters, "description"),
	}

	created, err := store.CreateType(t)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"type": typeToMap(created),
	}, nil
}

// GetType retrieves a type definition by API ID and type name.
// GET /v1/apis/{apiId}/types/{typeName}
func (s *AppSyncService) GetType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	typeName := request.GetStringParam(req.Parameters, "typeName")
	if apiId == "" || typeName == "" {
		return nil, NewBadRequestException("apiId and typeName are required")
	}

	t, err := store.GetType(apiId, typeName)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"type": typeToMap(t),
	}, nil
}

// UpdateType updates an existing type definition.
// POST /v1/apis/{apiId}/types/{typeName}
func (s *AppSyncService) UpdateType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	typeName := request.GetStringParam(req.Parameters, "typeName")
	if apiId == "" || typeName == "" {
		return nil, NewBadRequestException("apiId and typeName are required")
	}

	format := request.GetStringParam(req.Parameters, "format")
	if format == "" {
		return nil, NewBadRequestException("format is required")
	}

	t := &appsyncstore.Type{
		ApiId:       apiId,
		Name:        typeName,
		Definition:  request.GetStringParam(req.Parameters, "definition"),
		Format:      format,
		Description: request.GetStringParam(req.Parameters, "description"),
	}

	updated, err := store.UpdateType(t)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"type": typeToMap(updated),
	}, nil
}

// DeleteType removes a type definition.
// DELETE /v1/apis/{apiId}/types/{typeName}
func (s *AppSyncService) DeleteType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	typeName := request.GetStringParam(req.Parameters, "typeName")
	if apiId == "" || typeName == "" {
		return nil, NewBadRequestException("apiId and typeName are required")
	}

	if err := store.DeleteType(apiId, typeName); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// ListTypes returns a paginated list of type definitions for a GraphQL API.
// GET /v1/apis/{apiId}/types
func (s *AppSyncService) ListTypes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	opts := parsePaginationOptions(req)
	types, nextToken, err := store.ListTypes(apiId, opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(types))
	for _, t := range types {
		items = append(items, typeToMap(t))
	}

	response := map[string]interface{}{
		"types": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

// typeToMap converts a Type struct to a response map with correct wire format.
func typeToMap(t *appsyncstore.Type) map[string]interface{} {
	m := map[string]interface{}{
		"name":   t.Name,
		"format": t.Format,
	}

	if t.Arn != "" {
		m["arn"] = t.Arn
	}
	if t.Definition != "" {
		m["definition"] = t.Definition
	}
	if t.Description != "" {
		m["description"] = t.Description
	}

	return m
}
