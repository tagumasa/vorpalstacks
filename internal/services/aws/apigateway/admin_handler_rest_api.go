package apigateway

import (
	"context"

	"connectrpc.com/connect"

	tagutil "vorpalstacks/internal/common/tags"
	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// GetRestApis returns all REST APIs.
func (h *AdminHandler) GetRestApis(ctx context.Context, req *connect.Request[pb.GetRestApisRequest]) (*connect.Response[pb.RestApis], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	result, err := h.service.listRestApisCore(stores, int(req.Msg.GetLimit()), req.Msg.Position)
	if err != nil {
		return nil, storeErr(err)
	}

	items := make([]*pb.RestApi, 0, len(result.Items))
	for _, api := range result.Items {
		items = append(items, toPbRestApi(api))
	}

	resp := &pb.RestApis{Items: items}
	if result.NextMarker != "" {
		resp.Position = result.NextMarker
	}
	return connect.NewResponse(resp), nil
}

// GetRestApi returns a single REST API by ID.
func (h *AdminHandler) GetRestApi(ctx context.Context, req *connect.Request[pb.GetRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	api, err := h.service.getRestApiCore(stores, req.Msg.Restapiid)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbRestApi(api)), nil
}

// CreateRestApi creates a new REST API.
func (h *AdminHandler) CreateRestApi(ctx context.Context, req *connect.Request[pb.CreateRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	api := &apigatewaystore.RestApi{
		Name:               req.Msg.Name,
		Description:        req.Msg.Description,
		Version:            req.Msg.Version,
		BinaryMediaTypes:   req.Msg.Binarymediatypes,
		ApiKeySource:       apiKeySourceFromPb(req.Msg.Apikeysource),
		Policy:             req.Msg.Policy,
		SecurityPolicy:     securityPolicyFromPb(req.Msg.Securitypolicy),
		EndpointAccessMode: endpointAccessModeFromPb(req.Msg.Endpointaccessmode),
	}
	if req.Msg.Disableexecuteapiendpoint != nil {
		api.DisableExecuteApiEndpoint = *req.Msg.Disableexecuteapiendpoint
	}
	if req.Msg.Minimumcompressionsize != nil {
		v := *req.Msg.Minimumcompressionsize
		api.MinimumCompressionSize = &v
	}
	if req.Msg.Endpointconfiguration != nil {
		types := make([]string, len(req.Msg.Endpointconfiguration.Types))
		for i, t := range req.Msg.Endpointconfiguration.Types {
			types[i] = fromPbEndpointType(t)
		}
		api.EndpointConfiguration = &apigatewaystore.EndpointConfiguration{Types: types}
	}
	if len(req.Msg.Tags) > 0 {
		api.Tags = tagutil.MapToTags(req.Msg.Tags)
	}

	created, err := h.service.createRestApiCore(stores, api, req.Msg.Clonefrom)
	if err != nil {
		return nil, storeErr(err)
	}

	return connect.NewResponse(toPbRestApi(created)), nil
}

// DeleteRestApi removes a REST API.
func (h *AdminHandler) DeleteRestApi(ctx context.Context, req *connect.Request[pb.DeleteRestApiRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := h.service.deleteRestApiCore(stores, req.Msg.Restapiid); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// UpdateRestApi modifies an existing REST API.
func (h *AdminHandler) UpdateRestApi(ctx context.Context, req *connect.Request[pb.UpdateRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	patches := make([]PatchOperation, 0, len(req.Msg.Patchoperations))
	for _, po := range req.Msg.Patchoperations {
		patches = append(patches, PatchOperation{
			Op:    opFromPb(po.Op),
			Path:  po.Path,
			Value: po.Value,
		})
	}

	api, err := h.service.updateRestApiCore(stores, req.Msg.Restapiid, patches)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbRestApi(api)), nil
}

// opFromPb maps a protobuf patch op enum to its canonical string form.
func opFromPb(op pb.Op) string {
	switch op {
	case pb.Op_OP_ADD:
		return "add"
	case pb.Op_OP_REMOVE:
		return "remove"
	case pb.Op_OP_REPLACE:
		return "replace"
	default:
		return ""
	}
}
