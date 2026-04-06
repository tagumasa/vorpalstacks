package appsync

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/appsync"
	appsyncconnect "vorpalstacks/internal/pb/aws/appsync/appsyncconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// AdminHandler implements the AppSync gRPC-Web admin console handler. It exposes
// list and describe operations for the Flutter management UI, delegating data
// access to the AppSyncStore via the shared AppSyncService cache.
type AdminHandler struct {
	appsyncconnect.UnimplementedAppSyncServiceHandler
	service        *AppSyncService
	storageManager *storage.RegionStorageManager
}

var _ appsyncconnect.AppSyncServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new AppSync admin console handler backed by the
// given service instance, ensuring the same per-region cached stores are used
// as the HTTP API handlers.
func NewAdminHandler(svc *AppSyncService, sm *storage.RegionStorageManager) *AdminHandler {
	return &AdminHandler{service: svc, storageManager: sm}
}

func (h *AdminHandler) getStoreByHeader(header http.Header) (*appsyncstore.AppSyncStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	rs, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return h.service.GetStoreForRegion(region, rs), nil
}

// ListApis returns a paginated list of AppSync APIs in the requested region.
func (h *AdminHandler) ListApis(ctx context.Context, req *connect.Request[pb.ListApisRequest]) (*connect.Response[pb.ListApisResponse], error) {
	store, err := h.getStoreByHeader(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	limit := int(req.Msg.Maxresults)
	if limit <= 0 {
		limit = 25
	}

	opts := storecommon.ListOptions{
		MaxItems: limit,
		Marker:   req.Msg.Nexttoken,
	}

	apis, nextToken, err := store.ListApis(opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbApis := make([]*pb.Api, len(apis))
	for i, a := range apis {
		pbApis[i] = &pb.Api{
			Apiid:        a.ApiId,
			Name:         a.Name,
			Apiarn:       a.Arn,
			Dns:          a.Dns,
			Tags:         a.Tags,
			Xrayenabled:  a.XrayEnabled,
			Wafwebaclarn: a.WafWebAclArn,
		}
	}

	return connect.NewResponse(&pb.ListApisResponse{
		Apis:      pbApis,
		Nexttoken: nextToken,
	}), nil
}

// ListGraphqlApis returns a paginated list of GraphQL APIs in the requested region.
func (h *AdminHandler) ListGraphqlApis(ctx context.Context, req *connect.Request[pb.ListGraphqlApisRequest]) (*connect.Response[pb.ListGraphqlApisResponse], error) {
	store, err := h.getStoreByHeader(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	limit := int(req.Msg.Maxresults)
	if limit <= 0 {
		limit = 25
	}

	opts := storecommon.ListOptions{
		MaxItems: limit,
		Marker:   req.Msg.Nexttoken,
	}

	graphqlApis, nextToken, err := store.ListGraphqlApis(opts, "")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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

func toPbGraphqlApi(a *appsyncstore.GraphqlApi) *pb.GraphqlApi {
	return &pb.GraphqlApi{
		Name:         a.Name,
		Apiid:        a.ApiId,
		Arn:          a.Arn,
		Uris:         a.Uris,
		Tags:         a.Tags,
		Xrayenabled:  a.XrayEnabled,
		Wafwebaclarn: a.WafWebAclArn,
	}
}
