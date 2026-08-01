package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

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
			Deletemarker: proto.Bool(result.IsDeleteMarker),
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
		encResult, encErr := h.encryptionManager.Encrypt(req.Msg.Body, encType, bucketEncryption, req.Msg.Bucket, req.Msg.Key, "")
		if encErr != nil {
			return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("encryption failed: %w", encErr))
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
		Size:      proto.Int64(obj.Size),
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

// DeleteObjects removes multiple objects from a bucket via the admin console.
func (h *AdminHandler) DeleteObjects(ctx context.Context, req *connect.Request[pb.DeleteObjectsRequest]) (*connect.Response[pb.DeleteObjectsOutput], error) {
	if req.Msg.Bucket == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("bucket is required"))
	}
	if req.Msg.Delete == nil || len(req.Msg.Delete.Objects) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no objects specified for deletion"))
	}

	_, objectStore := h.getStores(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	var deleted []*pb.DeletedObject
	var errors []*pb.Error

	for _, obj := range req.Msg.Delete.Objects {
		if obj.Key == "" {
			errors = append(errors, &pb.Error{
				Key:     obj.Key,
				Code:    "InvalidArgument",
				Message: "object key is required",
			})
			continue
		}

		if obj.Versionid != "" {
			result, err := objectStore.DeleteWithVersion(ctx, req.Msg.Bucket, obj.Key, obj.Versionid)
			if err != nil {
				errors = append(errors, &pb.Error{
					Key:     obj.Key,
					Code:    "InternalError",
					Message: err.Error(),
				})
				continue
			}
			deletedObj := &pb.DeletedObject{
				Key:       obj.Key,
				Versionid: obj.Versionid,
			}
			if result != nil {
				deletedObj.Deletemarker = proto.Bool(true)
				deletedObj.Deletemarkerversionid = result.VersionID
			}
			deleted = append(deleted, deletedObj)
		} else {
			err := objectStore.Delete(ctx, req.Msg.Bucket, obj.Key)
			if err != nil {
				errors = append(errors, &pb.Error{
					Key:     obj.Key,
					Code:    "InternalError",
					Message: err.Error(),
				})
				continue
			}
			deleted = append(deleted, &pb.DeletedObject{
				Key: obj.Key,
			})
		}
	}

	return connect.NewResponse(&pb.DeleteObjectsOutput{
		Deleted: deleted,
		Errors:  errors,
	}), nil
}

// CopyObject copies an object to another location via the admin console.
func (h *AdminHandler) CopyObject(ctx context.Context, req *connect.Request[pb.CopyObjectRequest]) (*connect.Response[pb.CopyObjectOutput], error) {
	if req.Msg.Bucket == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("destination bucket is required"))
	}
	if req.Msg.Key == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("destination key is required"))
	}
	if req.Msg.Copysource == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("copy source is required"))
	}

	bucketStore, objectStore := h.getStores(req.Header())
	if objectStore == nil {
		return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("storage unavailable"))
	}

	srcBucket, srcKey, _, err := parseCopySource(req.Msg.Copysource)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid copy source: %w", err))
	}

	// Verify source object exists and retrieve metadata.
	srcObj, err := objectStore.GetMetadata(srcBucket, srcKey)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("source object not found: %s/%s", srcBucket, srcKey))
	}

	if srcObj.Size > maxCopyObjectSize {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("source object size %d exceeds maximum copy size %d", srcObj.Size, maxCopyObjectSize))
	}

	// Determine content type for the copy.
	contentType := req.Msg.Contenttype
	if contentType == "" {
		contentType = srcObj.ContentType
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Determine metadata for the copy.
	metadata := srcObj.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	var obj *s3store.Object

	// Handle encryption for the destination bucket.
	var bucketEncryption *s3store.EncryptionConfig
	if bucketStore != nil {
		bucketEncryption, _ = bucketStore.GetEncryptionConfiguration(req.Msg.Bucket)
	}

	encType := EncryptionTypeNone
	if h.encryptionManager != nil {
		encType = h.encryptionManager.DetermineEncryptionType(EncryptionTypeNone, bucketEncryption)
	}

	if h.encryptionManager != nil && h.encryptionManager.ShouldEncrypt(encType, bucketEncryption) {
		// Read source data, encrypt, and write to destination.
		reader, _, readErr := objectStore.Get(ctx, srcBucket, srcKey)
		if readErr != nil {
			return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("failed to read source object: %w", readErr))
		}
		srcData, readErr := io.ReadAll(reader)
		reader.Close()
		if readErr != nil {
			return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("failed to read source data: %w", readErr))
		}

		encResult, encErr := h.encryptionManager.Encrypt(srcData, encType, bucketEncryption, req.Msg.Bucket, req.Msg.Key, "")
		if encErr != nil {
			return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("encryption failed: %w", encErr))
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
		// Use store-level Copy for the unencrypted case.
		obj, err = objectStore.Copy(ctx, srcBucket, srcKey, req.Msg.Bucket, req.Msg.Key)
	}
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CopyObjectOutput{
		Copyobjectresult: &pb.CopyObjectResult{
			Etag:         formatETag(obj.ETag),
			Lastmodified: obj.LastModified.Format(timeutils.ISO8601UTCFormat),
		},
	}), nil
}
