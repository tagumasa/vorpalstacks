package s3

import (
	"context"
	"fmt"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/s3"
	s3connect "vorpalstacks/internal/pb/aws/s3/s3connect"
)

// AdminHandler implements the S3 admin console gRPC-Web handler.
type AdminHandler struct {
	s3connect.UnimplementedS3ServiceHandler
	s3Store   S3StoreProvider
	accountId string
}

var _ s3connect.S3ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new S3 admin handler backed by the shared S3 store.
func NewAdminHandler(s3Store S3StoreProvider, accountId string) *AdminHandler {
	return &AdminHandler{
		s3Store:   s3Store,
		accountId: accountId,
	}
}

// ListBuckets retrieves all S3 buckets from the regional store.
func (h *AdminHandler) ListBuckets(ctx context.Context, req *connect.Request[pb.ListBucketsRequest]) (*connect.Response[pb.ListBucketsOutput], error) {
	bucketStore := h.getBucketStoreFromHeaders(req.Header())
	if bucketStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	buckets, err := bucketStore.List()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var bucketInfos []*pb.Bucket
	for _, b := range buckets {
		bucketInfos = append(bucketInfos, &pb.Bucket{
			Name:         b.Name,
			Creationdate: b.CreationDate.Format(timeutils.ISO8601UTCFormat),
		})
	}

	return connect.NewResponse(&pb.ListBucketsOutput{
		Buckets: bucketInfos,
	}), nil
}

// CreateBucket creates a new S3 bucket via the admin console.
func (h *AdminHandler) CreateBucket(ctx context.Context, req *connect.Request[pb.CreateBucketRequest]) (*connect.Response[pb.CreateBucketOutput], error) {
	if req.Msg.Bucket == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("bucket name is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	bucketStore := h.s3Store.Buckets(region)
	if bucketStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	bucket, err := bucketStore.Create(req.Msg.Bucket, region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateBucketOutput{
		Location: bucket.Name,
	}), nil
}

// DeleteBucket deletes an S3 bucket via the admin console.
func (h *AdminHandler) DeleteBucket(ctx context.Context, req *connect.Request[pb.DeleteBucketRequest]) (*connect.Response[common.Empty], error) {
	if req.Msg.Bucket == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("bucket name is required"))
	}

	bucketStore := h.getBucketStoreFromHeaders(req.Header())
	if bucketStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	if err := bucketStore.Delete(req.Msg.Bucket); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&common.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the S3 admin console.
func NewConnectHandler(s3Store S3StoreProvider, accountID string) (string, http.Handler) {
	return s3connect.NewS3ServiceHandler(NewAdminHandler(s3Store, accountID))
}
