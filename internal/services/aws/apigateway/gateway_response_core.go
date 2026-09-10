package apigateway

import (
	"fmt"
	"strings"

	"vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/store/aws/common"
)

// GatewayResponseInput carries the parsed wire members of a
// PutGatewayResponse request.
type GatewayResponseInput struct {
	StatusCode         string
	ResponseParameters map[string]string
	ResponseTemplates  map[string]string
}

// gatewayResponseTypes is the GatewayResponseType enum of the model; a
// gateway response only exists for a modelled response type.
var gatewayResponseTypes = map[string]bool{
	"ACCESS_DENIED": true, "API_CONFIGURATION_ERROR": true,
	"AUTHORIZER_CONFIGURATION_ERROR": true, "AUTHORIZER_FAILURE": true,
	"BAD_REQUEST_BODY": true, "BAD_REQUEST_PARAMETERS": true,
	"DEFAULT_4XX": true, "DEFAULT_5XX": true, "EXPIRED_TOKEN": true,
	"INTEGRATION_FAILURE": true, "INTEGRATION_TIMEOUT": true,
	"INVALID_API_KEY": true, "INVALID_SIGNATURE": true,
	"MISSING_AUTHENTICATION_TOKEN": true, "QUOTA_EXCEEDED": true,
	"REQUEST_TOO_LARGE": true, "RESOURCE_NOT_FOUND": true,
	"THROTTLED": true, "UNAUTHORIZED": true,
	"UNSUPPORTED_MEDIA_TYPE": true, "WAF_FILTERED": true,
}

// validateGatewayResponseType rejects an unmodelled response type.
func validateGatewayResponseType(responseType string) *ApiGatewayError {
	if !gatewayResponseTypes[responseType] {
		return NewBadRequestException(fmt.Sprintf(
			"Invalid responseType: %s; must be a modelled GatewayResponseType", responseType))
	}
	return nil
}

// putGatewayResponseCore creates or replaces the gateway response of a
// response type. A response customised through the API is never the
// platform-generated default one.
func (s *APIGatewayService) putGatewayResponseCore(
	stores *apiGatewayStores,
	apiId, responseType string,
	in *GatewayResponseInput,
) (*apigateway.GatewayResponse, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if err := validateGatewayResponseType(responseType); err != nil {
		return nil, err
	}
	if in.StatusCode != "" {
		if err := validateStatusCode(in.StatusCode); err != nil {
			return nil, err
		}
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	resp := &apigateway.GatewayResponse{
		ResponseType:       responseType,
		StatusCode:         in.StatusCode,
		ResponseParameters: in.ResponseParameters,
		ResponseTemplates:  in.ResponseTemplates,
		DefaultResponse:    false,
	}
	created, err := stores.restApis.PutGatewayResponse(apiId, resp)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getGatewayResponseCore retrieves the gateway response of a response type.
func (s *APIGatewayService) getGatewayResponseCore(
	stores *apiGatewayStores, apiId, responseType string,
) (*apigateway.GatewayResponse, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if err := validateGatewayResponseType(responseType); err != nil {
		return nil, err
	}
	resp, err := stores.restApis.GetGatewayResponse(apiId, responseType)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return resp, nil
}

// deleteGatewayResponseCore removes the gateway response of a response type.
func (s *APIGatewayService) deleteGatewayResponseCore(
	stores *apiGatewayStores, apiId, responseType string,
) error {
	if apiId == "" {
		return NewBadRequestException("restApiId is required")
	}
	if err := validateGatewayResponseType(responseType); err != nil {
		return err
	}
	if err := stores.restApis.DeleteGatewayResponse(apiId, responseType); err != nil {
		return toApiGatewayError(err)
	}
	return nil
}

// listGatewayResponsesCore returns a page of the API's gateway responses.
func (s *APIGatewayService) listGatewayResponsesCore(
	stores *apiGatewayStores, limit int, marker, apiId string,
) (*common.ListResult[apigateway.GatewayResponse], error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	resolved, err := resolvePageLimit(limit)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	all, err := stores.restApis.ListGatewayResponses(apiId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return paginateList(all, marker, resolved, func(r *apigateway.GatewayResponse) string {
		return r.ResponseType
	})
}

// updateGatewayResponseCore applies the documented patch surface of a
// gateway response: replace of /statusCode, and add/replace/remove on
// /responseParameters and /responseTemplates.
func (s *APIGatewayService) updateGatewayResponseCore(
	stores *apiGatewayStores,
	apiId, responseType string,
	patches []PatchOperation,
) (*apigateway.GatewayResponse, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}
	if err := validateGatewayResponseType(responseType); err != nil {
		return nil, err
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	resp, err := stores.restApis.GetGatewayResponse(apiId, responseType)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	if resp.ResponseParameters == nil {
		resp.ResponseParameters = map[string]string{}
	}
	if resp.ResponseTemplates == nil {
		resp.ResponseTemplates = map[string]string{}
	}

	for _, po := range patches {
		handled := false
		switch {
		case po.Path == "/statusCode":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if err := validateStatusCode(po.Value); err != nil {
				return nil, err
			}
			resp.StatusCode = po.Value
		case po.Path == "/responseParameters":
			handled = true
			// Whole-member row of the official UpdateGatewayResponse
			// patch table: add, replace and remove are all supported.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if err := applyWholeStringMapPatch(&resp.ResponseParameters, po, nil, nil); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/responseParameters/"):
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if resp.ResponseParameters == nil {
				resp.ResponseParameters = make(map[string]string)
			}
			if err := applyMapPatch(resp.ResponseParameters, po, "/responseParameters/", nil, nil); err != nil {
				return nil, err
			}
		case po.Path == "/responseTemplates":
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if err := applyWholeStringMapPatch(&resp.ResponseTemplates, po, nil, nil); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/responseTemplates/"):
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if resp.ResponseTemplates == nil {
				resp.ResponseTemplates = make(map[string]string)
			}
			if err := applyMapPatch(resp.ResponseTemplates, po, "/responseTemplates/", nil, nil); err != nil {
				return nil, err
			}
		}
		if !handled {
			return nil, NewBadRequestException("unsupported patch path " + po.Path)
		}
	}

	updated, err := stores.restApis.PutGatewayResponse(apiId, resp)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return updated, nil
}
