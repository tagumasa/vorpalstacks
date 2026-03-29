// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/utils/crypto"
)

// S3Handler handles S3 HTTP requests.
type S3Handler struct {
	svc            *S3Service
	bucketOps      *BucketOperations
	objectOps      *ObjectOperations
	region         string
	storageManager *storage.RegionStorageManager
}

// NewS3Handler creates a new S3 handler.
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
	ctx := request.NewRequestContext(r.Context(), h.storageManager, h.svc.accountID, h.region)
	ctx.SourceIP = extractSourceIP(r)
	ctx.UserAgent = r.UserAgent()
	ctx.Principal = h.svc.accountID
	ctx.PrincipalID = h.svc.accountID
	ctx.PrincipalType = request.PrincipalTypeUser
	return ctx
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

// ServeHTTP implements the http.Handler interface to handle S3 HTTP requests.
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
		h.writeError(w, err, bucket, key)
		return
	}

	for k, v := range header {
		w.Header()[k] = v
	}

	switch v := result.(type) {
	case *SelectObjectContentOutput:
		if v.Payload != nil {
			defer v.Payload.Close()
			if _, err := io.Copy(w, v.Payload); err != nil {
				logs.Error("S3: failed to stream SelectObjectContent payload", logs.Err(err))
			}
		}
		return
	case *GetObjectOutput:
		if v.AcceptRanges != "" {
			w.Header().Set("Accept-Ranges", v.AcceptRanges)
		}
		if v.ContentRange != "" {
			w.Header().Set("Content-Range", v.ContentRange)
		}
		if v.IsPartial {
			w.WriteHeader(http.StatusPartialContent)
		} else {
			w.WriteHeader(statusCode)
		}
		if v.Body != nil {
			defer v.Body.Close()
			if _, err := io.Copy(w, v.Body); err != nil {
				logs.Error("S3: failed to stream GetObject body", logs.Err(err))
			}
		}
		return
	case *HeadObjectOutput:
		w.WriteHeader(statusCode)
		return
	case *ListBucketsOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *ListObjectsOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *ListObjectsV2Output:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *ListObjectVersionsOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *ListMultipartUploadsOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *ListPartsOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *GetBucketAclOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *GetObjectAclOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
		_, _ = w.Write([]byte(v.ToXML()))
	case *CopyObjectOutput:
		h.writeXMLResponse(w, "CopyObjectResult", v.CopyObjectResult, statusCode, "")
	case *UploadPartCopyOutput:
		h.writeXMLResponse(w, "CopyPartResult", v.CopyPartResult, statusCode, "")
	case *DeleteObjectsOutput:
		h.writeXMLResponse(w, "DeleteResult", v, statusCode, "")
	case *GetBucketNotificationOutput:
		h.writeXMLResponse(w, "NotificationConfiguration", v.NotificationConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
	case *GetBucketLoggingOutput:
		if v.LoggingConfiguration != nil {
			h.writeXMLResponse(w, "BucketLoggingStatus", v, statusCode, "")
		} else {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><BucketLoggingStatus xmlns="http://s3.amazonaws.com/doc/2006-03-01/"></BucketLoggingStatus>`))
		}
	case *GetObjectAttributesOutput:
		h.writeXMLResponse(w, "GetObjectAttributesOutput", v, statusCode, "")
	case *GetBucketEncryptionOutput:
		if v.ServerSideEncryptionConfiguration != nil {
			h.writeXMLResponse(w, "ServerSideEncryptionConfiguration", v.ServerSideEncryptionConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
		} else {
			w.WriteHeader(statusCode)
		}
	case *CreateBucketOutput:
		w.Header().Set("Location", v.Location)
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><CreateBucketResult><Location>` + v.Location + `</Location></CreateBucketResult>`))
	case *HeadBucketOutput:
		w.Header().Set("x-amz-bucket-region", v.BucketRegion)
		w.WriteHeader(statusCode)
	case nil:
		w.WriteHeader(statusCode)
	default:
		h.writeXMLResponse(w, "", v, statusCode, "")
	}
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
	verifier := crypto.NewPresignedURLVerifier(h.svc.credentialsProvider)
	return verifier.VerifyPresignedURL(r, bucket, h.region)
}
