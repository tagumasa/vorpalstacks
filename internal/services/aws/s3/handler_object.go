package s3

import (
	"io"
	"net/http"

	"vorpalstacks/internal/common/request"
)

func (h *S3Handler) handleObjectRequest(ctx *request.RequestContext, r *http.Request, bucket, key string, actions []s3ActionSpec) (interface{}, http.Header, int, error) {
	stores, err := h.svc.store(ctx)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	if len(actions) > 0 {
		for _, a := range actions {
			if err := h.checkAccess(ctx, r, stores, a.Action, a.Bucket, a.Key); err != nil {
				return nil, nil, http.StatusForbidden, err
			}
		}
	}

	return h.objectOps.HandleRequest(r.Context(), ctx, stores, r, bucket, key)
}

func (h *S3Handler) handleDeleteObjects(ctx *request.RequestContext, r *http.Request, bucket string, body io.Reader) (interface{}, int, error) {
	stores, err := h.svc.store(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := h.checkAccess(ctx, r, stores, "s3:DeleteObject", bucket, ""); err != nil {
		return nil, http.StatusForbidden, err
	}

	var deleteReq Delete
	if err := request.NewSafeXMLDecoder(body).Decode(&deleteReq); err != nil {
		return nil, http.StatusBadRequest, err
	}

	for _, obj := range deleteReq.Objects {
		if err := h.checkAccess(ctx, r, stores, "s3:DeleteObject", bucket, obj.Key); err != nil {
			return nil, http.StatusForbidden, err
		}
	}

	result, err := h.objectOps.DeleteObjects(r.Context(), ctx, stores, &DeleteObjectsInput{
		Bucket:                    bucket,
		Delete:                    &deleteReq,
		BypassGovernanceRetention: r.Header.Get("x-amz-bypass-governance-retention") == "true",
	})
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return result, http.StatusOK, nil
}
