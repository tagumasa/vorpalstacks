package s3

import (
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
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
	StorageClass         string
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
	SSECustomerKey       string
	SSECustomerKeyMD5    string
	ACLHeaders           aclHeaders
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
func (o *ObjectOperations) CreateMultipartUpload(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *CreateMultipartUploadInput) (*CreateMultipartUploadOutput, error) {
	return o.svc.createMultipartUploadCore(ctx, stores, input)
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
	ContentLength        int64
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
func (o *ObjectOperations) UploadPart(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *UploadPartInput) (*UploadPartOutput, error) {
	return o.svc.uploadPartCore(ctx, stores, input)
}

// UploadPartCopyInput contains the parameters for uploading a part by copying from an existing object.
type UploadPartCopyInput struct {
	Bucket                      string
	Key                         string
	UploadId                    string
	PartNumber                  int
	CopySource                  string
	CopySourceRange             string
	CopySourceVersionId         string
	CopySourceIfMatch           string
	CopySourceIfNoneMatch       string
	CopySourceIfModifiedSince   *time.Time
	CopySourceIfUnmodifiedSince *time.Time
	CopySourceSSECustomerAlgo   string
	CopySourceSSECustomerKey    string
	CopySourceSSECustomerMD5    string
	SSECustomerAlgorithm        string
	SSECustomerKey              string
	SSECustomerKeyMD5           string
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
func (o *ObjectOperations) UploadPartCopy(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *UploadPartCopyInput) (*UploadPartCopyOutput, error) {
	return o.svc.uploadPartCopyCore(ctx, stores, input)
}

func parseCopyRange(rangeHeader string, objectSize int64) (offset, length int64, err error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, NewInvalidArgumentError("invalid range header")
	}
	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.SplitN(rangeSpec, "-", 2)
	if len(parts) != 2 {
		return 0, 0, NewInvalidArgumentError("invalid range header")
	}
	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	var start, end int64
	if startStr != "" {
		start, err = strconv.ParseInt(startStr, 10, 64)
		if err != nil || start < 0 {
			return 0, 0, NewInvalidArgumentError("invalid range start")
		}
	}
	if endStr != "" {
		end, err = strconv.ParseInt(endStr, 10, 64)
		if err != nil || end < 0 {
			return 0, 0, NewInvalidArgumentError("invalid range end")
		}
	} else {
		end = objectSize - 1
	}

	// RFC 7233 §4.2: if last-byte-pos is absent or greater than or equal
	// to the representation length, it is reduced to one less than the
	// length.  AWS S3 follows this for x-amz-copy-source-range.
	if end >= objectSize {
		end = objectSize - 1
	}

	// After clamping, if start exceeds the (now clamped) end, the range
	// is genuinely unsatisfiable — the requested start byte does not
	// exist in the object.  This returns HTTP 416, matching AWS.
	if start > end {
		return 0, 0, ErrInvalidRange
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
func (o *ObjectOperations) ListParts(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *ListPartsInput) (*ListPartsOutput, error) {
	return o.svc.listPartsCore(ctx, stores, input)
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
func (o *ObjectOperations) CompleteMultipartUpload(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *CompleteMultipartUploadInput) (*CompleteMultipartUploadOutput, error) {
	return o.svc.completeMultipartUploadCore(ctx, reqCtx, stores, input)
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
func (o *ObjectOperations) AbortMultipartUpload(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *AbortMultipartUploadInput) error {
	return o.svc.abortMultipartUploadCore(ctx, stores, input)
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
	writeCommonPrefixesXML(&result, o.CommonPrefixes, "")
	result.WriteString(`</ListMultipartUploadsResult>`)
	return result.String()
}

// ListMultipartUploads lists the in-progress multipart uploads for a bucket.
// Returns uploads that have been initiated but not yet completed or aborted.
func (o *ObjectOperations) ListMultipartUploads(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *ListMultipartUploadsInput) (*ListMultipartUploadsOutput, error) {
	return o.svc.listMultipartUploadsCore(stores, input)
}
