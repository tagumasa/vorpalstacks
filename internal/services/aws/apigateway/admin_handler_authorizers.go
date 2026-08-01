package apigateway

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// CreateAuthorizer creates a new authorizer for a REST API.
func (h *AdminHandler) CreateAuthorizer(ctx context.Context, req *connect.Request[pb.CreateAuthorizerRequest]) (*connect.Response[pb.Authorizer], error) {
	if req.Msg.Restapiid == "" || req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id and name are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	authorizer := &apigatewaystore.Authorizer{
		Name:                         req.Msg.Name,
		Type:                         fromPbAuthorizerType(req.Msg.Type),
		AuthType:                     req.Msg.Authtype,
		AuthorizerUri:                req.Msg.Authorizeruri,
		AuthorizerCredentials:        req.Msg.Authorizercredentials,
		IdentitySource:               req.Msg.Identitysource,
		IdentityValidationExpression: req.Msg.Identityvalidationexpression,
		AuthorizerResultTtlInSeconds: req.Msg.Authorizerresultttlinseconds,
		ProviderArns:                 req.Msg.Providerarns,
	}
	if !validateAuthorizerType(authorizer.Type) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid authorizer type: %s", authorizer.Type))
	}
	if authorizer.Type == "" {
		authorizer.Type = "TOKEN"
	}
	if authorizer.AuthorizerResultTtlInSeconds < 0 {
		authorizer.AuthorizerResultTtlInSeconds = 300
	}

	created, err := stores.restApis.CreateAuthorizer(req.Msg.Restapiid, authorizer)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbAuthorizer(created)), nil
}

// GetAuthorizers returns all authorizers for a REST API.
func (h *AdminHandler) GetAuthorizers(ctx context.Context, req *connect.Request[pb.GetAuthorizersRequest]) (*connect.Response[pb.Authorizers], error) {
	if req.Msg.Restapiid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	authorizers, err := stores.restApis.ListAuthorizers(req.Msg.Restapiid)
	if err != nil {
		return nil, storeErr(err)
	}

	limit := int(req.Msg.Limit)
	start, end, nextPos, ok := paginateAdminList(len(authorizers), req.Msg.Position, limit)
	if !ok {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid position: %s", req.Msg.Position))
	}

	items := make([]*pb.Authorizer, 0, end-start)
	for _, a := range authorizers[start:end] {
		items = append(items, toPbAuthorizer(a))
	}
	resp := &pb.Authorizers{Items: items}
	if nextPos != "" {
		resp.Position = nextPos
	}
	return connect.NewResponse(resp), nil
}

// GetAuthorizer returns a single authorizer by ID.
func (h *AdminHandler) GetAuthorizer(ctx context.Context, req *connect.Request[pb.GetAuthorizerRequest]) (*connect.Response[pb.Authorizer], error) {
	if req.Msg.Restapiid == "" || req.Msg.Authorizerid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id and authorizer_id are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	a, err := stores.restApis.GetAuthorizer(req.Msg.Restapiid, req.Msg.Authorizerid)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbAuthorizer(a)), nil
}

// DeleteAuthorizer removes an authorizer from a REST API.
func (h *AdminHandler) DeleteAuthorizer(ctx context.Context, req *connect.Request[pb.DeleteAuthorizerRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Restapiid == "" || req.Msg.Authorizerid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id and authorizer_id are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := stores.restApis.DeleteAuthorizer(req.Msg.Restapiid, req.Msg.Authorizerid); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}
