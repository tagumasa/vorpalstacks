package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateType creates a new type definition for a GraphQL API.
// POST /v1/apis/{apiId}/types
func (s *AppSyncService) CreateType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	t, err := s.createTypeCore(store, createTypeInput{
		ApiId:       request.GetStringParam(req.Parameters, "apiId"),
		Definition:  request.GetStringParam(req.Parameters, "definition"),
		Format:      request.GetStringParam(req.Parameters, "format"),
		Description: request.GetStringParam(req.Parameters, "description"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type": typeToMap(t),
	}, nil
}

// GetType retrieves a type definition by API ID and type name.
// GET /v1/apis/{apiId}/types/{typeName}
func (s *AppSyncService) GetType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	t, err := s.getTypeCore(store, request.GetStringParam(req.Parameters, "apiId"), request.GetStringParam(req.Parameters, "typeName"))
	if err != nil {
		return nil, err
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
		return mapStoreError(err)
	}

	t, err := s.updateTypeCore(store, updateTypeInput{
		ApiId:          request.GetStringParam(req.Parameters, "apiId"),
		TypeName:       request.GetStringParam(req.Parameters, "typeName"),
		Format:         request.GetStringParam(req.Parameters, "format"),
		Definition:     request.GetStringParam(req.Parameters, "definition"),
		HasDefinition:  request.HasParam(req.Parameters, "definition"),
		Description:    request.GetStringParam(req.Parameters, "description"),
		HasDescription: request.HasParam(req.Parameters, "description"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type": typeToMap(t),
	}, nil
}

// DeleteType removes a type definition.
// DELETE /v1/apis/{apiId}/types/{typeName}
func (s *AppSyncService) DeleteType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	if err := s.deleteTypeCore(store, request.GetStringParam(req.Parameters, "apiId"), request.GetStringParam(req.Parameters, "typeName")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListTypes returns a paginated list of type definitions for a GraphQL API.
// GET /v1/apis/{apiId}/types
func (s *AppSyncService) ListTypes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")

	types, nextToken, err := s.listTypesCore(store, apiId, request.GetStringParam(req.Parameters, "format"), request.GetIntParam(req.Parameters, "maxResults"), request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
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
