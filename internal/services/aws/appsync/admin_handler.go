// Package appsync implements AWS AppSync API operations including GraphQL API
// management, data sources, resolvers, schema, real-time subscriptions, and
// code/mapping template evaluation.
package appsync

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/appsync"
	appsyncconnect "vorpalstacks/internal/pb/aws/appsync/appsyncconnect"
)

// AdminHandler implements the AppSync gRPC-Web admin console handler.
type AdminHandler struct {
	appsyncconnect.UnimplementedAppSyncServiceHandler
	service *AppSyncService
}

var _ appsyncconnect.AppSyncServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new AppSync admin console handler.
func NewAdminHandler(svc *AppSyncService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListApis returns a paginated list of AppSync APIs in the requested region.
func (h *AdminHandler) ListApis(ctx context.Context, req *connect.Request[pb.ListApisRequest]) (*connect.Response[pb.ListApisResponse], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	apis, nextToken, err := h.service.listApisCore(store, int(req.Msg.GetMaxresults()), req.Msg.Nexttoken)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	pbApis := make([]*pb.Api, len(apis))
	for i, a := range apis {
		pbApis[i] = toPbApi(a)
	}

	return connect.NewResponse(&pb.ListApisResponse{
		Apis:      pbApis,
		Nexttoken: nextToken,
	}), nil
}

// ListGraphqlApis returns a paginated list of GraphQL APIs in the requested region.
func (h *AdminHandler) ListGraphqlApis(ctx context.Context, req *connect.Request[pb.ListGraphqlApisRequest]) (*connect.Response[pb.ListGraphqlApisResponse], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	graphqlApis, nextToken, err := h.service.listGraphqlApisCore(store, int(req.Msg.GetMaxresults()), req.Msg.Nexttoken, "")
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	pbApis := make([]*pb.GraphqlApi, len(graphqlApis))
	for i, a := range graphqlApis {
		pbApis[i] = toPbGraphqlApi(a)
	}

	return connect.NewResponse(&pb.ListGraphqlApisResponse{
		Graphqlapis: pbApis,
		Nexttoken:   nextToken,
	}), nil
}

// CreateGraphqlApi creates a new AppSync GraphQL API via the admin console.
func (h *AdminHandler) CreateGraphqlApi(ctx context.Context, req *connect.Request[pb.CreateGraphqlApiRequest]) (*connect.Response[pb.CreateGraphqlApiResponse], error) {
	if req.Msg.GetName() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	authType, err := pbAuthTypeToString(req.Msg.GetAuthenticationtype())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.createGraphqlApiCore(store, createGraphqlApiInput{
		Name:               req.Msg.GetName(),
		AuthenticationType: authType,
		Tags:               req.Msg.GetTags(),
		XrayEnabled:        req.Msg.GetXrayenabled(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateGraphqlApiResponse{
		Graphqlapi: toPbGraphqlApi(result),
	}), nil
}

// DeleteGraphqlApi deletes an AppSync GraphQL API via the admin console.
func (h *AdminHandler) DeleteGraphqlApi(ctx context.Context, req *connect.Request[pb.DeleteGraphqlApiRequest]) (*connect.Response[pb.DeleteGraphqlApiResponse], error) {
	if req.Msg.GetApiid() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("ApiId is required"))
	}

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := h.service.deleteGraphqlApiCore(store, req.Msg.GetApiid()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteGraphqlApiResponse{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Appsync admin console.
func NewConnectHandler(svc *AppSyncService) (string, http.Handler) {
	return appsyncconnect.NewAppSyncServiceHandler(NewAdminHandler(svc))
}
