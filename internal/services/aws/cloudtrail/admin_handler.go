package cloudtrail

import (
	"context"
	"net/http"
	"vorpalstacks/internal/common/defaults"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/cloudtrail"
	cloudtrailconnect "vorpalstacks/internal/pb/aws/cloudtrail/cloudtrailconnect"
)

// AdminHandler implements the CloudTrail admin console gRPC-Web handler.
// It delegates to the shared CloudTrailService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
type AdminHandler struct {
	cloudtrailconnect.UnimplementedCloudTrailServiceHandler
	service *CloudTrailService
}

var _ cloudtrailconnect.CloudTrailServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudTrail admin handler backed by the given
// service instance, ensuring the same per-region cached stores are used as
// the HTTP API handlers.
func NewAdminHandler(svc *CloudTrailService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreFromHeader(header http.Header) (StoreInterface, error) {
	region := defaults.GetRegionFromHeader(header)
	return h.service.GetStoreForRegion(region)
}

// ListTrails retrieves CloudTrail trails with pagination support.
func (h *AdminHandler) ListTrails(ctx context.Context, req *connect.Request[pb.ListTrailsRequest]) (*connect.Response[pb.ListTrailsResponse], error) {
	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.listTrailsCore(store, ListTrailsInput{
		NextToken: req.Msg.GetNexttoken(),
		MaxItems:  100,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	var trailInfos []*pb.TrailInfo
	for _, t := range result.Items {
		trailInfos = append(trailInfos, &pb.TrailInfo{
			Name:       t.Name,
			Trailarn:   t.TrailARN,
			Homeregion: t.HomeRegion,
		})
	}

	return connect.NewResponse(&pb.ListTrailsResponse{
		Trails:    trailInfos,
		Nexttoken: result.NextToken,
	}), nil
}

// CreateTrail creates a new CloudTrail trail via the admin console.
func (h *AdminHandler) CreateTrail(ctx context.Context, req *connect.Request[pb.CreateTrailRequest]) (*connect.Response[pb.CreateTrailResponse], error) {
	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	in := CreateTrailInput{
		Name:                      req.Msg.GetName(),
		S3BucketName:              req.Msg.GetS3Bucketname(),
		S3KeyPrefix:               req.Msg.GetS3Keyprefix(),
		SnsTopicName:              req.Msg.GetSnstopicname(),
		CloudWatchLogsLogGroupARN: req.Msg.GetCloudwatchlogsloggrouparn(),
		CloudWatchLogsRoleARN:     req.Msg.GetCloudwatchlogsrolearn(),
		KMSKeyID:                  req.Msg.GetKmskeyid(),
		Region:                    defaults.GetRegionFromHeader(req.Header()),
	}
	if v := req.Msg.Includeglobalserviceevents; v != nil {
		in.IncludeGlobalServiceEvents = v
	}
	if v := req.Msg.Ismultiregiontrail; v != nil {
		in.IsMultiRegionTrail = v
	}
	if v := req.Msg.Isorganizationtrail; v != nil {
		in.IsOrganizationTrail = v
	}
	if v := req.Msg.Enablelogfilevalidation; v != nil {
		in.EnableLogFileValidation = v
	}

	created, err := h.service.createTrailCore(store, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateTrailResponse{
		Name:                       created.Name,
		Trailarn:                   created.TrailARN,
		S3Bucketname:               created.S3BucketName,
		S3Keyprefix:                created.S3KeyPrefix,
		Snstopicname:               created.SnsTopicName,
		Snstopicarn:                created.SnsTopicARN,
		Includeglobalserviceevents: proto.Bool(created.IncludeGlobalServiceEvents),
		Ismultiregiontrail:         proto.Bool(created.IsMultiRegionTrail),
		Isorganizationtrail:        proto.Bool(created.IsOrganizationTrail),
		Logfilevalidationenabled:   proto.Bool(created.LogFileValidationEnabled),
		Cloudwatchlogsloggrouparn:  created.CloudWatchLogsLogGroupARN,
		Cloudwatchlogsrolearn:      created.CloudWatchLogsRoleARN,
		Kmskeyid:                   created.KMSKeyID,
	}), nil
}

// DeleteTrail deletes a CloudTrail trail via the admin console.
func (h *AdminHandler) DeleteTrail(ctx context.Context, req *connect.Request[pb.DeleteTrailRequest]) (*connect.Response[pb.DeleteTrailResponse], error) {
	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.deleteTrailCore(store, DeleteTrailInput{NameOrARN: req.Msg.GetName()}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteTrailResponse{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the CloudTrail admin console.
func NewConnectHandler(svc *CloudTrailService) (string, http.Handler) {
	return cloudtrailconnect.NewCloudTrailServiceHandler(NewAdminHandler(svc))
}
