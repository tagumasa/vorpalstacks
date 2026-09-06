package apigateway

import (
	"strings"

	integration "vorpalstacks/internal/services/aws/apigateway/runtime/integration"
	"vorpalstacks/internal/store/aws/apigateway"
)

// MethodInput is the transport-agnostic input for creating or replacing a
// method. Both the data-plane HTTP handler and the admin gRPC handler
// build this struct from their respective request formats and pass it to
// putMethodCore.
type MethodInput struct {
	AuthorizationType   string
	AuthorizerId        string
	ApiKeyRequired      bool
	RequestValidatorId  string
	OperationName       string
	RequestParameters   map[string]bool
	RequestModels       map[string]string
	AuthorizationScopes []string
}

// IntegrationInput is the transport-agnostic input for creating or
// replacing an integration. Fields are populated from either an HTTP query
// parameter map or a protobuf message, then validated and persisted by
// putIntegrationCore.
type IntegrationInput struct {
	Type                  string
	IntegrationHttpMethod string
	Uri                   string
	Credentials           string
	PassthroughBehavior   string
	ContentHandling       string
	CacheNamespace        string
	CacheKeyParameters    []string
	TimeoutInMillis       int32
	ConnectionType        string
	ConnectionId          string
	RequestParameters     map[string]string
	RequestTemplates      map[string]string
	TlsConfig             *apigateway.TlsConfig
	ResponseTransferMode  string
	IntegrationTarget     string
}

// IntegrationResponseInput is the transport-agnostic input for an
// integration response.
type IntegrationResponseInput struct {
	SelectionPattern   string
	ContentHandling    string
	ResponseParameters map[string]string
	ResponseTemplates  map[string]string
}

// putMethodCore validates the input, builds a Method struct, and persists
// it via the store. The caller is responsible for resolving the HTTP
// method value before invocation.
func (s *APIGatewayService) putMethodCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod string,
	in *MethodInput,
) (*apigateway.Method, error) {
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}
	if !validateHTTPMethod(httpMethod) {
		return nil, NewBadRequestException("Invalid HTTP method: " + httpMethod)
	}
	authType := in.AuthorizationType
	if authType == "" {
		authType = "NONE"
	}
	if !validateAuthorizationType(authType) {
		return nil, NewBadRequestException("Invalid authorization type: " + authType)
	}
	if authType == "COGNITO_USER_POOLS" && in.AuthorizerId == "" {
		return nil, NewBadRequestException("authorizerId is required when authorizationType is COGNITO_USER_POOLS")
	}

	method := &apigateway.Method{
		HttpMethod:          httpMethod,
		AuthorizationType:   authType,
		AuthorizerId:        in.AuthorizerId,
		ApiKeyRequired:      in.ApiKeyRequired,
		RequestValidatorId:  in.RequestValidatorId,
		OperationName:       in.OperationName,
		RequestParameters:   in.RequestParameters,
		RequestModels:       in.RequestModels,
		AuthorizationScopes: in.AuthorizationScopes,
	}

	created, err := stores.restApis.PutMethod(apiId, resourceId, method)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getMethodCore retrieves a method.
func (s *APIGatewayService) getMethodCore(stores *apiGatewayStores, apiId, resourceId, httpMethod string) (*apigateway.Method, error) {
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}
	method, err := stores.restApis.GetMethod(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return method, nil
}

// deleteMethodCore removes a method.
func (s *APIGatewayService) deleteMethodCore(stores *apiGatewayStores, apiId, resourceId, httpMethod string) error {
	if apiId == "" || resourceId == "" {
		return NewBadRequestException("restApiId and resourceId are required")
	}
	if httpMethod == "" {
		return NewBadRequestException("httpMethod is required")
	}
	if err := stores.restApis.DeleteMethod(apiId, resourceId, httpMethod); err != nil {
		return toApiGatewayError(err)
	}
	return nil
}

// putIntegrationCore validates the input, builds an Integration struct, and
// persists it via the store. Validation centralises the type, URI,
// integrationHttpMethod, passthroughBehaviour, contentHandling,
// connectionType, responseTransferMode and timeout range checks that the
// admin handler previously skipped or duplicated.
func (s *APIGatewayService) putIntegrationCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod string,
	in *IntegrationInput,
) (*apigateway.Integration, error) {
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}
	if in.Type == "" {
		return nil, NewBadRequestException("type is required")
	}
	if !validateIntegrationType(in.Type) {
		return nil, NewBadRequestException("Invalid integration type: " + in.Type)
	}
	if in.Type != "MOCK" && in.Uri == "" {
		return nil, NewBadRequestException("uri is required for " + in.Type + " integration")
	}
	if in.Type == "AWS" && in.IntegrationHttpMethod == "" {
		return nil, NewBadRequestException("integrationHttpMethod is required for AWS integration")
	}
	if in.PassthroughBehavior != "" && !validatePassthroughBehavior(in.PassthroughBehavior) {
		return nil, NewBadRequestException("Invalid passthroughBehavior: " + in.PassthroughBehavior)
	}
	if in.ContentHandling != "" && !validateContentHandling(in.ContentHandling) {
		return nil, NewBadRequestException("Invalid contentHandling: " + in.ContentHandling)
	}
	if in.ConnectionType != "" && !validateConnectionType(in.ConnectionType) {
		return nil, NewBadRequestException("Invalid connectionType: " + in.ConnectionType)
	}
	// A VPC_LINK integration routes through a VpcLink, whose target is a
	// Network Load Balancer; no such substrate or consumer exists on this
	// platform, so the connection type rejects instead of storing inert
	// configuration.
	if in.ConnectionType == "VPC_LINK" {
		return nil, NewBadRequestException("connectionType VPC_LINK is not supported: VpcLink targets a Network Load Balancer, which this platform does not provide")
	}
	if in.ResponseTransferMode != "" && !validateResponseTransferMode(in.ResponseTransferMode) {
		return nil, NewBadRequestException("Invalid responseTransferMode: must be BUFFERED or STREAM")
	}

	timeout := in.TimeoutInMillis
	if timeout <= 0 {
		timeout = integration.DefaultIntegrationTimeoutMillis
	}
	if !validateTimeoutInMillis(timeout) {
		return nil, NewBadRequestException("timeoutInMillis must be between 50 and 30000")
	}

	integration := &apigateway.Integration{
		Type:                  in.Type,
		IntegrationHttpMethod: in.IntegrationHttpMethod,
		Uri:                   in.Uri,
		Credentials:           in.Credentials,
		PassthroughBehavior:   in.PassthroughBehavior,
		ContentHandling:       in.ContentHandling,
		CacheNamespace:        in.CacheNamespace,
		CacheKeyParameters:    in.CacheKeyParameters,
		TimeoutInMillis:       timeout,
		ConnectionType:        in.ConnectionType,
		ConnectionId:          in.ConnectionId,
		RequestParameters:     in.RequestParameters,
		RequestTemplates:      in.RequestTemplates,
		TlsConfig:             in.TlsConfig,
		ResponseTransferMode:  in.ResponseTransferMode,
		IntegrationTarget:     in.IntegrationTarget,
	}

	created, err := stores.restApis.PutIntegration(apiId, resourceId, httpMethod, integration)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getIntegrationCore retrieves an integration.
func (s *APIGatewayService) getIntegrationCore(stores *apiGatewayStores, apiId, resourceId, httpMethod string) (*apigateway.Integration, error) {
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}
	integration, err := stores.restApis.GetIntegration(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return integration, nil
}

// deleteIntegrationCore removes an integration.
func (s *APIGatewayService) deleteIntegrationCore(stores *apiGatewayStores, apiId, resourceId, httpMethod string) error {
	if apiId == "" || resourceId == "" {
		return NewBadRequestException("restApiId and resourceId are required")
	}
	if httpMethod == "" {
		return NewBadRequestException("httpMethod is required")
	}
	if err := stores.restApis.DeleteIntegration(apiId, resourceId, httpMethod); err != nil {
		return toApiGatewayError(err)
	}
	return nil
}

// putIntegrationResponseCore creates or replaces an integration response.
func (s *APIGatewayService) putIntegrationResponseCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod, statusCode string,
	in *IntegrationResponseInput,
) (*apigateway.IntegrationResponse, error) {
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}
	if statusCode == "" {
		return nil, NewBadRequestException("statusCode is required")
	}
	if in.ContentHandling != "" && !validateContentHandling(in.ContentHandling) {
		return nil, NewBadRequestException("Invalid contentHandling: " + in.ContentHandling)
	}

	response := &apigateway.IntegrationResponse{
		StatusCode:         statusCode,
		SelectionPattern:   in.SelectionPattern,
		ContentHandling:    in.ContentHandling,
		ResponseParameters: in.ResponseParameters,
		ResponseTemplates:  in.ResponseTemplates,
	}

	created, err := stores.restApis.PutIntegrationResponse(apiId, resourceId, httpMethod, statusCode, response)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return created, nil
}

// getIntegrationResponseCore retrieves an integration response.
func (s *APIGatewayService) getIntegrationResponseCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod, statusCode string,
) (*apigateway.IntegrationResponse, error) {
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}
	if statusCode == "" {
		return nil, NewBadRequestException("statusCode is required")
	}
	response, err := stores.restApis.GetIntegrationResponse(apiId, resourceId, httpMethod, statusCode)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return response, nil
}

// updateMethodCore applies JSON Patch operations to a method under the api
// key lock.
func (s *APIGatewayService) updateMethodCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod string,
	ops []PatchOperation,
) (*apigateway.Method, error) {
	if apiId == "" || resourceId == "" || httpMethod == "" {
		return nil, NewBadRequestException("restApiId, resourceId, and httpMethod are required")
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	method, err := stores.restApis.GetMethod(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range ops {
		handled := false
		switch {
		case po.Path == "/authorizationType":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if !validateAuthorizationType(po.Value) {
				return nil, NewBadRequestException("Invalid authorization type: " + po.Value)
			}
			method.AuthorizationType = po.Value
		case po.Path == "/authorizerId":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			method.AuthorizerId = po.Value
		case po.Path == "/apiKeyRequired":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			method.ApiKeyRequired = po.Value == "true"
		case po.Path == "/requestValidatorId":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			method.RequestValidatorId = po.Value
		case po.Path == "/operationName":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			method.OperationName = po.Value
		case po.Path == "/requestParameters":
			handled = true
			// Whole-member row of the official UpdateMethod patch table:
			// add, replace — except for MOCK integrations — and remove.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if po.Op == "replace" {
				if integration, err := stores.restApis.GetIntegration(apiId, resourceId, httpMethod); err == nil && integration.Type == "MOCK" {
					return nil, NewBadRequestException("requestParameters replace is not supported for MOCK integrations")
				}
			}
			if err := applyWholeBoolMapPatch(&method.RequestParameters, po); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/requestParameters/"):
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if method.RequestParameters == nil {
				method.RequestParameters = make(map[string]bool)
			}
			if err := applyBoolMapPatch(method.RequestParameters, po, "/requestParameters/", nil, nil); err != nil {
				return nil, err
			}
		case po.Path == "/requestModels":
			handled = true
			// Whole-member row: add, replace and remove are all supported.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if err := applyWholeStringMapPatch(&method.RequestModels, po, nil, nil); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/requestModels/"):
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if method.RequestModels == nil {
				method.RequestModels = make(map[string]string)
			}
			if err := applyMapPatch(method.RequestModels, po, "/requestModels/", nil, nil); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/authorizationScopes"):
			handled = true
			rest := strings.TrimPrefix(po.Path, "/authorizationScopes")
			switch {
			case rest == "":
				// Whole-member form of the official UpdateMethod patch
				// table: add appends the value, remove clears the list;
				// replace is not supported there.
				if err := requirePatchOp(po, opAdd|opRemove); err != nil {
					return nil, err
				}
				if po.Op == "remove" {
					method.AuthorizationScopes = nil
				} else {
					method.AuthorizationScopes = append(method.AuthorizationScopes, po.Value)
				}
			default:
				// Numeric index addressing appears nowhere in the official
				// patch tables — the documented list ops address the whole
				// member only.
				return nil, unknownPatchPathError(po)
			}
		}
		if !handled {
			return nil, unknownPatchPathError(po)
		}
	}

	_, err = stores.restApis.PutMethod(apiId, resourceId, method)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return method, nil
}

// updateIntegrationCore applies JSON Patch operations to an integration
// under the api key lock.
func (s *APIGatewayService) updateIntegrationCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod string,
	ops []PatchOperation,
) (*apigateway.Integration, error) {
	if apiId == "" || resourceId == "" || httpMethod == "" {
		return nil, NewBadRequestException("restApiId, resourceId, and httpMethod are required")
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	integrationRec, err := stores.restApis.GetIntegration(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range ops {
		handled := false
		switch {
		case po.Path == "/uri":
			handled = true
			// The row documents replace, except for MOCK integrations.
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if integrationRec.Type == "MOCK" {
				return nil, NewBadRequestException("uri replace is not supported for MOCK integrations")
			}
			integrationRec.Uri = po.Value
		case po.Path == "/type":
			handled = true
			// The row marks every operation Not supported: the
			// integration type is immutable through patching.
			return nil, unknownPatchPathError(po)
		case po.Path == "/httpMethod":
			handled = true
			// The row documents replace, except for MOCK integrations.
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if integrationRec.Type == "MOCK" {
				return nil, NewBadRequestException("httpMethod replace is not supported for MOCK integrations")
			}
			if po.Value != "" && !validateHTTPMethod(po.Value) {
				return nil, NewBadRequestException("Invalid integration HTTP method: " + po.Value)
			}
			integrationRec.IntegrationHttpMethod = po.Value
		case po.Path == "/integrationTarget":
			handled = true
			// The row documents replace only.
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			integrationRec.IntegrationTarget = po.Value
		case po.Path == "/responseTransferMode":
			handled = true
			// The row documents replace only.
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if po.Value != "" && !validateResponseTransferMode(po.Value) {
				return nil, NewBadRequestException("Invalid responseTransferMode: must be BUFFERED or STREAM")
			}
			integrationRec.ResponseTransferMode = po.Value
		case po.Path == "/credentials":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			integrationRec.Credentials = po.Value
		case po.Path == "/passthroughBehavior":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if !validatePassthroughBehavior(po.Value) {
				return nil, NewBadRequestException("Invalid passthroughBehavior: " + po.Value)
			}
			integrationRec.PassthroughBehavior = po.Value
		case po.Path == "/contentHandling":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if !validateContentHandling(po.Value) {
				return nil, NewBadRequestException("Invalid contentHandling: " + po.Value)
			}
			integrationRec.ContentHandling = po.Value
		case po.Path == "/cacheNamespace":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			integrationRec.CacheNamespace = po.Value
		case po.Path == "/connectionType":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if !validateConnectionType(po.Value) {
				return nil, NewBadRequestException("Invalid connectionType: " + po.Value)
			}
			// Same substrate bound as the create path: a VPC_LINK
			// integration would reference a VpcLink over a Network Load
			// Balancer this platform does not provide.
			if po.Value == "VPC_LINK" {
				return nil, NewBadRequestException("connectionType VPC_LINK is not supported: VpcLink targets a Network Load Balancer, which this platform does not provide")
			}
			integrationRec.ConnectionType = po.Value
		case po.Path == "/connectionId":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			integrationRec.ConnectionId = po.Value
		case po.Path == "/timeoutInMillis":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			v, err := parseInt32(po.Value)
			if err != nil {
				return nil, NewBadRequestException("invalid timeoutInMillis: not a number")
			}
			if v <= 0 {
				v = integration.DefaultIntegrationTimeoutMillis
			}
			if !validateTimeoutInMillis(v) {
				return nil, NewBadRequestException("timeoutInMillis must be between 50 and 30000")
			}
			integrationRec.TimeoutInMillis = v
		case po.Path == "/requestParameters":
			handled = true
			// Whole-member row: add, replace and remove are all supported.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if err := applyWholeStringMapPatch(&integrationRec.RequestParameters, po, nil, nil); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/requestParameters/"):
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if integrationRec.RequestParameters == nil {
				integrationRec.RequestParameters = make(map[string]string)
			}
			if err := applyMapPatch(integrationRec.RequestParameters, po, "/requestParameters/", nil, nil); err != nil {
				return nil, err
			}
		case po.Path == "/requestTemplates":
			handled = true
			// Whole-member row: add, replace and remove are all supported.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if err := applyWholeStringMapPatch(&integrationRec.RequestTemplates, po, nil, nil); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/requestTemplates/"):
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if integrationRec.RequestTemplates == nil {
				integrationRec.RequestTemplates = make(map[string]string)
			}
			if err := applyMapPatch(integrationRec.RequestTemplates, po, "/requestTemplates/", nil, nil); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/cacheKeyParameters"):
			handled = true
			rest := strings.TrimPrefix(po.Path, "/cacheKeyParameters")
			switch {
			case rest == "":
				// Whole-member row: add appends the value, remove clears
				// the list, replace sets it from a JSON array of strings.
				if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
					return nil, err
				}
				switch po.Op {
				case "remove":
					integrationRec.CacheKeyParameters = nil
				case "add":
					integrationRec.CacheKeyParameters = append(integrationRec.CacheKeyParameters, po.Value)
				default:
					parsed, err := parseWholeStringListValue(po)
					if err != nil {
						return nil, err
					}
					integrationRec.CacheKeyParameters = parsed
				}
			default:
				// Numeric index addressing appears nowhere in the official
				// patch tables — the documented list ops address the whole
				// member only.
				return nil, unknownPatchPathError(po)
			}
		case po.Path == "/tlsConfig/insecureSkipVerification":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if integrationRec.TlsConfig == nil {
				integrationRec.TlsConfig = &apigateway.TlsConfig{}
			}
			integrationRec.TlsConfig.InsecureSkipVerification = po.Value == "true"
		}
		if !handled {
			return nil, unknownPatchPathError(po)
		}
	}

	if err := stores.restApis.UpdateIntegration(apiId, resourceId, httpMethod, integrationRec); err != nil {
		return nil, toApiGatewayError(err)
	}

	return integrationRec, nil
}

// updateIntegrationResponseCore applies JSON Patch operations to an
// integration response under the api key lock.
func (s *APIGatewayService) updateIntegrationResponseCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod, statusCode string,
	ops []PatchOperation,
) (*apigateway.IntegrationResponse, error) {
	if apiId == "" || resourceId == "" || httpMethod == "" || statusCode == "" {
		return nil, NewBadRequestException("missing required parameters")
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	intResp, err := stores.restApis.GetIntegrationResponse(apiId, resourceId, httpMethod, statusCode)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	for _, po := range ops {
		handled := false
		switch {
		case po.Path == "/selectionPattern":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			intResp.SelectionPattern = po.Value
		case po.Path == "/contentHandling":
			handled = true
			if err := requirePatchOp(po, opReplace); err != nil {
				return nil, err
			}
			if !validateContentHandling(po.Value) {
				return nil, NewBadRequestException("Invalid contentHandling: " + po.Value)
			}
			intResp.ContentHandling = po.Value
		case po.Path == "/responseParameters":
			handled = true
			// Whole-member row: add, replace and remove are all supported.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if err := applyWholeStringMapPatch(&intResp.ResponseParameters, po, nil, nil); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/responseParameters/"):
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if intResp.ResponseParameters == nil {
				intResp.ResponseParameters = make(map[string]string)
			}
			if err := applyMapPatch(intResp.ResponseParameters, po, "/responseParameters/", nil, nil); err != nil {
				return nil, err
			}
		case po.Path == "/responseTemplates":
			handled = true
			// Whole-member row: add, replace and remove are all supported.
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if err := applyWholeStringMapPatch(&intResp.ResponseTemplates, po, nil, nil); err != nil {
				return nil, err
			}
		case strings.HasPrefix(po.Path, "/responseTemplates/"):
			handled = true
			if err := requirePatchOp(po, opAdd|opReplace|opRemove); err != nil {
				return nil, err
			}
			if intResp.ResponseTemplates == nil {
				intResp.ResponseTemplates = make(map[string]string)
			}
			if err := applyMapPatch(intResp.ResponseTemplates, po, "/responseTemplates/", nil, nil); err != nil {
				return nil, err
			}
		}
		if !handled {
			return nil, unknownPatchPathError(po)
		}
	}

	if err := stores.restApis.UpdateIntegrationResponse(apiId, resourceId, httpMethod, statusCode, intResp); err != nil {
		return nil, toApiGatewayError(err)
	}

	return intResp, nil
}

// deleteIntegrationResponseCore removes an integration response.
func (s *APIGatewayService) deleteIntegrationResponseCore(
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod, statusCode string,
) error {
	if apiId == "" || resourceId == "" {
		return NewBadRequestException("restApiId and resourceId are required")
	}
	if httpMethod == "" {
		return NewBadRequestException("httpMethod is required")
	}
	if statusCode == "" {
		return NewBadRequestException("statusCode is required")
	}
	if err := stores.restApis.DeleteIntegrationResponse(apiId, resourceId, httpMethod, statusCode); err != nil {
		return toApiGatewayError(err)
	}
	return nil
}
