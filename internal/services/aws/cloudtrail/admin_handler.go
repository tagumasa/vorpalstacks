package cloudtrail

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/common/defaults"
	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"

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

// adminListTrailsPageSize is the admin console's ListTrails page size — a
// platform paging choice, not an AWS limit.
const adminListTrailsPageSize = 100

// ListTrails retrieves CloudTrail trails with pagination support.
func (h *AdminHandler) ListTrails(ctx context.Context, req *connect.Request[pb.ListTrailsRequest]) (*connect.Response[pb.ListTrailsResponse], error) {
	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.listTrailsCore(store, ListTrailsInput{
		NextToken: req.Msg.GetNexttoken(),
		MaxItems:  adminListTrailsPageSize,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	var trailInfos []*pb.TrailInfo
	for _, t := range result.Items {
		trailInfos = append(trailInfos, &pb.TrailInfo{
			Name:       proto.String(t.Name),
			Trailarn:   proto.String(t.TrailARN),
			Homeregion: proto.String(t.HomeRegion),
		})
	}

	return connect.NewResponse(&pb.ListTrailsResponse{
		Trails:    trailInfos,
		Nexttoken: proto.String(result.NextToken),
	}), nil
}

// CreateTrail creates a new CloudTrail trail via the admin console.
func (h *AdminHandler) CreateTrail(ctx context.Context, req *connect.Request[pb.CreateTrailRequest]) (*connect.Response[pb.CreateTrailResponse], error) {
	store, err := h.getStoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	// The CloudWatchLogsRoleArn trust validation runs inside the Core; the
	// admin plane builds the same validator the HTTP request context
	// provides.
	var iamValidator *iam.IAMValidator
	if rp := h.service.RoleProvider(); rp != nil {
		iamValidator = iam.NewIAMValidator(rp, h.service.AccountID())
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
		IAMValidator:              iamValidator,
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

	created, err := h.service.createTrailCore(ctx, store, in)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateTrailResponse{
		Name:                       proto.String(created.Name),
		Trailarn:                   proto.String(created.TrailARN),
		S3Bucketname:               proto.String(created.S3BucketName),
		S3Keyprefix:                proto.String(created.S3KeyPrefix),
		Snstopicname:               proto.String(created.SnsTopicName),
		Snstopicarn:                proto.String(created.SnsTopicARN),
		Includeglobalserviceevents: proto.Bool(created.IncludeGlobalServiceEvents),
		Ismultiregiontrail:         proto.Bool(created.IsMultiRegionTrail),
		Isorganizationtrail:        proto.Bool(created.IsOrganizationTrail),
		Logfilevalidationenabled:   proto.Bool(created.LogFileValidationEnabled),
		Cloudwatchlogsloggrouparn:  proto.String(created.CloudWatchLogsLogGroupARN),
		Cloudwatchlogsrolearn:      proto.String(created.CloudWatchLogsRoleARN),
		Kmskeyid:                   proto.String(created.KMSKeyID),
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
