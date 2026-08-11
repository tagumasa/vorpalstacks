package apigateway

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
)

// GetResources returns all resources for a REST API.
func (h *AdminHandler) GetResources(ctx context.Context, req *connect.Request[pb.GetResourcesRequest]) (*connect.Response[pb.Resources], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	resources, err := h.service.listResourcesCore(stores, req.Msg.Restapiid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	limit := int(req.Msg.GetLimit())
	start, end, nextPos, ok := paginateAdminList(len(resources), req.Msg.Position, limit)
	if !ok {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid position: %s", req.Msg.Position))
	}

	items := make([]*pb.Resource, 0, end-start)
	for _, r := range resources[start:end] {
		items = append(items, toPbResource(r))
	}
	resp := &pb.Resources{Items: items}
	if nextPos != "" {
		resp.Position = nextPos
	}
	return connect.NewResponse(resp), nil
}

// GetResource returns a single resource by ID.
func (h *AdminHandler) GetResource(ctx context.Context, req *connect.Request[pb.GetResourceRequest]) (*connect.Response[pb.Resource], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	r, err := h.service.getResourceCore(stores, req.Msg.Restapiid, req.Msg.Resourceid)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbResource(r)), nil
}

// CreateResource creates a child resource under a parent.
func (h *AdminHandler) CreateResource(ctx context.Context, req *connect.Request[pb.CreateResourceRequest]) (*connect.Response[pb.Resource], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	created, err := h.service.createResourceCore(stores, req.Msg.Restapiid, req.Msg.Parentid, req.Msg.Pathpart)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(toPbResource(created)), nil
}

// DeleteResource removes a resource from a REST API.
func (h *AdminHandler) DeleteResource(ctx context.Context, req *connect.Request[pb.DeleteResourceRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	if err := h.service.deleteResourceCore(stores, req.Msg.Restapiid, req.Msg.Resourceid); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}
