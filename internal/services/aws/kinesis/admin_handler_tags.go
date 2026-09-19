package kinesis

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	svcerrors "vorpalstacks/internal/common/errors"

	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/kinesis"
)

// TagResource adds or overwrites tags on the resource an ARN addresses,
// via the admin console gRPC-Web interface. Thin adapters over the same
// ARN-addressed tag cores the HTTP tag framework serves — one validation
// and persistence path on both planes.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceInput]) (*connect.Response[pbcommon.Empty], error) {
	reqCtx, err := h.reqCtx(ctx, req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.tagResourceCore(reqCtx, req.Msg.GetResourcearn(), toTags(req.Msg.GetTags())); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// UntagResource removes tag keys from the resource an ARN addresses, via
// the admin console gRPC-Web interface.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceInput]) (*connect.Response[pbcommon.Empty], error) {
	reqCtx, err := h.reqCtx(ctx, req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.untagResourceCore(reqCtx, req.Msg.GetResourcearn(), req.Msg.GetTagkeys()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// ListTagsForResource lists the tags of the resource an ARN addresses,
// via the admin console gRPC-Web interface.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceInput]) (*connect.Response[pb.ListTagsForResourceOutput], error) {
	reqCtx, err := h.reqCtx(ctx, req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	tagList, err := h.service.listResourceTagsCore(reqCtx, req.Msg.GetResourcearn())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	out := &pb.ListTagsForResourceOutput{Tags: make([]*pb.Tag, len(tagList))}
	for i, t := range tagList {
		out.Tags[i] = &pb.Tag{Key: t.Key, Value: proto.String(t.Value)}
	}
	return connect.NewResponse(out), nil
}
