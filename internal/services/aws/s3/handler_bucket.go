package s3

import (
	"fmt"
	"net/http"

	"vorpalstacks/internal/common/request"
)

// handleServiceRequest handles service-level S3 operations (e.g. ListBuckets).
func (h *S3Handler) handleServiceRequest(ctx *request.RequestContext, r *http.Request) (interface{}, int, error) {
	if r.Method == "GET" {
		result, err := h.bucketOps.ListBuckets(ctx, &ListBucketsInput{})
		return result, http.StatusOK, err
	}
	return nil, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed")
}

// handleBucketRequest dispatches bucket-level S3 operations based on HTTP method.
// Resolves the bucket store once, then delegates to method-specific dispatchers.
func (h *S3Handler) handleBucketRequest(ctx *request.RequestContext, r *http.Request, bucket string) (interface{}, int, error) {
	stores, err := h.svc.store(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	switch r.Method {
	case "PUT":
		return h.dispatchPutBucket(ctx, r, bucket, stores)
	case "GET":
		return h.dispatchGetBucket(ctx, r, bucket, stores)
	case "HEAD":
		return h.headBucket(ctx, r, bucket, stores)
	case "DELETE":
		return h.dispatchDeleteBucket(ctx, r, bucket, stores)
	default:
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed")
	}
}
