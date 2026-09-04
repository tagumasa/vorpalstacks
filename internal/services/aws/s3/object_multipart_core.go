package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// Core functions for the multipart-upload operations. The ObjectOperations
// methods in object_multipart.go are thin adapters that receive the
// per-region store bundle and delegate all validation and persistence here.

type ssePartEncryptResult struct {
	Reader        io.Reader
	EncryptedSize int64
	PlainSize     int64
	ContentNonce  []byte
	DataKey       []byte
}

func (s *S3Service) encryptPartData(data []byte, upload *s3store.MultipartUpload, inputBucket, inputKey, sseCustomerKey, sseCustomerKeyMD5 string, store *s3Stores) (*ssePartEncryptResult, error) {
	result := &ssePartEncryptResult{
		Reader:    bytes.NewReader(data),
		PlainSize: int64(len(data)),
	}
	if upload.SSEType == "" {
		return result, nil
	}
	var encResult *EncryptionResult
	var err error
	switch upload.SSEType {
	case s3store.SSETypeAES256:
		if upload.PlaintextDataKey != nil && upload.SSEMetadata != nil {
			encResult, err = s.encryptionManager.EncryptWithPlaintextKey(data, EncryptionTypeSSE_S3, inputBucket, upload.PlaintextDataKey, "", nil)
		} else {
			encResult, err = s.encryptionManager.Encrypt(data, EncryptionTypeSSE_S3, nil, inputBucket, inputKey, "")
		}
	case s3store.SSETypeKMS, s3store.SSETypeDSSEKMS:
		if upload.PlaintextDataKey != nil && upload.SSEMetadata != nil {
			encResult, err = s.encryptionManager.EncryptWithPlaintextKey(data, EncryptionTypeSSE_KMS, inputBucket, upload.PlaintextDataKey, upload.KMSKeyID, upload.SSEMetadata.EncryptedDataKey)
		} else {
			bucketEncryption, _ := store.buckets.GetEncryptionConfiguration(inputBucket)
			encResult, err = s.encryptionManager.Encrypt(data, EncryptionTypeSSE_KMS, bucketEncryption, inputBucket, inputKey, upload.KMSKeyID)
		}
	case s3store.SSETypeCustomer:
		if sseCustomerKey == "" {
			return nil, ErrInvalidSSECustomerKey
		}
		if upload.CustomerKeyMD5 != "" && upload.CustomerKeyMD5 != sseCustomerKeyMD5 {
			return nil, ErrInvalidSSECustomerKey
		}
		customerKey, parseErr := s.encryptionManager.ParseCustomerKey(sseCustomerKey, sseCustomerKeyMD5)
		if parseErr != nil {
			return nil, parseErr
		}
		encResult, err = s.encryptionManager.EncryptWithCustomerKey(data, EncryptionTypeSSE_C, nil, inputBucket, inputKey, "", customerKey)
	}
	if err != nil {
		return nil, err
	}
	if encResult != nil {
		result.Reader = bytes.NewReader(encResult.EncryptedData)
		result.EncryptedSize = int64(len(encResult.EncryptedData))
		result.ContentNonce = encResult.ContentNonce
		result.DataKey = encResult.EncryptedDataKey
	}
	return result, nil
}

func (s *S3Service) createMultipartUploadCore(ctx context.Context, stores *s3Stores, input *CreateMultipartUploadInput) (*CreateMultipartUploadOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	if err := validateStorageClass(input.StorageClass); err != nil {
		return nil, err
	}

	bucketEncryption, _ := stores.buckets.GetEncryptionConfiguration(input.Bucket)

	var sseType s3store.SSEType
	var kmsKeyID string
	var customerKeyMD5 string
	var effectiveEncryptionType EncryptionType

	if input.ServerSideEncryption != "" {
		effectiveEncryptionType = EncryptionType(input.ServerSideEncryption)
	} else if input.SSECustomerAlgorithm != "" {
		effectiveEncryptionType = EncryptionTypeSSE_C
	} else {
		effectiveEncryptionType = s.encryptionManager.DetermineEncryptionType(EncryptionTypeNone, bucketEncryption)
	}

	switch effectiveEncryptionType {
	case EncryptionTypeSSE_S3:
		sseType = s3store.SSETypeAES256
	case EncryptionTypeSSE_KMS:
		sseType = s3store.SSETypeKMS
		kmsKeyID = input.SSEKMSKeyId
		if kmsKeyID == "" && bucketEncryption != nil {
			kmsKeyID = bucketEncryption.KMSMasterKeyID
		}
	case EncryptionTypeSSE_DSSE_KMS:
		sseType = s3store.SSETypeDSSEKMS
		kmsKeyID = input.SSEKMSKeyId
		if kmsKeyID == "" && bucketEncryption != nil {
			kmsKeyID = bucketEncryption.KMSMasterKeyID
		}
	case EncryptionTypeSSE_C:
		sseType = s3store.SSETypeCustomer
		if input.SSECustomerKey != "" {
			_, err := s.encryptionManager.ParseCustomerKey(input.SSECustomerKey, input.SSECustomerKeyMD5)
			if err != nil {
				return nil, err
			}
			customerKeyMD5 = input.SSECustomerKeyMD5
		}
	}

	var sseMetadata *s3store.SSEObjectMetadata
	var plaintextDataKey []byte
	if sseType != "" && sseType != s3store.SSETypeCustomer {
		genKey, err := s.encryptionManager.GenerateKey(effectiveEncryptionType, bucketEncryption, input.Bucket, input.Key, kmsKeyID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate SSE key: %w", err)
		}
		if genKey != nil {
			sseMetadata = &s3store.SSEObjectMetadata{
				EncryptionType:   sseType,
				EncryptedDataKey: genKey.EncryptedDataKey,
				ContentNonce:     genKey.ContentNonce,
				KMSKeyID:         genKey.KMSKeyID,
			}
			plaintextDataKey = genKey.PlaintextKey
		}
	}

	sc := s3store.ObjectStorageClass(input.StorageClass)
	if sc == "" {
		sc = s3store.StorageClassStandard
	}

	acl, err := s.resolveUploadACL(ctx, stores, input.Bucket, input.ACLHeaders)
	if err != nil {
		return nil, err
	}

	upload, err := stores.objects.CreateMultipartUpload(ctx, input.Bucket, input.Key, input.ContentType, input.Metadata, sseType, kmsKeyID, customerKeyMD5, sseMetadata, plaintextDataKey, sc, acl)
	if err != nil {
		return nil, err
	}

	output := &CreateMultipartUploadOutput{
		Bucket:   input.Bucket,
		Key:      input.Key,
		UploadId: upload.UploadID,
	}

	if sseType != "" {
		output.ServerSideEncryption = string(sseType)
		if sseType == s3store.SSETypeKMS && kmsKeyID != "" {
			output.SSEKMSKeyId = kmsKeyID
		}
		if sseType == s3store.SSETypeCustomer {
			output.SSECustomerAlgorithm = "AES256"
		}
	}

	return output, nil
}

func (s *S3Service) uploadPartCore(ctx context.Context, stores *s3Stores, input *UploadPartInput) (*UploadPartOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	if err := validatePartNumber(input.PartNumber); err != nil {
		return nil, err
	}

	upload, err := stores.objects.GetMultipartUpload(input.UploadId)
	if err != nil {
		return nil, err
	}

	var reader io.Reader = input.Body
	var encryptedSize int64
	var plainSize int64
	var contentNonce, dataKey []byte

	if upload.SSEType != "" {
		data, err := io.ReadAll(input.Body)
		if err != nil {
			return nil, err
		}
		encResult, encErr := s.encryptPartData(data, upload, input.Bucket, input.Key, input.SSECustomerKey, input.SSECustomerKeyMD5, stores)
		if encErr != nil {
			return nil, encErr
		}
		reader = encResult.Reader
		encryptedSize = encResult.EncryptedSize
		plainSize = encResult.PlainSize
		contentNonce = encResult.ContentNonce
		dataKey = encResult.DataKey
	} else if input.ContentLength >= 0 {
		reader = input.Body
		plainSize = input.ContentLength
	} else {
		data, err := io.ReadAll(input.Body)
		if err != nil {
			return nil, err
		}
		plainSize = int64(len(data))
		reader = bytes.NewReader(data)
	}

	part, err := stores.objects.UploadPart(ctx, input.Bucket, input.Key, input.UploadId, input.PartNumber, reader, encryptedSize, plainSize, contentNonce, dataKey)
	if err != nil {
		return nil, err
	}

	output := &UploadPartOutput{
		ETag: formatETag(part.ETag),
	}
	if upload.SSEType == s3store.SSETypeCustomer {
		output.SSECustomerAlgorithm = "AES256"
	}

	return output, nil
}

func (s *S3Service) uploadPartCopyCore(ctx context.Context, stores *s3Stores, input *UploadPartCopyInput) (*UploadPartCopyOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}
	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}
	if err := validatePartNumber(input.PartNumber); err != nil {
		return nil, err
	}

	srcBucket, srcKey, srcVersionId, err := parseCopySource(input.CopySource)
	if err != nil {
		return nil, err
	}
	if input.CopySourceVersionId != "" {
		srcVersionId = input.CopySourceVersionId
	}

	if err := s.validateBucketExists(stores, srcBucket); err != nil {
		return nil, ErrInvalidCopySource
	}

	upload, err := stores.objects.GetMultipartUpload(input.UploadId)
	if err != nil {
		return nil, err
	}

	var srcObj *s3store.Object
	if srcVersionId != "" {
		srcObj, err = stores.objects.HeadWithVersion(ctx, srcBucket, srcKey, srcVersionId)
	} else {
		srcObj, err = stores.objects.GetMetadata(srcBucket, srcKey)
	}
	if err != nil {
		return nil, ErrInvalidCopySource
	}

	if err := checkCopySourcePreconditions(srcObj, input.CopySourceIfMatch, input.CopySourceIfNoneMatch, input.CopySourceIfModifiedSince, input.CopySourceIfUnmodifiedSince); err != nil {
		return nil, err
	}

	// A part copy reads the source object, so an archived source must have
	// a restored temporary copy available — the same rule CopyObject
	// enforces for its copy source.
	if isArchiveClass(srcObj.StorageClass) && !objectRestored(srcObj, time.Now()) {
		return nil, ErrObjectNotInActiveTier
	}

	if input.CopySourceRange == "" && srcObj.Size > maxCopyObjectSize {
		return nil, ErrEntityTooLarge
	}

	var data []byte
	if srcObj.SSEMetadata != nil || input.CopySourceSSECustomerKey != "" {
		coreResult, err := s.getObjectStreamCore(ctx, stores.objects, GetObjectStreamInput{
			Bucket:               srcBucket,
			Key:                  srcKey,
			VersionID:            srcVersionId,
			Range:                input.CopySourceRange,
			SSECustomerAlgorithm: input.CopySourceSSECustomerAlgo,
			SSECustomerKey:       input.CopySourceSSECustomerKey,
			SSECustomerKeyMD5:    input.CopySourceSSECustomerMD5,
		})
		if err != nil {
			return nil, err
		}
		data, err = io.ReadAll(coreResult.Body)
		coreResult.Body.Close()
		if err != nil {
			return nil, err
		}
	} else if input.CopySourceRange != "" {
		offset, length, rangeErr := parseCopyRange(input.CopySourceRange, srcObj.Size)
		if rangeErr != nil {
			return nil, rangeErr
		}
		var reader io.ReadCloser
		if srcVersionId != "" {
			reader, _, err = stores.objects.GetRangeWithVersion(ctx, srcBucket, srcKey, srcVersionId, offset, length)
		} else {
			reader, _, err = stores.objects.GetRange(ctx, srcBucket, srcKey, offset, length)
		}
		if err != nil {
			return nil, err
		}
		data, err = io.ReadAll(reader)
		reader.Close()
		if err != nil {
			return nil, err
		}
	} else {
		var reader io.ReadCloser
		if srcVersionId != "" {
			reader, _, err = stores.objects.GetWithVersion(ctx, srcBucket, srcKey, srcVersionId)
		} else {
			reader, _, err = stores.objects.Get(ctx, srcBucket, srcKey)
		}
		if err != nil {
			return nil, err
		}
		data, err = io.ReadAll(reader)
		reader.Close()
		if err != nil {
			return nil, err
		}
	}

	var reader io.Reader = bytes.NewReader(data)
	var encryptedSize int64
	var plainSize int64 = int64(len(data))
	var contentNonce, dataKey []byte

	if upload.SSEType != "" {
		encResult, encErr := s.encryptPartData(data, upload, input.Bucket, input.Key, input.SSECustomerKey, input.SSECustomerKeyMD5, stores)
		if encErr != nil {
			return nil, encErr
		}
		reader = encResult.Reader
		encryptedSize = encResult.EncryptedSize
		contentNonce = encResult.ContentNonce
		dataKey = encResult.DataKey
	}

	part, err := stores.objects.UploadPart(ctx, input.Bucket, input.Key, input.UploadId, input.PartNumber, reader, encryptedSize, plainSize, contentNonce, dataKey)
	if err != nil {
		return nil, err
	}

	return &UploadPartCopyOutput{
		CopyPartResult: &CopyPartResult{
			ETag:         formatETag(part.ETag),
			LastModified: part.LastModified,
		},
	}, nil
}

func (s *S3Service) listPartsCore(ctx context.Context, stores *s3Stores, input *ListPartsInput) (*ListPartsOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	maxParts := input.MaxParts

	partNumberMarker := 0
	if input.PartNumberMarker != "" {
		var err error
		partNumberMarker, err = strconv.Atoi(input.PartNumberMarker)
		if err != nil {
			return nil, NewInvalidArgumentError("Provided part-number-marker not an integer")
		}
	}

	parts, nextPartNumberMarker, isTruncated, err := stores.objects.ListParts(ctx, input.Bucket, input.Key, input.UploadId, partNumberMarker, maxParts)
	if err != nil {
		if errors.Is(err, s3store.ErrUploadNotFound) {
			return nil, ErrNoSuchUpload
		}
		return nil, err
	}

	var outputParts []*Part
	now := time.Now().UTC()
	for _, p := range parts {
		lastModified := p.LastModified
		if lastModified.IsZero() {
			lastModified = now
		}
		outputParts = append(outputParts, &Part{
			PartNumber:   p.PartNumber,
			ETag:         formatETag(p.ETag),
			Size:         p.Size,
			LastModified: lastModified,
		})
	}

	output := &ListPartsOutput{
		Bucket:       input.Bucket,
		Key:          input.Key,
		UploadId:     input.UploadId,
		Parts:        outputParts,
		MaxParts:     maxParts,
		StorageClass: "STANDARD",
		IsTruncated:  isTruncated,
	}
	if nextPartNumberMarker > 0 {
		output.NextPartNumberMarker = strconv.Itoa(nextPartNumberMarker)
	}
	return output, nil
}

func (s *S3Service) completeMultipartUploadCore(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *CompleteMultipartUploadInput) (*CompleteMultipartUploadOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if len(input.Parts) == 0 {
		return nil, NewInvalidArgumentError("at least one part is required")
	}

	for _, p := range input.Parts {
		if err := validatePartNumber(p.PartNumber); err != nil {
			return nil, err
		}
	}

	// Build the parts list. AWS accepts duplicate part numbers and uses
	// the value from the last occurrence, so we build a map keyed by
	// PartNumber that naturally keeps the last value.
	partsMap := make(map[int]string)
	var orderedPartNumbers []int
	for _, p := range input.Parts {
		if _, exists := partsMap[p.PartNumber]; !exists {
			orderedPartNumbers = append(orderedPartNumbers, p.PartNumber)
		}
		partsMap[p.PartNumber] = strings.Trim(p.ETag, "\"")
	}

	// AWS requires parts to be assembled in ascending part-number order.
	// Sort once so that both the ETag validation loop, the size check, and
	// the final parts slice passed to CompleteMultipartUpload are correct.
	sort.Ints(orderedPartNumbers)

	var parts []s3store.ObjectPart
	for _, pn := range orderedPartNumbers {
		parts = append(parts, s3store.ObjectPart{
			PartNumber: pn,
			ETag:       partsMap[pn],
		})
	}

	upload, err := stores.objects.GetMultipartUpload(input.UploadId)
	if err != nil {
		return nil, ErrNoSuchUpload
	}

	for _, p := range parts {
		idx, exists := upload.FindPart(p.PartNumber)
		if !exists {
			return nil, ErrInvalidPart
		}
		expectedETag := strings.Trim(upload.Parts[idx].ETag, "\"")
		if p.ETag != expectedETag {
			return nil, ErrInvalidPart
		}
	}

	// Validate part sizes: all parts except the last must be at least
	// minPartSize (5 MiB). The last part (highest part number) may be
	// smaller.
	if len(orderedPartNumbers) > 0 {
		maxPartNum := orderedPartNumbers[len(orderedPartNumbers)-1]
		for _, pn := range orderedPartNumbers {
			if pn == maxPartNum {
				continue
			}
			idx, _ := upload.FindPart(pn)
			if idx >= 0 && upload.Parts[idx].Size < int64(minPartSize) {
				return nil, ErrEntityTooSmall
			}
		}
	}

	obj, err := stores.objects.CompleteMultipartUpload(ctx, input.Bucket, input.Key, input.UploadId, parts)
	if err != nil {
		return nil, err
	}

	// The ACL requested by the CreateMultipartUpload headers is applied to
	// the completed object.
	if upload.ACL != nil {
		obj.ACL = upload.ACL
		if err := stores.objects.SetACL(input.Bucket, input.Key, upload.ACL); err != nil {
			return nil, err
		}
	}

	s.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectCreatedCompleteMultipartUpload)
	s.launchObjectReplication(reqCtx, stores, input.Bucket, input.Key, obj)

	output := &CompleteMultipartUploadOutput{
		Location:  fmt.Sprintf("http://%s.s3.amazonaws.com/%s", input.Bucket, input.Key),
		Bucket:    input.Bucket,
		Key:       input.Key,
		ETag:      formatETag(obj.ETag),
		VersionId: obj.VersionID,
	}

	if obj.SSEMetadata != nil {
		output.ServerSideEncryption = string(obj.SSEMetadata.EncryptionType)
		if obj.SSEMetadata.KMSKeyID != "" {
			output.SSEKMSKeyId = obj.SSEMetadata.KMSKeyID
		}
	}

	return output, nil
}

func (s *S3Service) abortMultipartUploadCore(ctx context.Context, stores *s3Stores, input *AbortMultipartUploadInput) error {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	return stores.objects.AbortMultipartUpload(ctx, input.Bucket, input.Key, input.UploadId)
}

func (s *S3Service) listMultipartUploadsCore(stores *s3Stores, input *ListMultipartUploadsInput) (*ListMultipartUploadsOutput, error) {
	if err := s.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	result, err := stores.objects.ListMultipartUploads(input.Bucket, input.Prefix, input.KeyMarker, input.UploadIdMarker, input.MaxUploads)
	if err != nil {
		return nil, err
	}

	var uploads []*Upload
	for _, u := range result.Uploads {
		uploads = append(uploads, &Upload{
			Key:          u.Key,
			UploadId:     u.UploadID,
			Initiated:    u.Initiated,
			StorageClass: string(u.StorageClass),
			Initiator: &Owner{
				ID:          u.Initiator,
				DisplayName: u.Initiator,
			},
			Owner: &Owner{
				ID:          u.Owner,
				DisplayName: u.Owner,
			},
		})
	}

	return &ListMultipartUploadsOutput{
		Bucket:             input.Bucket,
		KeyMarker:          input.KeyMarker,
		UploadIdMarker:     input.UploadIdMarker,
		NextKeyMarker:      result.NextKeyMarker,
		NextUploadIdMarker: result.NextUploadIDMarker,
		MaxUploads:         input.MaxUploads,
		IsTruncated:        result.IsTruncated,
		Prefix:             input.Prefix,
		Delimiter:          input.Delimiter,
		Uploads:            uploads,
	}, nil
}
