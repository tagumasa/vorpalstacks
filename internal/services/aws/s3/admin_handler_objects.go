package s3

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"
	pb "vorpalstacks/internal/pb/aws/s3"
)

// ListObjectsV2 retrieves objects in a bucket for the admin console.
func (h *AdminHandler) ListObjectsV2(ctx context.Context, req *connect.Request[pb.ListObjectsV2Request]) (*connect.Response[pb.ListObjectsV2Output], error) {
	if err := requireBucket(req.Msg.Bucket); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	objectStore := h.getObjectStoreFromHeaders(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	input := pbToListObjectsInput(req.Msg)
	result, err := h.service.listObjectsCore(objectStore, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	maxKeys := input.MaxKeys
	if maxKeys <= 0 {
		maxKeys = 1000
	}
	if maxKeys > 1000 {
		maxKeys = 1000
	}

	return connect.NewResponse(listObjectsResultToPb(result, input, maxKeys)), nil
}

// HeadObject retrieves metadata for an object without returning the body.
func (h *AdminHandler) HeadObject(ctx context.Context, req *connect.Request[pb.HeadObjectRequest]) (*connect.Response[pb.HeadObjectOutput], error) {
	if err := requireBucket(req.Msg.Bucket); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if err := requireKey(req.Msg.Key); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	objectStore := h.getObjectStoreFromHeaders(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	input := pbToHeadObjectInput(req.Msg)
	result, err := h.service.headObjectCore(ctx, objectStore, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(headObjectResultToPb(result)), nil
}

// GetObject retrieves an object body for the admin console.
func (h *AdminHandler) GetObject(ctx context.Context, req *connect.Request[pb.GetObjectRequest]) (*connect.Response[pb.GetObjectOutput], error) {
	if err := requireBucket(req.Msg.Bucket); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if err := requireKey(req.Msg.Key); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	objectStore := h.getObjectStoreFromHeaders(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	input := pbToGetObjectInput(req.Msg)
	result, err := h.service.getObjectCore(ctx, objectStore, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(getObjectResultToPb(result)), nil
}
