package s3

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// GetObjectInput contains the input parameters for the GetObject operation.
type GetObjectInput struct {
	Bucket               string
	Key                  string
	VersionId            string
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

// GetObjectOutput contains the output from the GetObject operation.
type GetObjectOutput struct {
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
	VersionId            string
	Restore              string
	ContentRange         string
	IsPartial            bool
	AcceptRanges         string
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
	SSECustomerKeyMD5    string
	ReplicationStatus    string
	PartsCount           int32
}

// GetObject retrieves an object from S3.
func (o *ObjectOperations) GetObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *GetObjectInput) (*GetObjectOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	coreResult, err := o.svc.getObjectStreamCore(ctx, stores.objects, GetObjectStreamInput{
		Bucket:               input.Bucket,
		Key:                  input.Key,
		VersionID:            input.VersionId,
		IfMatch:              input.IfMatch,
		IfNoneMatch:          input.IfNoneMatch,
		IfModifiedSince:      input.IfModifiedSince,
		IfUnmodifiedSince:    input.IfUnmodifiedSince,
		Range:                input.Range,
		PartNumber:           input.PartNumber,
		SSECustomerAlgorithm: input.SSECustomerAlgorithm,
		SSECustomerKey:       input.SSECustomerKey,
		SSECustomerKeyMD5:    input.SSECustomerKeyMD5,
	})
	if err != nil {
		return nil, err
	}

	return &GetObjectOutput{
		Body:                 coreResult.Body,
		ContentLength:        coreResult.ContentLength,
		ContentType:          coreResult.ContentType,
		ContentEncoding:      coreResult.ContentEncoding,
		ContentLanguage:      coreResult.ContentLanguage,
		ContentDisposition:   coreResult.ContentDisposition,
		CacheControl:         coreResult.CacheControl,
		ETag:                 coreResult.ETag,
		LastModified:         coreResult.LastModified,
		Metadata:             coreResult.Metadata,
		StorageClass:         coreResult.StorageClass,
		VersionId:            coreResult.VersionID,
		Restore:              coreResult.Restore,
		ContentRange:         coreResult.ContentRange,
		IsPartial:            coreResult.IsPartial,
		AcceptRanges:         coreResult.AcceptRanges,
		ServerSideEncryption: coreResult.ServerSideEncryption,
		SSEKMSKeyId:          coreResult.SSEKMSKeyId,
		SSECustomerAlgorithm: coreResult.SSECustomerAlgorithm,
		SSECustomerKeyMD5:    coreResult.SSECustomerKeyMD5,
		ReplicationStatus:    coreResult.ReplicationStatus,
		PartsCount:           coreResult.PartsCount,
	}, nil
}

// HeadObjectInput contains the input parameters for the HeadObject operation.
type HeadObjectInput struct {
	Bucket               string
	Key                  string
	VersionId            string
	IfMatch              string
	IfNoneMatch          string
	IfModifiedSince      *time.Time
	IfUnmodifiedSince    *time.Time
	Range                string
	PartNumber           int
	SSECustomerAlgorithm string
	SSECustomerKey       string
	SSECustomerKeyMD5    string
	ReplicationStatus    string
}

// HeadObjectOutput contains the output from the HeadObject operation.
type HeadObjectOutput struct {
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
	VersionId            string
	Restore              string
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
	SSECustomerKeyMD5    string
	ReplicationStatus    string
	ContentRange         string
	IsPartial            bool
	PartsCount           int32
}

// HeadObject retrieves metadata for an object without returning the object itself.
func (o *ObjectOperations) HeadObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *HeadObjectInput) (*HeadObjectOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	coreResult, err := o.svc.headObjectCore(ctx, stores.objects, AdminHeadObjectInput{
		Bucket:    input.Bucket,
		Key:       input.Key,
		VersionID: input.VersionId,
	})
	if err != nil {
		return nil, err
	}
	obj := coreResult.Object

	if input.IfMatch != "" || input.IfNoneMatch != "" || input.IfModifiedSince != nil || input.IfUnmodifiedSince != nil {
		if err := checkObjectPreconditions(obj, input.IfMatch, input.IfNoneMatch, input.IfModifiedSince, input.IfUnmodifiedSince); err != nil {
			return nil, err
		}
	}

	contentLength := obj.Size
	if obj.SSEMetadata != nil {
		contentLength = obj.SSEMetadata.UnencryptedSize
	}

	output := &HeadObjectOutput{
		ContentLength:      contentLength,
		ContentType:        obj.ContentType,
		ContentEncoding:    obj.ContentEncoding,
		ContentLanguage:    obj.ContentLanguage,
		ContentDisposition: obj.ContentDisposition,
		CacheControl:       obj.CacheControl,
		ETag:               formatETag(obj.ETag),
		LastModified:       obj.LastModified,
		Metadata:           obj.Metadata,
		StorageClass:       string(obj.StorageClass),
		VersionId:          obj.VersionID,
		Restore:            restoreHeaderValue(obj, time.Now()),
		ReplicationStatus:  obj.ReplicationStatus,
	}

	// A HEAD Range is resolved against metadata only: the response carries
	// the Content-Range and partial Content-Length the corresponding GET
	// would produce, without a body.  A partNumber request is resolved to
	// the part's window first, with an optional Range evaluated within it.
	if input.PartNumber > 0 {
		start, end, count, err := resolvePartRange(obj.Parts, input.PartNumber, contentLength, input.Range)
		if err != nil {
			return nil, err
		}
		input.Range = fmt.Sprintf("bytes=%d-%d", start, end)
		output.PartsCount = count
	}
	if input.Range != "" {
		ranges, rangeErr := parseRangeHeader(input.Range)
		if rangeErr != nil {
			return nil, rangeErr
		}
		firstRange := ranges[0]
		var offset, length int64
		totalSize := contentLength

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
			return nil, ErrInvalidRange
		}

		actualEnd := offset + length - 1
		if actualEnd >= totalSize {
			actualEnd = totalSize - 1
			length = totalSize - offset
		}

		output.ContentLength = length
		output.ContentRange = fmt.Sprintf("bytes %d-%d/%d", offset, actualEnd, totalSize)
		output.IsPartial = true
	}

	sseCRequested := sseCustomerRequested(input.SSECustomerAlgorithm, input.SSECustomerKey, input.SSECustomerKeyMD5)
	if sseCRequested && (obj.SSEMetadata == nil || obj.SSEMetadata.EncryptionType != s3store.SSETypeCustomer) {
		return nil, NewInvalidRequestError("The encryption parameters are not applicable to this object.")
	}

	if obj.SSEMetadata != nil {
		if obj.SSEMetadata.EncryptionType == s3store.SSETypeCustomer {
			if input.SSECustomerKey == "" {
				return nil, awserrors.NewAWSError("InvalidRequest", "The object was stored using a form of Server Side Encryption. The correct parameters must be provided to retrieve the object.", http.StatusBadRequest)
			}
			output.SSECustomerAlgorithm = "AES256"
			output.SSECustomerKeyMD5 = input.SSECustomerKeyMD5
		} else {
			output.ServerSideEncryption = string(obj.SSEMetadata.EncryptionType)
			output.SSEKMSKeyId = obj.SSEMetadata.KMSKeyID
		}
	}

	return output, nil
}

// GetObjectAttributesInput contains the input parameters for the GetObjectAttributes operation.
type GetObjectAttributesInput struct {
	Bucket           string
	Key              string
	VersionId        string
	MaxParts         int32
	PartNumberMarker string
	ObjectAttributes []string
}

// GetObjectAttributesOutput contains the output from the GetObjectAttributes operation.
type GetObjectAttributesOutput struct {
	XMLName      xml.Name                     `xml:"GetObjectAttributesOutput"`
	VersionId    string                       `xml:"VersionId,omitempty"`
	ETag         string                       `xml:"ETag,omitempty"`
	ObjectSize   int64                        `xml:"ObjectSize,omitempty"`
	StorageClass string                       `xml:"StorageClass,omitempty"`
	LastModified s3Timestamp                  `xml:"LastModified,omitempty"`
	ObjectParts  *GetObjectAttributesParts    `xml:"ObjectParts,omitempty"`
	Checksum     *GetObjectAttributesChecksum `xml:"Checksum,omitempty"`
}

// s3Timestamp marshals time.Time as RFC1123 for S3 XML responses that require
// HTTP-date format (e.g. GetObjectAttributes LastModified field).
// Implements encoding/xml.Marshaler.
type s3Timestamp time.Time

// MarshalXML encodes the timestamp in S3's preferred format (RFC1123 with
// GMT timezone suffix) for XML serialisation.
func (t s3Timestamp) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(strings.Replace(time.Time(t).UTC().Format(time.RFC1123), "UTC", "GMT", 1), start)
}

// GetObjectAttributesParts contains information about the parts of an object.
type GetObjectAttributesParts struct {
	IsTruncated          bool                      `xml:"IsTruncated"`
	MaxParts             int32                     `xml:"MaxParts"`
	NextPartNumberMarker string                    `xml:"NextPartNumberMarker,omitempty"`
	PartNumberMarker     string                    `xml:"PartNumberMarker,omitempty"`
	Parts                []GetObjectAttributesPart `xml:"Parts>Part,omitempty"`
	TotalPartsCount      int32                     `xml:"PartsCount"`
}

// GetObjectAttributesPart contains information about a specific part.
type GetObjectAttributesPart struct {
	PartNumber   int32  `xml:"PartNumber"`
	Size         int64  `xml:"Size"`
	ETag         string `xml:"ETag,omitempty"`
	LastModified string `xml:"LastModified,omitempty"`
}

// GetObjectAttributesChecksum contains checksum information for an object.
type GetObjectAttributesChecksum struct {
	ChecksumCRC32  string `xml:"ChecksumCRC32,omitempty"`
	ChecksumCRC32C string `xml:"ChecksumCRC32C,omitempty"`
	ChecksumSHA1   string `xml:"ChecksumSHA1,omitempty"`
	ChecksumSHA256 string `xml:"ChecksumSHA256,omitempty"`
}

// GetObjectAttributes retrieves attributes of an object.
func (o *ObjectOperations) GetObjectAttributes(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *GetObjectAttributesInput) (*GetObjectAttributesOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	obj, err := stores.objects.HeadWithVersion(ctx, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, mapVersionLookupError(err, input.VersionId)
	}

	objectSize := obj.Size
	if obj.SSEMetadata != nil && obj.SSEMetadata.UnencryptedSize > 0 {
		objectSize = obj.SSEMetadata.UnencryptedSize
	}

	output := &GetObjectAttributesOutput{
		VersionId:    obj.VersionID,
		ETag:         formatETag(obj.ETag),
		ObjectSize:   objectSize,
		StorageClass: string(obj.StorageClass),
		LastModified: s3Timestamp(obj.LastModified),
	}

	for _, attr := range input.ObjectAttributes {
		switch attr {
		case "ETag":
			output.ETag = formatETag(obj.ETag)
		case "ObjectSize":
			output.ObjectSize = objectSize
		case "StorageClass":
			output.StorageClass = string(obj.StorageClass)
		case "ObjectParts":
			if obj.SSEMetadata != nil && len(obj.SSEMetadata.PartEncryptionInfos) > 0 {
				partInfos := obj.SSEMetadata.PartEncryptionInfos
				totalParts := int32(len(partInfos))

				partNumberStart := int32(0)
				if input.PartNumberMarker != "" {
					if parsed, pErr := strconv.ParseInt(input.PartNumberMarker, 10, 32); pErr == nil && parsed > 0 {
						partNumberStart = int32(parsed)
					}
				}

				maxParts := input.MaxParts
				if maxParts <= 0 {
					maxParts = 1000
				}

				var filteredParts []GetObjectAttributesPart
				for i, pi := range partInfos {
					pn := int32(i + 1)
					if pn <= partNumberStart {
						continue
					}
					if int32(len(filteredParts)) >= maxParts {
						break
					}
					filteredParts = append(filteredParts, GetObjectAttributesPart{
						PartNumber: pn,
						Size:       pi.PlainSize,
					})
				}

				isTruncated := int32(len(partInfos)) > partNumberStart+int32(len(filteredParts))
				var nextMarker string
				if isTruncated && len(filteredParts) > 0 {
					nextMarker = strconv.FormatInt(int64(filteredParts[len(filteredParts)-1].PartNumber), 10)
				}

				output.ObjectParts = &GetObjectAttributesParts{
					IsTruncated:          isTruncated,
					MaxParts:             maxParts,
					NextPartNumberMarker: nextMarker,
					PartNumberMarker:     input.PartNumberMarker,
					Parts:                filteredParts,
					TotalPartsCount:      totalParts,
				}
			}
		case "Checksum":
			output.Checksum = &GetObjectAttributesChecksum{}
		}
	}

	return output, nil
}
