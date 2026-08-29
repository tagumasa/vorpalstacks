package kinesis

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	svcerrors "vorpalstacks/internal/common/errors"

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

// NewAdminHandler creates a new Kinesis admin handler backed by the given
// service instance, ensuring the same per-region cached stores are used as
// the HTTP API handlers.
func NewAdminHandler(svc *KinesisService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListStreams returns a list of Kinesis streams via the admin console
// gRPC-Web interface.
func (h *AdminHandler) ListStreams(ctx context.Context, req *connect.Request[pb.ListStreamsInput]) (*connect.Response[pb.ListStreamsOutput], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.listStreamsCore(stores, ListStreamsInput{
		ExclusiveStartStreamName: req.Msg.Exclusivestartstreamname,
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
		Nexttoken:       result.NextMarker,
	}), nil
}

// DescribeStream returns detailed information about a Kinesis stream via
// the admin console gRPC-Web interface.
func (h *AdminHandler) DescribeStream(ctx context.Context, req *connect.Request[pb.DescribeStreamInput]) (*connect.Response[pb.DescribeStreamOutput], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.describeStreamCore(stores, DescribeStreamInput{
		StreamName: req.Msg.Streamname,
		StreamARN:  req.Msg.Streamarn,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DescribeStreamOutput{
		Streamdescription: toPbStreamDescription(result.Stream, result.Shards),
	}), nil
}

// CreateStream creates a new Kinesis stream via the admin console.
func (h *AdminHandler) CreateStream(ctx context.Context, req *connect.Request[pb.CreateStreamInput]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if _, err := h.service.createStreamCore(stores, CreateStreamInput{
		StreamName: req.Msg.GetStreamname(),
		ShardCount: req.Msg.GetShardcount(),
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DeleteStream deletes a Kinesis stream via the admin console.
func (h *AdminHandler) DeleteStream(ctx context.Context, req *connect.Request[pb.DeleteStreamInput]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.deleteStreamCore(stores, DeleteStreamInput{
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
