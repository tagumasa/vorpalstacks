package s3

import (
	"net/http"

	"vorpalstacks/internal/common/request"
)

// dispatchDeleteBucket handles DELETE requests targeting a bucket by dispatching
// to the appropriate sub-resource operation based on the query string.
// Falls through to DeleteBucket when no query parameter matches.
func (h *S3Handler) dispatchDeleteBucket(ctx *request.RequestContext, r *http.Request, bucket string, stores *s3Stores) (interface{}, int, error) {
	query := r.URL.Query()

	if query.Has("encryption") {
		action := "s3:DeleteEncryptionConfiguration"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.DeleteBucketEncryption(ctx, &DeleteBucketEncryptionInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("policy") {
		action := "s3:DeleteBucketPolicy"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.DeleteBucketPolicy(ctx, &DeleteBucketPolicyInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("cors") {
		action := "s3:DeleteBucketCORS"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.DeleteBucketCORS(ctx, &DeleteBucketCORSInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("tagging") {
		action := "s3:DeleteBucketTagging"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.DeleteBucketTagging(ctx, &DeleteBucketTaggingInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("lifecycle") {
		action := "s3:DeleteLifecycleConfiguration"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.DeleteBucketLifecycleConfiguration(ctx, &DeleteBucketLifecycleConfigurationInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("website") {
		action := "s3:DeleteBucketWebsite"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.DeleteBucketWebsite(ctx, &DeleteBucketWebsiteInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("ownershipControls") {
		action := "s3:DeleteBucketOwnershipControls"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.DeleteBucketOwnershipControls(ctx, &DeleteBucketOwnershipControlsInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("publicAccessBlock") {
		action := "s3:DeleteBucketPublicAccessBlock"
		if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
			return nil, http.StatusForbidden, err
		}
		err := h.bucketOps.DeletePublicAccessBlock(ctx, &DeletePublicAccessBlockInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}

	action := "s3:DeleteBucket"
	if err := h.checkAccess(ctx, r, stores, action, bucket, ""); err != nil {
		return nil, http.StatusForbidden, err
	}
	err := h.bucketOps.DeleteBucket(ctx, &DeleteBucketInput{Bucket: bucket})
	return nil, http.StatusNoContent, err
}
