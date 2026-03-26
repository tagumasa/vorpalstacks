package s3

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// HandleRequest processes HTTP requests for object-level operations such as
// get, put, delete, copy, multipart uploads, tagging, ACL, and legal hold.
func (o *ObjectOperations) HandleRequest(ctx context.Context, reqCtx *request.RequestContext, r *http.Request, bucket, key string) (interface{}, http.Header, int, error) {
	method := r.Method
	query := r.URL.Query()
	header := make(http.Header)

	switch {
	case method == "GET" && key == "" && query.Has("versions"):
		input := &ListObjectVersionsInput{
			Bucket:          bucket,
			Delimiter:       query.Get("delimiter"),
			Prefix:          query.Get("prefix"),
			KeyMarker:       query.Get("key-marker"),
			VersionIdMarker: query.Get("version-id-marker"),
		}
		if maxKeys := query.Get("max-keys"); maxKeys != "" {
			mk, err := strconv.Atoi(maxKeys)
			if err != nil {
				return nil, header, http.StatusBadRequest, NewInvalidArgumentError("Provided max-keys not an integer")
			}
			input.MaxKeys = mk
		}
		result, err := o.ListObjectVersions(ctx, reqCtx, input)
		return result, header, http.StatusOK, err

	case method == "GET" && key == "" && query.Has("uploads"):
		input := &ListMultipartUploadsInput{
			Bucket:         bucket,
			Delimiter:      query.Get("delimiter"),
			Prefix:         query.Get("prefix"),
			KeyMarker:      query.Get("key-marker"),
			UploadIdMarker: query.Get("upload-id-marker"),
		}
		if maxUploads := query.Get("max-uploads"); maxUploads != "" {
			mu, err := strconv.Atoi(maxUploads)
			if err != nil {
				return nil, header, http.StatusBadRequest, NewInvalidArgumentError("Provided max-uploads not an integer")
			}
			input.MaxUploads = mu
		}
		result, err := o.ListMultipartUploads(ctx, reqCtx, input)
		return result, header, http.StatusOK, err

	case method == "POST" && query.Has("restore"):
		input := &RestoreObjectInput{
			Bucket:    bucket,
			Key:       key,
			VersionId: r.Header.Get("x-amz-version-id"),
			Body:      r.Body,
		}
		result, err := o.RestoreObject(ctx, reqCtx, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		return result, header, http.StatusAccepted, nil

	case method == "POST" && query.Has("uploads"):
		input := &CreateMultipartUploadInput{
			Bucket:               bucket,
			Key:                  key,
			ContentType:          r.Header.Get("Content-Type"),
			ServerSideEncryption: r.Header.Get("x-amz-server-side-encryption"),
			SSEKMSKeyId:          r.Header.Get("x-amz-server-side-encryption-aws-kms-key-id"),
			SSECustomerAlgorithm: r.Header.Get("x-amz-server-side-encryption-customer-algorithm"),
			SSECustomerKey:       r.Header.Get("x-amz-server-side-encryption-customer-key"),
			SSECustomerKeyMD5:    r.Header.Get("x-amz-server-side-encryption-customer-key-MD5"),
		}
		result, err := o.CreateMultipartUpload(ctx, reqCtx, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		if result.ServerSideEncryption != "" {
			header.Set("x-amz-server-side-encryption", result.ServerSideEncryption)
		}
		if result.SSEKMSKeyId != "" {
			header.Set("x-amz-server-side-encryption-aws-kms-key-id", result.SSEKMSKeyId)
		}
		if result.SSECustomerAlgorithm != "" {
			header.Set("x-amz-server-side-encryption-customer-algorithm", result.SSECustomerAlgorithm)
		}
		return result, header, http.StatusOK, nil

	case method == "PUT" && query.Has("uploadId") && query.Has("partNumber") && r.Header.Get("x-amz-copy-source") != "":
		partNumber, err := strconv.Atoi(query.Get("partNumber"))
		if err != nil {
			return nil, nil, http.StatusBadRequest, fmt.Errorf("invalid partNumber: %w", err)
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
		result, err := o.UploadPartCopy(ctx, reqCtx, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		return result, header, http.StatusOK, nil

	case method == "PUT" && query.Has("uploadId") && query.Has("partNumber"):
		partNumber, err := strconv.Atoi(query.Get("partNumber"))
		if err != nil {
			return nil, nil, http.StatusBadRequest, fmt.Errorf("invalid partNumber: %w", err)
		}
		input := &UploadPartInput{
			Bucket:               bucket,
			Key:                  key,
			UploadId:             query.Get("uploadId"),
			PartNumber:           partNumber,
			Body:                 r.Body,
			SSECustomerAlgorithm: r.Header.Get("x-amz-server-side-encryption-customer-algorithm"),
			SSECustomerKey:       r.Header.Get("x-amz-server-side-encryption-customer-key"),
			SSECustomerKeyMD5:    r.Header.Get("x-amz-server-side-encryption-customer-key-MD5"),
		}
		result, err := o.UploadPart(ctx, reqCtx, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		header.Set("ETag", result.ETag)
		if result.SSECustomerAlgorithm != "" {
			header.Set("x-amz-server-side-encryption-customer-algorithm", result.SSECustomerAlgorithm)
		}
		return result, header, http.StatusOK, nil

	case method == "GET" && query.Has("uploadId"):
		input := &ListPartsInput{
			Bucket:   bucket,
			Key:      key,
			UploadId: query.Get("uploadId"),
		}
		if maxParts := query.Get("max-parts"); maxParts != "" {
			input.MaxParts, _ = strconv.Atoi(maxParts)
		}
		if partNumberMarker := query.Get("part-number-marker"); partNumberMarker != "" {
			input.PartNumberMarker = partNumberMarker
		}
		result, err := o.ListParts(ctx, reqCtx, input)
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
		result, err := o.CompleteMultipartUpload(ctx, reqCtx, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		if result.ServerSideEncryption != "" {
			header.Set("x-amz-server-side-encryption", result.ServerSideEncryption)
		}
		if result.SSEKMSKeyId != "" {
			header.Set("x-amz-server-side-encryption-aws-kms-key-id", result.SSEKMSKeyId)
		}
		return result, header, http.StatusOK, nil

	case method == "DELETE" && query.Has("uploadId"):
		err := o.AbortMultipartUpload(ctx, reqCtx, &AbortMultipartUploadInput{
			Bucket:   bucket,
			Key:      key,
			UploadId: query.Get("uploadId"),
		})
		return nil, header, http.StatusNoContent, err

	case method == "GET" && key == "" && query.Has("list-type") && query.Get("list-type") == "2":
		input := &ListObjectsV2Input{
			Bucket:            bucket,
			Delimiter:         query.Get("delimiter"),
			Prefix:            query.Get("prefix"),
			ContinuationToken: query.Get("continuation-token"),
			StartAfter:        query.Get("start-after"),
		}
		if maxKeys := query.Get("max-keys"); maxKeys != "" {
			input.MaxKeys, _ = strconv.Atoi(maxKeys)
		}
		result, err := o.ListObjectsV2(ctx, reqCtx, input)
		return result, header, http.StatusOK, err

	case method == "GET" && key == "":
		input := &ListObjectsInput{
			Bucket:    bucket,
			Delimiter: query.Get("delimiter"),
			Prefix:    query.Get("prefix"),
			Marker:    query.Get("marker"),
		}
		if maxKeys := query.Get("max-keys"); maxKeys != "" {
			input.MaxKeys, _ = strconv.Atoi(maxKeys)
		}
		result, err := o.ListObjects(ctx, reqCtx, input)
		return result, header, http.StatusOK, err

	case method == "GET" && query.Has("tagging"):
		result, err := o.GetObjectTagging(ctx, reqCtx, &GetObjectTaggingInput{Bucket: bucket, Key: key})
		return result, header, http.StatusOK, err

	case method == "GET" && query.Has("legal-hold"):
		result, err := o.GetObjectLegalHold(ctx, reqCtx, &GetObjectLegalHoldInput{
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
		err := o.PutObjectLegalHold(ctx, reqCtx, &PutObjectLegalHoldInput{
			Bucket:    bucket,
			Key:       key,
			VersionId: query.Get("versionId"),
			LegalHold: &legalHoldReq,
		})
		return nil, header, http.StatusOK, err

	case method == "GET" && query.Has("retention"):
		result, err := o.GetObjectRetention(ctx, reqCtx, &GetObjectRetentionInput{
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
		err := o.PutObjectRetention(ctx, reqCtx, &PutObjectRetentionInput{
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
		err := o.PutObjectAcl(ctx, reqCtx, input)
		return nil, header, http.StatusOK, err

	case method == "GET" && query.Has("acl"):
		result, err := o.GetObjectAcl(ctx, reqCtx, bucket, key, query.Get("versionId"))
		return result, header, http.StatusOK, err

	case method == "GET" && query.Has("attributes"):
		objectAttributes := r.Header["X-Amz-Object-Attributes"]
		maxParts := int32(1000)
		if mp := query.Get("max-parts"); mp != "" {
			if parsed, err := strconv.ParseInt(mp, 10, 32); err == nil {
				maxParts = int32(parsed)
			}
		}
		result, err := o.GetObjectAttributes(ctx, reqCtx, &GetObjectAttributesInput{
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
		result, err := o.GetObject(ctx, reqCtx, input)
		if err != nil {
			return nil, header, http.StatusNotFound, err
		}
		header.Set("ETag", result.ETag)
		header.Set("Content-Length", strconv.FormatInt(result.ContentLength, 10))
		header.Set("Content-Type", result.ContentType)
		header.Set("Last-Modified", result.LastModified.Format(http.TimeFormat))
		if result.VersionId != "" && result.VersionId != "null" {
			header.Set("x-amz-version-id", result.VersionId)
		}
		if result.SSECustomerAlgorithm != "" {
			header.Set("x-amz-server-side-encryption-customer-algorithm", result.SSECustomerAlgorithm)
		}
		if result.SSECustomerKeyMD5 != "" {
			header.Set("x-amz-server-side-encryption-customer-key-MD5", result.SSECustomerKeyMD5)
		}
		if result.ServerSideEncryption != "" {
			header.Set("x-amz-server-side-encryption", result.ServerSideEncryption)
			if result.SSEKMSKeyId != "" {
				header.Set("x-amz-server-side-encryption-aws-kms-key-id", result.SSEKMSKeyId)
			}
		}
		for k, v := range result.Metadata {
			header.Set("x-amz-meta-"+k, v)
		}
		return result, header, http.StatusOK, nil

	case method == "HEAD":
		result, err := o.HeadObject(ctx, reqCtx, &HeadObjectInput{Bucket: bucket, Key: key, VersionId: query.Get("versionId")})
		if err != nil {
			return nil, header, http.StatusNotFound, err
		}
		header.Set("ETag", result.ETag)
		header.Set("Content-Length", strconv.FormatInt(result.ContentLength, 10))
		header.Set("Content-Type", result.ContentType)
		header.Set("Last-Modified", result.LastModified.Format(http.TimeFormat))
		if result.VersionId != "" && result.VersionId != "null" {
			header.Set("x-amz-version-id", result.VersionId)
		}
		if result.ServerSideEncryption != "" {
			header.Set("x-amz-server-side-encryption", result.ServerSideEncryption)
			if result.SSEKMSKeyId != "" {
				header.Set("x-amz-server-side-encryption-aws-kms-key-id", result.SSEKMSKeyId)
			}
		}
		if result.CacheControl != "" {
			header.Set("Cache-Control", result.CacheControl)
		}
		if result.ContentDisposition != "" {
			header.Set("Content-Disposition", result.ContentDisposition)
		}
		if result.ContentEncoding != "" {
			header.Set("Content-Encoding", result.ContentEncoding)
		}
		if result.ContentLanguage != "" {
			header.Set("Content-Language", result.ContentLanguage)
		}
		if result.StorageClass != "" {
			header.Set("x-amz-storage-class", result.StorageClass)
		}
		for k, v := range result.Metadata {
			header.Set("x-amz-meta-"+k, v)
		}
		return result, header, http.StatusOK, nil

	case method == "PUT" && query.Has("tagging"):
		var tagSet struct {
			Tags []Tag `xml:"TagSet>Tag"`
		}
		if err := request.NewSafeXMLDecoder(r.Body).Decode(&tagSet); err != nil {
			return nil, header, http.StatusBadRequest, err
		}
		err := o.PutObjectTagging(ctx, reqCtx, &PutObjectTaggingInput{
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
			ServerSideEncryption:      r.Header.Get("x-amz-server-side-encryption"),
			SSEKMSKeyId:               r.Header.Get("x-amz-server-side-encryption-aws-kms-key-id"),
			SSECustomerAlgorithm:      r.Header.Get("x-amz-server-side-encryption-customer-algorithm"),
			SSECustomerKey:            r.Header.Get("x-amz-server-side-encryption-customer-key"),
			SSECustomerKeyMD5:         r.Header.Get("x-amz-server-side-encryption-customer-key-MD5"),
			CopySourceSSECustomerAlgo: r.Header.Get("x-amz-copy-source-server-side-encryption-customer-algorithm"),
			CopySourceSSECustomerKey:  r.Header.Get("x-amz-copy-source-server-side-encryption-customer-key"),
			CopySourceSSECustomerMD5:  r.Header.Get("x-amz-copy-source-server-side-encryption-customer-key-MD5"),
		}
		result, err := o.CopyObject(ctx, reqCtx, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		if result.ServerSideEncryption != "" {
			header.Set("x-amz-server-side-encryption", result.ServerSideEncryption)
		}
		if result.SSEKMSKeyId != "" {
			header.Set("x-amz-server-side-encryption-aws-kms-key-id", result.SSEKMSKeyId)
		}
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
			if parseErr != nil {
				return nil, header, http.StatusBadRequest, fmt.Errorf("invalid Content-Length: %s", contentLengthStr)
			}
		}

		metadata := make(map[string]string)
		for k, v := range r.Header {
			if strings.HasPrefix(k, "X-Amz-Meta-") {
				metaKey := strings.TrimPrefix(k, "X-Amz-Meta-")
				metadata[metaKey] = v[0]
			}
		}
		input := &PutObjectInput{
			Bucket:               bucket,
			Key:                  key,
			Body:                 r.Body,
			ContentLength:        contentLength,
			ContentType:          contentType,
			Metadata:             metadata,
			IfMatch:              r.Header.Get("If-Match"),
			IfNoneMatch:          r.Header.Get("If-None-Match"),
			ServerSideEncryption: r.Header.Get("x-amz-server-side-encryption"),
			SSEKMSKeyId:          r.Header.Get("x-amz-server-side-encryption-aws-kms-key-id"),
			SSECustomerAlgorithm: r.Header.Get("x-amz-server-side-encryption-customer-algorithm"),
			SSECustomerKey:       r.Header.Get("x-amz-server-side-encryption-customer-key"),
			SSECustomerKeyMD5:    r.Header.Get("x-amz-server-side-encryption-customer-key-md5"),
		}
		result, err := o.PutObject(ctx, reqCtx, input)
		if err != nil {
			return nil, header, http.StatusInternalServerError, err
		}
		header.Set("ETag", result.ETag)
		if result.VersionId != "" && result.VersionId != "null" {
			header.Set("x-amz-version-id", result.VersionId)
		}
		if input.SSECustomerAlgorithm != "" {
			header.Set("x-amz-server-side-encryption-customer-algorithm", input.SSECustomerAlgorithm)
		}
		if input.SSECustomerKeyMD5 != "" {
			header.Set("x-amz-server-side-encryption-customer-key-MD5", input.SSECustomerKeyMD5)
		}
		if result.ServerSideEncryption != "" {
			header.Set("x-amz-server-side-encryption", result.ServerSideEncryption)
			if result.SSEKMSKeyId != "" {
				header.Set("x-amz-server-side-encryption-aws-kms-key-id", result.SSEKMSKeyId)
			}
		}
		for k, v := range metadata {
			header.Set("x-amz-meta-"+k, v)
		}
		return result, header, http.StatusOK, nil

	case method == "DELETE" && query.Has("tagging"):
		err := o.DeleteObjectTagging(ctx, reqCtx, &DeleteObjectTaggingInput{Bucket: bucket, Key: key})
		return nil, header, http.StatusNoContent, err

	case method == "DELETE":
		result, err := o.DeleteObject(ctx, reqCtx, &DeleteObjectInput{Bucket: bucket, Key: key, VersionId: query.Get("versionId")})
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
		return nil, header, http.StatusMethodNotAllowed, fmt.Errorf("unsupported operation")
	}
}
