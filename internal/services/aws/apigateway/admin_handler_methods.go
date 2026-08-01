package apigateway

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// PutMethod creates or replaces a method on a resource.
func (h *AdminHandler) PutMethod(ctx context.Context, req *connect.Request[pb.PutMethodRequest]) (*connect.Response[pb.Method], error) {
	if req.Msg.Restapiid == "" || req.Msg.Resourceid == "" || req.Msg.Httpmethod == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id, resource_id, and http_method are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	httpMethod := req.Msg.Httpmethod
	if !validateHTTPMethod(httpMethod) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid HTTP method: %s", httpMethod))
	}

	authType := req.Msg.Authorizationtype
	if authType == "" {
		authType = "NONE"
	}
	if !validateAuthorizationType(authType) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid authorization type: %s", authType))
	}
	if authType == "COGNITO_USER_POOLS" && req.Msg.Authorizerid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("authorizerId is required when authorizationType is COGNITO_USER_POOLS"))
	}

	method := &apigatewaystore.Method{
		HttpMethod:         req.Msg.Httpmethod,
		AuthorizationType:  authType,
		AuthorizerId:       req.Msg.Authorizerid,
		ApiKeyRequired:     req.Msg.GetApikeyrequired(),
		RequestValidatorId: req.Msg.Requestvalidatorid,
		OperationName:      req.Msg.Operationname,
		RequestParameters:  req.Msg.Requestparameters,
		RequestModels:      req.Msg.Requestmodels,
	}

	created, err := stores.restApis.PutMethod(req.Msg.Restapiid, req.Msg.Resourceid, method)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbMethod(created)), nil
}

// GetMethod returns a method configuration.
func (h *AdminHandler) GetMethod(ctx context.Context, req *connect.Request[pb.GetMethodRequest]) (*connect.Response[pb.Method], error) {
	if req.Msg.Restapiid == "" || req.Msg.Resourceid == "" || req.Msg.Httpmethod == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id, resource_id, and http_method are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	m, err := stores.restApis.GetMethod(req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbMethod(m)), nil
}

// DeleteMethod removes a method from a resource.
func (h *AdminHandler) DeleteMethod(ctx context.Context, req *connect.Request[pb.DeleteMethodRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Restapiid == "" || req.Msg.Resourceid == "" || req.Msg.Httpmethod == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id, resource_id, and http_method are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := stores.restApis.DeleteMethod(req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// PutIntegration creates or replaces an integration on a method.
func (h *AdminHandler) PutIntegration(ctx context.Context, req *connect.Request[pb.PutIntegrationRequest]) (*connect.Response[pb.Integration], error) {
	if req.Msg.Restapiid == "" || req.Msg.Resourceid == "" || req.Msg.Httpmethod == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id, resource_id, and http_method are required"))
	}
	integrationType := fromPbIntegrationType(req.Msg.Type)
	if integrationType == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("type is required and must be a valid integration type"))
	}

	// URI is required for all types except MOCK, matching the HTTP path.
	if integrationType != "MOCK" && req.Msg.Uri == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("uri is required for %s integration", integrationType))
	}

	// integrationHttpMethod is required for AWS (non-proxy) integrations.
	if integrationType == "AWS" && req.Msg.Integrationhttpmethod == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("integrationHttpMethod is required for AWS integration"))
	}

	// Validate passthroughBehaviour if provided.
	passthroughBehavior := req.Msg.Passthroughbehavior
	if passthroughBehavior != "" && !validatePassthroughBehavior(passthroughBehavior) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid passthroughBehavior: %s", passthroughBehavior))
	}

	// Validate and default timeoutInMillis, matching the HTTP path.
	timeoutInMillis := req.Msg.GetTimeoutinmillis()
	if timeoutInMillis > 0 && (timeoutInMillis < 50 || timeoutInMillis > 30000) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("timeoutInMillis must be between 50 and 30000"))
	}
	if timeoutInMillis <= 0 {
		timeoutInMillis = 29000
	}

	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	integration := &apigatewaystore.Integration{
		Type:                  integrationType,
		IntegrationHttpMethod: req.Msg.Integrationhttpmethod,
		Uri:                   req.Msg.Uri,
		Credentials:           req.Msg.Credentials,
		PassthroughBehavior:   passthroughBehavior,
		CacheNamespace:        req.Msg.Cachenamespace,
		ConnectionType:        fromPbConnectionType(req.Msg.Connectiontype),
		ConnectionId:          req.Msg.Connectionid,
		TimeoutInMillis:       timeoutInMillis,
		RequestParameters:     req.Msg.Requestparameters,
		RequestTemplates:      req.Msg.Requesttemplates,
		CacheKeyParameters:    req.Msg.Cachekeyparameters,
	}
	ch := fromPbContentHandling(req.Msg.Contenthandling)
	if ch != "" {
		integration.ContentHandling = ch
	}

	created, err := stores.restApis.PutIntegration(req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod, integration)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbIntegration(created)), nil
}

// GetIntegration returns a single integration.
func (h *AdminHandler) GetIntegration(ctx context.Context, req *connect.Request[pb.GetIntegrationRequest]) (*connect.Response[pb.Integration], error) {
	if req.Msg.Restapiid == "" || req.Msg.Resourceid == "" || req.Msg.Httpmethod == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id, resource_id, and http_method are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	i, err := stores.restApis.GetIntegration(req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbIntegration(i)), nil
}

// DeleteIntegration removes an integration from a method.
func (h *AdminHandler) DeleteIntegration(ctx context.Context, req *connect.Request[pb.DeleteIntegrationRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Restapiid == "" || req.Msg.Resourceid == "" || req.Msg.Httpmethod == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id, resource_id, and http_method are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := stores.restApis.DeleteIntegration(req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}
