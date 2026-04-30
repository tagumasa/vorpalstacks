package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// CreateMultipartUploadInput contains the parameters for initiating a multipart upload.
// Bucket is the name of the S3 bucket.
// Key is the object key to upload.
// ContentType, ContentEncoding, ContentDisposition, CacheControl specify object metadata.
// Metadata is custom key-value pairs.
// ServerSideEncryption, SSEKMSKeyId, SSECustomerAlgorithm, SSECustomerKey, SSECustomerKeyMD5 specify encryption settings.
type CreateMultipartUploadInput struct {
	Bucket               string
	Key                  string
	ContentType          string
	ContentEncoding      string
	ContentDisposition   string
	CacheControl         string
	Metadata             map[string]string
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
	SSECustomerKey       string
	SSECustomerKeyMD5    string
}

// CreateMultipartUploadOutput contains the result of initiating a multipart upload.
// Bucket is the target bucket name.
// Key is the target object key.
// UploadId is the unique identifier for this multipart upload (required for subsequent part uploads).
// ServerSideEncryption, SSEKMSKeyId, SSECustomerAlgorithm contain encryption settings.
type CreateMultipartUploadOutput struct {
	Bucket               string `xml:"Bucket"`
	Key                  string `xml:"Key"`
	UploadId             string `xml:"UploadId"`
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
}

// CreateMultipartUpload initiates a multipart upload for an object.
// Returns an UploadId that must be used for subsequent UploadPart and CompleteMultipartUpload calls.
func (o *ObjectOperations) CreateMultipartUpload(ctx context.Context, reqCtx *request.RequestContext, input *CreateMultipartUploadInput) (*CreateMultipartUploadOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	bucketEncryption, _ := store.buckets.GetEncryptionConfiguration(input.Bucket)

	var sseType s3store.SSEType
	var kmsKeyID string
	var customerKeyMD5 string
	var effectiveEncryptionType EncryptionType

	if input.ServerSideEncryption != "" {
		effectiveEncryptionType = EncryptionType(input.ServerSideEncryption)
	} else if input.SSECustomerAlgorithm != "" {
		effectiveEncryptionType = EncryptionTypeSSE_C
	} else {
		effectiveEncryptionType = o.svc.encryptionManager.DetermineEncryptionType(EncryptionTypeNone, bucketEncryption)
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
	case EncryptionTypeSSE_C:
		sseType = s3store.SSETypeCustomer
		if input.SSECustomerKey != "" {
			_, err := o.svc.encryptionManager.ParseCustomerKey(input.SSECustomerKey, input.SSECustomerKeyMD5)
			if err != nil {
				return nil, err
			}
			customerKeyMD5 = input.SSECustomerKeyMD5
		}
	}

	var sseMetadata *s3store.SSEObjectMetadata
	var plaintextDataKey []byte
	if sseType != "" && sseType != s3store.SSETypeCustomer {
		genKey, err := o.svc.encryptionManager.GenerateKey(effectiveEncryptionType, bucketEncryption, input.Bucket, input.Key, kmsKeyID)
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

	upload, err := store.objects.CreateMultipartUpload(ctx, input.Bucket, input.Key, input.ContentType, input.Metadata, sseType, kmsKeyID, customerKeyMD5, sseMetadata, plaintextDataKey)
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

// UploadPartInput contains the parameters for uploading a part in a multipart upload.
// Bucket is the name of the S3 bucket.
// Key is the object key.
// UploadId is the multipart upload identifier returned by CreateMultipartUpload.
// PartNumber is the part's sequence number (1-10000).
// Body is the part data.
// SSECustomerAlgorithm, SSECustomerKey, SSECustomerKeyMD5 specify customer-provided encryption.
type UploadPartInput struct {
	Bucket               string
	Key                  string
	UploadId             string
	PartNumber           int
	Body                 io.Reader
	SSECustomerAlgorithm string
	SSECustomerKey       string
	SSECustomerKeyMD5    string
}

// UploadPartOutput contains the result of uploading a part.
// ETag is the entity tag of the uploaded part.
// SSECustomerAlgorithm is set when using customer-provided encryption.
type UploadPartOutput struct {
	ETag                 string
	SSECustomerAlgorithm string
}

// UploadPart uploads a single part of a multipart upload.
// The part number must be unique within the upload and between 1 and 10000.
func (o *ObjectOperations) UploadPart(ctx context.Context, reqCtx *request.RequestContext, input *UploadPartInput) (*UploadPartOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	if err := validatePartNumber(input.PartNumber); err != nil {
		return nil, err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	upload, err := store.objects.GetMultipartUpload(input.UploadId)
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
		plainSize = int64(len(data))

		var encResult *EncryptionResult
		switch upload.SSEType {
		case s3store.SSETypeAES256:
			if upload.PlaintextDataKey != nil && upload.SSEMetadata != nil {
				encResult, err = o.svc.encryptionManager.EncryptWithPlaintextKey(data, EncryptionTypeSSE_S3, input.Bucket, upload.PlaintextDataKey, "", nil)
			} else {
				encResult, err = o.svc.encryptionManager.Encrypt(data, EncryptionTypeSSE_S3, nil, input.Bucket, input.Key, "")
			}
		case s3store.SSETypeKMS:
			if upload.PlaintextDataKey != nil && upload.SSEMetadata != nil {
				encResult, err = o.svc.encryptionManager.EncryptWithPlaintextKey(data, EncryptionTypeSSE_KMS, input.Bucket, upload.PlaintextDataKey, upload.KMSKeyID, upload.SSEMetadata.EncryptedDataKey)
			} else {
				bucketEncryption, _ := store.buckets.GetEncryptionConfiguration(input.Bucket)
				encResult, err = o.svc.encryptionManager.Encrypt(data, EncryptionTypeSSE_KMS, bucketEncryption, input.Bucket, input.Key, upload.KMSKeyID)
			}
		case s3store.SSETypeCustomer:
			if input.SSECustomerKey == "" {
				return nil, ErrInvalidSSECustomerKey
			}
			if upload.CustomerKeyMD5 != "" && upload.CustomerKeyMD5 != input.SSECustomerKeyMD5 {
				return nil, ErrInvalidSSECustomerKey
			}
			customerKey, parseErr := o.svc.encryptionManager.ParseCustomerKey(input.SSECustomerKey, input.SSECustomerKeyMD5)
			if parseErr != nil {
				return nil, parseErr
			}
			encResult, err = o.svc.encryptionManager.EncryptWithCustomerKey(data, EncryptionTypeSSE_C, nil, input.Bucket, input.Key, "", customerKey)
		}
		if err != nil {
			return nil, err
		}
		if encResult != nil {
			reader = bytes.NewReader(encResult.EncryptedData)
			encryptedSize = int64(len(encResult.EncryptedData))
			contentNonce = encResult.ContentNonce
			dataKey = encResult.EncryptedDataKey
		}
	}

	part, err := store.objects.UploadPart(ctx, input.Bucket, input.Key, input.UploadId, input.PartNumber, reader, encryptedSize, plainSize, contentNonce, dataKey)
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

// UploadPartCopyInput contains the parameters for uploading a part by copying from an existing object.
type UploadPartCopyInput struct {
	Bucket                    string
	Key                       string
	UploadId                  string
	PartNumber                int
	CopySource                string
	CopySourceRange           string
	CopySourceVersionId       string
	CopySourceSSECustomerAlgo string
	CopySourceSSECustomerKey  string
	CopySourceSSECustomerMD5  string
	SSECustomerAlgorithm      string
	SSECustomerKey            string
	SSECustomerKeyMD5         string
}

// UploadPartCopyOutput contains the result of an UploadPartCopy operation.
type UploadPartCopyOutput struct {
	CopyPartResult *CopyPartResult `xml:"CopyPartResult"`
}

// CopyPartResult contains the ETag and last modified time of a copied part.
type CopyPartResult struct {
	ETag         string    `xml:"ETag"`
	LastModified time.Time `xml:"LastModified"`
}

// UploadPartCopy uploads a part by copying bytes from an existing S3 object.
func (o *ObjectOperations) UploadPartCopy(ctx context.Context, reqCtx *request.RequestContext, input *UploadPartCopyInput) (*UploadPartCopyOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
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

	if err := o.validateBucketExists(reqCtx, srcBucket); err != nil {
		return nil, ErrInvalidCopySource
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	upload, err := store.objects.GetMultipartUpload(input.UploadId)
	if err != nil {
		return nil, err
	}

	var srcObj *s3store.Object
	if srcVersionId != "" {
		srcObj, err = store.objects.HeadWithVersion(ctx, srcBucket, srcKey, srcVersionId)
	} else {
		srcObj, err = store.objects.GetMetadata(srcBucket, srcKey)
	}
	if err != nil {
		return nil, ErrInvalidCopySource
	}

	var data []byte
	if srcObj.SSEMetadata != nil || input.CopySourceSSECustomerKey != "" {
		getInput := &GetObjectInput{
			Bucket:               srcBucket,
			Key:                  srcKey,
			VersionId:            srcVersionId,
			SSECustomerAlgorithm: input.CopySourceSSECustomerAlgo,
			SSECustomerKey:       input.CopySourceSSECustomerKey,
			SSECustomerKeyMD5:    input.CopySourceSSECustomerMD5,
		}
		getOutput, err := o.GetObject(ctx, reqCtx, getInput)
		if err != nil {
			return nil, err
		}
		data, err = io.ReadAll(getOutput.Body)
		getOutput.Body.Close()
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
			reader, _, err = store.objects.GetRangeWithVersion(ctx, srcBucket, srcKey, srcVersionId, offset, length)
		} else {
			reader, _, err = store.objects.GetRange(ctx, srcBucket, srcKey, offset, length)
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
			reader, _, err = store.objects.GetWithVersion(ctx, srcBucket, srcKey, srcVersionId)
		} else {
			reader, _, err = store.objects.Get(ctx, srcBucket, srcKey)
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
		var encResult *EncryptionResult
		switch upload.SSEType {
		case s3store.SSETypeAES256:
			if upload.PlaintextDataKey != nil && upload.SSEMetadata != nil {
				encResult, err = o.svc.encryptionManager.EncryptWithPlaintextKey(data, EncryptionTypeSSE_S3, input.Bucket, upload.PlaintextDataKey, "", nil)
			} else {
				encResult, err = o.svc.encryptionManager.Encrypt(data, EncryptionTypeSSE_S3, nil, input.Bucket, input.Key, "")
			}
		case s3store.SSETypeKMS:
			if upload.PlaintextDataKey != nil && upload.SSEMetadata != nil {
				encResult, err = o.svc.encryptionManager.EncryptWithPlaintextKey(data, EncryptionTypeSSE_KMS, input.Bucket, upload.PlaintextDataKey, upload.KMSKeyID, upload.SSEMetadata.EncryptedDataKey)
			} else {
				bucketEncryption, _ := store.buckets.GetEncryptionConfiguration(input.Bucket)
				encResult, err = o.svc.encryptionManager.Encrypt(data, EncryptionTypeSSE_KMS, bucketEncryption, input.Bucket, input.Key, upload.KMSKeyID)
			}
		case s3store.SSETypeCustomer:
			if input.SSECustomerKey == "" {
				return nil, ErrInvalidSSECustomerKey
			}
			if upload.CustomerKeyMD5 != "" && upload.CustomerKeyMD5 != input.SSECustomerKeyMD5 {
				return nil, ErrInvalidSSECustomerKey
			}
			customerKey, parseErr := o.svc.encryptionManager.ParseCustomerKey(input.SSECustomerKey, input.SSECustomerKeyMD5)
			if parseErr != nil {
				return nil, parseErr
			}
			encResult, err = o.svc.encryptionManager.EncryptWithCustomerKey(data, EncryptionTypeSSE_C, nil, input.Bucket, input.Key, "", customerKey)
		}
		if err != nil {
			return nil, err
		}
		if encResult != nil {
			reader = bytes.NewReader(encResult.EncryptedData)
			encryptedSize = int64(len(encResult.EncryptedData))
			contentNonce = encResult.ContentNonce
			dataKey = encResult.EncryptedDataKey
		}
	}

	part, err := store.objects.UploadPart(ctx, input.Bucket, input.Key, input.UploadId, input.PartNumber, reader, encryptedSize, plainSize, contentNonce, dataKey)
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

func parseCopyRange(rangeHeader string, objectSize int64) (offset, length int64, err error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, fmt.Errorf("invalid range header")
	}
	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.SplitN(rangeSpec, "-", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range header")
	}
	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	var start, end int64
	if startStr != "" {
		start, err = strconv.ParseInt(startStr, 10, 64)
		if err != nil || start < 0 {
			return 0, 0, fmt.Errorf("invalid range start")
		}
	}
	if endStr != "" {
		end, err = strconv.ParseInt(endStr, 10, 64)
		if err != nil || end < 0 {
			return 0, 0, fmt.Errorf("invalid range end")
		}
	} else {
		end = objectSize - 1
	}

	if start > end {
		return 0, 0, fmt.Errorf("invalid range: start > end")
	}
	if end >= objectSize {
		end = objectSize - 1
	}

	return start, end - start + 1, nil
}

// Bucket is the name of the S3 bucket.
// Key is the object key.
// UploadId is the multipart upload identifier.
// MaxParts limits the number of parts returned.
// PartNumberMarker specifies where to start in the part list.
type ListPartsInput struct {
	Bucket           string
	Key              string
	UploadId         string
	MaxParts         int
	PartNumberMarker string
}

// ListPartsOutput contains the result of listing uploaded parts.
// Bucket, Key, UploadId identify the multipart upload.
// Parts contains the list of uploaded parts.
// NextPartNumberMarker is used for pagination.
// IsTruncated indicates if more parts exist.
// MaxParts, PartNumberMarker, StorageClass contain request and response metadata.
type ListPartsOutput struct {
	Bucket               string
	Key                  string
	UploadId             string
	Parts                []*Part
	NextPartNumberMarker string
	IsTruncated          bool
	MaxParts             int
	PartNumberMarker     string
	StorageClass         string
}

// Part represents an uploaded part in a multipart upload.
// PartNumber is the part's sequence number.
// ETag is the entity tag of the part.
// Size is the part size in bytes.
// LastModified is when the part was uploaded.
type Part struct {
	PartNumber   int       `xml:"PartNumber"`
	ETag         string    `xml:"ETag"`
	Size         int64     `xml:"Size"`
	LastModified time.Time `xml:"LastModified"`
}

// ToXML serialises the ListPartsResult to XML format for S3 API response.
func (o *ListPartsOutput) ToXML() string {
	var result strings.Builder
	result.WriteString(`<ListPartsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">`)
	result.WriteString(`<Bucket>`)
	result.WriteString(xmlEscape(o.Bucket))
	result.WriteString(`</Bucket><Key>`)
	result.WriteString(xmlEscape(o.Key))
	result.WriteString(`</Key><UploadId>`)
	result.WriteString(xmlEscape(o.UploadId))
	result.WriteString(`</UploadId><StorageClass>`)
	result.WriteString(o.StorageClass)
	result.WriteString(`</StorageClass>`)
	if o.PartNumberMarker != "" {
		result.WriteString(`<PartNumberMarker>`)
		result.WriteString(xmlEscape(o.PartNumberMarker))
		result.WriteString(`</PartNumberMarker>`)
	}
	result.WriteString(`<MaxParts>`)
	result.WriteString(strconv.Itoa(o.MaxParts))
	result.WriteString(`</MaxParts><IsTruncated>`)
	result.WriteString(strconv.FormatBool(o.IsTruncated))
	result.WriteString(`</IsTruncated>`)
	if o.NextPartNumberMarker != "" {
		result.WriteString(`<NextPartNumberMarker>`)
		result.WriteString(xmlEscape(o.NextPartNumberMarker))
		result.WriteString(`</NextPartNumberMarker>`)
	}
	for _, p := range o.Parts {
		result.WriteString(`<Part><PartNumber>`)
		result.WriteString(strconv.Itoa(p.PartNumber))
		result.WriteString(`</PartNumber><ETag>`)
		result.WriteString(xmlEscape(p.ETag))
		result.WriteString(`</ETag><Size>`)
		result.WriteString(strconv.FormatInt(p.Size, 10))
		result.WriteString(`</Size><LastModified>`)
		result.WriteString(p.LastModified.Format(time.RFC3339))
		result.WriteString(`</LastModified></Part>`)
	}
	result.WriteString(`</ListPartsResult>`)
	return result.String()
}

// ListParts returns the list of parts that have been uploaded for a multipart upload.
func (o *ObjectOperations) ListParts(ctx context.Context, reqCtx *request.RequestContext, input *ListPartsInput) (*ListPartsOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	maxParts := input.MaxParts
	if maxParts <= 0 {
		maxParts = 1000
	}

	partNumberMarker := 0
	if input.PartNumberMarker != "" {
		partNumberMarker, _ = strconv.Atoi(input.PartNumberMarker)
	}

	parts, nextPartNumberMarker, isTruncated, err := store.objects.ListParts(ctx, input.Bucket, input.Key, input.UploadId, partNumberMarker, maxParts)
	if err != nil {
		return nil, ErrNoSuchUpload
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

// CompletedPart represents a part to be used when completing a multipart upload.
// PartNumber is the part's sequence number.
// ETag is the entity tag returned by UploadPart.
type CompletedPart struct {
	PartNumber int    `xml:"PartNumber"`
	ETag       string `xml:"ETag"`
}

// CompleteMultipartUploadInput contains the parameters for completing a multipart upload.
// Bucket is the name of the S3 bucket.
// Key is the object key.
// UploadId is the multipart upload identifier.
// Parts is the list of uploaded parts in the order they should be assembled.
type CompleteMultipartUploadInput struct {
	Bucket   string
	Key      string
	UploadId string
	Parts    []CompletedPart
}

// CompleteMultipartUploadOutput contains the result of completing a multipart upload.
// Location is the URI of the newly created object.
// Bucket, Key identify the completed object.
// ETag is the entity tag of the completed object.
// VersionId is the version ID if versioning is enabled.
// ServerSideEncryption, SSEKMSKeyId contain encryption settings.
type CompleteMultipartUploadOutput struct {
	Location             string `xml:"Location"`
	Bucket               string `xml:"Bucket"`
	Key                  string `xml:"Key"`
	ETag                 string `xml:"ETag"`
	VersionId            string `xml:"VersionId,omitempty"`
	ServerSideEncryption string
	SSEKMSKeyId          string
}

// CompleteMultipartUpload assembles the uploaded parts into a complete object.
// At least one part is required. Parts are assembled in the order specified.
func (o *ObjectOperations) CompleteMultipartUpload(ctx context.Context, reqCtx *request.RequestContext, input *CompleteMultipartUploadInput) (*CompleteMultipartUploadOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	if len(input.Parts) == 0 {
		return nil, fmt.Errorf("at least one part is required")
	}

	for _, p := range input.Parts {
		if err := validatePartNumber(p.PartNumber); err != nil {
			return nil, err
		}
	}

	var parts []s3store.ObjectPart
	for _, p := range input.Parts {
		etag := strings.Trim(p.ETag, "\"")
		parts = append(parts, s3store.ObjectPart{
			PartNumber: p.PartNumber,
			ETag:       etag,
		})
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	obj, err := store.objects.CompleteMultipartUpload(ctx, input.Bucket, input.Key, input.UploadId, parts)
	if err != nil {
		return nil, err
	}

	o.svc.publishObjectNotification(ctx, reqCtx, input.Bucket, input.Key, obj.Size, obj.VersionID, obj.ETag, eventbus.S3ObjectCreatedCompleteMultipartUpload)

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

// AbortMultipartUploadInput contains the parameters for aborting a multipart upload.
// Bucket is the name of the S3 bucket.
// Key is the object key.
// UploadId is the multipart upload identifier.
type AbortMultipartUploadInput struct {
	Bucket   string
	Key      string
	UploadId string
}

// AbortMultipartUpload aborts a multipart upload and removes all uploaded parts.
// After aborting, the UploadId is no longer valid.
func (o *ObjectOperations) AbortMultipartUpload(ctx context.Context, reqCtx *request.RequestContext, input *AbortMultipartUploadInput) error {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return err
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return err
	}

	return store.objects.AbortMultipartUpload(ctx, input.Bucket, input.Key, input.UploadId)
}

// ListMultipartUploadsInput contains the parameters for listing multipart uploads.
// Bucket is the name of the S3 bucket.
// Delimiter groups keys by common prefix.
// Prefix filters keys by prefix.
// KeyMarker, UploadIdMarker specify where to start (for pagination).
// MaxUploads limits the number of uploads returned.
type ListMultipartUploadsInput struct {
	Bucket         string
	Delimiter      string
	Prefix         string
	KeyMarker      string
	UploadIdMarker string
	MaxUploads     int
}

// ListMultipartUploadsOutput contains the result of listing multipart uploads.
// Bucket is the bucket name.
// KeyMarker, UploadIdMarker are the starting points for this response.
// NextKeyMarker, NextUploadIdMarker are the starting points for the next response.
// MaxUploads is the maximum number of uploads requested.
// IsTruncated indicates if more uploads exist.
// Prefix, Delimiter are the request parameters.
// Uploads contains the list of in-progress multipart uploads.
// CommonPrefixes contains grouped keys (when delimiter is specified).
type ListMultipartUploadsOutput struct {
	Bucket             string
	KeyMarker          string
	UploadIdMarker     string
	NextKeyMarker      string
	NextUploadIdMarker string
	MaxUploads         int
	IsTruncated        bool
	Prefix             string
	Delimiter          string
	Uploads            []*Upload
	CommonPrefixes     []CommonPrefix
}

// Upload represents an in-progress multipart upload.
// Key is the object key.
// UploadId is the multipart upload identifier.
// Initiator is who initiated the upload.
// Owner is the object owner.
// StorageClass is the storage class (e.g., STANDARD).
// Initiated is when the upload was initiated.
type Upload struct {
	Key          string    `xml:"Key"`
	UploadId     string    `xml:"UploadId"`
	Initiator    *Owner    `xml:"Initiator"`
	Owner        *Owner    `xml:"Owner"`
	StorageClass string    `xml:"StorageClass"`
	Initiated    time.Time `xml:"Initiated"`
}

// ToXML serialises the ListMultipartUploadsResult to XML format for S3 API response.
func (o *ListMultipartUploadsOutput) ToXML() string {
	var result strings.Builder
	result.WriteString(`<ListMultipartUploadsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">`)
	result.WriteString(`<Bucket>`)
	result.WriteString(xmlEscape(o.Bucket))
	result.WriteString(`</Bucket>`)
	if o.KeyMarker != "" {
		result.WriteString(`<KeyMarker>`)
		result.WriteString(xmlEscape(o.KeyMarker))
		result.WriteString(`</KeyMarker>`)
	}
	if o.UploadIdMarker != "" {
		result.WriteString(`<UploadIdMarker>`)
		result.WriteString(xmlEscape(o.UploadIdMarker))
		result.WriteString(`</UploadIdMarker>`)
	}
	if o.NextKeyMarker != "" {
		result.WriteString(`<NextKeyMarker>`)
		result.WriteString(xmlEscape(o.NextKeyMarker))
		result.WriteString(`</NextKeyMarker>`)
	}
	if o.NextUploadIdMarker != "" {
		result.WriteString(`<NextUploadIdMarker>`)
		result.WriteString(xmlEscape(o.NextUploadIdMarker))
		result.WriteString(`</NextUploadIdMarker>`)
	}
	result.WriteString(`<MaxUploads>`)
	result.WriteString(strconv.Itoa(o.MaxUploads))
	result.WriteString(`</MaxUploads><IsTruncated>`)
	result.WriteString(strconv.FormatBool(o.IsTruncated))
	result.WriteString(`</IsTruncated>`)
	if o.Prefix != "" {
		result.WriteString(`<Prefix>`)
		result.WriteString(xmlEscape(o.Prefix))
		result.WriteString(`</Prefix>`)
	}
	if o.Delimiter != "" {
		result.WriteString(`<Delimiter>`)
		result.WriteString(xmlEscape(o.Delimiter))
		result.WriteString(`</Delimiter>`)
	}
	for _, u := range o.Uploads {
		result.WriteString(`<Upload><Key>`)
		result.WriteString(xmlEscape(u.Key))
		result.WriteString(`</Key><UploadId>`)
		result.WriteString(xmlEscape(u.UploadId))
		result.WriteString(`</UploadId><StorageClass>`)
		result.WriteString(u.StorageClass)
		result.WriteString(`</StorageClass><Initiated>`)
		result.WriteString(u.Initiated.Format(time.RFC3339))
		result.WriteString(`</Initiated></Upload>`)
	}
	for _, p := range o.CommonPrefixes {
		result.WriteString(`<CommonPrefixes><Prefix>`)
		result.WriteString(xmlEscape(p.Prefix))
		result.WriteString(`</Prefix></CommonPrefixes>`)
	}
	result.WriteString(`</ListMultipartUploadsResult>`)
	return result.String()
}

// ListMultipartUploads lists the in-progress multipart uploads for a bucket.
// Returns uploads that have been initiated but not yet completed or aborted.
func (o *ObjectOperations) ListMultipartUploads(ctx context.Context, reqCtx *request.RequestContext, input *ListMultipartUploadsInput) (*ListMultipartUploadsOutput, error) {
	if err := o.validateBucketExists(reqCtx, input.Bucket); err != nil {
		return nil, err
	}

	if input.MaxUploads <= 0 {
		input.MaxUploads = 1000
	}

	store, err := o.svc.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.objects.ListMultipartUploads(input.Bucket, input.Prefix, input.KeyMarker, input.UploadIdMarker, input.MaxUploads)
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
