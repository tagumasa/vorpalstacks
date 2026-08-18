package s3

import (
	"net/http"

	"vorpalstacks/internal/common/request"
)

// dispatchGetBucket handles GET requests targeting a bucket by dispatching
// to the appropriate sub-resource operation based on the query string.
// Falls through to ListObjectsV1 when no query parameter matches.
// IAM action checks run in handleBucketRequest before dispatch.
func (h *S3Handler) dispatchGetBucket(ctx *request.RequestContext, r *http.Request, bucket string, stores *s3Stores) (interface{}, int, error) {
	query := r.URL.Query()

	if query.Has("acl") {
		result, err := h.bucketOps.GetBucketAcl(ctx, bucket)
		return result, http.StatusOK, err
	}
	if query.Has("versioning") {
		result, err := h.bucketOps.GetBucketVersioning(ctx, &GetBucketVersioningInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("encryption") {
		result, err := h.bucketOps.GetBucketEncryption(ctx, &GetBucketEncryptionInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("policy") {
		result, err := h.bucketOps.GetBucketPolicy(ctx, &GetBucketPolicyInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("policyStatus") {
		result, err := h.bucketOps.GetBucketPolicyStatus(ctx, &GetBucketPolicyStatusInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("cors") {
		result, err := h.bucketOps.GetBucketCORS(ctx, &GetBucketCORSInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("tagging") {
		result, err := h.bucketOps.GetBucketTagging(ctx, &GetBucketTaggingInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("lifecycle") {
		result, err := h.bucketOps.GetBucketLifecycleConfiguration(ctx, &GetBucketLifecycleConfigurationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("website") {
		result, err := h.bucketOps.GetBucketWebsite(ctx, &GetBucketWebsiteInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("replication") {
		result, err := h.bucketOps.GetBucketReplication(ctx, &GetBucketReplicationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("object-lock") {
		result, err := h.bucketOps.GetObjectLockConfiguration(ctx, &GetObjectLockConfigurationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("notification") {
		result, err := h.bucketOps.GetBucketNotificationConfiguration(ctx, &GetBucketNotificationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("logging") {
		result, err := h.bucketOps.GetBucketLogging(ctx, &GetBucketLoggingInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("ownershipControls") {
		result, err := h.bucketOps.GetBucketOwnershipControls(ctx, &GetBucketOwnershipControlsInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("requestPayment") {
		result, err := h.bucketOps.GetBucketRequestPayment(ctx, &GetBucketRequestPaymentInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("accelerate") {
		result, err := h.bucketOps.GetBucketAccelerateConfiguration(ctx, &GetBucketAccelerateConfigurationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("location") {
		result, err := h.bucketOps.GetBucketLocation(ctx, &GetBucketLocationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("publicAccessBlock") {
		result, err := h.bucketOps.GetPublicAccessBlock(ctx, &GetPublicAccessBlockInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("versions") {
		input := &ListObjectVersionsInput{
			Bucket:          bucket,
			Delimiter:       query.Get("delimiter"),
			Prefix:          query.Get("prefix"),
			KeyMarker:       query.Get("key-marker"),
			VersionIdMarker: query.Get("version-id-marker"),
			EncodingType:    query.Get("encoding-type"),
		}
		maxKeys, err := parseListLimit(query, "max-keys", s3MaxKeys)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		input.MaxKeys = maxKeys
		result, err := h.objectOps.ListObjectVersions(r.Context(), ctx, stores, input)
		return result, http.StatusOK, err
	}
	if query.Has("uploads") {
		input := &ListMultipartUploadsInput{
			Bucket:         bucket,
			Delimiter:      query.Get("delimiter"),
			Prefix:         query.Get("prefix"),
			KeyMarker:      query.Get("key-marker"),
			UploadIdMarker: query.Get("upload-id-marker"),
		}
		maxUploads, err := parseListLimit(query, "max-uploads", s3MaxUploads)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		input.MaxUploads = maxUploads
		result, err := h.objectOps.ListMultipartUploads(r.Context(), ctx, stores, input)
		return result, http.StatusOK, err
	}
	if query.Has("list-type") && query.Get("list-type") == "2" {
		input := &ListObjectsV2Input{
			Bucket:            bucket,
			Delimiter:         query.Get("delimiter"),
			Prefix:            query.Get("prefix"),
			ContinuationToken: query.Get("continuation-token"),
			StartAfter:        query.Get("start-after"),
			EncodingType:      query.Get("encoding-type"),
		}
		maxKeys, err := parseListLimit(query, "max-keys", s3MaxKeys)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
		input.MaxKeys = maxKeys
		result, err := h.objectOps.ListObjectsV2(r.Context(), ctx, stores, input)
		return result, http.StatusOK, err
	}

	input := &ListObjectsInput{
		Bucket:       bucket,
		Delimiter:    query.Get("delimiter"),
		Prefix:       query.Get("prefix"),
		Marker:       query.Get("marker"),
		EncodingType: query.Get("encoding-type"),
	}
	maxKeys, err := parseListLimit(query, "max-keys", s3MaxKeys)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	input.MaxKeys = maxKeys
	result, err := h.objectOps.ListObjects(r.Context(), ctx, stores, input)
	return result, http.StatusOK, err
}

// headBucket handles HEAD requests for bucket existence checks.
func (h *S3Handler) headBucket(ctx *request.RequestContext, r *http.Request, bucket string, stores *s3Stores) (interface{}, int, error) {
	result, err := h.bucketOps.HeadBucket(ctx, &HeadBucketInput{Bucket: bucket})
	return result, http.StatusOK, err
}
