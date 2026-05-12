package apigateway

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// GetResources returns all resources for a REST API.
func (h *AdminHandler) GetResources(ctx context.Context, req *connect.Request[pb.GetResourcesRequest]) (*connect.Response[pb.Resources], error) {
	if req.Msg.Restapiid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	resources, err := stores.restApis.ListResources(req.Msg.Restapiid)
	if err != nil {
		return nil, storeErr(err)
	}

	items := make([]*pb.Resource, 0, len(resources))
	for _, r := range resources {
		items = append(items, toPbResource(r))
	}
	return connect.NewResponse(&pb.Resources{Items: items}), nil
}

// GetResource returns a single resource by ID.
func (h *AdminHandler) GetResource(ctx context.Context, req *connect.Request[pb.GetResourceRequest]) (*connect.Response[pb.Resource], error) {
	if req.Msg.Restapiid == "" || req.Msg.Resourceid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id and resource_id are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	r, err := stores.restApis.GetResource(req.Msg.Restapiid, req.Msg.Resourceid)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbResource(r)), nil
}

// CreateResource creates a child resource under a parent.
func (h *AdminHandler) CreateResource(ctx context.Context, req *connect.Request[pb.CreateResourceRequest]) (*connect.Response[pb.Resource], error) {
	if req.Msg.Restapiid == "" || req.Msg.Parentid == "" || req.Msg.Pathpart == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id, parent_id, and path_part are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	parent, err := stores.restApis.GetResource(req.Msg.Restapiid, req.Msg.Parentid)
	if err != nil {
		return nil, storeErr(err)
	}

	path := parent.Path
	if path != "/" {
		path += "/"
	}
	path += req.Msg.Pathpart

	resource := &apigatewaystore.Resource{
		ParentId:        req.Msg.Parentid,
		Path:            path,
		PathPart:        req.Msg.Pathpart,
		ResourceMethods: make(map[string]*apigatewaystore.Method),
	}

	created, err := stores.restApis.CreateResource(req.Msg.Restapiid, resource)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbResource(created)), nil
}

// DeleteResource removes a resource from a REST API.
func (h *AdminHandler) DeleteResource(ctx context.Context, req *connect.Request[pb.DeleteResourceRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Restapiid == "" || req.Msg.Resourceid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id and resource_id are required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := stores.restApis.DeleteResource(req.Msg.Restapiid, req.Msg.Resourceid); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}
