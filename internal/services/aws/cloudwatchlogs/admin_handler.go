package cloudwatchlogs

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/cloudwatchlogs"
	cloudwatchlogsconnect "vorpalstacks/internal/pb/aws/cloudwatchlogs/cloudwatchlogsconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
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
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	limit := int(req.Msg.Limit)
	if limit <= 0 {
		limit = 50
	}

	groups, nextToken, err := store.ListLogGroups(req.Msg.Loggroupnamepattern, req.Msg.Nexttoken, limit)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	summaries := make([]*pb.LogGroupSummary, len(groups))
	for i, g := range groups {
		summaries[i] = &pb.LogGroupSummary{
			Loggrouparn:   g.ARN,
			Loggroupname:  g.Name,
			Loggroupclass: pb.LogGroupClass_LOG_GROUP_CLASS_STANDARD,
		}
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
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	limit := int(req.Msg.Limit)
	if limit <= 0 {
		limit = 50
	}

	streams, nextToken, err := store.ListLogStreams(req.Msg.Loggroupname, req.Msg.Logstreamnameprefix, req.Msg.Nexttoken, limit)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
