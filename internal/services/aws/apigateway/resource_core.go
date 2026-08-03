package apigateway

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"vorpalstacks/internal/store/aws/apigateway"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// createResourceCore creates a child resource under parentId. The store
// performs the duplicate-path check via ErrResourceAlreadyExists; both
// transport handlers map that to a ConflictException.
func (s *APIGatewayService) createResourceCore(
	stores *apiGatewayStores,
	apiId string,
	parentId string,
	pathPart string,
) (*apigateway.Resource, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if parentId == "" {
		return nil, NewBadRequestException("parentId is required")
	}
	if pathPart == "" {
		return nil, NewBadRequestException("pathPart is required")
	}
	if strings.Contains(pathPart, "/") {
		return nil, NewBadRequestException("pathPart must not contain '/'")
	}
	if !validatePathPart(pathPart) {
		return nil, NewBadRequestException("pathPart has unbalanced braces or invalid path parameter syntax")
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

	resource := &apigateway.Resource{
		ParentId:        parentId,
		Path:            path,
		PathPart:        pathPart,
		ResourceMethods: make(map[string]*apigateway.Method),
	}

	created, err := stores.restApis.CreateResource(apiId, resource)
	if err != nil {
		if errors.Is(err, apigateway.ErrResourceAlreadyExists) {
			return nil, NewConflictException("Resource already exists")
		}
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getResourceCore returns a single resource by id.
func (s *APIGatewayService) getResourceCore(stores *apiGatewayStores, apiId, resourceId string) (*apigateway.Resource, error) {
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}
	resource, err := stores.restApis.GetResource(apiId, resourceId)
	if err != nil {
		return nil, ErrNotFoundException
	}
	return resource, nil
}

// deleteResourceCore removes a resource. Maps store-level errors to the
// API Gateway exception types.
func (s *APIGatewayService) deleteResourceCore(stores *apiGatewayStores, apiId, resourceId string) error {
	if apiId == "" || resourceId == "" {
		return NewBadRequestException("restApiId and resourceId are required")
	}
	err := stores.restApis.DeleteResource(apiId, resourceId)
	if err == nil {
		return nil
	}
	var storeErr *storecommon.StoreError
	if errors.As(err, &storeErr) {
		msg := storeErr.Err.Error()
		if strings.Contains(msg, "cannot delete resource with child resources") {
			return NewBadRequestException("Resource has child resources")
		}
		if strings.Contains(msg, "cannot delete the root resource") {
			return NewBadRequestException("Cannot delete the root resource")
		}
	}
	if storecommon.IsNotFound(err) {
		return ErrNotFoundException
	}
	return NewApiGatewayError("InternalServerError", fmt.Sprintf("Failed to delete resource: %v", err), http.StatusInternalServerError)
}

// listResourcesCore returns all resources for an API; pagination is applied
// by the caller.
func (s *APIGatewayService) listResourcesCore(stores *apiGatewayStores, apiId string) ([]*apigateway.Resource, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	return stores.restApis.ListResources(apiId)
}
