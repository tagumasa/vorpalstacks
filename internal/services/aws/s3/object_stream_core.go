package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// Streaming object engines: versioned put, ranged get with SSE decryption,
// server-side copy and in-place encryption rewrite. Shared by the AWS API
// handlers and the admin console through the Core functions.

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
	// ACL is the resolved object ACL for the upload; nil leaves the object
	// without an ACL policy.
	ACL *s3store.AccessControlPolicy
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
	PartNumber           int
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
	Restore              string
	ContentRange         string
	IsPartial            bool
	AcceptRanges         string
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
	SSECustomerKeyMD5    string
	ReplicationStatus    string
	Expiration           string
	SSEMetadata          *s3store.SSEObjectMetadata
	Tags                 []types.Tag
	// PartsCount carries x-amz-mp-parts-count for partNumber reads of
	// multipart-uploaded objects; zero omits the header.
	PartsCount int32
}

// CopyObjectStreamInput is the transport-agnostic input for CopyObject.
type CopyObjectStreamInput struct {
	Bucket                      string
	Key                         string
	CopySource                  string
	CopySourceVersionId         string
	CopySourceIfMatch           string
	CopySourceIfNoneMatch       string
	CopySourceIfModifiedSince   *time.Time
	CopySourceIfUnmodifiedSince *time.Time
	MetadataDirective           string
	ContentType                 string
	ContentEncoding             string
	ContentDisposition          string
	ContentLanguage             string
	CacheControl                string
	Metadata                    map[string]string
	StorageClass                string
	ServerSideEncryption        string
	SSEKMSKeyId                 string
	SSECustomerAlgorithm        string
	SSECustomerKey              string
	SSECustomerKeyMD5           string
	CopySourceSSECustomerAlgo   string
	CopySourceSSECustomerKey    string
	CopySourceSSECustomerMD5    string
	// ACL is the resolved object ACL for the copied object; nil leaves the
	// copy without an ACL policy.
	ACL *s3store.AccessControlPolicy
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

// getObjectStreamCore is the streaming variant of getObjectCore. It returns
// the object body as io.ReadCloser and handles conditional headers, SSE
// decryption (streaming for chunked encryption, materialise for others),
// SSE-C key parsing, and Range requests. HTTP and admin handlers share
// this method.
func (s *S3Service) getObjectStreamCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in GetObjectStreamInput) (*GetObjectStreamResult, error) {
	if in.IfMatch != "" || in.IfNoneMatch != "" || in.IfModifiedSince != nil || in.IfUnmodifiedSince != nil {
		obj, err := objectStore.HeadWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
		if err != nil {
			return nil, mapVersionReadError(err, obj, in.Key, in.VersionID)
		}
		if err := checkObjectPreconditions(obj, in.IfMatch, in.IfNoneMatch, in.IfModifiedSince, in.IfUnmodifiedSince); err != nil {
			return nil, err
		}
	}

	partsCount := int32(0)
	if in.PartNumber > 0 {
		meta, err := objectStore.HeadWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
		if err != nil {
			return nil, mapVersionReadError(err, meta, in.Key, in.VersionID)
		}
		totalSize := meta.Size
		if meta.SSEMetadata != nil {
			totalSize = meta.SSEMetadata.UnencryptedSize
		}
		start, end, count, err := resolvePartRange(meta.Parts, in.PartNumber, totalSize, in.Range)
		if err != nil {
			return nil, err
		}
		partsCount = count
		// The part selection behaves as a ranged GET of the part's bytes.
		in.Range = fmt.Sprintf("bytes=%d-%d", start, end)
		in.PartNumber = 0
	}

	reader, obj, err := objectStore.GetWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
	if err != nil {
		return nil, mapVersionReadError(err, obj, in.Key, in.VersionID)
	}

	if isArchiveClass(obj.StorageClass) && !objectRestored(obj, time.Now()) {
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
		Restore:            restoreHeaderValue(obj, time.Now()),
		ReplicationStatus:  obj.ReplicationStatus,
		SSEMetadata:        obj.SSEMetadata,
		Tags:               obj.Tags,
		Expiration:         s.objectExpirationHeaderFor(in.Bucket, obj),
	}

	var unencryptedSize int64

	// The range window resolves before any decryption work: an invalid
	// range fails the request outright, and the chunked path below fetches
	// only the window's chunks.
	var rangeOffset, rangeLength int64
	ranged := false
	if in.Range != "" {
		ranges, rangeErr := parseRangeHeader(in.Range)
		if rangeErr != nil {
			reader.Close()
			return nil, rangeErr
		}

		totalSize := obj.Size
		if obj.SSEMetadata != nil {
			totalSize = obj.SSEMetadata.UnencryptedSize
		}

		var winErr error
		rangeOffset, rangeLength, winErr = resolveRangeWindow(ranges[0], totalSize)
		if winErr != nil {
			reader.Close()
			return nil, winErr
		}
		ranged = true
	}
	actualEnd := rangeOffset + rangeLength - 1

	// Streaming decryption for chunked encrypted objects (no Range).
	if obj.SSEMetadata != nil && !ranged && len(obj.SSEMetadata.PartEncryptionInfos) > 0 {
		customerKey, keyErr := s.resolveEncryptedReadKey(obj, in, result)
		if keyErr != nil {
			reader.Close()
			return nil, keyErr
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
		customerKey, keyErr := s.resolveEncryptedReadKey(obj, in, result)
		if keyErr != nil {
			reader.Close()
			return nil, keyErr
		}
		unencryptedSize = obj.SSEMetadata.UnencryptedSize

		// A ranged read of a chunk-encrypted object decrypts only the
		// chunks the window overlaps instead of the whole ciphertext.
		if ranged && len(obj.SSEMetadata.PartEncryptionInfos) > 0 {
			reader.Close()
			window, rangeErr := s.encryptionManager.DecryptChunkRange(obj.SSEMetadata, in.Bucket, in.Key, customerKey, rangeOffset, rangeLength, func(encOffset, encLength int64) ([]byte, error) {
				blobReader, _, rErr := objectStore.GetRangeWithVersion(ctx, in.Bucket, in.Key, in.VersionID, encOffset, encLength)
				if rErr != nil {
					return nil, rErr
				}
				defer blobReader.Close()
				return io.ReadAll(blobReader)
			})
			if rangeErr != nil {
				return nil, rangeErr
			}

			result.Body = io.NopCloser(bytes.NewReader(window))
			result.ContentLength = rangeLength
			result.ContentRange = fmt.Sprintf("bytes %d-%d/%d", rangeOffset, actualEnd, unencryptedSize)
			result.IsPartial = true
			result.AcceptRanges = "bytes"
			result.PartsCount = partsCount
			return result, nil
		}

		encryptedData, readErr := io.ReadAll(reader)
		reader.Close()
		if readErr != nil {
			return nil, fmt.Errorf("failed to read encrypted data: %w", readErr)
		}

		decryptedData, _, err := s.decryptObjectData(encryptedData, obj.SSEMetadata, in.Bucket, in.Key, in.SSECustomerKey, in.SSECustomerKeyMD5)
		if err != nil {
			return nil, err
		}
		if ranged {
			start := rangeOffset
			end := start + rangeLength
			if end > int64(len(decryptedData)) {
				end = int64(len(decryptedData))
			}
			result.Body = io.NopCloser(bytes.NewReader(decryptedData[start:end]))
			result.ContentLength = rangeLength
			result.ContentRange = fmt.Sprintf("bytes %d-%d/%d", rangeOffset, actualEnd, unencryptedSize)
			result.IsPartial = true
			result.AcceptRanges = "bytes"
		} else {
			result.Body = io.NopCloser(bytes.NewReader(decryptedData))
			result.ContentLength = unencryptedSize
		}
		result.PartsCount = partsCount
		return result, nil
	}

	if ranged {
		reader.Close()
		rangeReader, _, err := objectStore.GetRangeWithVersion(ctx, in.Bucket, in.Key, in.VersionID, rangeOffset, rangeLength)
		if err != nil {
			return nil, err
		}

		result.Body = rangeReader
		result.ContentLength = rangeLength
		result.ContentRange = fmt.Sprintf("bytes %d-%d/%d", rangeOffset, actualEnd, obj.Size)
		result.IsPartial = true
		result.AcceptRanges = "bytes"
		result.PartsCount = partsCount
		return result, nil
	}

	result.PartsCount = partsCount
	return result, nil
}

// resolveEncryptedReadKey validates the SSE-C request parameters against an
// encrypted object for a read and stamps the encryption response fields on
// the result. It returns the parsed customer key for SSE-C objects and nil
// for the server-side encryption types.
func (s *S3Service) resolveEncryptedReadKey(obj *s3store.Object, in GetObjectStreamInput, result *GetObjectStreamResult) ([]byte, error) {
	if obj.SSEMetadata.EncryptionType == s3store.SSETypeCustomer {
		if in.SSECustomerKey == "" {
			return nil, awserrors.NewAWSError("InvalidRequest", "The object was stored using a form of Server Side Encryption. The correct parameters must be provided to retrieve the object.", http.StatusBadRequest)
		}
		customerKey, err := s.encryptionManager.ParseCustomerKey(in.SSECustomerKey, in.SSECustomerKeyMD5)
		if err != nil {
			return nil, ErrInvalidSSECustomerKey
		}
		result.SSECustomerAlgorithm = "AES256"
		result.SSECustomerKeyMD5 = in.SSECustomerKeyMD5
		return customerKey, nil
	}
	result.ServerSideEncryption = string(obj.SSEMetadata.EncryptionType)
	result.SSEKMSKeyId = obj.SSEMetadata.KMSKeyID
	return nil, nil
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

	// One head answers both the conditional-request preconditions and the
	// encryption-family overwrite check below.
	existingObj, headErr := objectStore.Head(ctx, in.Bucket, in.Key)
	objectExists := headErr == nil && existingObj != nil

	if in.IfMatch != "" || in.IfNoneMatch != "" {
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
	if objectExists {
		existingIsSSEC := existingObj.SSEMetadata != nil && existingObj.SSEMetadata.EncryptionType == s3store.SSETypeCustomer
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

	// The tagging header resolves before any bytes are written: a malformed
	// header or an invalid tag set must fail the request without leaving a
	// stored object behind, and the tag-set rules are the shared
	// validator's — the same ones PutObjectTagging applies.
	var parsedTags []types.Tag
	if in.Tagging != "" {
		parsedTags, err = parseTaggingHeader(in.Tagging)
		if err != nil {
			return nil, err
		}
		if err := validateTags(parsedTags); err != nil {
			return nil, err
		}
	}

	var obj *s3store.Object
	result := &PutObjectStreamResult{Bucket: bucket}

	if s.encryptionManager.ShouldEncrypt(encryptionType, bucket.EncryptionConfig) {
		sseMeta := &s3store.SSEObjectMetadata{}
		encReader, encErr := s.encryptionManager.NewChunkEncryptReader(in.Body, encryptionType, bucket.EncryptionConfig, in.Bucket, in.Key, in.SSEKMSKeyId, customerKey, sseMeta)
		if encErr != nil {
			return nil, fmt.Errorf("failed to encrypt data: %w", encErr)
		}
		obj, err = objectStore.PutEncryptedStreaming(ctx, in.Bucket, in.Key, encReader, contentType, metadata, sseMeta, storageClass, sysMeta)
		if err != nil {
			return nil, err
		}
		result.ServerSideEncryption = string(sseMeta.EncryptionType)
		result.SSEKMSKeyId = sseMeta.KMSKeyID
	} else {
		obj, err = objectStore.PutWithVersioning(ctx, in.Bucket, in.Key, in.Body, contentType, metadata, false, storageClass, sysMeta)
		if err != nil {
			return nil, err
		}
	}

	if len(parsedTags) > 0 {
		obj.Tags = parsedTags
		if err := objectStore.SetTags(in.Bucket, in.Key, "", parsedTags); err != nil {
			return nil, err
		}
	}

	if in.ACL != nil {
		obj.ACL = in.ACL
		if err := objectStore.SetACL(in.Bucket, in.Key, in.ACL); err != nil {
			return nil, err
		}
	}

	result.Object = obj
	return result, nil
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

	if err := checkCopySourcePreconditions(srcObj, in.CopySourceIfMatch, in.CopySourceIfNoneMatch, in.CopySourceIfModifiedSince, in.CopySourceIfUnmodifiedSince); err != nil {
		return nil, err
	}

	if isArchiveClass(srcObj.StorageClass) && !objectRestored(srcObj, time.Now()) {
		return nil, ErrObjectNotInActiveTier
	}

	if srcObj.Size > maxSingleUploadSize {
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

	// The metadata directive governs the content-* fields exactly as it
	// governs user metadata: COPY carries the source values over, REPLACE
	// installs the request's.
	sysMeta := &s3store.SystemMetadata{
		ContentEncoding:    srcObj.ContentEncoding,
		ContentLanguage:    srcObj.ContentLanguage,
		ContentDisposition: srcObj.ContentDisposition,
		CacheControl:       srcObj.CacheControl,
	}
	if in.MetadataDirective == "REPLACE" {
		sysMeta.ContentEncoding = in.ContentEncoding
		sysMeta.ContentLanguage = in.ContentLanguage
		sysMeta.ContentDisposition = in.ContentDisposition
		sysMeta.CacheControl = in.CacheControl
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

		sseMeta := &s3store.SSEObjectMetadata{}
		encReader, encErr := s.encryptionManager.NewChunkEncryptReader(srcReader, targetEncryptionType, bucketEncryption, in.Bucket, in.Key, targetKMSKeyID, customerKey, sseMeta)
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
		obj, err = objectStore.PutEncryptedStreaming(ctx, in.Bucket, in.Key, encReader, contentType, metadata, sseMeta, targetStorageClass, sysMeta)
		if err != nil {
			return nil, err
		}
		result.ServerSideEncryption = string(sseMeta.EncryptionType)
		if sseMeta.KMSKeyID != "" {
			result.SSEKMSKeyId = sseMeta.KMSKeyID
		}
	} else {
		// The metadata directive selects what the copy replaces: REPLACE
		// swaps in the request's content type and metadata, while COPY
		// carries the source values over untouched.
		overrides := &s3store.CopyOverrides{
			VersionId:       srcVersionId,
			ReplaceMetadata: in.MetadataDirective == "REPLACE",
			Metadata:        metadata,
			StorageClass:    s3store.ObjectStorageClass(in.StorageClass),
		}
		if overrides.ReplaceMetadata {
			overrides.ContentType = contentType
			overrides.ContentEncoding = in.ContentEncoding
			overrides.ContentDisposition = in.ContentDisposition
			overrides.ContentLanguage = in.ContentLanguage
			overrides.CacheControl = in.CacheControl
		}
		obj, err = objectStore.CopyObject(ctx, srcBucket, srcKey, in.Bucket, in.Key, overrides)
		if err != nil {
			return nil, err
		}
	}

	if in.ACL != nil {
		obj.ACL = in.ACL
		if err := objectStore.SetACL(in.Bucket, in.Key, in.ACL); err != nil {
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
			return versionLookupError(in.Key, in.VersionID)
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
