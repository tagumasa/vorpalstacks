package apigateway

import (
	"context"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"

	tagutil "vorpalstacks/internal/common/tags"
	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
)

// TagResource adds tags to an API Gateway resource.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.tagResourceCore(stores, req.Msg.Resourcearn, req.Msg.Tags); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// UntagResource removes tags from an API Gateway resource.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.untagResourceCore(stores, req.Msg.Resourcearn, req.Msg.Tagkeys); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// GetTags lists tags for an API Gateway resource.
func (h *AdminHandler) GetTags(ctx context.Context, req *connect.Request[pb.GetTagsRequest]) (*connect.Response[pb.Tags], error) {
	stores, err := h.getStores(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	tagsList, err := h.service.getResourceTagsCore(stores, req.Msg.Resourcearn)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.Tags{
		Tags: tagutil.ToMap(tagsList),
	}), nil
}
