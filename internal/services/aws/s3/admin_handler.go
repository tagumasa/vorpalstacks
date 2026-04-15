package s3

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/s3"
	s3connect "vorpalstacks/internal/pb/aws/s3/s3connect"
	s3store "vorpalstacks/internal/store/aws/s3"
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

func (h *AdminHandler) getBucketStore(req *connect.Request[pb.ListBucketsRequest]) s3store.BucketStoreInterface {
	region := svccommon.GetRegionFromHeader(req.Header())
	return h.s3Store.Buckets(region)
}

// ListBuckets retrieves all S3 buckets from the regional store.
func (h *AdminHandler) ListBuckets(ctx context.Context, req *connect.Request[pb.ListBucketsRequest]) (*connect.Response[pb.ListBucketsOutput], error) {
	bucketStore := h.getBucketStore(req)
	if bucketStore == nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("storage unavailable"))
	}

	buckets, err := bucketStore.List()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var bucketInfos []*pb.Bucket
	for _, b := range buckets {
		bucketInfos = append(bucketInfos, &pb.Bucket{
			Name:         b.Name,
			Creationdate: b.CreationDate.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	return connect.NewResponse(&pb.ListBucketsOutput{
		Buckets: bucketInfos,
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the S3 admin console.
func NewConnectHandler(s3Store S3StoreProvider, accountID string) (string, http.Handler) {
	return s3connect.NewS3ServiceHandler(NewAdminHandler(s3Store, accountID))
}
