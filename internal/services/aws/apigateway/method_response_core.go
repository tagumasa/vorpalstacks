package apigateway

import (
	"strings"

	"vorpalstacks/internal/store/aws/apigateway"
)

// MethodResponseInput carries the parsed wire members of a PutMethodResponse
// request.
type MethodResponseInput struct {
	StatusCode         string
	ResponseParameters map[string]bool
	ResponseModels     map[string]string
}

// putMethodResponseCore creates or replaces a method response.
func (s *APIGatewayService) putMethodResponseCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod, statusCode string,
	in *MethodResponseInput,
) (*apigateway.MethodResponse, error) {
	if apiId == "" || resourceId == "" || httpMethod == "" || statusCode == "" {
		return nil, NewBadRequestException("missing required parameters")
	}

	response := &apigateway.MethodResponse{
		StatusCode:         statusCode,
		ResponseParameters: in.ResponseParameters,
		ResponseModels:     in.ResponseModels,
	}

	result, err := stores.restApis.PutMethodResponse(apiId, resourceId, httpMethod, statusCode, response)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return result, nil
}

// getMethodResponseCore retrieves a method response.
func (s *APIGatewayService) getMethodResponseCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod, statusCode string,
) (*apigateway.MethodResponse, error) {
	if apiId == "" || resourceId == "" || httpMethod == "" || statusCode == "" {
		return nil, NewBadRequestException("missing required parameters")
	}

	result, err := stores.restApis.GetMethodResponse(apiId, resourceId, httpMethod, statusCode)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return result, nil
}

// deleteMethodResponseCore removes a method response.
func (s *APIGatewayService) deleteMethodResponseCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod, statusCode string,
) error {
	if apiId == "" || resourceId == "" || httpMethod == "" || statusCode == "" {
		return NewBadRequestException("missing required parameters")
	}

	if err := stores.restApis.DeleteMethodResponse(apiId, resourceId, httpMethod, statusCode); err != nil {
		return toApiGatewayError(err)
	}

	return nil
}

// updateMethodResponseCore applies JSON Patch operations to a method
// response.
func (s *APIGatewayService) updateMethodResponseCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod, statusCode string,
	ops []PatchOperation,
) (*apigateway.MethodResponse, error) {
	if apiId == "" || resourceId == "" || httpMethod == "" || statusCode == "" {
		return nil, NewBadRequestException("missing required parameters")
	}

	methodResp, err := stores.restApis.GetMethodResponse(apiId, resourceId, httpMethod, statusCode)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range ops {
		handled := false
		switch {
		case po.Path == "/responseParameters":
			handled = true
			// Whole-member row of the official UpdateMethodResponse patch
			// table: add, replace and remove are all supported.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if err := applyWholeBoolMapPatch(&methodResp.ResponseParameters, po); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/responseParameters/"):
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if methodResp.ResponseParameters == nil {
				methodResp.ResponseParameters = make(map[string]bool)
			}
			if err := applyBoolMapPatch(methodResp.ResponseParameters, po, "/responseParameters/", nil, nil); err != nil {
				return nil, err
			}
		case po.Path == "/responseModels":
			handled = true
			// Whole-member row: add, replace and remove are all supported.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if err := applyWholeStringMapPatch(&methodResp.ResponseModels, po, nil, nil); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/responseModels/"):
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if methodResp.ResponseModels == nil {
				methodResp.ResponseModels = make(map[string]string)
			}
			if err := applyMapPatch(methodResp.ResponseModels, po, "/responseModels/", nil, nil); err != nil {
				return nil, err
			}
		}
		if !handled {
			return nil, unknownPatchPathError(po)
		}
	}

	updatedResp, err := stores.restApis.PutMethodResponse(apiId, resourceId, httpMethod, statusCode, methodResp)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return updatedResp, nil
}
