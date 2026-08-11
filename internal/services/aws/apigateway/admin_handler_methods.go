package apigateway

import (
	"context"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
)

// PutMethod creates or replaces a method on a resource.
func (h *AdminHandler) PutMethod(ctx context.Context, req *connect.Request[pb.PutMethodRequest]) (*connect.Response[pb.Method], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	in := &MethodInput{
		AuthorizationType:  req.Msg.Authorizationtype,
		AuthorizerId:       req.Msg.Authorizerid,
		ApiKeyRequired:     req.Msg.GetApikeyrequired(),
		RequestValidatorId: req.Msg.Requestvalidatorid,
		OperationName:      req.Msg.Operationname,
		RequestParameters:  req.Msg.Requestparameters,
		RequestModels:      req.Msg.Requestmodels,
	}
	for _, scope := range req.Msg.Authorizationscopes {
		in.AuthorizationScopes = append(in.AuthorizationScopes, scope)
	}

	created, err := h.service.putMethodCore(stores, req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbMethod(created)), nil
}

// GetMethod returns a method configuration.
func (h *AdminHandler) GetMethod(ctx context.Context, req *connect.Request[pb.GetMethodRequest]) (*connect.Response[pb.Method], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	m, err := h.service.getMethodCore(stores, req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbMethod(m)), nil
}

// DeleteMethod removes a method from a resource.
func (h *AdminHandler) DeleteMethod(ctx context.Context, req *connect.Request[pb.DeleteMethodRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if err := h.service.deleteMethodCore(stores, req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// PutIntegration creates or replaces an integration on a method.
func (h *AdminHandler) PutIntegration(ctx context.Context, req *connect.Request[pb.PutIntegrationRequest]) (*connect.Response[pb.Integration], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	in := &IntegrationInput{
		Type:                  fromPbIntegrationType(req.Msg.Type),
		IntegrationHttpMethod: req.Msg.Integrationhttpmethod,
		Uri:                   req.Msg.Uri,
		Credentials:           req.Msg.Credentials,
		PassthroughBehavior:   req.Msg.Passthroughbehavior,
		ContentHandling:       fromPbContentHandling(req.Msg.Contenthandling),
		CacheNamespace:        req.Msg.Cachenamespace,
		CacheKeyParameters:    req.Msg.Cachekeyparameters,
		TimeoutInMillis:       req.Msg.GetTimeoutinmillis(),
		ConnectionType:        fromPbConnectionType(req.Msg.Connectiontype),
		ConnectionId:          req.Msg.Connectionid,
		RequestParameters:     req.Msg.Requestparameters,
		RequestTemplates:      req.Msg.Requesttemplates,
		ResponseTransferMode:  fromPbResponseTransferMode(req.Msg.GetResponsetransfermode()),
		IntegrationTarget:     req.Msg.Integrationtarget,
	}
	if req.Msg.Tlsconfig != nil {
		in.TlsConfig = fromPbTlsConfig(req.Msg.Tlsconfig)
	}

	created, err := h.service.putIntegrationCore(stores, req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbIntegration(created)), nil
}

// GetIntegration returns a single integration.
func (h *AdminHandler) GetIntegration(ctx context.Context, req *connect.Request[pb.GetIntegrationRequest]) (*connect.Response[pb.Integration], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	i, err := h.service.getIntegrationCore(stores, req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbIntegration(i)), nil
}

// DeleteIntegration removes an integration from a method.
func (h *AdminHandler) DeleteIntegration(ctx context.Context, req *connect.Request[pb.DeleteIntegrationRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if err := h.service.deleteIntegrationCore(stores, req.Msg.Restapiid, req.Msg.Resourceid, req.Msg.Httpmethod); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// fromPbResponseTransferMode converts a protobuf ResponseTransferMode enum
// to the canonical short string ("STREAM" / "BUFFERED") used in storage.
// Returns "" for the zero/unset value, which the core treats as "not
// specified".
func fromPbResponseTransferMode(t pb.ResponseTransferMode) string {
	switch t {
	case pb.ResponseTransferMode_RESPONSE_TRANSFER_MODE_BUFFERED:
		return "BUFFERED"
	case pb.ResponseTransferMode_RESPONSE_TRANSFER_MODE_STREAM:
		return "STREAM"
	default:
		return ""
	}
}
