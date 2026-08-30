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

// updateResourceCore applies JSON Patch operations to a resource under the
// api key lock, recomputing the full path when pathPart changes.
func (s *APIGatewayService) updateResourceCore(
	stores *apiGatewayStores,
	apiId, resourceId string,
	ops []PatchOperation,
) (*apigateway.Resource, error) {
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	resource, err := stores.restApis.GetResource(apiId, resourceId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range ops {
		switch po.Path {
		case "/pathPart":
			// The root resource's path part is the API's root path ("/") and
			// is immutable in AWS; DeleteResource already rejects deleting it.
			if resource.ParentId == "" {
				return nil, NewBadRequestException("cannot modify the pathPart of the root resource")
			}
			if po.Value == "" {
				return nil, NewBadRequestException("pathPart cannot be empty")
			}
			if !validatePathPart(po.Value) {
				return nil, NewBadRequestException("Invalid pathPart: malformed path parameter")
			}
			// Check for path collision with siblings under the same parent.
			siblings, err := stores.restApis.ListResources(apiId)
			if err != nil {
				return nil, err
			}
			for _, sib := range siblings {
				if sib.Id != resourceId && sib.ParentId == resource.ParentId && sib.PathPart == po.Value {
					return nil, NewConflictException("pathPart already exists under this parent")
				}
			}
			resource.PathPart = po.Value
			// A non-root resource always has a parent; a failed parent
			// lookup means the resource tree is inconsistent, and
			// persisting the new path part with the stale full path would
			// corrupt it further.
			parent, err := stores.restApis.GetResource(apiId, resource.ParentId)
			if err != nil {
				return nil, toApiGatewayError(err)
			}
			parentPath := strings.TrimRight(parent.Path, "/")
			resource.Path = parentPath + "/" + po.Value
		}
	}

	if err := stores.restApis.UpdateResourceCascade(apiId, resource); err != nil {
		return nil, toApiGatewayError(err)
	}

	return resource, nil
}

// listResourcesCore returns all resources for an API; pagination is applied
// by the caller.
func (s *APIGatewayService) listResourcesCore(stores *apiGatewayStores, apiId string) ([]*apigateway.Resource, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	resources, err := stores.restApis.ListResources(apiId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return resources, nil
}
