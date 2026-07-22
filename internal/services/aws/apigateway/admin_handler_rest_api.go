package apigateway

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"

	tagutil "vorpalstacks/internal/common/tags"
	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// GetRestApis returns all REST APIs.
func (h *AdminHandler) GetRestApis(ctx context.Context, req *connect.Request[pb.GetRestApisRequest]) (*connect.Response[pb.RestApis], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}

	limit := int(req.Msg.Limit)
	if limit <= 0 {
		limit = 100
	}
	result, err := stores.restApis.List(storecommon.ListOptions{
		Marker:   req.Msg.Position,
		MaxItems: limit,
	})
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
	if req.Msg.Restapiid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	api, err := stores.restApis.Get(req.Msg.Restapiid)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbRestApi(api)), nil
}

// CreateRestApi creates a new REST API.
func (h *AdminHandler) CreateRestApi(ctx context.Context, req *connect.Request[pb.CreateRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}
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
	if req.Msg.Minimumcompressionsize > 0 {
		v := req.Msg.Minimumcompressionsize
		api.MinimumCompressionSize = &v
	}
	if req.Msg.Endpointconfiguration != nil {
		types := make([]string, len(req.Msg.Endpointconfiguration.Types))
		for i, t := range req.Msg.Endpointconfiguration.Types {
			types[i] = t.String()
		}
		api.EndpointConfiguration = &apigatewaystore.EndpointConfiguration{Types: types}
	}
	if len(req.Msg.Tags) > 0 {
		api.Tags = tagutil.MapToTags(req.Msg.Tags)
	}

	created, err := stores.restApis.Create(api)
	if err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbRestApi(created)), nil
}

// DeleteRestApi removes a REST API.
func (h *AdminHandler) DeleteRestApi(ctx context.Context, req *connect.Request[pb.DeleteRestApiRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Restapiid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	if err := stores.restApis.Delete(req.Msg.Restapiid); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// UpdateRestApi modifies an existing REST API.
func (h *AdminHandler) UpdateRestApi(ctx context.Context, req *connect.Request[pb.UpdateRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	if req.Msg.Restapiid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id is required"))
	}
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, storeErr(err)
	}
	api, err := stores.restApis.Get(req.Msg.Restapiid)
	if err != nil {
		return nil, storeErr(err)
	}

	for _, po := range req.Msg.Patchoperations {
		switch po.Path {
		case "/name":
			api.Name = po.Value
		case "/description":
			api.Description = po.Value
		case "/version":
			api.Version = po.Value
		case "/apiKeySource":
			api.ApiKeySource = po.Value
		case "/policy":
			api.Policy = po.Value
		case "/disableExecuteApiEndpoint":
			api.DisableExecuteApiEndpoint = po.Value == "true"
		case "/securityPolicy":
			api.SecurityPolicy = po.Value
		case "/endpointAccessMode":
			api.EndpointAccessMode = po.Value
		case "/minimumCompressionSize":
			v, err := parseInt32(po.Value)
			if err != nil {
				return nil, NewBadRequestException("invalid minimumCompressionSize: not a number")
			}
			api.MinimumCompressionSize = &v
		}

		if strings.HasPrefix(po.Path, "/binaryMediaTypes/") {
		MediaTypeLoop:
			for i, mt := range api.BinaryMediaTypes {
				if mt == po.Value {
					if po.Op == pb.Op_OP_REMOVE {
						api.BinaryMediaTypes = append(api.BinaryMediaTypes[:i], api.BinaryMediaTypes[i+1:]...)
					}
					break MediaTypeLoop
				}
			}
			if (po.Op == pb.Op_OP_ADD || po.Op == pb.Op_OP_REPLACE) && !containsAny(api.BinaryMediaTypes, po.Value) {
				api.BinaryMediaTypes = append(api.BinaryMediaTypes, po.Value)
			}
		}
	}

	if err := stores.restApis.Update(api); err != nil {
		return nil, storeErr(err)
	}
	return connect.NewResponse(toPbRestApi(api)), nil
}
