package s3

import (
	"context"
	"errors"
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

const (
	s3MaxKeys    = 1000
	s3MaxParts   = 1000
	s3MaxUploads = 1000
)

func errorStatusCode(err error, fallback int) int {
	var awsErr interface{ StatusCode() int }
	if errors.As(err, &awsErr) {
		return awsErr.StatusCode()
	}
	return fallback
}

func clampInt(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func setSSEHeaders(header http.Header, customerAlgo, customerKeyMD5, sseType, kmsKeyID string) {
	if customerAlgo != "" {
		header.Set("x-amz-server-side-encryption-customer-algorithm", customerAlgo)
	}
	if customerKeyMD5 != "" {
		header.Set("x-amz-server-side-encryption-customer-key-MD5", customerKeyMD5)
	}
	if sseType != "" {
		header.Set("x-amz-server-side-encryption", sseType)
		if kmsKeyID != "" {
			header.Set("x-amz-server-side-encryption-aws-kms-key-id", kmsKeyID)
		}
	}
}

type objectResponseHeaders struct {
	ETag                 string
	ContentType          string
	ContentLength        int64
	LastModified         time.Time
	VersionId            string
	SSECustomerAlgorithm string
	SSECustomerKeyMD5    string
	ServerSideEncryption string
	SSEKMSKeyId          string
	CacheControl         string
	ContentDisposition   string
	ContentEncoding      string
	ContentLanguage      string
	StorageClass         string
	ReplicationStatus    string
	Metadata             map[string]string
}

func setObjectResponseHeaders(header http.Header, h objectResponseHeaders) {
	header.Set("ETag", h.ETag)
	header.Set("Content-Length", strconv.FormatInt(h.ContentLength, 10))
	header.Set("Content-Type", h.ContentType)
	header.Set("Last-Modified", h.LastModified.Format(http.TimeFormat))
	if h.VersionId != "" && h.VersionId != "null" {
		header.Set("x-amz-version-id", h.VersionId)
	}
	setSSEHeaders(header, h.SSECustomerAlgorithm, h.SSECustomerKeyMD5, h.ServerSideEncryption, h.SSEKMSKeyId)
	if h.CacheControl != "" {
		header.Set("Cache-Control", h.CacheControl)
	}
	if h.ContentDisposition != "" {
		header.Set("Content-Disposition", h.ContentDisposition)
	}
	if h.ContentEncoding != "" {
		header.Set("Content-Encoding", h.ContentEncoding)
	}
	if h.ContentLanguage != "" {
		header.Set("Content-Language", h.ContentLanguage)
	}
	if h.StorageClass != "" {
		header.Set("x-amz-storage-class", h.StorageClass)
	}
	if h.ReplicationStatus != "" {
		header.Set("x-amz-replication-status", h.ReplicationStatus)
	}
	for k, v := range h.Metadata {
		header.Set("x-amz-meta-"+k, v)
	}
}

// HandleRequest processes HTTP requests for object-level operations such as
// get, put, delete, copy, multipart uploads, tagging, ACL, and legal hold.
func (o *ObjectOperations) HandleRequest(ctx context.Context, reqCtx *request.RequestContext, stores *s3Stores, r *http.Request, bucket, key string) (interface{}, http.Header, int, error) {
	method := r.Method
	query := r.URL.Query()
	header := make(http.Header)

	switch {
	case method == "POST" && query.Has("restore"):
		input := &RestoreObjectInput{
			Bucket:    bucket,
			Key:       key,
			VersionId: r.Header.Get("x-amz-version-id"),
			Body:      r.Body,
		}
		result, err := o.RestoreObject(ctx, reqCtx, stores, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		return result, header, http.StatusAccepted, nil

	case method == "POST" && query.Has("uploads"):
		input := &CreateMultipartUploadInput{
			Bucket:               bucket,
			Key:                  key,
			ContentType:          r.Header.Get("Content-Type"),
			ContentEncoding:      r.Header.Get("Content-Encoding"),
			ContentDisposition:   r.Header.Get("Content-Disposition"),
			CacheControl:         r.Header.Get("Cache-Control"),
			StorageClass:         r.Header.Get("x-amz-storage-class"),
			ServerSideEncryption: r.Header.Get("x-amz-server-side-encryption"),
			SSEKMSKeyId:          r.Header.Get("x-amz-server-side-encryption-aws-kms-key-id"),
			SSECustomerAlgorithm: r.Header.Get("x-amz-server-side-encryption-customer-algorithm"),
			SSECustomerKey:       r.Header.Get("x-amz-server-side-encryption-customer-key"),
			SSECustomerKeyMD5:    r.Header.Get("x-amz-server-side-encryption-customer-key-MD5"),
		}
		for k, v := range r.Header {
			if strings.HasPrefix(k, "X-Amz-Meta-") {
				if input.Metadata == nil {
					input.Metadata = make(map[string]string)
				}
				input.Metadata[strings.TrimPrefix(k, "X-Amz-Meta-")] = v[0]
			}
		}
		result, err := o.CreateMultipartUpload(ctx, reqCtx, stores, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		setSSEHeaders(header, result.SSECustomerAlgorithm, "", result.ServerSideEncryption, result.SSEKMSKeyId)
		return result, header, http.StatusOK, nil

	case method == "PUT" && query.Has("uploadId") && query.Has("partNumber") && r.Header.Get("x-amz-copy-source") != "":
		partNumber, err := strconv.Atoi(query.Get("partNumber"))
		if err != nil || partNumber < minPartNumber || partNumber > maxPartNumber {
			return nil, header, http.StatusBadRequest, NewInvalidArgumentError(fmt.Sprintf("invalid partNumber: must be between %d and %d", minPartNumber, maxPartNumber))
		}
		input := &UploadPartCopyInput{
			Bucket:                    bucket,
			Key:                       key,
			UploadId:                  query.Get("uploadId"),
			PartNumber:                partNumber,
			CopySource:                r.Header.Get("x-amz-copy-source"),
			CopySourceRange:           r.Header.Get("x-amz-copy-source-range"),
			CopySourceVersionId:       r.Header.Get("x-amz-copy-source-version-id"),
			CopySourceSSECustomerAlgo: r.Header.Get("x-amz-copy-source-server-side-encryption-customer-algorithm"),
			CopySourceSSECustomerKey:  r.Header.Get("x-amz-copy-source-server-side-encryption-customer-key"),
			CopySourceSSECustomerMD5:  r.Header.Get("x-amz-copy-source-server-side-encryption-customer-key-MD5"),
			SSECustomerAlgorithm:      r.Header.Get("x-amz-server-side-encryption-customer-algorithm"),
			SSECustomerKey:            r.Header.Get("x-amz-server-side-encryption-customer-key"),
			SSECustomerKeyMD5:         r.Header.Get("x-amz-server-side-encryption-customer-key-MD5"),
		}
		result, err := o.UploadPartCopy(ctx, reqCtx, stores, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		return result, header, http.StatusOK, nil

	case method == "PUT" && query.Has("uploadId") && query.Has("partNumber"):
		partNumber, err := strconv.Atoi(query.Get("partNumber"))
		if err != nil || partNumber < minPartNumber || partNumber > maxPartNumber {
			return nil, header, http.StatusBadRequest, NewInvalidArgumentError(fmt.Sprintf("invalid partNumber: must be between %d and %d", minPartNumber, maxPartNumber))
		}

		var partBody io.Reader = r.Body
		if isAwsChunkedRequest(r) {
			partBody = decodeAwsChunkedBody(r.Body)
		}

		partContentLength := int64(-1)
		if !isAwsChunkedRequest(r) {
			if cl := r.Header.Get("Content-Length"); cl != "" {
				if n, parseErr := strconv.ParseInt(cl, 10, 64); parseErr == nil && n >= 0 {
					partContentLength = n
				}
			}
		}

		input := &UploadPartInput{
			Bucket:               bucket,
			Key:                  key,
			UploadId:             query.Get("uploadId"),
			PartNumber:           partNumber,
			Body:                 partBody,
			ContentLength:        partContentLength,
			SSECustomerAlgorithm: r.Header.Get("x-amz-server-side-encryption-customer-algorithm"),
			SSECustomerKey:       r.Header.Get("x-amz-server-side-encryption-customer-key"),
			SSECustomerKeyMD5:    r.Header.Get("x-amz-server-side-encryption-customer-key-MD5"),
		}
		result, err := o.UploadPart(ctx, reqCtx, stores, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		header.Set("ETag", result.ETag)
		setSSEHeaders(header, result.SSECustomerAlgorithm, "", "", "")
		return result, header, http.StatusOK, nil

	case method == "GET" && query.Has("uploadId"):
		input := &ListPartsInput{
			Bucket:   bucket,
			Key:      key,
			UploadId: query.Get("uploadId"),
		}
		if maxParts := query.Get("max-parts"); maxParts != "" {
			mp, err := strconv.Atoi(maxParts)
			if err != nil {
				return nil, header, http.StatusBadRequest, NewInvalidArgumentError("Provided max-parts not an integer")
			}
			input.MaxParts = clampInt(mp, 0, s3MaxParts)
		}
		if partNumberMarker := query.Get("part-number-marker"); partNumberMarker != "" {
			input.PartNumberMarker = partNumberMarker
		}
		result, err := o.ListParts(ctx, reqCtx, stores, input)
		return result, header, http.StatusOK, err

	case method == "POST" && query.Has("uploadId"):
		var completeReq struct {
			Parts []CompletedPart `xml:"Part"`
		}
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&completeReq); err != nil {
			return nil, header, http.StatusBadRequest, err
		}
		input := &CompleteMultipartUploadInput{
			Bucket:   bucket,
			Key:      key,
			UploadId: query.Get("uploadId"),
			Parts:    completeReq.Parts,
		}
		result, err := o.CompleteMultipartUpload(ctx, reqCtx, stores, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		setSSEHeaders(header, "", "", result.ServerSideEncryption, result.SSEKMSKeyId)
		return result, header, http.StatusOK, nil

	case method == "DELETE" && query.Has("uploadId"):
		err := o.AbortMultipartUpload(ctx, reqCtx, stores, &AbortMultipartUploadInput{
			Bucket:   bucket,
			Key:      key,
			UploadId: query.Get("uploadId"),
		})
		return nil, header, http.StatusNoContent, err

	case method == "GET" && query.Has("tagging"):
		result, err := o.GetObjectTagging(ctx, reqCtx, stores, &GetObjectTaggingInput{Bucket: bucket, Key: key})
		return result, header, http.StatusOK, err

	case method == "GET" && query.Has("legal-hold"):
		result, err := o.GetObjectLegalHold(ctx, reqCtx, stores, &GetObjectLegalHoldInput{
			Bucket:    bucket,
			Key:       key,
			VersionId: query.Get("versionId"),
		})
		return result, header, http.StatusOK, err

	case method == "PUT" && query.Has("legal-hold"):
		var legalHoldReq LegalHoldInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&legalHoldReq); err != nil {
			return nil, header, http.StatusBadRequest, err
		}
		err := o.PutObjectLegalHold(ctx, reqCtx, stores, &PutObjectLegalHoldInput{
			Bucket:    bucket,
			Key:       key,
			VersionId: query.Get("versionId"),
			LegalHold: &legalHoldReq,
		})
		return nil, header, http.StatusOK, err

	case method == "GET" && query.Has("retention"):
		result, err := o.GetObjectRetention(ctx, reqCtx, stores, &GetObjectRetentionInput{
			Bucket:    bucket,
			Key:       key,
			VersionId: query.Get("versionId"),
		})
		return result, header, http.StatusOK, err

	case method == "PUT" && query.Has("retention"):
		var retentionReq RetentionInput
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&retentionReq); err != nil {
			return nil, header, http.StatusBadRequest, err
		}
		err := o.PutObjectRetention(ctx, reqCtx, stores, &PutObjectRetentionInput{
			Bucket:    bucket,
			Key:       key,
			VersionId: query.Get("versionId"),
			Retention: &retentionReq,
		})
		return nil, header, http.StatusOK, err

	case method == "PUT" && query.Has("acl"):
		input := &PutObjectAclInput{
			Bucket:           bucket,
			Key:              key,
			VersionId:        query.Get("versionId"),
			ACL:              r.Header.Get("x-amz-acl"),
			GrantFullControl: r.Header.Get("x-amz-grant-full-control"),
			GrantRead:        r.Header.Get("x-amz-grant-read"),
			GrantReadACP:     r.Header.Get("x-amz-grant-read-acp"),
			GrantWrite:       r.Header.Get("x-amz-grant-write"),
			GrantWriteACP:    r.Header.Get("x-amz-grant-write-acp"),
		}
		if input.ACL == "" && input.GrantFullControl == "" && input.GrantRead == "" && input.GrantWrite == "" {
			var acp s3store.AccessControlPolicy
			if err := request.NewSafeXMLDecoder(r.Body).Decode(&acp); err == nil {
				input.AccessControlPolicy = &acp
			}
		}
		err := o.PutObjectAcl(ctx, reqCtx, stores, input)
		return nil, header, http.StatusOK, err

	case method == "GET" && query.Has("acl"):
		result, err := o.GetObjectAcl(ctx, reqCtx, stores, bucket, key, query.Get("versionId"))
		return result, header, http.StatusOK, err

	case method == "GET" && query.Has("attributes"):
		objectAttributes := r.Header["X-Amz-Object-Attributes"]
		maxParts := int32(1000)
		if mp := query.Get("max-parts"); mp != "" {
			if parsed, err := strconv.ParseInt(mp, 10, 32); err == nil {
				maxParts = int32(parsed)
			}
		}
		result, err := o.GetObjectAttributes(ctx, reqCtx, stores, &GetObjectAttributesInput{
			Bucket:           bucket,
			Key:              key,
			VersionId:        query.Get("versionId"),
			MaxParts:         maxParts,
			PartNumberMarker: query.Get("part-number-marker"),
			ObjectAttributes: objectAttributes,
		})
		return result, header, http.StatusOK, err

	case method == "GET":
		input := &GetObjectInput{
			Bucket:               bucket,
			Key:                  key,
			Range:                r.Header.Get("Range"),
			VersionId:            query.Get("versionId"),
			IfMatch:              r.Header.Get("If-Match"),
			IfNoneMatch:          r.Header.Get("If-None-Match"),
			SSECustomerAlgorithm: r.Header.Get("x-amz-server-side-encryption-customer-algorithm"),
			SSECustomerKey:       r.Header.Get("x-amz-server-side-encryption-customer-key"),
			SSECustomerKeyMD5:    r.Header.Get("x-amz-server-side-encryption-customer-key-md5"),
		}
		if modSince := r.Header.Get("If-Modified-Since"); modSince != "" {
			if t, err := time.Parse(http.TimeFormat, modSince); err == nil {
				input.IfModifiedSince = &t
			}
		}
		if unmodSince := r.Header.Get("If-Unmodified-Since"); unmodSince != "" {
			if t, err := time.Parse(http.TimeFormat, unmodSince); err == nil {
				input.IfUnmodifiedSince = &t
			}
		}
		result, err := o.GetObject(ctx, reqCtx, stores, input)
		if err != nil {
			return nil, header, errorStatusCode(err, http.StatusNotFound), err
		}
		setObjectResponseHeaders(header, objectResponseHeaders{
			ETag: result.ETag, ContentType: result.ContentType, ContentLength: result.ContentLength, LastModified: result.LastModified,
			VersionId: result.VersionId, SSECustomerAlgorithm: result.SSECustomerAlgorithm, SSECustomerKeyMD5: result.SSECustomerKeyMD5,
			ServerSideEncryption: result.ServerSideEncryption, SSEKMSKeyId: result.SSEKMSKeyId,
			CacheControl: result.CacheControl, ContentDisposition: result.ContentDisposition, ContentEncoding: result.ContentEncoding, ContentLanguage: result.ContentLanguage,
			StorageClass: result.StorageClass, ReplicationStatus: result.ReplicationStatus, Metadata: result.Metadata,
		})
		return result, header, http.StatusOK, nil

	case method == "HEAD":
		result, err := o.HeadObject(ctx, reqCtx, stores, &HeadObjectInput{Bucket: bucket, Key: key, VersionId: query.Get("versionId"), SSECustomerKey: r.Header.Get("x-amz-server-side-encryption-customer-key"), SSECustomerKeyMD5: r.Header.Get("x-amz-server-side-encryption-customer-key-MD5")})
		if err != nil {
			return nil, header, errorStatusCode(err, http.StatusNotFound), err
		}
		setObjectResponseHeaders(header, objectResponseHeaders{
			ETag: result.ETag, ContentType: result.ContentType, ContentLength: result.ContentLength, LastModified: result.LastModified,
			VersionId: result.VersionId, SSECustomerAlgorithm: result.SSECustomerAlgorithm, SSECustomerKeyMD5: result.SSECustomerKeyMD5,
			ServerSideEncryption: result.ServerSideEncryption, SSEKMSKeyId: result.SSEKMSKeyId,
			CacheControl: result.CacheControl, ContentDisposition: result.ContentDisposition, ContentEncoding: result.ContentEncoding, ContentLanguage: result.ContentLanguage,
			StorageClass: result.StorageClass, ReplicationStatus: result.ReplicationStatus, Metadata: result.Metadata,
		})
		return result, header, http.StatusOK, nil

	case method == "PUT" && query.Has("tagging"):
		var tagSet struct {
			Tags []Tag `xml:"TagSet>Tag"`
		}
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&tagSet); err != nil {
			return nil, header, http.StatusBadRequest, err
		}
		err := o.PutObjectTagging(ctx, reqCtx, stores, &PutObjectTaggingInput{
			Bucket: bucket,
			Key:    key,
			Tags:   tagSet.Tags,
		})
		return nil, header, http.StatusOK, err

	case method == "PUT" && r.Header.Get("x-amz-copy-source") != "" && !query.Has("uploadId"):
		input := &CopyObjectInput{
			Bucket:                    bucket,
			Key:                       key,
			CopySource:                r.Header.Get("x-amz-copy-source"),
			MetadataDirective:         r.Header.Get("x-amz-metadata-directive"),
			ContentType:               r.Header.Get("Content-Type"),
			StorageClass:              r.Header.Get("x-amz-storage-class"),
			ServerSideEncryption:      r.Header.Get("x-amz-server-side-encryption"),
			SSEKMSKeyId:               r.Header.Get("x-amz-server-side-encryption-aws-kms-key-id"),
			SSECustomerAlgorithm:      r.Header.Get("x-amz-server-side-encryption-customer-algorithm"),
			SSECustomerKey:            r.Header.Get("x-amz-server-side-encryption-customer-key"),
			SSECustomerKeyMD5:         r.Header.Get("x-amz-server-side-encryption-customer-key-MD5"),
			CopySourceSSECustomerAlgo: r.Header.Get("x-amz-copy-source-server-side-encryption-customer-algorithm"),
			CopySourceSSECustomerKey:  r.Header.Get("x-amz-copy-source-server-side-encryption-customer-key"),
			CopySourceSSECustomerMD5:  r.Header.Get("x-amz-copy-source-server-side-encryption-customer-key-MD5"),
		}
		for k, v := range r.Header {
			if strings.HasPrefix(k, "X-Amz-Meta-") {
				if input.Metadata == nil {
					input.Metadata = make(map[string]string)
				}
				input.Metadata[strings.TrimPrefix(k, "X-Amz-Meta-")] = v[0]
			}
		}
		result, err := o.CopyObject(ctx, reqCtx, stores, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		setSSEHeaders(header, "", "", result.ServerSideEncryption, result.SSEKMSKeyId)
		return result, header, http.StatusOK, nil

	case method == "PUT":
		contentType := r.Header.Get("Content-Type")
		contentLengthStr := r.Header.Get("Content-Length")
		transferEncoding := r.Header.Get("Transfer-Encoding")

		if contentLengthStr == "" && transferEncoding != "chunked" {
			return nil, header, http.StatusLengthRequired, ErrMissingContentLength
		}

		var contentLength int64
		if contentLengthStr != "" {
			var parseErr error
			contentLength, parseErr = strconv.ParseInt(contentLengthStr, 10, 64)
			if parseErr != nil || contentLength < 0 {
				return nil, header, http.StatusBadRequest, NewInvalidArgumentError(fmt.Sprintf("invalid Content-Length: %s", contentLengthStr))
			}
		}

		if decodedContentLengthStr := r.Header.Get("X-Amz-Decoded-Content-Length"); decodedContentLengthStr != "" {
			if dcl, err := strconv.ParseInt(decodedContentLengthStr, 10, 64); err == nil {
				contentLength = dcl
			}
		}

		metadata := make(map[string]string)
		for k, v := range r.Header {
			if strings.HasPrefix(k, "X-Amz-Meta-") {
				metaKey := strings.TrimPrefix(k, "X-Amz-Meta-")
				metadata[metaKey] = v[0]
			}
		}

		var body io.Reader = r.Body
		if isAwsChunkedRequest(r) {
			body = decodeAwsChunkedBody(r.Body)
		}

		input := &PutObjectInput{
			Bucket:               bucket,
			Key:                  key,
			Body:                 body,
			ContentLength:        contentLength,
			ContentType:          contentType,
			ContentEncoding:      r.Header.Get("Content-Encoding"),
			ContentDisposition:   r.Header.Get("Content-Disposition"),
			CacheControl:         r.Header.Get("Cache-Control"),
			ContentLanguage:      r.Header.Get("Content-Language"),
			Metadata:             metadata,
			StorageClass:         r.Header.Get("x-amz-storage-class"),
			IfMatch:              r.Header.Get("If-Match"),
			IfNoneMatch:          r.Header.Get("If-None-Match"),
			ServerSideEncryption: r.Header.Get("x-amz-server-side-encryption"),
			SSEKMSKeyId:          r.Header.Get("x-amz-server-side-encryption-aws-kms-key-id"),
			SSECustomerAlgorithm: r.Header.Get("x-amz-server-side-encryption-customer-algorithm"),
			SSECustomerKey:       r.Header.Get("x-amz-server-side-encryption-customer-key"),
			SSECustomerKeyMD5:    r.Header.Get("x-amz-server-side-encryption-customer-key-md5"),
			Tagging:              r.Header.Get("x-amz-tagging"),
		}
		result, err := o.PutObject(ctx, reqCtx, stores, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		header.Set("ETag", result.ETag)
		if result.VersionId != "" && result.VersionId != "null" {
			header.Set("x-amz-version-id", result.VersionId)
		}
		setSSEHeaders(header, input.SSECustomerAlgorithm, input.SSECustomerKeyMD5, result.ServerSideEncryption, result.SSEKMSKeyId)
		for k, v := range metadata {
			header.Set("x-amz-meta-"+k, v)
		}
		return result, header, http.StatusOK, nil

	case method == "DELETE" && query.Has("tagging"):
		err := o.DeleteObjectTagging(ctx, reqCtx, stores, &DeleteObjectTaggingInput{Bucket: bucket, Key: key})
		return nil, header, http.StatusNoContent, err

	case method == "DELETE":
		result, err := o.DeleteObject(ctx, reqCtx, stores, &DeleteObjectInput{
			Bucket:                    bucket,
			Key:                       key,
			VersionId:                 query.Get("versionId"),
			BypassGovernanceRetention: r.Header.Get("x-amz-bypass-governance-retention") == "true",
		})
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		if result.DeleteMarker {
			header.Set("x-amz-delete-marker", "true")
		}
		if result.VersionId != "" {
			header.Set("x-amz-version-id", result.VersionId)
		}
		return nil, header, http.StatusNoContent, nil

	default:
		return nil, header, http.StatusNotImplemented, awserrors.NewAWSError("NotImplemented", "unsupported operation", http.StatusNotImplemented)
	}
}
