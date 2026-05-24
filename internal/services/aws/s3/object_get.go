package s3

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
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
	ContentRange         string
	IsPartial            bool
	AcceptRanges         string
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
	SSECustomerKeyMD5    string
}

// GetObject retrieves an object from S3.
func (o *ObjectOperations) GetObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *GetObjectInput) (*GetObjectOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	if input.IfMatch != "" || input.IfNoneMatch != "" || input.IfModifiedSince != nil || input.IfUnmodifiedSince != nil {
		obj, err := stores.objects.HeadWithVersion(ctx, input.Bucket, input.Key, input.VersionId)
		if err != nil {
			return nil, err
		}

		if input.IfMatch != "" {
			if strings.Trim(obj.ETag, "\"") != strings.Trim(input.IfMatch, "\"") {
				return nil, ErrPreconditionFailed
			}
		}
		if input.IfNoneMatch != "" {
			if strings.Trim(obj.ETag, "\"") == strings.Trim(input.IfNoneMatch, "\"") {
				return nil, ErrNotModified
			}
		}
		if input.IfUnmodifiedSince != nil {
			if obj.LastModified.After(*input.IfUnmodifiedSince) {
				return nil, ErrPreconditionFailed
			}
		}
		if input.IfModifiedSince != nil {
			if !obj.LastModified.After(*input.IfModifiedSince) {
				return nil, ErrNotModified
			}
		}
	}

	reader, obj, err := stores.objects.GetWithVersion(ctx, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, err
	}

	output := &GetObjectOutput{
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
		VersionId:          obj.VersionID,
	}

	var decryptedData []byte
	var unencryptedSize int64

	if obj.SSEMetadata != nil {
		encryptedData, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read encrypted data: %w", err)
		}

		decryptedData, unencryptedSize, err = o.decryptObjectData(encryptedData, obj.SSEMetadata, input)
		if err != nil {
			return nil, err
		}

		if obj.SSEMetadata.EncryptionType == s3store.SSETypeCustomer {
			output.SSECustomerAlgorithm = "AES256"
			output.SSECustomerKeyMD5 = input.SSECustomerKeyMD5
		} else {
			output.ServerSideEncryption = string(obj.SSEMetadata.EncryptionType)
			output.SSEKMSKeyId = obj.SSEMetadata.KMSKeyID
		}
	}

	if input.Range != "" {
		ranges, err := parseRangeHeader(input.Range)
		if err != nil {
			if obj.SSEMetadata == nil {
				reader.Close()
			}
			return nil, err
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
			var getErr error
			rangeReader, _, getErr = stores.objects.GetRangeWithVersion(ctx, input.Bucket, input.Key, input.VersionId, offset, length)
			if getErr != nil {
				return nil, getErr
			}
		}

		return &GetObjectOutput{
			Body:                 rangeReader,
			ContentLength:        length,
			ContentType:          obj.ContentType,
			ETag:                 formatETag(obj.ETag),
			LastModified:         obj.LastModified,
			Metadata:             obj.Metadata,
			ContentRange:         fmt.Sprintf("bytes %d-%d/%d", offset, actualEnd, totalSize),
			IsPartial:            true,
			AcceptRanges:         "bytes",
			VersionId:            obj.VersionID,
			ServerSideEncryption: output.ServerSideEncryption,
			SSEKMSKeyId:          output.SSEKMSKeyId,
			SSECustomerAlgorithm: output.SSECustomerAlgorithm,
			SSECustomerKeyMD5:    output.SSECustomerKeyMD5,
			ContentEncoding:      obj.ContentEncoding,
			ContentLanguage:      obj.ContentLanguage,
			ContentDisposition:   obj.ContentDisposition,
			CacheControl:         obj.CacheControl,
			StorageClass:         string(obj.StorageClass),
		}, nil
	}

	if obj.SSEMetadata != nil {
		output.Body = io.NopCloser(bytes.NewReader(decryptedData))
		output.ContentLength = unencryptedSize
	}

	return output, nil
}

// HeadObjectInput contains the input parameters for the HeadObject operation.
type HeadObjectInput struct {
	Bucket            string
	Key               string
	VersionId         string
	SSECustomerKey    string
	SSECustomerKeyMD5 string
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
	ServerSideEncryption string
	SSEKMSKeyId          string
	SSECustomerAlgorithm string
	SSECustomerKeyMD5    string
}

// HeadObject retrieves metadata for an object without returning the object itself.
func (o *ObjectOperations) HeadObject(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, input *HeadObjectInput) (*HeadObjectOutput, error) {
	if err := o.validateBucketExists(stores, input.Bucket); err != nil {
		return nil, err
	}

	if err := validateObjectKey(input.Key); err != nil {
		return nil, err
	}

	obj, err := stores.objects.HeadWithVersion(ctx, input.Bucket, input.Key, input.VersionId)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	objectSize := obj.Size
	if obj.SSEMetadata != nil && obj.SSEMetadata.UnencryptedSize > 0 {
		objectSize = obj.SSEMetadata.UnencryptedSize
	}

	output := &GetObjectAttributesOutput{
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
				output.ObjectParts = &GetObjectAttributesParts{
					IsTruncated:     false,
					MaxParts:        input.MaxParts,
					TotalPartsCount: totalParts,
				}
				for i, pi := range partInfos {
					output.ObjectParts.Parts = append(output.ObjectParts.Parts, GetObjectAttributesPart{
						PartNumber: int32(i + 1),
						Size:       pi.PlainSize,
					})
				}
			}
		case "Checksum":
			output.Checksum = &GetObjectAttributesChecksum{}
		}
	}

	return output, nil
}

func (o *ObjectOperations) decryptMultipartParts(combinedEncrypted []byte, sseMeta *s3store.SSEObjectMetadata, bucket, key string) ([]byte, error) {
	var result bytes.Buffer
	offset := int64(0)

	for _, partInfo := range sseMeta.PartEncryptionInfos {
		encSize := partInfo.EncryptedSize
		if encSize == 0 {
			continue
		}
		if offset+encSize > int64(len(combinedEncrypted)) {
			return nil, fmt.Errorf("encrypted data truncated at part boundary")
		}
		partEncData := combinedEncrypted[offset : offset+encSize]
		offset += encSize

		partMeta := &s3store.SSEObjectMetadata{
			EncryptionType:   sseMeta.EncryptionType,
			EncryptedDataKey: partInfo.DataKey,
			ContentNonce:     partInfo.ContentNonce,
			KMSKeyID:         sseMeta.KMSKeyID,
		}

		decResult, err := o.svc.encryptionManager.Decrypt(partEncData, partMeta, bucket, key)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt part at offset %d: %w", offset-encSize, err)
		}
		result.Write(decResult.DecryptedData)
	}

	if offset < int64(len(combinedEncrypted)) {
		remaining := combinedEncrypted[offset:]
		decResult, err := o.svc.encryptionManager.Decrypt(remaining, sseMeta, bucket, key)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt trailing data: %w", err)
		}
		result.Write(decResult.DecryptedData)
	}

	return result.Bytes(), nil
}

func (o *ObjectOperations) decryptObjectData(encryptedData []byte, sseMeta *s3store.SSEObjectMetadata, input *GetObjectInput) ([]byte, int64, error) {
	unencryptedSize := sseMeta.UnencryptedSize

	if sseMeta.EncryptionType == s3store.SSETypeCustomer {
		if input.SSECustomerKey == "" {
			return nil, 0, awserrors.NewAWSError("InvalidRequest", "The object was stored using a form of Server Side Encryption. The correct parameters must be provided to retrieve the object.", http.StatusBadRequest)
		}
		customerKey, err := o.svc.encryptionManager.ParseCustomerKey(input.SSECustomerKey, input.SSECustomerKeyMD5)
		if err != nil {
			return nil, 0, ErrInvalidSSECustomerKey
		}

		var plainData []byte
		if len(sseMeta.PartEncryptionInfos) > 0 {
			plainData, err = o.svc.encryptionManager.DecryptChunked(encryptedData, sseMeta, input.Bucket, input.Key, customerKey)
		} else {
			decResult, decErr := o.svc.encryptionManager.DecryptWithCustomerKey(encryptedData, sseMeta, input.Bucket, input.Key, customerKey)
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
		plainData, err = o.svc.encryptionManager.DecryptChunked(encryptedData, sseMeta, input.Bucket, input.Key, nil)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to decrypt chunked data: %w", err)
		}
	} else {
		decResult, decErr := o.svc.encryptionManager.Decrypt(encryptedData, sseMeta, input.Bucket, input.Key)
		if decErr != nil {
			return nil, 0, fmt.Errorf("failed to decrypt data: %w", decErr)
		}
		plainData = decResult.DecryptedData
	}
	return plainData, unencryptedSize, nil
}
