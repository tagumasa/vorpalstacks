package s3

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
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
