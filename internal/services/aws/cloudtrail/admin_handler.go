package cloudtrail

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/cloudtrail"
	cloudtrailconnect "vorpalstacks/internal/pb/aws/cloudtrail/cloudtrailconnect"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// AdminHandler provides CloudTrail service administration functionality.
// It implements the CloudTrailServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	cloudtrailconnect.UnimplementedCloudTrailServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ cloudtrailconnect.CloudTrailServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudTrail AdminHandler.
// It initialises the handler with the provided storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

// getCloudTrailStore retrieves the CloudTrail store for the request.
// It extracts the region from the request header and creates a new CloudTrailStore instance.
func (h *AdminHandler) getCloudTrailStore(req *connect.Request[pb.ListTrailsRequest]) (*cloudtrailstore.CloudTrailStore, error) {
	region := req.Header().Get("X-Aws-Region")
	if region == "" {
		region = "us-east-1"
	}
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return cloudtrailstore.NewCloudTrailStore(regionStorage, h.accountId, region), nil
}

// ListTrails lists trails in CloudTrail.
func (h *AdminHandler) ListTrails(ctx context.Context, req *connect.Request[pb.ListTrailsRequest]) (*connect.Response[pb.ListTrailsResponse], error) {
	store, err := h.getCloudTrailStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	trails, err := store.ListTrails()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var trailInfos []*pb.TrailInfo
	for _, trail := range trails {
		trailInfos = append(trailInfos, &pb.TrailInfo{
			Name:       trail.Name,
			Trailarn:   trail.TrailARN,
			Homeregion: trail.HomeRegion,
		})
	}

	return connect.NewResponse(&pb.ListTrailsResponse{
		Trails: trailInfos,
	}), nil
}

// AddTags adds tags to a trail or event data store.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AddTags(ctx context.Context, req *connect.Request[pb.AddTagsRequest]) (*connect.Response[pb.AddTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CancelQuery cancels a CloudTrail Lake query.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CancelQuery(ctx context.Context, req *connect.Request[pb.CancelQueryRequest]) (*connect.Response[pb.CancelQueryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateChannel creates a new channel for ingesting events.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateChannel(ctx context.Context, req *connect.Request[pb.CreateChannelRequest]) (*connect.Response[pb.CreateChannelResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateDashboard creates a new CloudTrail Lake dashboard.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDashboard(ctx context.Context, req *connect.Request[pb.CreateDashboardRequest]) (*connect.Response[pb.CreateDashboardResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateEventDataStore creates a new event data store for ingesting events.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateEventDataStore(ctx context.Context, req *connect.Request[pb.CreateEventDataStoreRequest]) (*connect.Response[pb.CreateEventDataStoreResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateTrail creates a new trail in CloudTrail.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateTrail(ctx context.Context, req *connect.Request[pb.CreateTrailRequest]) (*connect.Response[pb.CreateTrailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteChannel deletes a channel.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteChannel(ctx context.Context, req *connect.Request[pb.DeleteChannelRequest]) (*connect.Response[pb.DeleteChannelResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteDashboard deletes a CloudTrail Lake dashboard.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDashboard(ctx context.Context, req *connect.Request[pb.DeleteDashboardRequest]) (*connect.Response[pb.DeleteDashboardResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteEventDataStore deletes an event data store.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteEventDataStore(ctx context.Context, req *connect.Request[pb.DeleteEventDataStoreRequest]) (*connect.Response[pb.DeleteEventDataStoreResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteResourcePolicy deletes a resource policy from a trail or event data store.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteResourcePolicy(ctx context.Context, req *connect.Request[pb.DeleteResourcePolicyRequest]) (*connect.Response[pb.DeleteResourcePolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteTrail deletes a trail.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteTrail(ctx context.Context, req *connect.Request[pb.DeleteTrailRequest]) (*connect.Response[pb.DeleteTrailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
