package dispatcher

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/store/api"
)

// StreamableResponse is an interface for responses that can be streamed.
type StreamableResponse interface {
	GetStream() io.Reader
	GetStreamHeaders() http.Header
}

func (d *Dispatcher) writeResponse(w http.ResponseWriter, r *http.Request, operation *api.Operation, opNameOverride string, response interface{}) {
	if streamable, ok := response.(StreamableResponse); ok {
		headers := streamable.GetStreamHeaders()
		for k, v := range headers {
			w.Header()[k] = v
		}
		w.WriteHeader(http.StatusOK)
		if reader := streamable.GetStream(); reader != nil {
			if rc, ok := reader.(io.ReadCloser); ok {
				defer rc.Close()
			}
			if _, err := io.Copy(w, reader); err != nil {
				log.Printf("Failed to stream response: %v", err)
			}
		}
		return
	}

	contentType := ""
	xAmzTarget := ""
	var xAmznQueryMode bool
	smithyProtocol := ""
	opName := opNameOverride

	if operation != nil {
		if operation.ContentType != nil {
			contentType = *operation.ContentType
		}
		if operation.XAmzTarget != nil {
			xAmzTarget = *operation.XAmzTarget
		}
		if operation.SmithyProtocol != nil {
			smithyProtocol = *operation.SmithyProtocol
		}
		xAmznQueryMode = operation.XAmznQueryMode
		if operation.Name != "" {
			opName = operation.Name
		}
	}

	reqContentType := r.Header.Get("Content-Type")
	acceptHeader := r.Header.Get("Accept")
	reqAmzTarget := r.Header.Get("X-Amz-Target")
	isJSONRequest := reqAmzTarget != "" || strings.Contains(reqContentType, "application/x-amz-json") || strings.Contains(reqContentType, "application/x-amzn-json")

	if smithyProtocol == "" && strings.Contains(reqContentType, "application/x-www-form-urlencoded") {
		hasAction := false
		if _, ok := r.URL.Query()["Action"]; ok {
			hasAction = true
		} else if r.Method == "POST" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				if err := r.ParseForm(); err == nil {
					_, hasAction = r.Form["Action"]
				}
				r.Form = nil
				r.PostForm = nil
				r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}
		if hasAction {
			smithyProtocol = "aws.protocols#awsQuery"
		}
	}

	if smithyProtocol == "" {
		if strings.HasPrefix(r.URL.Path, "/2013-04-01/") {
			smithyProtocol = "aws.protocols#restXml"
		} else if strings.HasPrefix(r.URL.Path, "/2020-05-31/") {
			smithyProtocol = "aws.protocols#restXml"
		}
	}

	protocol.SetProtocolHeaders(w, contentType, xAmzTarget, smithyProtocol, xAmznQueryMode)

	if isJSONRequest {
		if err := protocol.EncodeJSONResponse(w, response); err != nil {
			log.Printf("Failed to encode JSON response: %v", err)
		}
		return
	}

	if smithyProtocol == "aws.protocols#restXml" || smithyProtocol == "aws.protocols#ec2Query" || smithyProtocol == "aws.protocols#awsQuery" {
		var err error
		if smithyProtocol == "aws.protocols#awsQuery" || smithyProtocol == "aws.protocols#ec2Query" {
			err = protocol.EncodeQueryXMLResponse(w, opName, response)
		} else if strings.HasPrefix(r.URL.Path, "/2020-05-31/") {
			if isCloudFrontPayloadOperation(opName) {
				payloadRoot := getCloudFrontPayloadRoot(opName)
				extractCloudFrontETag(w, response, payloadRoot)
				err = protocol.EncodeRestXMLPayloadResponse(w, payloadRoot, response)
			} else {
				err = protocol.EncodeRestXMLResponse(w, opName, response)
			}
		} else {
			err = protocol.EncodeRestXMLResponse(w, opName, response)
		}
		if err != nil {
			log.Printf("Failed to encode XML response: %v", err)
		}
		return
	}

	if strings.Contains(acceptHeader, "application/cbor") || strings.Contains(reqContentType, "application/cbor") {
		if err := protocol.EncodeCBORResponse(w, response); err != nil {
			log.Printf("Failed to encode CBOR response: %v", err)
		}
		return
	}

	if err := protocol.EncodeJSONResponse(w, response); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func (d *Dispatcher) writeResponseWithOpName(w http.ResponseWriter, r *http.Request, serviceName, opName string, response interface{}) {
	op, _ := d.operationStore.Get(serviceName, opName)
	d.writeResponse(w, r, op, opName, response)
}

func (d *Dispatcher) handleError(w http.ResponseWriter, r *http.Request, operation *api.Operation, err error) {
	contentType := d.getErrorContentType(r, operation)

	if customErr, ok := err.(awserrors.CustomJSONMarshaler); ok {
		awserrors.WriteCustomJSONError(w, customErr, contentType)
		return
	}

	awsErr := d.extractAWSError(err)
	if awsErr != nil {
		awserrors.WriteAWSError(w, awsErr, contentType)
		return
	}

	awserrors.WriteAWSError(w, awserrors.ErrInternal, contentType)
}

func (d *Dispatcher) handleErrorForRequest(w http.ResponseWriter, r *http.Request, err error) {
	contentType := d.getErrorContentTypeForRequest(r)

	awsErr := d.extractAWSError(err)
	if awsErr != nil {
		awserrors.WriteAWSError(w, awsErr, contentType)
		return
	}

	if customErr, ok := err.(awserrors.CustomJSONMarshaler); ok {
		awserrors.WriteCustomJSONError(w, customErr, contentType)
		return
	}

	log.Printf("ERROR: unhandled error type=%T: %v", err, err)
	awserrors.WriteAWSError(w, awserrors.ErrInternal, contentType)
}

func (d *Dispatcher) extractAWSError(err error) *awserrors.AWSError {
	if err == nil {
		return nil
	}

	var awsErr *awserrors.AWSError
	if errors.As(err, &awsErr) {
		return awsErr
	}

	for err != nil {
		type awsErrorHolder interface {
			GetAWSError() *awserrors.AWSError
		}
		var holder awsErrorHolder
		if errors.As(err, &holder) {
			if awsErr := holder.GetAWSError(); awsErr != nil {
				return awsErr
			}
		}

		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			break
		}
		err = unwrapped
		if errors.As(err, &awsErr) {
			return awsErr
		}
	}

	return nil
}

func (d *Dispatcher) getErrorContentType(r *http.Request, operation *api.Operation) string {
	if operation != nil && operation.ContentType != nil {
		return *operation.ContentType
	}

	if operation != nil && operation.SmithyProtocol != nil {
		if *operation.SmithyProtocol == "aws.protocols#awsQuery" || *operation.SmithyProtocol == "aws.protocols#ec2Query" {
			return "text/xml"
		}
	}

	return d.getErrorContentTypeForRequest(r)
}

func (d *Dispatcher) getErrorContentTypeForRequest(r *http.Request) string {
	reqContentType := r.Header.Get("Content-Type")
	reqAmzTarget := r.Header.Get("X-Amz-Target")
	if reqAmzTarget != "" || strings.Contains(reqContentType, "application/x-amz-json") || strings.Contains(reqContentType, "application/x-amzn-json") {
		return "application/x-amz-json-1.0"
	}
	if strings.Contains(reqContentType, "application/x-www-form-urlencoded") {
		return "text/xml"
	}
	if strings.Contains(reqContentType, "application/cbor") {
		return "application/cbor"
	}
	if strings.HasPrefix(r.URL.Path, "/2020-05-31/") || strings.HasPrefix(r.URL.Path, "/2013-04-01/") {
		return "application/xml"
	}
	return "application/json"
}
