package s3

import (
	"net/http"

	"vorpalstacks/internal/common/request"
)

// dispatchDeleteBucket handles DELETE requests targeting a bucket by dispatching
// to the appropriate sub-resource operation based on the query string.
// Falls through to DeleteBucket when no query parameter matches.
// IAM action checks run in handleBucketRequest before dispatch.
func (h *S3Handler) dispatchDeleteBucket(ctx *request.RequestContext, r *http.Request, bucket string, stores *s3Stores) (interface{}, int, error) {
	query := r.URL.Query()

	if query.Has("encryption") {
		err := h.bucketOps.DeleteBucketEncryption(ctx, &DeleteBucketEncryptionInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("policy") {
		err := h.bucketOps.DeleteBucketPolicy(ctx, &DeleteBucketPolicyInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("cors") {
		err := h.bucketOps.DeleteBucketCORS(ctx, &DeleteBucketCORSInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("tagging") {
		err := h.bucketOps.DeleteBucketTagging(ctx, &DeleteBucketTaggingInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("lifecycle") {
		err := h.bucketOps.DeleteBucketLifecycleConfiguration(ctx, &DeleteBucketLifecycleConfigurationInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("website") {
		err := h.bucketOps.DeleteBucketWebsite(ctx, &DeleteBucketWebsiteInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("ownershipControls") {
		err := h.bucketOps.DeleteBucketOwnershipControls(ctx, &DeleteBucketOwnershipControlsInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("publicAccessBlock") {
		err := h.bucketOps.DeletePublicAccessBlock(ctx, &DeletePublicAccessBlockInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}
	if query.Has("replication") {
		err := h.bucketOps.DeleteBucketReplication(ctx, &DeleteBucketReplicationInput{Bucket: bucket})
		return nil, http.StatusNoContent, err
	}

	err := h.bucketOps.DeleteBucket(ctx, &DeleteBucketInput{Bucket: bucket})
	return nil, http.StatusNoContent, err
}
