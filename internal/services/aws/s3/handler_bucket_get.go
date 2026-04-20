package s3

import (
	"net/http"
	"strconv"

	"vorpalstacks/internal/common/request"
)

// dispatchGetBucket handles GET requests targeting a bucket by dispatching
// to the appropriate sub-resource operation based on the query string.
// Falls through to ListObjectsV1 when no query parameter matches.
func (h *S3Handler) dispatchGetBucket(ctx *request.RequestContext, r *http.Request, bucket string, stores *s3Stores) (interface{}, int, error) {
	query := r.URL.Query()

	if query.Has("acl") {
		action := "s3:GetBucketAcl"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketAcl(ctx, bucket)
		return result, http.StatusOK, err
	}
	if query.Has("versioning") {
		action := "s3:GetBucketVersioning"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketVersioning(ctx, &GetBucketVersioningInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("encryption") {
		action := "s3:GetEncryptionConfiguration"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketEncryption(ctx, &GetBucketEncryptionInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("policy") {
		action := "s3:GetBucketPolicy"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketPolicy(ctx, &GetBucketPolicyInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("cors") {
		action := "s3:GetBucketCORS"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketCORS(ctx, &GetBucketCORSInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("tagging") {
		action := "s3:GetBucketTagging"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketTagging(ctx, &GetBucketTaggingInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("lifecycle") {
		action := "s3:GetLifecycleConfiguration"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketLifecycleConfiguration(ctx, &GetBucketLifecycleConfigurationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("website") {
		action := "s3:GetBucketWebsite"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketWebsite(ctx, &GetBucketWebsiteInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("replication") {
		return nil, http.StatusNotFound, ErrNoSuchReplication
	}
	if query.Has("object-lock") {
		action := "s3:GetBucketObjectLockConfiguration"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetObjectLockConfiguration(ctx, &GetObjectLockConfigurationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("notification") {
		action := "s3:GetBucketNotification"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketNotificationConfiguration(ctx, &GetBucketNotificationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("logging") {
		action := "s3:GetBucketLogging"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketLogging(ctx, &GetBucketLoggingInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("ownershipControls") {
		action := "s3:GetBucketOwnershipControls"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketOwnershipControls(ctx, &GetBucketOwnershipControlsInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("requestPayment") {
		action := "s3:GetBucketRequestPayment"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketRequestPayment(ctx, &GetBucketRequestPaymentInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("accelerate") {
		action := "s3:GetAccelerateConfiguration"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketAccelerateConfiguration(ctx, &GetBucketAccelerateConfigurationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("location") {
		action := "s3:GetBucketLocation"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetBucketLocation(ctx, &GetBucketLocationInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("publicAccessBlock") {
		action := "s3:GetBucketPublicAccessBlock"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		result, err := h.bucketOps.GetPublicAccessBlock(ctx, &GetPublicAccessBlockInput{Bucket: bucket})
		return result, http.StatusOK, err
	}
	if query.Has("versions") {
		action := "s3:ListBucketVersions"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
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
				return nil, http.StatusBadRequest, NewInvalidArgumentError("Provided max-keys not an integer")
			}
			input.MaxKeys = mk
		}
		result, err := h.objectOps.ListObjectVersions(r.Context(), ctx, input)
		return result, http.StatusOK, err
	}
	if query.Has("uploads") {
		action := "s3:ListBucketMultipartUploads"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
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
				return nil, http.StatusBadRequest, NewInvalidArgumentError("Provided max-uploads not an integer")
			}
			input.MaxUploads = mu
		}
		result, err := h.objectOps.ListMultipartUploads(r.Context(), ctx, input)
		return result, http.StatusOK, err
	}
	if query.Has("list-type") && query.Get("list-type") == "2" {
		action := "s3:ListBucket"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		input := &ListObjectsV2Input{
			Bucket:            bucket,
			Delimiter:         query.Get("delimiter"),
			Prefix:            query.Get("prefix"),
			ContinuationToken: query.Get("continuation-token"),
			StartAfter:        query.Get("start-after"),
		}
		if maxKeys := query.Get("max-keys"); maxKeys != "" {
			mk, err := strconv.Atoi(maxKeys)
			if err != nil {
				return nil, http.StatusBadRequest, NewInvalidArgumentError("Provided max-keys not an integer")
			}
			input.MaxKeys = mk
		}
		result, err := h.objectOps.ListObjectsV2(r.Context(), ctx, input)
		return result, http.StatusOK, err
	}

	action := "s3:ListBucket"
	if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
		return nil, http.StatusForbidden, err
	}
	input := &ListObjectsInput{
		Bucket:    bucket,
		Delimiter: query.Get("delimiter"),
		Prefix:    query.Get("prefix"),
		Marker:    query.Get("marker"),
	}
	if maxKeys := query.Get("max-keys"); maxKeys != "" {
		mk, err := strconv.Atoi(maxKeys)
		if err != nil {
			return nil, http.StatusBadRequest, NewInvalidArgumentError("Provided max-keys not an integer")
		}
		input.MaxKeys = mk
	}
	result, err := h.objectOps.ListObjects(r.Context(), ctx, input)
	return result, http.StatusOK, err
}

// headBucket handles HEAD requests for bucket existence checks.
func (h *S3Handler) headBucket(ctx *request.RequestContext, r *http.Request, bucket string, stores *s3Stores) (interface{}, int, error) {
	action := "s3:ListBucket"
	if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
		return nil, http.StatusForbidden, err
	}
	result, err := h.bucketOps.HeadBucket(ctx, &HeadBucketInput{Bucket: bucket})
	return result, http.StatusOK, err
}
