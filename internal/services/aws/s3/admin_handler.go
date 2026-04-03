package s3

import (
	"context"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/s3"
	s3connect "vorpalstacks/internal/pb/aws/s3/s3connect"
	svccommon "vorpalstacks/internal/services/aws/common"
	s3store "vorpalstacks/internal/store/aws/s3"
)

type AdminHandler struct {
	s3connect.UnimplementedS3ServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ s3connect.S3ServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getBucketStore(req *connect.Request[pb.ListBucketsRequest]) (*s3store.BucketStore, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return s3store.NewBucketStore(regionStorage, h.accountId, region), nil
}

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
