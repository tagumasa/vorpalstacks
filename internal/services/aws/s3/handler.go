package s3

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"vorpalstacks/internal/common/audit"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/utils/crypto"
)

// S3Handler routes incoming S3 HTTP requests to the appropriate bucket or
// object operation based on the request path and query parameters.
type S3Handler struct {
	svc            *S3Service
	bucketOps      *BucketOperations
	objectOps      *ObjectOperations
	region         string
	storageManager *storage.RegionStorageManager
	auditRecorder  request.AuditRecorder
}

// NewS3Handler creates a new S3Handler for the given service, region, and storage manager.
func NewS3Handler(svc *S3Service, region string, storageMgr *storage.RegionStorageManager) *S3Handler {
	return &S3Handler{
		svc:            svc,
		bucketOps:      NewBucketOperations(svc),
		objectOps:      NewObjectOperations(svc),
		region:         region,
		storageManager: storageMgr,
	}
}

func (h *S3Handler) newRequestContext(r *http.Request) *request.RequestContext {
	region := request.ExtractRegionFromAuth(r.Header.Get("Authorization"))
	if region == "" {
		region = h.region
	}
	ctx := request.NewRequestContext(r.Context(), h.storageManager, h.svc.accountID, region)
	ctx.SourceIP = extractSourceIP(r)
	ctx.UserAgent = r.UserAgent()
	ctx.Principal = h.svc.accountID
	ctx.PrincipalID = h.svc.accountID
	ctx.PrincipalType = request.PrincipalTypeUser
	return ctx
}

// SetAuditRecorder sets the audit recorder for CloudTrail event logging.
func (h *S3Handler) SetAuditRecorder(recorder request.AuditRecorder) {
	h.auditRecorder = recorder
}

func extractSourceIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	idx := strings.LastIndex(r.RemoteAddr, ":")
	if idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

// ServeHTTP implements http.Handler, dispatching S3 API requests to bucket or
// object handlers and writing the response.
func (h *S3Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := h.newRequestContext(r)

	requestID := fmt.Sprintf("%016X", time.Now().UnixNano())
	w.Header().Set("x-amz-request-id", requestID)
	w.Header().Set("x-amz-id-2", requestID)

	bucket, key := parseS3Path(r.URL.Path)
	query := r.URL.Query()

	if h.isPresignedURLRequest(query) {
		if err := h.verifyPresignedURL(r, bucket); err != nil {
			h.writeError(w, err, bucket, key)
			return
		}
	}

	var result interface{}
	var statusCode int
	var err error
	var header http.Header

	switch {
	case query.Has("delete") && bucket != "" && key == "":
		result, statusCode, err = h.handleDeleteObjects(reqCtx, r, bucket, r.Body)
	case bucket == "":
		result, statusCode, err = h.handleServiceRequest(reqCtx, r)
	case key == "":
		result, statusCode, err = h.handleBucketRequest(reqCtx, r, bucket)
	default:
		result, header, statusCode, err = h.handleObjectRequest(reqCtx, r, bucket, key)
	}

	if err != nil {
		h.recordAudit(determineS3EventName(r, bucket, key), reqCtx, r, nil, err)
		h.writeError(w, err, bucket, key)
		return
	}

	for k, v := range header {
		w.Header()[k] = v
	}

	h.recordAudit(determineS3EventName(r, bucket, key), reqCtx, r, result, nil)
	h.writeResult(w, result, statusCode)
}

func parseS3Path(urlPath string) (bucket, key string) {
	urlPath = strings.TrimPrefix(urlPath, "/")
	if urlPath == "" {
		return "", ""
	}

	parts := strings.SplitN(urlPath, "/", 2)
	bucket = parts[0]
	if len(parts) > 1 {
		key = parts[1]
	}
	return bucket, key
}

func (h *S3Handler) isPresignedURLRequest(query url.Values) bool {
	return query.Has("X-Amz-Algorithm") &&
		query.Has("X-Amz-Credential") &&
		query.Has("X-Amz-Signature")
}

func (h *S3Handler) verifyPresignedURL(r *http.Request, bucket string) error {
	if h.svc.credentialsProvider == nil {
		return nil
	}
	verifier := crypto.NewPresignedURLVerifier(h.svc.credentialsProvider)
	return verifier.VerifyPresignedURL(r, bucket, h.region)
}

func (h *S3Handler) recordAudit(eventName string, reqCtx *request.RequestContext, r *http.Request, response interface{}, err error) {
	if h.auditRecorder == nil {
		return
	}
	builder := audit.NewEventBuilder("s3", eventName, reqCtx, nil)
	event := builder.Build(response, err)
	if recorder, ok := h.auditRecorder.(audit.Recorder); ok {
		_ = recorder.RecordEvent(event)
	}
}
