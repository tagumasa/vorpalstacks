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
)

// AdminHandler provides CloudWatch Logs service administration functionality.
// It delegates to the shared LogsService Core methods to ensure the same
// validation, business logic, and per-region cached stores are used as the
// HTTP API handlers. Per architecture rule #29, this handler does NOT import
// any store package directly — all data access flows through Core methods
// and conversion helpers in admin_handler_convert.go.
type AdminHandler struct {
	cloudwatchlogsconnect.UnimplementedCloudWatchLogsServiceHandler
	service *LogsService
}

// NewAdminHandler creates a new AdminHandler bound to the given service.
func NewAdminHandler(svc *LogsService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListLogGroups lists log groups in CloudWatch Logs via the admin console.
func (h *AdminHandler) ListLogGroups(ctx context.Context, req *connect.Request[pb.ListLogGroupsRequest]) (*connect.Response[pb.ListLogGroupsResponse], error) {
	region := svccommon.GetRegionFromHeader(req.Header())

	limit, err := validateListLimit(int32(req.Msg.GetLimit()), 50, 1000)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	input := ListLogGroupsInput{
		LogGroupNamePrefix: req.Msg.Loggroupnamepattern,
		NextToken:          req.Msg.Nexttoken,
		Limit:              limit,
		Region:             region,
	}

	result, err := h.service.listLogGroupsCore(input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	summaries := make([]*pb.LogGroupSummary, len(result.LogGroups))
	for i, lg := range result.LogGroups {
		summaries[i] = toPbLogGroupSummary(lg)
	}

	return connect.NewResponse(&pb.ListLogGroupsResponse{
		Loggroups: summaries,
		Nexttoken: result.NextToken,
	}), nil
}

// DescribeLogStreams describes log streams in CloudWatch Logs via the admin console.
func (h *AdminHandler) DescribeLogStreams(ctx context.Context, req *connect.Request[pb.DescribeLogStreamsRequest]) (*connect.Response[pb.DescribeLogStreamsResponse], error) {
	region := svccommon.GetRegionFromHeader(req.Header())

	input := DescribeLogStreamsInput{
		LogGroupName:        req.Msg.Loggroupname,
		LogStreamNamePrefix: req.Msg.Logstreamnameprefix,
		NextToken:           req.Msg.Nexttoken,
		Limit:               int32(req.Msg.GetLimit()),
		Region:              region,
	}

	result, err := h.service.describeLogStreamsCore(input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	pbStreams := make([]*pb.LogStream, len(result.LogStreams))
	for i, ls := range result.LogStreams {
		pbStreams[i] = toPbLogStream(ls)
	}

	return connect.NewResponse(&pb.DescribeLogStreamsResponse{
		Logstreams: pbStreams,
		Nexttoken:  result.NextToken,
	}), nil
}

// CreateLogGroup creates a new CloudWatch Logs log group via the admin console.
func (h *AdminHandler) CreateLogGroup(ctx context.Context, req *connect.Request[pb.CreateLogGroupRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Loggroupname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("LogGroupName is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())

	input := CreateLogGroupInput{
		LogGroupName:              req.Msg.Loggroupname,
		LogGroupClass:             pbLogGroupClassToString(req.Msg.Loggroupclass),
		Tags:                      req.Msg.Tags,
		DeletionProtectionEnabled: req.Msg.GetDeletionprotectionenabled(),
		Region:                    region,
	}

	if _, err := h.service.createLogGroupCore(input); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DeleteLogGroup deletes a CloudWatch Logs log group by name via the admin console.
func (h *AdminHandler) DeleteLogGroup(ctx context.Context, req *connect.Request[pb.DeleteLogGroupRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Loggroupname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("LogGroupName is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())

	input := DeleteLogGroupInput{
		LogGroupName: req.Msg.Loggroupname,
		Region:       region,
	}

	if err := h.service.deleteLogGroupCore(input); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the CloudWatch Logs admin console.
func NewConnectHandler(svc *LogsService) (string, http.Handler) {
	return cloudwatchlogsconnect.NewCloudWatchLogsServiceHandler(NewAdminHandler(svc))
}
