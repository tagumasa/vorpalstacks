package s3

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/logs"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// #nosec G705
func (h *S3Handler) writeXMLResponse(w http.ResponseWriter, rootElement string, data interface{}, statusCode int, xmlns string, requestID string) {
	var xmlData []byte
	var err error

	if rootElement != "" {
		inner, marshalErr := xml.Marshal(data)
		if marshalErr != nil {
			h.writeError(w, marshalErr, "", "", requestID)
			return
		}
		innerStr := xmlStripOuterTag(string(inner))
		if xmlns != "" {
			xmlData = []byte(fmt.Sprintf(`<%s xmlns="%s">%s</%s>`, rootElement, xmlns, innerStr, rootElement))
		} else {
			xmlData = []byte(fmt.Sprintf(`<%s>%s</%s>`, rootElement, innerStr, rootElement))
		}
	} else {
		xmlData, err = xml.Marshal(data)
	}

	if err != nil {
		h.writeError(w, err, "", "", requestID)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	_, _ = w.Write(xmlData)
}

// xmlStripOuterTag removes the outermost XML element produced by xml.Marshal.
// xml.Marshal wraps the value in a tag named after the Go struct (e.g.
// <GetBucketTaggingOutput>...</GetBucketTaggingOutput>). We strip that outer
// wrapper so the caller can re-wrap with the correct S3 root element name.
// Because xml.Marshal always produces well-formed XML with the struct name as
// the outermost tag, splitting at the first ">" and last "<" is safe — inner
// content is escaped (&gt;, &lt;) and never breaks the boundary markers.
func xmlStripOuterTag(xmlStr string) string {
	if idx := strings.Index(xmlStr, ">"); idx != -1 {
		xmlStr = xmlStr[idx+1:]
	}
	if idx := strings.LastIndex(xmlStr, "<"); idx != -1 {
		xmlStr = xmlStr[:idx]
	}
	return xmlStr
}

// #nosec G705
func (h *S3Handler) writeError(w http.ResponseWriter, err error, bucket, key, requestID string) {
	var awsErr *awserrors.AWSError

	switch {
	case errors.Is(err, storecommon.ErrNotFound):
		if bucket != "" && key == "" {
			awsErr = NewNoSuchBucketError(bucket)
		} else if key != "" {
			awsErr = NewNoSuchKeyError(key)
		} else {
			awsErr = ErrNoSuchBucket
		}
	case errors.Is(err, storecommon.ErrAlreadyExists):
		awsErr = NewBucketAlreadyExistsError(bucket)
	case errors.Is(err, storecommon.ErrConflict):
		awsErr = ErrBucketNotEmpty
	case errors.Is(err, storecommon.ErrInvalidInput):
		awsErr = ErrInvalidRequest
	case errors.Is(err, ErrPreconditionFailed):
		awsErr = ErrPreconditionFailed
	default:
		var castErr *awserrors.AWSError
		if errors.As(err, &castErr) {
			awsErr = castErr
		} else {
			awsErr = awserrors.NewAWSError("InternalError", err.Error(), http.StatusInternalServerError)
		}
	}

	// A delete-marker read carries the marker's response headers alongside
	// the error: the documented 405 (explicit versionId) sets Last-Modified,
	// and both the 405 and the latest-marker 404 report the marker
	// identification headers.
	var markerErr *deleteMarkerReadError
	if errors.As(err, &markerErr) {
		if awsErr.HTTPStatus == http.StatusMethodNotAllowed {
			w.Header().Set("Last-Modified", markerErr.marker.LastModified.UTC().Format(http.TimeFormat))
		}
		w.Header().Set("x-amz-delete-marker", "true")
		if markerErr.marker.VersionID != "" {
			w.Header().Set("x-amz-version-id", markerErr.marker.VersionID)
		}
	}

	resource := bucket
	if key != "" {
		resource = bucket + "/" + key
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(awsErr.HTTPStatus)
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	_, _ = w.Write([]byte(fmt.Sprintf(`<Error><Code>%s</Code><Message>%s</Message><Resource>%s</Resource><RequestId>%s</RequestId></Error>`,
		xmlEscape(awsErr.Code), xmlEscape(awsErr.Message), xmlEscape(resource), requestID)))
}

// writeToXMLResult emits an output that renders itself through ToXML.
func writeToXMLResult(w http.ResponseWriter, v interface{ ToXML() string }, statusCode int) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	_, _ = w.Write([]byte(v.ToXML()))
}

// writeResult serialises the operation result into the HTTP response.
// Streaming operations (GetObject) write directly; types that render
// themselves through ToXML fall through to the shared writer; all other
// results are rendered as XML or plain status codes.
func (h *S3Handler) writeResult(w http.ResponseWriter, result interface{}, statusCode int, requestID string) {
	switch v := result.(type) {
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
	case *HeadObjectOutput:
		if v.ContentRange != "" {
			w.Header().Set("Content-Range", v.ContentRange)
		}
		if v.IsPartial {
			w.WriteHeader(http.StatusPartialContent)
		} else {
			w.WriteHeader(statusCode)
		}
	case *CopyObjectOutput:
		h.writeXMLResponse(w, "CopyObjectResult", v.CopyObjectResult, statusCode, "", requestID)
	case *UploadPartCopyOutput:
		h.writeXMLResponse(w, "CopyPartResult", v.CopyPartResult, statusCode, "", requestID)
	case *DeleteObjectsOutput:
		h.writeXMLResponse(w, "DeleteResult", v, statusCode, "", requestID)
	case *GetBucketNotificationOutput:
		h.writeXMLResponse(w, "NotificationConfiguration", v.NotificationConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *GetBucketLoggingOutput:
		if v.LoggingConfiguration != nil {
			h.writeXMLResponse(w, "BucketLoggingStatus", v, statusCode, "", requestID)
		} else {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><BucketLoggingStatus xmlns="http://s3.amazonaws.com/doc/2006-03-01/"></BucketLoggingStatus>`))
		}
	case *GetObjectAttributesOutput:
		// VersionId and LastModified are header-bound members (the
		// operation's model): x-amz-version-id and Last-Modified. The
		// body carries only the requested attributes.
		if v.VersionId != "" && v.VersionId != "null" {
			w.Header().Set("x-amz-version-id", v.VersionId)
		}
		if !time.Time(v.LastModified).IsZero() {
			w.Header().Set("Last-Modified", time.Time(v.LastModified).UTC().Format(http.TimeFormat))
		}
		h.writeXMLResponse(w, "GetObjectAttributesOutput", v, statusCode, "", requestID)
	case *GetBucketEncryptionOutput:
		if v.ServerSideEncryptionConfiguration != nil {
			h.writeXMLResponse(w, "ServerSideEncryptionConfiguration", v.ServerSideEncryptionConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetBucketCORSOutput:
		if v.CORSConfiguration != nil {
			h.writeXMLResponse(w, "CORSConfiguration", v.CORSConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetBucketOwnershipControlsOutput:
		if v.OwnershipControls != nil {
			h.writeXMLResponse(w, "OwnershipControls", v.OwnershipControls, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetBucketRequestPaymentOutput:
		if v.RequestPaymentConfiguration != nil {
			h.writeXMLResponse(w, "RequestPaymentConfiguration", v.RequestPaymentConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetBucketAccelerateConfigurationOutput:
		if v.AccelerateConfiguration != nil {
			h.writeXMLResponse(w, "AccelerateConfiguration", v.AccelerateConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetBucketReplicationOutput:
		if v.ReplicationConfiguration != nil {
			h.writeXMLResponse(w, "ReplicationConfiguration", v.ReplicationConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetBucketPolicyOutput:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(v.Policy))
	case *GetBucketPolicyStatusOutput:
		h.writeXMLResponse(w, "PolicyStatus", v.PolicyStatus, statusCode, "", requestID)
	case *GetBucketLocationOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">%s</LocationConstraint>`, v.LocationConstraint)))
	case *GetObjectLockConfigurationOutput:
		if v.ObjectLockConfiguration != nil {
			h.writeXMLResponse(w, "ObjectLockConfiguration", v.ObjectLockConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetPublicAccessBlockOutput:
		if v.PublicAccessBlockConfiguration != nil {
			h.writeXMLResponse(w, "PublicAccessBlockConfiguration", v.PublicAccessBlockConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetObjectLegalHoldOutput:
		if v.LegalHold != nil {
			h.writeXMLResponse(w, "LegalHold", v.LegalHold, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetObjectRetentionOutput:
		if v.Retention != nil {
			h.writeXMLResponse(w, "Retention", v.Retention, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
		} else {
			w.WriteHeader(statusCode)
		}
	case *CreateBucketOutput:
		w.Header().Set("Location", v.Location)
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><CreateBucketResult><Location>` + xmlEscape(v.Location) + `</Location></CreateBucketResult>`))
	case *HeadBucketOutput:
		w.Header().Set("x-amz-bucket-region", v.BucketRegion)
		w.WriteHeader(statusCode)
	case *GetBucketTaggingOutput:
		h.writeXMLResponse(w, "Tagging", v, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *GetBucketVersioningOutput:
		h.writeXMLResponse(w, "VersioningConfiguration", v, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *GetBucketInventoryConfigurationOutput:
		h.writeXMLResponse(w, "InventoryConfiguration", v.InventoryConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *ListBucketInventoryConfigurationsOutput:
		h.writeXMLResponse(w, "ListInventoryConfigurationsResult", v, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *GetBucketMetricsConfigurationOutput:
		h.writeXMLResponse(w, "MetricsConfiguration", v.MetricsConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *ListBucketMetricsConfigurationsOutput:
		h.writeXMLResponse(w, "ListMetricsConfigurationsResult", v, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *GetBucketLifecycleConfigurationOutput:
		h.writeXMLResponse(w, "LifecycleConfiguration", v, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *GetBucketWebsiteOutput:
		h.writeXMLResponse(w, "WebsiteConfiguration", v, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *GetObjectTaggingOutput:
		h.writeXMLResponse(w, "Tagging", v, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *CreateMultipartUploadOutput:
		h.writeXMLResponse(w, "InitiateMultipartUploadResult", v, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case *CompleteMultipartUploadOutput:
		// The complete response carries its per-object facts as headers —
		// the Smithy model binds VersionId to x-amz-version-id and the SSE
		// fields mirror every other write response — never as body
		// elements.
		if v.VersionId != "" && v.VersionId != "null" {
			w.Header().Set("x-amz-version-id", v.VersionId)
		}
		setSSEHeaders(w.Header(), "", "", v.ServerSideEncryption, v.SSEKMSKeyId)
		h.writeXMLResponse(w, "CompleteMultipartUploadResult", v, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/", requestID)
	case nil:
		w.WriteHeader(statusCode)
	case interface{ ToXML() string }:
		writeToXMLResult(w, v, statusCode)
	default:
		h.writeXMLResponse(w, "", v, statusCode, "", requestID)
	}
}
