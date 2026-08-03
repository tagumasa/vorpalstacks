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
	return s.toResourceResponse(created), nil
}

// GetResource retrieves a resource from API Gateway.
func (s *APIGatewayService) GetResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	resource, err := s.getResourceCore(stores, apiId, resourceId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toResourceResponse(resource), nil
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
		items = append(items, s.toResourceResponse(r))
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

func (s *APIGatewayService) toResourceResponse(r *store.Resource) map[string]interface{} {
	response := map[string]interface{}{
		"id":       r.Id,
		"parentId": r.ParentId,
		"path":     r.Path,
		"pathPart": r.PathPart,
	}

	if len(r.ResourceMethods) > 0 {
		methods := make(map[string]interface{})
		for method, m := range r.ResourceMethods {
			methods[method] = s.toMethodResponse(m)
		}
		response["resourceMethods"] = methods
	}

	return response
}
