package kinesis

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"

	"vorpalstacks/internal/common/defaults"

	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/kinesis"
	kinesisconnect "vorpalstacks/internal/pb/aws/kinesis/kinesisconnect"
)

// AdminHandler implements the Kinesis admin console gRPC-Web handler.
// It is a thin adapter that delegates all operations to service-layer
// Core methods and converts results to proto types via the helpers in
// admin_handler_convert.go. This file does not import any store package.
type AdminHandler struct {
	kinesisconnect.UnimplementedKinesisServiceHandler
	service *KinesisService
}

var _ kinesisconnect.KinesisServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new admin handler backed by the given
// service instance, ensuring the same per-region cached stores are used as
// the HTTP API handlers.
func NewAdminHandler(svc *KinesisService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// reqCtx builds the request context for the admin plane: the console's
// requests carry the region in headers rather than in the caller context,
// and the context resolves the same per-region cached stores as the HTTP
// API plane.
func (h *AdminHandler) reqCtx(ctx context.Context, headers http.Header) (*request.RequestContext, error) {
	region := defaults.GetRegionFromHeader(headers)
	if h.service.storageManager == nil {
		return nil, fmt.Errorf("kinesis storage manager not initialised for region %s", region)
	}
	return request.NewRequestContext(ctx, h.service.storageManager, h.service.accountID, region), nil
}

// ListStreams returns a list of Kinesis streams via the admin console
// gRPC-Web interface.
func (h *AdminHandler) ListStreams(ctx context.Context, req *connect.Request[pb.ListStreamsInput]) (*connect.Response[pb.ListStreamsOutput], error) {
	reqCtx, err := h.reqCtx(ctx, req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.listStreamsCore(reqCtx, ListStreamsInput{
		ExclusiveStartStreamName: req.Msg.GetExclusivestartstreamname(),
		Limit:                    int(req.Msg.GetLimit()),
		// The console's zero limit means "no limit chosen", matching the
		// proto default for an absent member.
		HasLimit: req.Msg.GetLimit() > 0,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	streamNames := make([]string, len(result.Streams))
	summaries := make([]*pb.StreamSummary, len(result.Streams))
	for i, s := range result.Streams {
		streamNames[i] = s.StreamName
		summaries[i] = toPbStreamSummary(s)
	}

	return connect.NewResponse(&pb.ListStreamsOutput{
		Streamnames:     streamNames,
		Streamsummaries: summaries,
		Hasmorestreams:  proto.Bool(result.IsTruncated),
		Nexttoken:       proto.String(result.NextMarker),
	}), nil
}

// DescribeStream returns detailed information about a Kinesis stream via
// the admin console gRPC-Web interface. The shard page travels with the
// request: Limit and ExclusiveStartShardId page the core's documented
// hundred-shard window, and the response's HasMoreShards is the core's
// truth — a console client can walk a stream of any shard count.
func (h *AdminHandler) DescribeStream(ctx context.Context, req *connect.Request[pb.DescribeStreamInput]) (*connect.Response[pb.DescribeStreamOutput], error) {
	reqCtx, err := h.reqCtx(ctx, req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	limit := req.Msg.GetLimit()
	result, err := h.service.describeStreamCore(reqCtx, DescribeStreamInput{
		StreamName:            req.Msg.GetStreamname(),
		StreamARN:             req.Msg.GetStreamarn(),
		Limit:                 int(limit),
		HasLimit:              limit > 0, // zero means no limit chosen, the proto default for an absent member
		ExclusiveStartShardId: req.Msg.GetExclusivestartshardid(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DescribeStreamOutput{
		Streamdescription: toPbStreamDescription(result.Stream, result.Shards, result.HasMoreShards),
	}), nil
}

// CreateStream creates a new Kinesis stream via the admin console. The
// console's mode selection travels with the request — an on-demand
// creation must not silently become a provisioned one — and create-time
// tags ride the same locked creation the HTTP plane uses.
func (h *AdminHandler) CreateStream(ctx context.Context, req *connect.Request[pb.CreateStreamInput]) (*connect.Response[pbcommon.Empty], error) {
	reqCtx, err := h.reqCtx(ctx, req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	modeDetails := req.Msg.GetStreammodedetails()
	if _, err := h.service.createStreamCore(reqCtx, CreateStreamInput{
		StreamName:           req.Msg.GetStreamname(),
		ShardCount:           req.Msg.GetShardcount(),
		HasShardCount:        req.Msg.Shardcount != nil,
		StreamMode:           fromPbStreamMode(modeDetails.GetStreammode()),
		HasStreamModeDetails: modeDetails != nil,
		Tags:                 toTags(req.Msg.GetTags()),
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DeleteStream deletes a Kinesis stream via the admin console.
func (h *AdminHandler) DeleteStream(ctx context.Context, req *connect.Request[pb.DeleteStreamInput]) (*connect.Response[pbcommon.Empty], error) {
	reqCtx, err := h.reqCtx(ctx, req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.deleteStreamCore(reqCtx, DeleteStreamInput{
		StreamName: req.Msg.GetStreamname(),
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Kinesis
// admin console.
func NewConnectHandler(svc *KinesisService) (string, http.Handler) {
	return kinesisconnect.NewKinesisServiceHandler(NewAdminHandler(svc))
}
