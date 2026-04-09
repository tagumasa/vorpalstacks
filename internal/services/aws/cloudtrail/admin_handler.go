package cloudtrail

import (
	"context"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/cloudtrail"
	cloudtrailconnect "vorpalstacks/internal/pb/aws/cloudtrail/cloudtrailconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// AdminHandler implements the CloudTrail admin console gRPC-Web handler.
type AdminHandler struct {
	cloudtrailconnect.UnimplementedCloudTrailServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ cloudtrailconnect.CloudTrailServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudTrail admin handler with the given storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getCloudTrailStore(req *connect.Request[pb.ListTrailsRequest]) (*cloudtrailstore.CloudTrailStore, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return cloudtrailstore.NewCloudTrailStore(regionStorage, h.accountId, region), nil
}

// ListTrails retrieves all CloudTrail trails from the regional store.
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
