package apigateway

import (
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
		return nil, ErrNotFoundException
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
		return ErrNotFoundException
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
	if in.ResponseTransferMode != "" && !validateResponseTransferMode(in.ResponseTransferMode) {
		return nil, NewBadRequestException("Invalid responseTransferMode: must be BUFFERED or STREAM")
	}

	timeout := in.TimeoutInMillis
	if timeout <= 0 {
		timeout = 29000
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
		return nil, ErrNotFoundException
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
		return ErrNotFoundException
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
		return nil, ErrNotFoundException
	}
	return response, nil
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
		return ErrNotFoundException
	}
	return nil
}
