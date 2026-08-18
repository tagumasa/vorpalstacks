package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"time"
	"vorpalstacks/internal/common/defaults"
	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/logs"
	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

// ---------------------------------------------------------------------------
// Input structs — transport-agnostic DTOs for the admin gRPC-Web handler.
// Prefixed with Admin to avoid collisions with the HTTP API input types.
// ---------------------------------------------------------------------------

// AdminListBucketsInput carries the fields needed for ListBuckets.
type AdminListBucketsInput struct {
	Prefix            string
	BucketRegion      string
	MaxBuckets        int
	ContinuationToken string
}

// AdminCreateBucketInput carries the fields needed for CreateBucket.
type AdminCreateBucketInput struct {
	Bucket                     string
	Region                     string
	ACL                        string
	ObjectLockEnabledForBucket bool
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
	Bucket       string
	Key          string
	CopySource   string
	ContentType  string
	StorageClass string
}

// ---------------------------------------------------------------------------
// Streaming input/result structs — used by both HTTP and admin handlers
// for operations involving io.Reader bodies.
// ---------------------------------------------------------------------------

// PutObjectStreamInput is the transport-agnostic input for PutObject.
type PutObjectStreamInput struct {
	Body                 io.Reader
	ContentLength        int64
	Bucket               string
	Key                  string
	ContentType          string
	ContentEncoding      string
	ContentLanguage      string
	ContentDisposition   string
	CacheControl         string
	Metadata             map[string]string
	StorageClass         string
	IfMatch              string
	IfNoneMatch          string
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
	SSECustomerKey       string
	SSECustomerKeyMD5    string
	Tagging              string
}

// PutObjectStreamResult holds the transport-agnostic result of PutObject.
type PutObjectStreamResult struct {
	Object               *s3store.Object
	Bucket               *s3store.Bucket
	ServerSideEncryption string
	SSEKMSKeyId          string
}

// GetObjectStreamInput is the transport-agnostic input for GetObject.
type GetObjectStreamInput struct {
	Bucket               string
	Key                  string
	VersionID            string
	IfMatch              string
	IfNoneMatch          string
	IfModifiedSince      *time.Time
	IfUnmodifiedSince    *time.Time
	Range                string
	SSECustomerAlgorithm string
	SSECustomerKey       string
	SSECustomerKeyMD5    string
}

// GetObjectStreamResult holds the transport-agnostic result of GetObject.
type GetObjectStreamResult struct {
	Body                 io.ReadCloser
	ContentLength        int64
	ContentType          string
	ContentEncoding      string
	ContentLanguage      string
	ContentDisposition   string
	CacheControl         string
	ETag                 string
	LastModified         time.Time
	Metadata             map[string]string
	StorageClass         string
	VersionID            string
	ContentRange         string
	IsPartial            bool
	AcceptRanges         string
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
	SSECustomerKeyMD5    string
	ReplicationStatus    string
	SSEMetadata          *s3store.SSEObjectMetadata
	Tags                 []types.Tag
}

// CopyObjectStreamInput is the transport-agnostic input for CopyObject.
type CopyObjectStreamInput struct {
	Bucket                    string
	Key                       string
	CopySource                string
	CopySourceVersionId       string
	MetadataDirective         string
	ContentType               string
	Metadata                  map[string]string
	StorageClass              string
	ServerSideEncryption      string
	SSEKMSKeyId               string
	SSECustomerAlgorithm      string
	SSECustomerKey            string
	SSECustomerKeyMD5         string
	CopySourceSSECustomerAlgo string
	CopySourceSSECustomerKey  string
	CopySourceSSECustomerMD5  string
}

// CopyObjectStreamResult holds the transport-agnostic result of CopyObject.
type CopyObjectStreamResult struct {
	Object               *s3store.Object
	ServerSideEncryption string
	SSEKMSKeyId          string
}

// ---------------------------------------------------------------------------
// Result structs
// ---------------------------------------------------------------------------

// AdminListBucketsResult holds the transport-agnostic result of ListBuckets.
// Returns raw store buckets so each transport layer can format as needed.
type AdminListBucketsResult struct {
	Buckets           []*s3store.Bucket
	ContinuationToken string
	IsTruncated       bool
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

// listBucketsCore returns all buckets in the regional store, optionally
// filtered by prefix and bucket-region, sorted by name, and paginated.
func (s *S3Service) listBucketsCore(bucketStore s3store.BucketStoreInterface, in AdminListBucketsInput) (*AdminListBucketsResult, error) {
	buckets, err := bucketStore.List()
	if err != nil {
		return nil, err
	}

	filtered := make([]*s3store.Bucket, 0, len(buckets))
	for _, b := range buckets {
		if in.Prefix != "" && !strings.HasPrefix(b.Name, in.Prefix) {
			continue
		}
		if in.BucketRegion != "" && b.Region != in.BucketRegion {
			continue
		}
		filtered = append(filtered, b)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name < filtered[j].Name
	})

	startIdx := 0
	if in.ContinuationToken != "" {
		for i, b := range filtered {
			if b.Name > in.ContinuationToken {
				startIdx = i
				break
			}
			if i == len(filtered)-1 {
				startIdx = len(filtered)
			}
		}
	}

	maxBuckets := in.MaxBuckets
	if maxBuckets <= 0 || maxBuckets > 10000 {
		maxBuckets = 10000
	}

	endIdx := startIdx + maxBuckets
	var nextToken string
	isTruncated := false
	if endIdx < len(filtered) {
		nextToken = filtered[endIdx-1].Name
		isTruncated = true
	} else {
		endIdx = len(filtered)
	}

	return &AdminListBucketsResult{
		Buckets:           filtered[startIdx:endIdx],
		ContinuationToken: nextToken,
		IsTruncated:       isTruncated,
	}, nil
}

// createBucketCore validates the bucket name, creates the bucket in the
// regional store, and optionally applies a canned ACL and/or Object Lock.
func (s *S3Service) createBucketCore(bucketStore s3store.BucketStoreInterface, in AdminCreateBucketInput) (*AdminCreateBucketResult, error) {
	if in.Bucket == "" {
		return nil, NewInvalidArgumentError("bucket name is required")
	}

	if err := validateBucketName(in.Bucket); err != nil {
		return nil, err
	}

	// Bucket names form a single global namespace and this platform has a
	// single account, so a duplicate name always collides with a bucket the
	// requester owns. S3 returns BucketAlreadyOwnedByYou for that case in
	// every Region except North Virginia, where legacy clients rely on a
	// 200 OK response that resets the existing bucket's ACLs.
	if bucketStore.Exists(in.Bucket) {
		if in.Region == defaults.DefaultRegion {
			existing, getErr := bucketStore.Get(in.Bucket)
			if getErr != nil {
				return nil, getErr
			}
			existing.ACL = nil
			if putErr := bucketStore.Put(existing); putErr != nil {
				return nil, putErr
			}
			return &AdminCreateBucketResult{Location: "/" + in.Bucket}, nil
		}
		return nil, ErrBucketAlreadyOwnedByYou
	}
	if s.s3Store != nil {
		if found, _ := s.s3Store.FindBucket(in.Bucket); found != nil {
			return nil, ErrBucketAlreadyOwnedByYou
		}
	}

	bucket, err := bucketStore.Create(in.Bucket, in.Region)
	if err != nil {
		if errors.Is(err, s3store.ErrBucketAlreadyExists) {
			return nil, ErrBucketAlreadyOwnedByYou
		}
		return nil, err
	}

	if in.ACL != "" {
		acp, err := CannedACLToPolicy(in.ACL, &s3store.ACLOwner{ID: s.accountID, DisplayName: s.accountID})
		if err != nil {
			return nil, err
		}
		bucket.ACL = acp
	}

	if in.ObjectLockEnabledForBucket {
		bucket.ObjectLockEnabled = true
	}

	if in.ACL != "" || in.ObjectLockEnabledForBucket {
		if err := bucketStore.Put(bucket); err != nil {
			return nil, err
		}
	}

	return &AdminCreateBucketResult{Location: "/" + in.Bucket}, nil
}

// deleteBucketCore verifies the bucket is empty (no objects, no incomplete
// multipart uploads), purges the encryption key, and deletes the bucket.
func (s *S3Service) deleteBucketCore(bucketStore s3store.BucketStoreInterface, objectStore s3store.ObjectStoreInterface, in AdminDeleteBucketInput) (*AdminDeleteBucketResult, error) {
	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}
	if bucket.ObjectLockEnabled {
		logs.Warn("s3: deleting bucket with Object Lock enabled", logs.String("bucket", in.Bucket))
	}

	count, err := objectStore.CountByBucket(in.Bucket)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrBucketNotEmpty
	}

	multipartCount, err := objectStore.CountMultipartUploadsByBucket(in.Bucket)
	if err != nil {
		return nil, err
	}
	if multipartCount > 0 {
		return nil, ErrBucketNotEmpty
	}

	if s.encryptionManager != nil {
		s.encryptionManager.DeleteBucketKey(in.Bucket)
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
	streamResult, err := s.getObjectStreamCore(ctx, objectStore, GetObjectStreamInput{
		Bucket:    in.Bucket,
		Key:       in.Key,
		VersionID: in.VersionID,
	})
	if err != nil {
		return nil, err
	}
	defer streamResult.Body.Close()

	data, err := io.ReadAll(streamResult.Body)
	if err != nil {
		return nil, err
	}

	return &AdminGetObjectResult{
		Object: &s3store.Object{
			ETag:               streamResult.ETag,
			LastModified:       streamResult.LastModified,
			ContentType:        streamResult.ContentType,
			Metadata:           streamResult.Metadata,
			Size:               streamResult.ContentLength,
			StorageClass:       s3store.ObjectStorageClass(streamResult.StorageClass),
			VersionID:          streamResult.VersionID,
			ContentEncoding:    streamResult.ContentEncoding,
			ContentLanguage:    streamResult.ContentLanguage,
			ContentDisposition: streamResult.ContentDisposition,
			CacheControl:       streamResult.CacheControl,
			ReplicationStatus:  streamResult.ReplicationStatus,
			SSEMetadata:        streamResult.SSEMetadata,
			Tags:               streamResult.Tags,
		},
		Body: data,
	}, nil
}

// getObjectStreamCore is the streaming variant of getObjectCore. It returns
// the object body as io.ReadCloser and handles conditional headers, SSE
// decryption (streaming for chunked encryption, materialise for others),
// SSE-C key parsing, and Range requests. HTTP and admin handlers share
// this method.
func (s *S3Service) getObjectStreamCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in GetObjectStreamInput) (*GetObjectStreamResult, error) {
	if in.IfMatch != "" || in.IfNoneMatch != "" || in.IfModifiedSince != nil || in.IfUnmodifiedSince != nil {
		obj, err := objectStore.HeadWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
		if err != nil {
			return nil, err
		}

		if in.IfMatch != "" {
			if in.IfMatch == "*" {
				// Wildcard: object must exist. HeadWithVersion succeeded.
			} else if strings.Trim(obj.ETag, "\"") != strings.Trim(in.IfMatch, "\"") {
				return nil, ErrPreconditionFailed
			}
		}
		if in.IfNoneMatch != "" {
			if in.IfNoneMatch == "*" {
				return nil, ErrNotModified
			} else if strings.Trim(obj.ETag, "\"") == strings.Trim(in.IfNoneMatch, "\"") {
				return nil, ErrNotModified
			}
		}
		if in.IfUnmodifiedSince != nil {
			if obj.LastModified.After(*in.IfUnmodifiedSince) {
				return nil, ErrPreconditionFailed
			}
		}
		if in.IfModifiedSince != nil {
			if !obj.LastModified.After(*in.IfModifiedSince) {
				return nil, ErrNotModified
			}
		}
	}

	reader, obj, err := objectStore.GetWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
	if err != nil {
		return nil, err
	}

	if isArchiveClass(obj.StorageClass) {
		reader.Close()
		return nil, ErrInvalidObjectState
	}

	sseCRequested := sseCustomerRequested(in.SSECustomerAlgorithm, in.SSECustomerKey, in.SSECustomerKeyMD5)
	if sseCRequested && (obj.SSEMetadata == nil || obj.SSEMetadata.EncryptionType != s3store.SSETypeCustomer) {
		reader.Close()
		return nil, NewInvalidRequestError("The encryption parameters are not applicable to this object.")
	}

	result := &GetObjectStreamResult{
		Body:               reader,
		ContentLength:      obj.Size,
		ContentType:        obj.ContentType,
		ContentEncoding:    obj.ContentEncoding,
		ContentLanguage:    obj.ContentLanguage,
		ContentDisposition: obj.ContentDisposition,
		CacheControl:       obj.CacheControl,
		ETag:               formatETag(obj.ETag),
		LastModified:       obj.LastModified,
		Metadata:           obj.Metadata,
		StorageClass:       string(obj.StorageClass),
		VersionID:          obj.VersionID,
		ReplicationStatus:  obj.ReplicationStatus,
		SSEMetadata:        obj.SSEMetadata,
		Tags:               obj.Tags,
	}

	var decryptedData []byte
	var unencryptedSize int64

	// Streaming decryption for chunked encrypted objects (no Range).
	if obj.SSEMetadata != nil && in.Range == "" && len(obj.SSEMetadata.PartEncryptionInfos) > 0 {
		var customerKey []byte
		if obj.SSEMetadata.EncryptionType == s3store.SSETypeCustomer {
			if in.SSECustomerKey == "" {
				reader.Close()
				return nil, awserrors.NewAWSError("InvalidRequest", "The object was stored using a form of Server Side Encryption. The correct parameters must be provided to retrieve the object.", http.StatusBadRequest)
			}
			customerKey, err = s.encryptionManager.ParseCustomerKey(in.SSECustomerKey, in.SSECustomerKeyMD5)
			if err != nil {
				reader.Close()
				return nil, ErrInvalidSSECustomerKey
			}
			result.SSECustomerAlgorithm = "AES256"
			result.SSECustomerKeyMD5 = in.SSECustomerKeyMD5
		} else {
			result.ServerSideEncryption = string(obj.SSEMetadata.EncryptionType)
			result.SSEKMSKeyId = obj.SSEMetadata.KMSKeyID
		}

		streamReader, streamErr := s.encryptionManager.NewChunkDecryptReader(reader, obj.SSEMetadata, in.Bucket, in.Key, customerKey)
		if streamErr != nil {
			reader.Close()
			return nil, streamErr
		}

		result.Body = streamReader
		result.ContentLength = obj.SSEMetadata.UnencryptedSize
		return result, nil
	}

	if obj.SSEMetadata != nil {
		encryptedData, readErr := io.ReadAll(reader)
		reader.Close()
		if readErr != nil {
			return nil, fmt.Errorf("failed to read encrypted data: %w", readErr)
		}

		decryptedData, unencryptedSize, err = s.decryptObjectData(encryptedData, obj.SSEMetadata, in.Bucket, in.Key, in.SSECustomerKey, in.SSECustomerKeyMD5)
		if err != nil {
			return nil, err
		}

		if obj.SSEMetadata.EncryptionType == s3store.SSETypeCustomer {
			result.SSECustomerAlgorithm = "AES256"
			result.SSECustomerKeyMD5 = in.SSECustomerKeyMD5
		} else {
			result.ServerSideEncryption = string(obj.SSEMetadata.EncryptionType)
			result.SSEKMSKeyId = obj.SSEMetadata.KMSKeyID
		}
	}

	if in.Range != "" {
		ranges, rangeErr := parseRangeHeader(in.Range)
		if rangeErr != nil {
			if obj.SSEMetadata == nil {
				reader.Close()
			}
			return nil, rangeErr
		}

		firstRange := ranges[0]
		var offset, length int64
		totalSize := obj.Size
		if obj.SSEMetadata != nil {
			totalSize = unencryptedSize
		}

		if firstRange.Start == -1 {
			length = firstRange.Length
			offset = totalSize - length
			if offset < 0 {
				offset = 0
				length = totalSize
			}
		} else {
			offset = firstRange.Start
			if firstRange.Length == -1 {
				length = totalSize - offset
				if length < 0 {
					length = 0
				}
			} else {
				length = firstRange.Length
			}
		}

		if offset >= totalSize {
			if obj.SSEMetadata == nil {
				reader.Close()
			}
			return nil, ErrInvalidRange
		}

		actualEnd := offset + length - 1
		if actualEnd >= totalSize {
			actualEnd = totalSize - 1
			length = totalSize - offset
		}

		var rangeReader io.ReadCloser
		if obj.SSEMetadata != nil {
			start := offset
			end := start + length
			if end > int64(len(decryptedData)) {
				end = int64(len(decryptedData))
			}
			rangeReader = io.NopCloser(bytes.NewReader(decryptedData[start:end]))
		} else {
			reader.Close()
			rangeReader, _, err = objectStore.GetRangeWithVersion(ctx, in.Bucket, in.Key, in.VersionID, offset, length)
			if err != nil {
				return nil, err
			}
		}

		return &GetObjectStreamResult{
			Body:                 rangeReader,
			ContentLength:        length,
			ContentType:          obj.ContentType,
			ETag:                 formatETag(obj.ETag),
			LastModified:         obj.LastModified,
			Metadata:             obj.Metadata,
			ContentRange:         fmt.Sprintf("bytes %d-%d/%d", offset, actualEnd, totalSize),
			IsPartial:            true,
			AcceptRanges:         "bytes",
			VersionID:            obj.VersionID,
			ServerSideEncryption: result.ServerSideEncryption,
			SSEKMSKeyId:          result.SSEKMSKeyId,
			SSECustomerAlgorithm: result.SSECustomerAlgorithm,
			SSECustomerKeyMD5:    result.SSECustomerKeyMD5,
			ContentEncoding:      obj.ContentEncoding,
			ContentLanguage:      obj.ContentLanguage,
			ContentDisposition:   obj.ContentDisposition,
			CacheControl:         obj.CacheControl,
			StorageClass:         string(obj.StorageClass),
		}, nil
	}

	if obj.SSEMetadata != nil {
		result.Body = io.NopCloser(bytes.NewReader(decryptedData))
		result.ContentLength = unencryptedSize
	}

	return result, nil
}

// decryptObjectData decrypts encrypted object data using the appropriate
// decryption method based on the SSE metadata. This is the Core-layer
// counterpart of ObjectOperations.decryptObjectData, taking individual
// parameters instead of the HTTP-specific GetObjectInput struct.
func (s *S3Service) decryptObjectData(encryptedData []byte, sseMeta *s3store.SSEObjectMetadata, bucket, key, sseCustomerKey, sseCustomerKeyMD5 string) ([]byte, int64, error) {
	unencryptedSize := sseMeta.UnencryptedSize

	if sseMeta.EncryptionType == s3store.SSETypeCustomer {
		if sseCustomerKey == "" {
			return nil, 0, awserrors.NewAWSError("InvalidRequest", "The object was stored using a form of Server Side Encryption. The correct parameters must be provided to retrieve the object.", http.StatusBadRequest)
		}
		customerKey, err := s.encryptionManager.ParseCustomerKey(sseCustomerKey, sseCustomerKeyMD5)
		if err != nil {
			return nil, 0, ErrInvalidSSECustomerKey
		}

		var plainData []byte
		if len(sseMeta.PartEncryptionInfos) > 0 {
			plainData, err = s.encryptionManager.DecryptChunked(encryptedData, sseMeta, bucket, key, customerKey)
			if err != nil {
				return nil, 0, ErrInvalidSSECustomerKey
			}
		} else {
			decResult, decErr := s.encryptionManager.DecryptWithCustomerKey(encryptedData, sseMeta, bucket, key, customerKey)
			if decErr != nil {
				return nil, 0, ErrInvalidSSECustomerKey
			}
			plainData = decResult.DecryptedData
		}
		return plainData, unencryptedSize, nil
	}

	var plainData []byte
	var err error
	if len(sseMeta.PartEncryptionInfos) > 0 {
		plainData, err = s.encryptionManager.DecryptChunked(encryptedData, sseMeta, bucket, key, nil)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to decrypt chunked data: %w", err)
		}
	} else {
		decResult, decErr := s.encryptionManager.Decrypt(encryptedData, sseMeta, bucket, key)
		if decErr != nil {
			return nil, 0, fmt.Errorf("failed to decrypt data: %w", decErr)
		}
		plainData = decResult.DecryptedData
	}
	return plainData, unencryptedSize, nil
}

// putObjectCore validates the upload, determines encryption settings, and
// stores the object.
func (s *S3Service) putObjectCore(ctx context.Context, bucketStore s3store.BucketStoreInterface, objectStore s3store.ObjectStoreInterface, in AdminPutObjectInput) (*AdminPutObjectResult, error) {
	streamResult, err := s.putObjectStreamCore(ctx, bucketStore, objectStore, PutObjectStreamInput{
		Body:          bytes.NewReader(in.Body),
		ContentLength: int64(len(in.Body)),
		Bucket:        in.Bucket,
		Key:           in.Key,
		ContentType:   in.ContentType,
		Metadata:      in.Metadata,
	})
	if err != nil {
		return nil, err
	}
	obj := streamResult.Object
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

// putObjectStreamCore is the streaming variant of putObjectCore. It accepts
// an io.Reader body and handles the full Put logic: conditional headers,
// encryption type determination, SSE-C key parsing, EncryptStream for
// encrypted path, PutWithVersioning for non-encrypted path, and tagging.
// HTTP and admin handlers share this method to avoid duplicated store
// interaction logic.
func (s *S3Service) putObjectStreamCore(ctx context.Context, bucketStore s3store.BucketStoreInterface, objectStore s3store.ObjectStoreInterface, in PutObjectStreamInput) (*PutObjectStreamResult, error) {
	if in.ContentLength > maxSingleUploadSize {
		return nil, ErrEntityTooLarge
	}

	if in.IfMatch != "" || in.IfNoneMatch != "" {
		existingObj, headErr := objectStore.Head(ctx, in.Bucket, in.Key)
		objectExists := headErr == nil && existingObj != nil

		if in.IfNoneMatch == "*" {
			if objectExists {
				return nil, ErrPreconditionFailed
			}
		} else if in.IfNoneMatch != "" {
			if objectExists && strings.Trim(existingObj.ETag, "\"") == strings.Trim(in.IfNoneMatch, "\"") {
				return nil, ErrPreconditionFailed
			}
		}

		if in.IfMatch != "" {
			if !objectExists {
				return nil, ErrPreconditionFailed
			}
			if strings.Trim(existingObj.ETag, "\"") != strings.Trim(in.IfMatch, "\"") {
				return nil, ErrPreconditionFailed
			}
		}
	}

	bucket, err := bucketStore.Get(in.Bucket)
	if err != nil {
		return nil, err
	}

	// An overwrite must carry the same encryption family as the object it
	// replaces: a PUT that supplies SSE-C parameters for a non-SSE-C object,
	// or omits them for an SSE-C object, is rejected rather than silently
	// re-encrypting under a different type.
	requestUsesSSEC := sseCustomerRequested(in.SSECustomerAlgorithm, in.SSECustomerKey, in.SSECustomerKeyMD5)
	if existing, headErr := objectStore.Head(ctx, in.Bucket, in.Key); headErr == nil && existing != nil {
		existingIsSSEC := existing.SSEMetadata != nil && existing.SSEMetadata.EncryptionType == s3store.SSETypeCustomer
		if existingIsSSEC != requestUsesSSEC {
			return nil, ErrEncryptionTypeMismatch
		}
	}

	metadata := in.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	contentType := in.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	var encryptionType EncryptionType
	var customerKey []byte
	if in.SSECustomerAlgorithm != "" {
		encryptionType = EncryptionTypeSSE_C
		customerKey, err = s.encryptionManager.ParseCustomerKey(in.SSECustomerKey, in.SSECustomerKeyMD5)
		if err != nil {
			return nil, NewInvalidArgumentError(fmt.Sprintf("invalid SSE-C customer key: %v", err))
		}
	} else {
		encryptionType = s.encryptionManager.DetermineEncryptionType(
			EncryptionType(in.ServerSideEncryption),
			bucket.EncryptionConfig,
		)
	}

	storageClass := s3store.ObjectStorageClass(in.StorageClass)
	if storageClass == "" {
		storageClass = s3store.StorageClassStandard
	}

	sysMeta := &s3store.SystemMetadata{
		ContentEncoding:    in.ContentEncoding,
		ContentLanguage:    in.ContentLanguage,
		ContentDisposition: in.ContentDisposition,
		CacheControl:       in.CacheControl,
	}

	var obj *s3store.Object
	result := &PutObjectStreamResult{Bucket: bucket}

	if s.encryptionManager.ShouldEncrypt(encryptionType, bucket.EncryptionConfig) {
		encResult, encErr := s.encryptionManager.EncryptStream(in.Body, encryptionType, bucket.EncryptionConfig, in.Bucket, in.Key, in.SSEKMSKeyId, customerKey)
		if encErr != nil {
			return nil, fmt.Errorf("failed to encrypt data: %w", encErr)
		}
		obj, err = objectStore.PutEncrypted(ctx, in.Bucket, in.Key, encResult.EncryptedData, contentType, metadata, encResult.SSEMetadata, storageClass, sysMeta)
		if err != nil {
			return nil, err
		}
		result.ServerSideEncryption = string(encResult.SSEMetadata.EncryptionType)
		result.SSEKMSKeyId = encResult.SSEMetadata.KMSKeyID
	} else {
		obj, err = objectStore.PutWithVersioning(ctx, in.Bucket, in.Key, in.Body, contentType, metadata, false, storageClass, sysMeta)
		if err != nil {
			return nil, err
		}
	}

	if in.Tagging != "" {
		parsedTags := parseTaggingHeader(in.Tagging)
		if len(parsedTags) > 0 {
			obj.Tags = parsedTags
			_ = objectStore.SetTags(in.Bucket, in.Key, parsedTags)
		}
	}

	result.Object = obj
	return result, nil
}

// deleteObjectCore deletes a single object. In a versioned bucket, deleting
// without a specific VersionID creates a delete marker.
func (s *S3Service) deleteObjectCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in AdminDeleteObjectInput) (*AdminDeleteObjectResult, error) {
	result, err := objectStore.DeleteWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return &AdminDeleteObjectResult{}, nil
	}
	return &AdminDeleteObjectResult{
		VersionID:      result.VersionID,
		IsDeleteMarker: result.IsDeleteMarker,
	}, nil
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
	streamResult, err := s.copyObjectStreamCore(ctx, bucketStore, objectStore, CopyObjectStreamInput{
		Bucket:       in.Bucket,
		Key:          in.Key,
		CopySource:   in.CopySource,
		ContentType:  in.ContentType,
		StorageClass: in.StorageClass,
	})
	if err != nil {
		return nil, err
	}
	obj := streamResult.Object
	return &AdminCopyObjectResult{
		ETag:         formatETag(obj.ETag),
		LastModified: obj.LastModified.Format(timeutils.ISO8601UTCFormat),
	}, nil
}

// isArchiveClass reports whether a storage class places objects in an
// archive tier that must be restored before the object data can be read.
// GLACIER_IR is excluded because it offers real-time retrieval.
func isArchiveClass(cls s3store.ObjectStorageClass) bool {
	return cls == s3store.StorageClassGlacier || cls == s3store.StorageClassDeepArchive
}

// copyObjectStreamCore is the streaming variant of copyObjectCore. It handles
// the full server-side copy logic: source object retrieval (with optional
// SSE-C decryption), target encryption determination, EncryptStream or
// store-level Copy, and metadata directive handling. HTTP and admin
// handlers share this method.
func (s *S3Service) copyObjectStreamCore(ctx context.Context, bucketStore s3store.BucketStoreInterface, objectStore s3store.ObjectStoreInterface, in CopyObjectStreamInput) (*CopyObjectStreamResult, error) {
	if err := validateStorageClass(in.StorageClass); err != nil {
		return nil, err
	}

	srcBucket, srcKey, srcVersionId, err := parseCopySource(in.CopySource)
	if err != nil {
		return nil, err
	}

	if in.CopySourceVersionId != "" {
		srcVersionId = in.CopySourceVersionId
	}

	var srcObj *s3store.Object
	if srcVersionId != "" {
		srcObj, err = objectStore.HeadWithVersion(ctx, srcBucket, srcKey, srcVersionId)
	} else {
		srcObj, err = objectStore.GetMetadata(srcBucket, srcKey)
	}
	if err != nil {
		if errors.Is(err, s3store.ErrObjectNotFound) {
			return nil, ErrNoSuchKey
		}
		return nil, err
	}

	if isArchiveClass(srcObj.StorageClass) {
		return nil, ErrObjectNotInActiveTier
	}

	if srcObj.Size > maxCopyObjectSize {
		return nil, ErrEntityTooLarge
	}

	var srcReader io.Reader
	if srcObj.SSEMetadata != nil || sseCustomerRequested(in.CopySourceSSECustomerAlgo, in.CopySourceSSECustomerKey, in.CopySourceSSECustomerMD5) {
		getResult, getErr := s.getObjectStreamCore(ctx, objectStore, GetObjectStreamInput{
			Bucket:               srcBucket,
			Key:                  srcKey,
			VersionID:            srcVersionId,
			SSECustomerAlgorithm: in.CopySourceSSECustomerAlgo,
			SSECustomerKey:       in.CopySourceSSECustomerKey,
			SSECustomerKeyMD5:    in.CopySourceSSECustomerMD5,
		})
		if getErr != nil {
			return nil, getErr
		}
		defer getResult.Body.Close()
		srcReader = getResult.Body
	} else {
		var reader io.ReadCloser
		if srcVersionId != "" {
			reader, _, err = objectStore.GetWithVersion(ctx, srcBucket, srcKey, srcVersionId)
		} else {
			reader, _, err = objectStore.Get(ctx, srcBucket, srcKey)
		}
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		srcReader = reader
	}

	bucketEncryption, err := bucketStore.GetEncryptionConfiguration(in.Bucket)
	if err != nil {
		return nil, err
	}

	var targetEncryptionType EncryptionType
	var targetKMSKeyID string

	if in.ServerSideEncryption != "" {
		targetEncryptionType = EncryptionType(in.ServerSideEncryption)
		targetKMSKeyID = in.SSEKMSKeyId
	} else if sseCustomerRequested(in.SSECustomerAlgorithm, in.SSECustomerKey, in.SSECustomerKeyMD5) {
		targetEncryptionType = EncryptionTypeSSE_C
	} else {
		targetEncryptionType = s.encryptionManager.DetermineEncryptionType(EncryptionTypeNone, bucketEncryption)
		if targetEncryptionType == EncryptionTypeSSE_KMS && bucketEncryption != nil {
			targetKMSKeyID = bucketEncryption.KMSMasterKeyID
		}
	}

	contentType := in.ContentType
	if contentType == "" {
		contentType = srcObj.ContentType
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	metadata := in.Metadata
	if in.MetadataDirective != "" && in.MetadataDirective != "COPY" && in.MetadataDirective != "REPLACE" {
		return nil, NewInvalidArgumentError(fmt.Sprintf("invalid MetadataDirective: %s (must be COPY or REPLACE)", in.MetadataDirective))
	}
	if in.MetadataDirective != "REPLACE" {
		metadata = srcObj.Metadata
	}

	var obj *s3store.Object
	result := &CopyObjectStreamResult{}

	if targetEncryptionType != EncryptionTypeNone {
		var customerKey []byte
		if in.SSECustomerKey != "" {
			customerKey, err = s.encryptionManager.ParseCustomerKey(in.SSECustomerKey, in.SSECustomerKeyMD5)
			if err != nil {
				return nil, err
			}
		}

		encResult, encErr := s.encryptionManager.EncryptStream(srcReader, targetEncryptionType, bucketEncryption, in.Bucket, in.Key, targetKMSKeyID, customerKey)
		if encErr != nil {
			return nil, encErr
		}

		targetStorageClass := s3store.ObjectStorageClass(in.StorageClass)
		if targetStorageClass == "" {
			targetStorageClass = srcObj.StorageClass
		}
		if targetStorageClass == "" {
			targetStorageClass = s3store.StorageClassStandard
		}
		obj, err = objectStore.PutEncrypted(ctx, in.Bucket, in.Key, encResult.EncryptedData, contentType, metadata, encResult.SSEMetadata, targetStorageClass, nil)
		if err != nil {
			return nil, err
		}
		result.ServerSideEncryption = string(encResult.SSEMetadata.EncryptionType)
		if encResult.SSEMetadata.KMSKeyID != "" {
			result.SSEKMSKeyId = encResult.SSEMetadata.KMSKeyID
		}
	} else {
		dstStorageClass := s3store.ObjectStorageClass(in.StorageClass)
		if srcVersionId != "" {
			if in.MetadataDirective == "REPLACE" {
				obj, err = objectStore.CopyWithVersionAndMetadata(ctx, srcBucket, srcKey, srcVersionId, in.Bucket, in.Key, contentType, metadata, dstStorageClass)
			} else {
				obj, err = objectStore.CopyWithVersion(ctx, srcBucket, srcKey, srcVersionId, in.Bucket, in.Key, dstStorageClass)
			}
		} else {
			if in.MetadataDirective == "REPLACE" {
				obj, err = objectStore.CopyWithMetadata(ctx, srcBucket, srcKey, in.Bucket, in.Key, contentType, metadata, dstStorageClass)
			} else {
				obj, err = objectStore.Copy(ctx, srcBucket, srcKey, in.Bucket, in.Key, dstStorageClass)
			}
		}
		if err != nil {
			return nil, err
		}
	}

	result.Object = obj
	return result, nil
}

// UpdateObjectEncryptionInput is the transport-agnostic input for updating
// the server-side encryption of an existing object.
type UpdateObjectEncryptionInput struct {
	Bucket    string
	Key       string
	VersionID string
	KMSKeyArn string
}

// updateObjectEncryptionCore re-encrypts an existing SSE-S3 or SSE-KMS object
// under the requested KMS key. Object data is rewritten in place — the ETag,
// timestamps, storage class, tags, ACL, lock state, and version identifier
// are preserved and no new version is created. Unencrypted sources and
// DSSE-KMS or SSE-C sources are rejected, as are objects protected by an
// active Object Lock, matching the S3 contract for this operation.
func (s *S3Service) updateObjectEncryptionCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in UpdateObjectEncryptionInput) error {
	if err := validateKMSKeyArn(in.KMSKeyArn); err != nil {
		return err
	}

	encryptedData, obj, err := objectStore.GetEncrypted(ctx, in.Bucket, in.Key, in.VersionID)
	if err != nil {
		if errors.Is(err, s3store.ErrObjectNotFound) {
			return NewNoSuchKeyError(in.Key)
		}
		return err
	}

	if obj.SSEMetadata == nil {
		return NewInvalidRequestError("The UpdateObjectEncryption operation doesn't support unencrypted source objects. Only source objects encrypted with SSE-S3 or SSE-KMS are supported.")
	}
	switch obj.SSEMetadata.EncryptionType {
	case s3store.SSETypeAES256, s3store.SSETypeKMS:
	default:
		return NewInvalidRequestError("The UpdateObjectEncryption operation doesn't support source objects with the encryption type DSSE-KMS or SSE-C. Only source objects encrypted with SSE-S3 or SSE-KMS are supported.")
	}

	if hold := obj.ObjectLockLegalHold; hold != nil && hold.Status == s3store.ObjectLockLegalHoldOn {
		return awserrors.NewAWSError("AccessDenied", "The encryption type for the specified object can't be updated because that object is protected by S3 Object Lock. If the object has a governance-mode retention period or a legal hold, you must first remove the Object Lock status on the object before you issue your UpdateObjectEncryption request.", http.StatusForbidden)
	}
	if ret := obj.ObjectLockRetention; ret != nil && ret.Mode != "" {
		if ret.RetainUntilDate == nil || ret.RetainUntilDate.After(time.Now().UTC()) {
			if ret.Mode == s3store.ObjectLockRetentionModeCompliance {
				return awserrors.NewAWSError("AccessDenied", "The encryption type for the specified object can't be updated because that object is protected by S3 Object Lock. You can't use the UpdateObjectEncryption operation with objects that have an Object Lock compliance mode retention period applied to them.", http.StatusForbidden)
			}
			return awserrors.NewAWSError("AccessDenied", "The encryption type for the specified object can't be updated because that object is protected by S3 Object Lock. If the object has a governance-mode retention period or a legal hold, you must first remove the Object Lock status on the object before you issue your UpdateObjectEncryption request.", http.StatusForbidden)
		}
	}

	if s.bus != nil {
		if invoker := s.bus.KMSInvoker(); invoker != nil && !invoker.KeyExists(ctx, in.KMSKeyArn) {
			return NewInvalidRequestError("Requests that modify an object's encryption type to SSE-KMS require a valid Amazon Web Services KMS key Amazon Resource Name (ARN). Confirm that you have a correctly formatted KMS key ARN in your request, and then try again.")
		}
	}

	plainData, _, err := s.decryptObjectData(encryptedData, obj.SSEMetadata, in.Bucket, in.Key, "", "")
	if err != nil {
		return err
	}

	encResult, err := s.encryptionManager.EncryptStream(bytes.NewReader(plainData), EncryptionTypeSSE_KMS, nil, in.Bucket, in.Key, in.KMSKeyArn, nil)
	if err != nil {
		return err
	}

	_, err = objectStore.UpdateObjectEncryption(ctx, in.Bucket, in.Key, obj.VersionID, encResult.EncryptedData, encResult.SSEMetadata)
	return err
}
