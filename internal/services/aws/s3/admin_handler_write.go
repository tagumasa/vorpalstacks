package s3

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"
	pb "vorpalstacks/internal/pb/aws/s3"
)

// DeleteObject removes an object from a bucket via the admin console.
func (h *AdminHandler) DeleteObject(ctx context.Context, req *connect.Request[pb.DeleteObjectRequest]) (*connect.Response[pb.DeleteObjectOutput], error) {
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

	input := pbToDeleteObjectInput(req.Msg)
	result, err := h.service.deleteObjectCore(ctx, objectStore, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(deleteObjectResultToPb(result)), nil
}

// PutObject uploads an object to a bucket via the admin console.
func (h *AdminHandler) PutObject(ctx context.Context, req *connect.Request[pb.PutObjectRequest]) (*connect.Response[pb.PutObjectOutput], error) {
	if err := requireBucket(req.Msg.Bucket); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if err := requireKey(req.Msg.Key); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	bucketStore, objectStore := h.getStoresFromHeaders(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	input := pbToPutObjectInput(req.Msg)
	result, err := h.service.putObjectCore(ctx, bucketStore, objectStore, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(putObjectResultToPb(result)), nil
}

// DeleteObjects removes multiple objects from a bucket via the admin console.
func (h *AdminHandler) DeleteObjects(ctx context.Context, req *connect.Request[pb.DeleteObjectsRequest]) (*connect.Response[pb.DeleteObjectsOutput], error) {
	if err := requireBucket(req.Msg.Bucket); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if req.Msg.Delete == nil || len(req.Msg.Delete.Objects) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no objects specified for deletion"))
	}

	objectStore := h.getObjectStoreFromHeaders(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	input := pbToDeleteObjectsInput(req.Msg)
	result, err := h.service.deleteObjectsCore(ctx, objectStore, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(deleteObjectsResultToPb(result)), nil
}

// CopyObject copies an object to another location via the admin console.
func (h *AdminHandler) CopyObject(ctx context.Context, req *connect.Request[pb.CopyObjectRequest]) (*connect.Response[pb.CopyObjectOutput], error) {
	if err := requireBucket(req.Msg.Bucket); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if err := requireKey(req.Msg.Key); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if req.Msg.Copysource == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("copy source is required"))
	}

	bucketStore, objectStore := h.getStoresFromHeaders(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	input := pbToCopyObjectInput(req.Msg)
	result, err := h.service.copyObjectCore(ctx, bucketStore, objectStore, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(copyObjectResultToPb(result)), nil
}
