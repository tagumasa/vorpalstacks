package s3

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/s3"
	s3connect "vorpalstacks/internal/pb/aws/s3/s3connect"
)

// AdminHandler implements the S3 admin console gRPC-Web handler.
// It delegates all business logic to S3Service Core methods and uses
// admin_handler_convert.go for proto<->DTO conversion.
type AdminHandler struct {
	s3connect.UnimplementedS3ServiceHandler
	service *S3Service
}

var _ s3connect.S3ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new S3 admin handler backed by the given service.
func NewAdminHandler(svc *S3Service) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListBuckets retrieves all S3 buckets from the regional store.
func (h *AdminHandler) ListBuckets(ctx context.Context, req *connect.Request[pb.ListBucketsRequest]) (*connect.Response[pb.ListBucketsOutput], error) {
	bucketStore := h.getBucketStoreFromHeaders(req.Header())
	if bucketStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	result, err := h.service.listBucketsCore(bucketStore, AdminListBucketsInput{})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(listBucketsResultToPb(result)), nil
}

// CreateBucket creates a new S3 bucket via the admin console.
func (h *AdminHandler) CreateBucket(ctx context.Context, req *connect.Request[pb.CreateBucketRequest]) (*connect.Response[pb.CreateBucketOutput], error) {
	if err := requireBucket(req.Msg.Bucket); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	bucketStore := h.getBucketStoreFromHeaders(req.Header())
	if bucketStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	input := pbToCreateBucketInput(req.Msg, region)
	result, err := h.service.createBucketCore(bucketStore, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(createBucketResultToPb(result)), nil
}

// DeleteBucket deletes an S3 bucket via the admin console.
func (h *AdminHandler) DeleteBucket(ctx context.Context, req *connect.Request[pb.DeleteBucketRequest]) (*connect.Response[pbcommon.Empty], error) {
	if err := requireBucket(req.Msg.Bucket); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	bucketStore, objectStore := h.getStoresFromHeaders(req.Header())
	if bucketStore == nil || objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	input := pbToDeleteBucketInput(req.Msg)
	_, err := h.service.deleteBucketCore(bucketStore, objectStore, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the S3 admin console.
func NewConnectHandler(svc *S3Service) (string, http.Handler) {
	return s3connect.NewS3ServiceHandler(NewAdminHandler(svc))
}
