package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/logs"
	storecommon "vorpalstacks/internal/store/aws/common"
	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/timeutils"
)

// ---------------------------------------------------------------------------
// Input structs — transport-agnostic DTOs for the admin gRPC-Web handler.
// Prefixed with Admin to avoid collisions with the HTTP API input types.
// ---------------------------------------------------------------------------

// AdminCreateBucketInput carries the fields needed for CreateBucket.
type AdminCreateBucketInput struct {
	Bucket string
	Region string
}

// AdminDeleteBucketInput carries the fields needed for DeleteBucket.
type AdminDeleteBucketInput struct {
	Bucket string
}

// AdminListObjectsInput carries the fields needed for ListObjectsV2.
type AdminListObjectsInput struct {
	Bucket    string
	Prefix    string
	Delimiter string
	Marker    string
	MaxKeys   int
}

// AdminHeadObjectInput carries the fields needed for HeadObject.
type AdminHeadObjectInput struct {
	Bucket    string
	Key       string
	VersionID string
}

// AdminGetObjectInput carries the fields needed for GetObject.
type AdminGetObjectInput struct {
	Bucket    string
	Key       string
	VersionID string
}

// AdminPutObjectInput carries the fields needed for PutObject.
type AdminPutObjectInput struct {
	Bucket      string
	Key         string
	Body        []byte
	ContentType string
	Metadata    map[string]string
}

// AdminDeleteObjectInput carries the fields needed for DeleteObject.
type AdminDeleteObjectInput struct {
	Bucket    string
	Key       string
	VersionID string
}

// AdminObjectIdentifier identifies a single object for bulk delete.
type AdminObjectIdentifier struct {
	Key       string
	VersionID string
}

// AdminDeleteObjectsInput carries the fields needed for DeleteObjects.
type AdminDeleteObjectsInput struct {
	Bucket  string
	Objects []AdminObjectIdentifier
}

// AdminCopyObjectInput carries the fields needed for CopyObject.
type AdminCopyObjectInput struct {
	Bucket      string
	Key         string
	CopySource  string
	ContentType string
}

// ---------------------------------------------------------------------------
// Result structs
// ---------------------------------------------------------------------------

// AdminBucketInfo holds the transport-agnostic representation of a bucket.
type AdminBucketInfo struct {
	Name         string
	Region       string
	CreationDate string
}

// AdminListBucketsResult holds the transport-agnostic result of ListBuckets.
type AdminListBucketsResult struct {
	Buckets []AdminBucketInfo
}

// AdminCreateBucketResult holds the transport-agnostic result of CreateBucket.
type AdminCreateBucketResult struct {
	Location string
}

// AdminDeleteBucketResult holds the transport-agnostic result of DeleteBucket.
type AdminDeleteBucketResult struct{}

// AdminListObjectsResult holds the transport-agnostic result of ListObjectsV2.
type AdminListObjectsResult struct {
	Objects        []*s3store.Object
	CommonPrefixes []string
	IsTruncated    bool
	NextMarker     string
}

// AdminHeadObjectResult holds the transport-agnostic result of HeadObject.
type AdminHeadObjectResult struct {
	Object *s3store.Object
}

// AdminGetObjectResult holds the transport-agnostic result of GetObject.
type AdminGetObjectResult struct {
	Object *s3store.Object
	Body   []byte
}

// AdminPutObjectResult holds the transport-agnostic result of PutObject.
type AdminPutObjectResult struct {
	ETag      string
	VersionID string
	Size      int64
	KMSKeyID  string
}

// AdminDeleteObjectResult holds the transport-agnostic result of DeleteObject.
type AdminDeleteObjectResult struct {
	VersionID      string
	IsDeleteMarker bool
}

// AdminDeletedObject holds info about a single successfully deleted object.
type AdminDeletedObject struct {
	Key                   string
	VersionID             string
	DeleteMarker          bool
	DeleteMarkerVersionID string
}

// AdminDeleteError holds info about a single failed deletion.
type AdminDeleteError struct {
	Key     string
	Code    string
	Message string
}

// AdminDeleteObjectsResult holds the transport-agnostic result of DeleteObjects.
type AdminDeleteObjectsResult struct {
	Deleted []AdminDeletedObject
	Errors  []AdminDeleteError
}

// AdminCopyObjectResult holds the transport-agnostic result of CopyObject.
type AdminCopyObjectResult struct {
	ETag         string
	LastModified string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path shared by admin
// handler only.  The HTTP API layer has its own operations in
// bucket_operations.go / object_*.go.
// ---------------------------------------------------------------------------

// listBucketsCore returns all buckets in the regional store.
func (s *S3Service) listBucketsCore(bucketStore s3store.BucketStoreInterface) (*AdminListBucketsResult, error) {
	buckets, err := bucketStore.List()
	if err != nil {
		return nil, err
	}

	result := &AdminListBucketsResult{}
	for _, b := range buckets {
		result.Buckets = append(result.Buckets, AdminBucketInfo{
			Name:         b.Name,
			Region:       b.Region,
			CreationDate: b.CreationDate.Format(timeutils.ISO8601UTCFormat),
		})
	}
	return result, nil
}

// createBucketCore validates the bucket name and creates the bucket in the
// regional store.
func (s *S3Service) createBucketCore(bucketStore s3store.BucketStoreInterface, in AdminCreateBucketInput) (*AdminCreateBucketResult, error) {
	if err := validateBucketName(in.Bucket); err != nil {
		return nil, awserrors.NewAWSError("InvalidBucketName",
			fmt.Sprintf("invalid bucket name: %s", in.Bucket), http.StatusBadRequest)
	}

	bucket, err := bucketStore.Create(in.Bucket, in.Region)
	if err != nil {
		return nil, err
	}
	return &AdminCreateBucketResult{Location: bucket.Name}, nil
}

// deleteBucketCore verifies the bucket is empty (no objects, no incomplete
// multipart uploads) before deleting it.
func (s *S3Service) deleteBucketCore(bucketStore s3store.BucketStoreInterface, objectStore s3store.ObjectStoreInterface, in AdminDeleteBucketInput) (*AdminDeleteBucketResult, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}
	if bucket.ObjectLockEnabled {
		logs.Warn("s3: admin deleting bucket with Object Lock enabled", logs.String("bucket", in.Bucket))
	}

	count, err := objectStore.CountByBucket(in.Bucket)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, awserrors.NewAWSError("BucketNotEmpty",
			fmt.Sprintf("bucket is not empty: contains %d object(s), delete all objects first", count),
			http.StatusConflict)
	}

	multipartCount, err := objectStore.CountMultipartUploadsByBucket(in.Bucket)
	if err != nil {
		return nil, err
	}
	if multipartCount > 0 {
		return nil, awserrors.NewAWSError("BucketNotEmpty",
			fmt.Sprintf("bucket has %d incomplete multipart upload(s)", multipartCount),
			http.StatusConflict)
	}

	if err := bucketStore.Delete(in.Bucket); err != nil {
		return nil, err
	}
	return &AdminDeleteBucketResult{}, nil
}

// listObjectsCore lists objects in a bucket with pagination.
func (s *S3Service) listObjectsCore(objectStore s3store.ObjectStoreInterface, in AdminListObjectsInput) (*AdminListObjectsResult, error) {
	maxKeys := in.MaxKeys
	if maxKeys <= 0 {
		maxKeys = 1000
	}
	if maxKeys > 1000 {
		maxKeys = 1000
	}

	result, err := objectStore.List(in.Bucket, in.Prefix, in.Delimiter, in.Marker, maxKeys)
	if err != nil {
		return nil, err
	}

	return &AdminListObjectsResult{
		Objects:        result.Objects,
		CommonPrefixes: result.CommonPrefixes,
		IsTruncated:    result.IsTruncated,
		NextMarker:     result.NextMarker,
	}, nil
}

// headObjectCore retrieves metadata for an object without returning the body.
func (s *S3Service) headObjectCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in AdminHeadObjectInput) (*AdminHeadObjectResult, error) {
	obj, err := objectStore.HeadWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
	if err != nil {
		return nil, err
	}
	return &AdminHeadObjectResult{Object: obj}, nil
}

// getObjectCore retrieves metadata and body for an object.
func (s *S3Service) getObjectCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in AdminGetObjectInput) (*AdminGetObjectResult, error) {
	obj, err := objectStore.HeadWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
	if err != nil {
		return nil, err
	}

	var reader io.ReadCloser
	if in.VersionID != "" {
		reader, _, err = objectStore.GetWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
	} else {
		reader, _, err = objectStore.Get(ctx, in.Bucket, in.Key)
	}
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &AdminGetObjectResult{Object: obj, Body: data}, nil
}

// putObjectCore validates the upload, determines encryption settings, and
// stores the object.
func (s *S3Service) putObjectCore(ctx context.Context, bucketStore s3store.BucketStoreInterface, objectStore s3store.ObjectStoreInterface, in AdminPutObjectInput) (*AdminPutObjectResult, error) {
	contentLength := int64(len(in.Body))
	if contentLength > maxSingleUploadSize {
		return nil, awserrors.NewAWSError("EntityTooLarge",
			fmt.Sprintf("object size %d exceeds maximum allowed size %d", contentLength, maxSingleUploadSize),
			http.StatusBadRequest)
	}

	metadata := in.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	contentType := in.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	var bucketEncryption *s3store.EncryptionConfig
	if bucketStore != nil {
		bucketEncryption, _ = bucketStore.GetEncryptionConfiguration(in.Bucket)
	}

	encType := EncryptionTypeNone
	if s.encryptionManager != nil {
		encType = s.encryptionManager.DetermineEncryptionType(EncryptionTypeNone, bucketEncryption)
	}

	var obj *s3store.Object
	var err error

	if s.encryptionManager != nil && s.encryptionManager.ShouldEncrypt(encType, bucketEncryption) {
		encResult, encErr := s.encryptionManager.Encrypt(in.Body, encType, bucketEncryption, in.Bucket, in.Key, "")
		if encErr != nil {
			return nil, fmt.Errorf("encryption failed: %w", encErr)
		}
		sseMeta := &s3store.SSEObjectMetadata{
			EncryptionType:   s3store.SSEType(encResult.EncryptionType),
			EncryptedDataKey: encResult.EncryptedDataKey,
			ContentNonce:     encResult.ContentNonce,
			KMSKeyID:         encResult.KMSKeyID,
			UnencryptedMD5:   encResult.UnencryptedMD5,
			UnencryptedSize:  encResult.UnencryptedSize,
		}
		obj, err = objectStore.PutEncrypted(ctx, in.Bucket, in.Key, encResult.EncryptedData, contentType, metadata, sseMeta, s3store.StorageClassStandard, nil)
	} else {
		reader := bytes.NewReader(in.Body)
		obj, err = objectStore.Put(ctx, in.Bucket, in.Key, reader, contentType, metadata)
	}
	if err != nil {
		return nil, err
	}

	result := &AdminPutObjectResult{
		ETag:      formatETag(obj.ETag),
		VersionID: obj.VersionID,
		Size:      obj.Size,
	}
	if obj.SSEMetadata != nil {
		result.KMSKeyID = obj.SSEMetadata.KMSKeyID
	}
	return result, nil
}

// deleteObjectCore deletes a single object, optionally with a specific version.
func (s *S3Service) deleteObjectCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in AdminDeleteObjectInput) (*AdminDeleteObjectResult, error) {
	if in.VersionID != "" {
		result, err := objectStore.DeleteWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
		if err != nil {
			return nil, err
		}
		return &AdminDeleteObjectResult{
			VersionID:      result.VersionID,
			IsDeleteMarker: result.IsDeleteMarker,
		}, nil
	}

	if err := objectStore.Delete(ctx, in.Bucket, in.Key); err != nil {
		return nil, err
	}
	return &AdminDeleteObjectResult{}, nil
}

// deleteObjectsCore deletes multiple objects, collecting per-object results.
func (s *S3Service) deleteObjectsCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in AdminDeleteObjectsInput) (*AdminDeleteObjectsResult, error) {
	result := &AdminDeleteObjectsResult{}

	for _, obj := range in.Objects {
		if obj.Key == "" {
			result.Errors = append(result.Errors, AdminDeleteError{
				Key:     obj.Key,
				Code:    "InvalidArgument",
				Message: "object key is required",
			})
			continue
		}

		if obj.VersionID != "" {
			delResult, err := objectStore.DeleteWithVersion(ctx, in.Bucket, obj.Key, obj.VersionID)
			if err != nil {
				result.Errors = append(result.Errors, AdminDeleteError{
					Key:     obj.Key,
					Code:    "InternalError",
					Message: err.Error(),
				})
				continue
			}
			deletedObj := AdminDeletedObject{
				Key:       obj.Key,
				VersionID: obj.VersionID,
			}
			if delResult != nil {
				deletedObj.DeleteMarker = true
				deletedObj.DeleteMarkerVersionID = delResult.VersionID
			}
			result.Deleted = append(result.Deleted, deletedObj)
		} else {
			if err := objectStore.Delete(ctx, in.Bucket, obj.Key); err != nil {
				result.Errors = append(result.Errors, AdminDeleteError{
					Key:     obj.Key,
					Code:    "InternalError",
					Message: err.Error(),
				})
				continue
			}
			result.Deleted = append(result.Deleted, AdminDeletedObject{
				Key: obj.Key,
			})
		}
	}

	return result, nil
}

// copyObjectCore copies an object, handling encryption for the destination.
func (s *S3Service) copyObjectCore(ctx context.Context, bucketStore s3store.BucketStoreInterface, objectStore s3store.ObjectStoreInterface, in AdminCopyObjectInput) (*AdminCopyObjectResult, error) {
	srcBucket, srcKey, _, err := parseCopySource(in.CopySource)
	if err != nil {
		return nil, awserrors.NewAWSError("InvalidArgument",
			fmt.Sprintf("invalid copy source: %v", err), http.StatusBadRequest)
	}

	srcObj, err := objectStore.GetMetadata(srcBucket, srcKey)
	if err != nil {
		if storecommon.IsNotFound(err) {
			return nil, awserrors.NewAWSError("NoSuchKey",
				fmt.Sprintf("source object not found: %s/%s", srcBucket, srcKey), http.StatusNotFound)
		}
		return nil, err
	}

	if srcObj.Size > maxCopyObjectSize {
		return nil, awserrors.NewAWSError("EntityTooLarge",
			fmt.Sprintf("source object size %d exceeds maximum copy size %d", srcObj.Size, maxCopyObjectSize),
			http.StatusBadRequest)
	}

	contentType := in.ContentType
	if contentType == "" {
		contentType = srcObj.ContentType
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	metadata := srcObj.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	var bucketEncryption *s3store.EncryptionConfig
	if bucketStore != nil {
		bucketEncryption, _ = bucketStore.GetEncryptionConfiguration(in.Bucket)
	}

	encType := EncryptionTypeNone
	if s.encryptionManager != nil {
		encType = s.encryptionManager.DetermineEncryptionType(EncryptionTypeNone, bucketEncryption)
	}

	var obj *s3store.Object

	if s.encryptionManager != nil && s.encryptionManager.ShouldEncrypt(encType, bucketEncryption) {
		reader, _, readErr := objectStore.Get(ctx, srcBucket, srcKey)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read source object: %w", readErr)
		}
		srcData, readErr := io.ReadAll(reader)
		reader.Close()
		if readErr != nil {
			return nil, fmt.Errorf("failed to read source data: %w", readErr)
		}

		encResult, encErr := s.encryptionManager.Encrypt(srcData, encType, bucketEncryption, in.Bucket, in.Key, "")
		if encErr != nil {
			return nil, fmt.Errorf("encryption failed: %w", encErr)
		}
		sseMeta := &s3store.SSEObjectMetadata{
			EncryptionType:   s3store.SSEType(encResult.EncryptionType),
			EncryptedDataKey: encResult.EncryptedDataKey,
			ContentNonce:     encResult.ContentNonce,
			KMSKeyID:         encResult.KMSKeyID,
			UnencryptedMD5:   encResult.UnencryptedMD5,
			UnencryptedSize:  encResult.UnencryptedSize,
		}
		obj, err = objectStore.PutEncrypted(ctx, in.Bucket, in.Key, encResult.EncryptedData, contentType, metadata, sseMeta, s3store.StorageClassStandard, nil)
	} else {
		obj, err = objectStore.Copy(ctx, srcBucket, srcKey, in.Bucket, in.Key)
	}
	if err != nil {
		return nil, err
	}

	return &AdminCopyObjectResult{
		ETag:         formatETag(obj.ETag),
		LastModified: obj.LastModified.Format(timeutils.ISO8601UTCFormat),
	}, nil
}
