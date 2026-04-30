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
func (h *S3Handler) writeXMLResponse(w http.ResponseWriter, rootElement string, data interface{}, statusCode int, xmlns string) {
	var xmlData []byte
	var err error

	if rootElement != "" {
		xmlData, err = xml.Marshal(data)
		if err != nil {
			h.writeError(w, err, "", "")
			return
		}
		xmlStr := string(xmlData)
		if idx := strings.Index(xmlStr, ">"); idx != -1 {
			xmlStr = xmlStr[idx+1:]
		}
		if idx := strings.LastIndex(xmlStr, "<"); idx != -1 {
			xmlStr = xmlStr[:idx]
		}
		if xmlns != "" {
			xmlData = []byte(fmt.Sprintf(`<%s xmlns="%s">%s</%s>`, rootElement, xmlns, xmlStr, rootElement))
		} else {
			xmlData = []byte(fmt.Sprintf(`<%s>%s</%s>`, rootElement, xmlStr, rootElement))
		}
	} else {
		xmlData, err = xml.Marshal(data)
	}

	if err != nil {
		h.writeError(w, err, "", "")
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	_, _ = w.Write(xmlData)
}

// #nosec G705
func (h *S3Handler) writeError(w http.ResponseWriter, err error, bucket, key string) {
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
	case errors.Is(err, storecommon.ErrConflict) || (err != nil && strings.Contains(err.Error(), "not empty")):
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

	resource := bucket
	if key != "" {
		resource = bucket + "/" + key
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(awsErr.HTTPStatus)
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	_, _ = w.Write([]byte(fmt.Sprintf(`<Error><Code>%s</Code><Message>%s</Message><Resource>%s</Resource><RequestId>%s</RequestId></Error>`,
		escapeXMLValue(awsErr.Code), escapeXMLValue(awsErr.Message), escapeXMLValue(resource), fmt.Sprintf("%016X", time.Now().UnixNano()))))
}

func escapeXMLValue(s string) string {
	var buf strings.Builder
	xml.Escape(&buf, []byte(s))
	return buf.String()
}

// writeResult serialises the operation result into the HTTP response.
// Streaming operations (GetObject, SelectObjectContent) write directly;
// all other results are rendered as XML or plain status codes.
func (h *S3Handler) writeResult(w http.ResponseWriter, result interface{}, statusCode int) {
	switch v := result.(type) {
	case *SelectObjectContentOutput:
		if v.Payload != nil {
			defer v.Payload.Close()
			if _, err := io.Copy(w, v.Payload); err != nil {
				logs.Error("S3: failed to stream SelectObjectContent payload", logs.Err(err))
			}
		}
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
		w.WriteHeader(statusCode)
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
	case *GetBucketCORSOutput:
		if v.CORSConfiguration != nil {
			h.writeXMLResponse(w, "CORSConfiguration", v.CORSConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetBucketOwnershipControlsOutput:
		if v.OwnershipControls != nil {
			h.writeXMLResponse(w, "OwnershipControls", v.OwnershipControls, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetBucketRequestPaymentOutput:
		if v.RequestPaymentConfiguration != nil {
			h.writeXMLResponse(w, "RequestPaymentConfiguration", v.RequestPaymentConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetBucketAccelerateConfigurationOutput:
		if v.AccelerateConfiguration != nil {
			h.writeXMLResponse(w, "AccelerateConfiguration", v.AccelerateConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetBucketPolicyOutput:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(v.Policy))
	case *GetBucketLocationOutput:
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">%s</LocationConstraint>`, v.LocationConstraint)))
	case *GetObjectLockConfigurationOutput:
		if v.ObjectLockConfiguration != nil {
			h.writeXMLResponse(w, "ObjectLockConfiguration", v.ObjectLockConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetPublicAccessBlockOutput:
		if v.PublicAccessBlockConfiguration != nil {
			h.writeXMLResponse(w, "PublicAccessBlockConfiguration", v.PublicAccessBlockConfiguration, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetObjectLegalHoldOutput:
		if v.LegalHold != nil {
			h.writeXMLResponse(w, "LegalHold", v.LegalHold, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
		} else {
			w.WriteHeader(statusCode)
		}
	case *GetObjectRetentionOutput:
		if v.Retention != nil {
			h.writeXMLResponse(w, "Retention", v.Retention, statusCode, "http://s3.amazonaws.com/doc/2006-03-01/")
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
