// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	store "vorpalstacks/internal/store/aws/apigateway"
)

func getApiIdAndResourceId(req *request.ParsedRequest) (string, string) {
	apiId := getRestApiId(req)
	resourceId := getResourceId(req)
	return apiId, resourceId
}

func getResourceId(req *request.ParsedRequest) string {
	resourceId := request.GetStringParam(req.Parameters, "resourceId")
	if resourceId == "" {
		resourceId = getPathParam(req, "resourceId")
	}
	return resourceId
}

// parseEmbedParam reads the embed query parameter: the model supports a
// single-valued list whose only allowed entry is "methods".
func parseEmbedParam(params map[string]interface{}) (bool, error) {
	raw, ok := params["embed"]
	if !ok || raw == nil {
		return false, nil
	}
	var values []string
	switch v := raw.(type) {
	case string:
		values = []string{v}
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok {
				values = append(values, s)
			}
		}
	}
	includeMethods := false
	for _, value := range values {
		if value == "methods" {
			includeMethods = true
			continue
		}
		return false, NewBadRequestException("the embed parameter supports only methods")
	}
	return includeMethods, nil
}

// CreateResource creates a new resource in API Gateway.
func (s *APIGatewayService) CreateResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	parentId := request.GetStringParam(req.Parameters, "parentId")
	if parentId == "" {
		parentId = request.GetStringParam(req.Parameters, "resourceId")
	}
	pathPart := request.GetStringParam(req.Parameters, "pathPart")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.createResourceCore(stores, apiId, parentId, pathPart)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toResourceResponse(created, false), nil
}

// GetResource retrieves a resource from API Gateway.
func (s *APIGatewayService) GetResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	includeMethods, err := parseEmbedParam(req.Parameters)
	if err != nil {
		return nil, err
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	resource, err := s.getResourceCore(stores, apiId, resourceId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toResourceResponse(resource, includeMethods), nil
}

// DeleteResource deletes a resource from API Gateway.
func (s *APIGatewayService) DeleteResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteResourceCore(stores, apiId, resourceId); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

// GetResources lists all resources for a REST API in API Gateway.
func (s *APIGatewayService) GetResources(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	position := request.GetStringParam(req.Parameters, "position")
	includeMethods, err := parseEmbedParam(req.Parameters)
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	resources, err := s.listResourcesCore(stores, apiId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	items := make([]interface{}, 0, len(resources))
	for _, r := range resources {
		items = append(items, s.toResourceResponse(r, includeMethods))
	}

	page, nextPos, found := paginateItems(items, position, limit)
	if !found {
		return nil, NewBadRequestException("Invalid position: " + position)
	}
	result := map[string]interface{}{
		"item": page,
	}
	if nextPos != "" {
		result["position"] = nextPos
	}
	return result, nil
}

// toResourceResponse renders a Resource; method summaries are only part of
// the response when the caller asked to embed them.
func (s *APIGatewayService) toResourceResponse(r *store.Resource, includeMethods bool) map[string]interface{} {
	response := map[string]interface{}{
		"id":       r.Id,
		"parentId": r.ParentId,
		"path":     r.Path,
		"pathPart": r.PathPart,
	}

	if includeMethods && len(r.ResourceMethods) > 0 {
		methods := make(map[string]interface{})
		for method, m := range r.ResourceMethods {
			methods[method] = s.toMethodResponse(m)
		}
		response["resourceMethods"] = methods
	}

	return response
}
