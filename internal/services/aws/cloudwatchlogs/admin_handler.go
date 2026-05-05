package cloudwatchlogs

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/cloudwatchlogs"
	cloudwatchlogsconnect "vorpalstacks/internal/pb/aws/cloudwatchlogs/cloudwatchlogsconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	cloudwatchlogsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// AdminHandler provides CloudWatch Logs service administration functionality.
// It delegates to the shared LogsService to ensure the same per-region cached
// stores are used as the HTTP API handlers.
type AdminHandler struct {
	cloudwatchlogsconnect.UnimplementedCloudWatchLogsServiceHandler
	service *LogsService
}

// NewAdminHandler creates a new CloudWatch Logs AdminHandler backed by the
// given service instance.
func NewAdminHandler(svc *LogsService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// getStoreByRegion retrieves the CloudWatch Logs store for the specified region
// via the shared service cache.
func (h *AdminHandler) getStoreByRegion(region string) (*cloudwatchlogsstore.Store, error) {
	return h.service.getLogsStoreByRegion(region)
}

// ListLogGroups lists log groups in CloudWatch Logs.
func (h *AdminHandler) ListLogGroups(ctx context.Context, req *connect.Request[pb.ListLogGroupsRequest]) (*connect.Response[pb.ListLogGroupsResponse], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getStoreByRegion(region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	limit := int(req.Msg.Limit)
	if limit <= 0 {
		limit = 50
	}

	groups, nextToken, err := store.ListLogGroups(req.Msg.Loggroupnamepattern, req.Msg.Nexttoken, limit)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	summaries := make([]*pb.LogGroupSummary, len(groups))
	for i, g := range groups {
		summary := &pb.LogGroupSummary{
			Loggrouparn:  g.ARN,
			Loggroupname: g.Name,
		}
		switch g.LogGroupClass {
		case "DELIVERY":
			summary.Loggroupclass = pb.LogGroupClass_LOG_GROUP_CLASS_DELIVERY
		case "INFREQUENT_ACCESS":
			summary.Loggroupclass = pb.LogGroupClass_LOG_GROUP_CLASS_INFREQUENT_ACCESS
		default:
			summary.Loggroupclass = pb.LogGroupClass_LOG_GROUP_CLASS_STANDARD
		}
		summaries[i] = summary
	}

	return connect.NewResponse(&pb.ListLogGroupsResponse{
		Loggroups: summaries,
		Nexttoken: nextToken,
	}), nil
}

// DescribeLogStreams describes log streams in CloudWatch Logs.
func (h *AdminHandler) DescribeLogStreams(ctx context.Context, req *connect.Request[pb.DescribeLogStreamsRequest]) (*connect.Response[pb.DescribeLogStreamsResponse], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getStoreByRegion(region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	limit := int(req.Msg.Limit)
	if limit <= 0 {
		limit = 50
	}

	streams, nextToken, err := store.ListLogStreams(req.Msg.Loggroupname, req.Msg.Logstreamnameprefix, req.Msg.Nexttoken, limit)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	pbStreams := make([]*pb.LogStream, len(streams))
	for i, s := range streams {
		pbStreams[i] = &pb.LogStream{
			Logstreamname:       s.Name,
			Arn:                 s.ARN,
			Creationtime:        s.CreatedAt.UnixMilli(),
			Firsteventtimestamp: s.FirstEventTs,
			Lasteventtimestamp:  s.LastEventTs,
			Lastingestiontime:   s.LastIngestionTs,
			Uploadsequencetoken: s.UploadSequenceToken,
		}
	}

	return connect.NewResponse(&pb.DescribeLogStreamsResponse{
		Logstreams: pbStreams,
		Nexttoken:  nextToken,
	}), nil
}

// CreateLogGroup creates a new CloudWatch Logs log group via the admin console.
func (h *AdminHandler) CreateLogGroup(ctx context.Context, req *connect.Request[pb.CreateLogGroupRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Loggroupname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("LogGroupName is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getStoreByRegion(region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	logGroupClass := "STANDARD"
	switch req.Msg.Loggroupclass {
	case pb.LogGroupClass_LOG_GROUP_CLASS_DELIVERY:
		logGroupClass = "DELIVERY"
	case pb.LogGroupClass_LOG_GROUP_CLASS_INFREQUENT_ACCESS:
		logGroupClass = "INFREQUENT_ACCESS"
	}

	if err := store.CreateLogGroup(&cloudwatchlogsstore.LogGroup{
		Name:          req.Msg.Loggroupname,
		Tags:          req.Msg.Tags,
		LogGroupClass: logGroupClass,
	}); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DeleteLogGroup deletes a CloudWatch Logs log group by name via the admin console.
func (h *AdminHandler) DeleteLogGroup(ctx context.Context, req *connect.Request[pb.DeleteLogGroupRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Loggroupname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("LogGroupName is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getStoreByRegion(region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := store.DeleteLogGroup(req.Msg.Loggroupname); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the CloudWatch Logs admin console.
func NewConnectHandler(svc *LogsService) (string, http.Handler) {
	return cloudwatchlogsconnect.NewCloudWatchLogsServiceHandler(NewAdminHandler(svc))
}
