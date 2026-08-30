package apigateway

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"

	tagutil "vorpalstacks/internal/common/tags"
	pb "vorpalstacks/internal/pb/aws/apigateway"
	pbcommon "vorpalstacks/internal/pb/aws/common"
)

// TagResource adds tags to an API Gateway resource.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[pbcommon.Empty], error) {
	if req.Msg.Resourcearn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("resource_arn is required"))
	}
	if len(req.Msg.Tags) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("tags is required"))
	}

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
	if req.Msg.Resourcearn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("resource_arn is required"))
	}
	if len(req.Msg.Tagkeys) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("tag_keys is required"))
	}

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
	if req.Msg.Resourcearn == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("resource_arn is required"))
	}

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
