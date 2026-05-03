package cloudtrail

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/cloudtrail"
	cloudtrailconnect "vorpalstacks/internal/pb/aws/cloudtrail/cloudtrailconnect"
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

func (h *AdminHandler) getStoreFromHeader(header http.Header) (*cloudtrailstore.CloudTrailStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return cloudtrailstore.NewCloudTrailStore(regionStorage, h.accountId, region), nil
}

// ListTrails retrieves all CloudTrail trails from the regional store.
func (h *AdminHandler) ListTrails(ctx context.Context, req *connect.Request[pb.ListTrailsRequest]) (*connect.Response[pb.ListTrailsResponse], error) {
	store, err := h.getStoreFromHeader(req.Header())
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

// CreateTrail creates a new CloudTrail trail via the admin console.
func (h *AdminHandler) CreateTrail(ctx context.Context, req *connect.Request[pb.CreateTrailRequest]) (*connect.Response[pb.CreateTrailResponse], error) {
	if req.Msg.GetName() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("Name is required"))
	}
	if req.Msg.GetS3Bucketname() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("S3BucketName is required"))
	}

	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	trail := &cloudtrailstore.Trail{
		Name:                       req.Msg.GetName(),
		S3BucketName:               req.Msg.GetS3Bucketname(),
		S3KeyPrefix:                req.Msg.GetS3Keyprefix(),
		SnsTopicName:               req.Msg.GetSnstopicname(),
		IncludeGlobalServiceEvents: req.Msg.GetIncludeglobalserviceevents(),
		IsMultiRegionTrail:         req.Msg.GetIsmultiregiontrail(),
		IsOrganizationTrail:        req.Msg.GetIsorganizationtrail(),
		LogFileValidationEnabled:   req.Msg.GetEnablelogfilevalidation(),
		CloudWatchLogsLogGroupARN:  req.Msg.GetCloudwatchlogsloggrouparn(),
		CloudWatchLogsRoleARN:      req.Msg.GetCloudwatchlogsrolearn(),
		KMSKeyID:                   req.Msg.GetKmskeyid(),
	}

	result, err := store.CreateTrail(trail)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.CreateTrailResponse{
		Name:                       result.Name,
		Trailarn:                   result.TrailARN,
		S3Bucketname:               result.S3BucketName,
		S3Keyprefix:                result.S3KeyPrefix,
		Snstopicname:               result.SnsTopicName,
		Snstopicarn:                result.SnsTopicARN,
		Includeglobalserviceevents: result.IncludeGlobalServiceEvents,
		Ismultiregiontrail:         result.IsMultiRegionTrail,
		Isorganizationtrail:        result.IsOrganizationTrail,
		Logfilevalidationenabled:   result.LogFileValidationEnabled,
		Cloudwatchlogsloggrouparn:  result.CloudWatchLogsLogGroupARN,
		Cloudwatchlogsrolearn:      result.CloudWatchLogsRoleARN,
		Kmskeyid:                   result.KMSKeyID,
	}), nil
}

// DeleteTrail deletes a CloudTrail trail via the admin console.
func (h *AdminHandler) DeleteTrail(ctx context.Context, req *connect.Request[pb.DeleteTrailRequest]) (*connect.Response[pb.DeleteTrailResponse], error) {
	if req.Msg.GetName() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("Name is required"))
	}

	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := store.DeleteTrail(req.Msg.GetName()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.DeleteTrailResponse{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the CloudTrail admin console.
func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return cloudtrailconnect.NewCloudTrailServiceHandler(NewAdminHandler(sm, accountID))
}
