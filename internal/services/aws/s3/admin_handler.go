package s3

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/s3"
	s3connect "vorpalstacks/internal/pb/aws/s3/s3connect"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// AdminHandler provides S3 service administration functionality.
// It implements the S3ServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	s3connect.UnimplementedS3ServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ s3connect.S3ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new S3 AdminHandler.
// It initialises the handler with the provided storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

// getBucketStore returns the S3 bucket store for the specified region.
// It extracts the region from the request headers and creates a new bucket store.
func (h *AdminHandler) getBucketStore(req *connect.Request[pb.ListBucketsRequest]) (*s3store.BucketStore, error) {
	region := req.Header().Get("X-Aws-Region")
	if region == "" {
		region = "us-east-1"
	}
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return s3store.NewBucketStore(regionStorage, h.accountId, region), nil
}

// ListBuckets lists all S3 buckets in the region.
// It returns the bucket names and their creation dates.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListBuckets(ctx context.Context, req *connect.Request[pb.ListBucketsRequest]) (*connect.Response[pb.ListBucketsOutput], error) {
	bucketStore, err := h.getBucketStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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

// CreateBucket creates a new S3 bucket in the specified region.
// The bucket name must be globally unique across all AWS accounts.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateBucket(ctx context.Context, req *connect.Request[pb.CreateBucketRequest]) (*connect.Response[pb.CreateBucketOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteBucket deletes an S3 bucket and all its objects.
// The bucket must be empty before it can be deleted.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteBucket(ctx context.Context, req *connect.Request[pb.DeleteBucketRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetObject retrieves an object from an S3 bucket.
// This returns the object's data and metadata.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetObject(ctx context.Context, req *connect.Request[pb.GetObjectRequest]) (*connect.Response[pb.GetObjectOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutObject uploads an object to an S3 bucket.
// This stores the provided data as an object in the specified bucket.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutObject(ctx context.Context, req *connect.Request[pb.PutObjectRequest]) (*connect.Response[pb.PutObjectOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteObject deletes an object from an S3 bucket.
// This permanently removes the object from the bucket.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteObject(ctx context.Context, req *connect.Request[pb.DeleteObjectRequest]) (*connect.Response[pb.DeleteObjectOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// HeadBucket retrieves metadata about an S3 bucket without returning the object keys.
// This includes information about bucket versioning, encryption, and permissions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) HeadBucket(ctx context.Context, req *connect.Request[pb.HeadBucketRequest]) (*connect.Response[pb.HeadBucketOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// HeadObject retrieves metadata about an S3 object without returning the object data.
// This includes information about object size, content type, and custom metadata.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) HeadObject(ctx context.Context, req *connect.Request[pb.HeadObjectRequest]) (*connect.Response[pb.HeadObjectOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListObjects lists the objects in an S3 bucket using the legacy API.
// This returns object keys and metadata for up to 1000 objects.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListObjects(ctx context.Context, req *connect.Request[pb.ListObjectsRequest]) (*connect.Response[pb.ListObjectsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListObjectsV2 lists the objects in an S3 bucket using the current API version.
// This returns object keys and metadata with improved performance and pagination.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListObjectsV2(ctx context.Context, req *connect.Request[pb.ListObjectsV2Request]) (*connect.Response[pb.ListObjectsV2Output], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
