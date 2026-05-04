// Package apigateway implements AWS API Gateway REST API operations including
// REST APIs, resources, methods, integrations, deployments, authorisers, and
// usage plans.
package apigateway

import (
	"context"
	"fmt"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/apigateway"
	apigatewayconnect "vorpalstacks/internal/pb/aws/apigateway/apigatewayconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/store/aws/common"
)

// AdminHandler implements the API Gateway admin console gRPC-Web handler.
type AdminHandler struct {
	apigatewayconnect.UnimplementedAPIGatewayServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ apigatewayconnect.APIGatewayServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new API Gateway admin handler with the given storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*apigatewaystore.RestApiStore, error) {
	region := svccommon.GetRegionFromHeader(headers)
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return apigatewaystore.NewRestApiStore(regionStorage, h.accountId, region), nil
}

// GetRestApis lists all REST APIs registered in the API Gateway service.
func (h *AdminHandler) GetRestApis(ctx context.Context, req *connect.Request[pb.GetRestApisRequest]) (*connect.Response[pb.RestApis], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	limit := int(req.Msg.Limit)
	if limit <= 0 {
		limit = 100
	}
	result, err := store.List(common.ListOptions{
		Marker:   req.Msg.Position,
		MaxItems: limit,
	})
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var items []*pb.RestApi
	for _, api := range result.Items {
		items = append(items, toPbRestApi(api))
	}

	return connect.NewResponse(&pb.RestApis{
		Items:    items,
		Position: result.NextMarker,
	}), nil
}

// CreateRestApi creates a new REST API via the admin console.
func (h *AdminHandler) CreateRestApi(ctx context.Context, req *connect.Request[pb.CreateRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	api := &apigatewaystore.RestApi{
		Name:                   req.Msg.Name,
		Description:            req.Msg.Description,
		Version:                req.Msg.Version,
		BinaryMediaTypes:       req.Msg.Binarymediatypes,
		MinimumCompressionSize: req.Msg.Minimumcompressionsize,
		Policy:                 req.Msg.Policy,
	}

	if req.Msg.Endpointconfiguration != nil {
		types := make([]string, len(req.Msg.Endpointconfiguration.Types))
		for i, t := range req.Msg.Endpointconfiguration.Types {
			types[i] = t.String()
		}
		api.EndpointConfiguration = &apigatewaystore.EndpointConfiguration{
			Types: types,
		}
	}

	created, err := store.Create(api)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(toPbRestApi(created)), nil
}

// DeleteRestApi deletes a REST API via the admin console.
func (h *AdminHandler) DeleteRestApi(ctx context.Context, req *connect.Request[pb.DeleteRestApiRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Restapiid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("rest_api_id is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := store.Delete(req.Msg.Restapiid); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// toPbRestApi converts a store RestApi to a proto RestApi.
func toPbRestApi(api *apigatewaystore.RestApi) *pb.RestApi {
	pbApi := &pb.RestApi{
		Id:          api.Id,
		Name:        api.Name,
		Description: api.Description,
		Version:     api.Version,
		Warnings:    api.Warnings,
		Createddate: api.CreatedDate.Format(timeutils.ISO8601UTCFormat),
	}
	if api.EndpointConfiguration != nil {
		types := make([]pb.EndpointType, len(api.EndpointConfiguration.Types))
		for i, t := range api.EndpointConfiguration.Types {
			switch t {
			case "PRIVATE":
				types[i] = pb.EndpointType_ENDPOINT_TYPE_PRIVATE
			case "REGIONAL":
				types[i] = pb.EndpointType_ENDPOINT_TYPE_REGIONAL
			case "EDGE":
				types[i] = pb.EndpointType_ENDPOINT_TYPE_EDGE
			}
		}
		pbApi.Endpointconfiguration = &pb.EndpointConfiguration{
			Types: types,
		}
	}
	return pbApi
}

// NewConnectHandler creates a gRPC-Web connect handler for the API Gateway admin console.
func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return apigatewayconnect.NewAPIGatewayServiceHandler(NewAdminHandler(sm, accountID))
}
