package s3

import (
	"bytes"
	"context"
	"fmt"
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

	bucketStore, objectStore := h.getStores(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	contentLength := int64(len(req.Msg.Body))
	if contentLength > maxSingleUploadSize {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("object size %d exceeds maximum allowed size %d", contentLength, maxSingleUploadSize))
	}

	metadata := req.Msg.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	contentType := req.Msg.Contenttype
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	var bucketEncryption *s3store.EncryptionConfig
	if bucketStore != nil {
		bucketEncryption, _ = bucketStore.GetEncryptionConfiguration(req.Msg.Bucket)
	}

	encType := EncryptionTypeNone
	if h.encryptionManager != nil {
		encType = h.encryptionManager.DetermineEncryptionType(EncryptionTypeNone, bucketEncryption)
	}

	var obj *s3store.Object
	var err error

	if h.encryptionManager != nil && h.encryptionManager.ShouldEncrypt(encType, bucketEncryption) {
		encResult, err := h.encryptionManager.Encrypt(req.Msg.Body, encType, bucketEncryption, req.Msg.Bucket, req.Msg.Key, "")
		if err != nil {
			return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("encryption failed: %w", err))
		}
		sseMeta := &s3store.SSEObjectMetadata{
			EncryptionType:   s3store.SSEType(encResult.EncryptionType),
			EncryptedDataKey: encResult.EncryptedDataKey,
			ContentNonce:     encResult.ContentNonce,
			KMSKeyID:         encResult.KMSKeyID,
			UnencryptedMD5:   encResult.UnencryptedMD5,
			UnencryptedSize:  encResult.UnencryptedSize,
		}
		obj, err = objectStore.PutEncrypted(ctx, req.Msg.Bucket, req.Msg.Key, encResult.EncryptedData, contentType, metadata, sseMeta, s3store.StorageClassStandard, nil)
	} else {
		reader := bytes.NewReader(req.Msg.Body)
		obj, err = objectStore.Put(ctx, req.Msg.Bucket, req.Msg.Key, reader, contentType, metadata)
	}
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
