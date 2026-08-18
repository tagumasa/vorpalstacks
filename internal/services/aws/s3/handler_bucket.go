package s3

import (
	"fmt"
	"net/http"
	"strconv"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
)

// handleServiceRequest handles service-level S3 operations (e.g. ListBuckets).
func (h *S3Handler) handleServiceRequest(ctx *request.RequestContext, r *http.Request) (interface{}, int, error) {
	if r.Method == "GET" {
		stores, err := h.svc.store(ctx)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if err := h.checkAccess(ctx, r, stores, "s3:ListAllMyBuckets", "", ""); err != nil {
			return nil, http.StatusForbidden, err
		}

		input := &ListBucketsInput{
			ContinuationToken: r.URL.Query().Get("continuation-token"),
			Prefix:            r.URL.Query().Get("prefix"),
			BucketRegion:      r.URL.Query().Get("bucket-region"),
		}
		if mb := r.URL.Query().Get("max-buckets"); mb != "" {
			n, err := strconv.Atoi(mb)
			if err != nil {
				return nil, http.StatusBadRequest, NewInvalidArgumentError(fmt.Sprintf("invalid max-buckets value: %s", mb))
			}
			input.MaxBuckets = n
		}
		result, err := h.bucketOps.ListBuckets(ctx, input)
		return result, http.StatusOK, err
	}
	return nil, http.StatusMethodNotAllowed, awserrors.NewAWSError("MethodNotAllowed", "The specified method is not allowed against this resource.", http.StatusMethodNotAllowed)
}

// handleBucketRequest dispatches bucket-level S3 operations based on HTTP method.
// Resolves the bucket store once, checks the IAM actions required by the
// classified operation, then delegates to method-specific dispatchers.
func (h *S3Handler) handleBucketRequest(ctx *request.RequestContext, r *http.Request, bucket string) (interface{}, int, error) {
	stores, err := h.svc.store(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if _, actions := classifyBucketRequest(r, bucket); len(actions) > 0 {
		for _, a := range actions {
			if err := h.checkAccess(ctx, r, stores, a.Action, a.Bucket, a.Key); err != nil {
				return nil, http.StatusForbidden, err
			}
		}
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
		return nil, http.StatusMethodNotAllowed, awserrors.NewAWSError("MethodNotAllowed", "The specified method is not allowed against this resource.", http.StatusMethodNotAllowed)
	}
}
