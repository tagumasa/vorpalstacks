// Package apigateway provides API Gateway service operations for vorpalstacks.
package apigateway

import (
	"context"
	"errors"
	"strings"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	store "vorpalstacks/internal/store/aws/apigateway"
	common "vorpalstacks/internal/store/aws/common"
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
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	parentId := request.GetStringParam(req.Parameters, "parentId")
	if parentId == "" {
		parentId = request.GetStringParam(req.Parameters, "resourceId")
	}
	pathPart := request.GetStringParam(req.Parameters, "pathPart")

	if parentId == "" {
		return nil, NewBadRequestException("parentId is required")
	}
	if pathPart == "" {
		return nil, NewBadRequestException("pathPart is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	parentResource, err := stores.restApis.GetResource(apiId, parentId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	path := parentResource.Path
	if path != "/" {
		path += "/"
	}
	path += pathPart

	resource := &store.Resource{
		ParentId:        parentId,
		Path:            path,
		PathPart:        pathPart,
		ResourceMethods: make(map[string]*store.Method),
	}

	created, err := stores.restApis.CreateResource(apiId, resource)
	if err != nil {
		return nil, err
	}

	return s.toResourceResponse(created), nil
}

// GetResource retrieves a resource from API Gateway.
func (s *APIGatewayService) GetResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	resource, err := stores.restApis.GetResource(apiId, resourceId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	return s.toResourceResponse(resource), nil
}

// DeleteResource deletes a resource from API Gateway.
func (s *APIGatewayService) DeleteResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.restApis.DeleteResource(apiId, resourceId); err != nil {
		var storeErr *common.StoreError
		if errors.As(err, &storeErr) && strings.Contains(storeErr.Err.Error(), "cannot delete resource with child resources") {
			return nil, NewBadRequestException("Resource has child resources")
		}
		return nil, ErrNotFoundException
	}

	return response.EmptyResponse(), nil
}

// GetResources lists all resources for a REST API in API Gateway.
func (s *APIGatewayService) GetResources(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	resources, err := stores.restApis.ListResources(apiId)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(resources))
	for _, r := range resources {
		items = append(items, s.toResourceResponse(r))
	}

	return map[string]interface{}{
		"item": items,
	}, nil
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
