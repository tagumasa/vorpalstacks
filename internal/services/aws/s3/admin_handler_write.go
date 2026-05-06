package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/s3"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// DeleteObject removes an object from a bucket via the admin console.
func (h *AdminHandler) DeleteObject(ctx context.Context, req *connect.Request[pb.DeleteObjectRequest]) (*connect.Response[pb.DeleteObjectOutput], error) {
	if req.Msg.Bucket == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("bucket is required"))
	}
	if req.Msg.Key == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("key is required"))
	}

	_, objectStore := h.getStores(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	if req.Msg.Versionid != "" {
		result, err := objectStore.DeleteWithVersion(ctx, req.Msg.Bucket, req.Msg.Key, req.Msg.Versionid)
		if err != nil {
			return nil, svcerrors.StoreErrorToGRPC(err)
		}
		output := &pb.DeleteObjectOutput{
			Versionid:    result.VersionID,
			Deletemarker: result.IsDeleteMarker,
		}
		return connect.NewResponse(output), nil
	}

	err := objectStore.Delete(ctx, req.Msg.Bucket, req.Msg.Key)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteObjectOutput{}), nil
}

// PutObject uploads an object to a bucket via the admin console.
func (h *AdminHandler) PutObject(ctx context.Context, req *connect.Request[pb.PutObjectRequest]) (*connect.Response[pb.PutObjectOutput], error) {
	if req.Msg.Bucket == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("bucket is required"))
	}
	if req.Msg.Key == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("key is required"))
	}

	_, objectStore := h.getStores(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	contentLength := int64(len(req.Msg.Body))
	if contentLength > maxSingleUploadSize {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("object size %d exceeds maximum allowed size %d", contentLength, maxSingleUploadSize))
	}

	var reader io.Reader = bytes.NewReader(req.Msg.Body)

	metadata := req.Msg.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	contentType := req.Msg.Contenttype
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	obj, err := objectStore.Put(ctx, req.Msg.Bucket, req.Msg.Key, reader, contentType, metadata)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	output := &pb.PutObjectOutput{
		Etag:      formatETag(obj.ETag),
		Versionid: obj.VersionID,
		Size:      obj.Size,
	}
	if obj.SSEMetadata != nil {
		output.Ssekmskeyid = obj.SSEMetadata.KMSKeyID
	}

	return connect.NewResponse(output), nil
}

func (h *AdminHandler) getBucketStoreFromHeaders(headers http.Header) s3store.BucketStoreInterface {
	region := svccommon.GetRegionFromHeader(headers)
	return h.s3Store.Buckets(region)
}
