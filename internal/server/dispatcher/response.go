package dispatcher

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/api"
)

func (d *Dispatcher) writeResponse(w http.ResponseWriter, r *http.Request, operation *api.Operation, opNameOverride string, resp interface{}) {
	if streamable, ok := resp.(response.StreamableResponse); ok {
		headers := streamable.GetStreamHeaders()
		for k, v := range headers {
			w.Header()[k] = v
		}
		statusCode := http.StatusOK
		if sc, ok := resp.(response.StatusCodeResponse); ok {
			if code := sc.GetStreamStatusCode(); code > 0 {
				statusCode = code
			}
		}
		w.WriteHeader(statusCode)
		if reader := streamable.GetStream(); reader != nil {
			if rc, ok := reader.(io.ReadCloser); ok {
				defer rc.Close()
			}
			if _, err := io.Copy(w, reader); err != nil {
				logs.Error("Failed to stream response", logs.Err(err))
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
	reqAmzTarget := r.Header.Get("X-Amz-Target")
	if r.Header.Get("X-Amzn-Query-Mode") == "true" {
		xAmznQueryMode = true
	}
	isJSONRequest := reqAmzTarget != "" || strings.Contains(reqContentType, "application/x-amz-json") || strings.Contains(reqContentType, "application/x-amzn-json")
	isCBORRequest := strings.Contains(reqContentType, "application/cbor") || r.Header.Get("smithy-protocol") == "rpc-v2-cbor"

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

	if contentType == "" && isJSONRequest && reqContentType != "" {
		contentType = reqContentType
	}

	protocol.SetProtocolHeaders(w, contentType, xAmzTarget, smithyProtocol, xAmznQueryMode)

	if isCBORRequest {
		if err := protocol.EncodeCBORResponse(w, resp); err != nil {
			logs.Error("Failed to encode CBOR response", logs.Err(err))
		}
		return
	}

	if isJSONRequest {
		if err := protocol.EncodeJSONResponse(w, resp); err != nil {
			logs.Error("Failed to encode JSON response", logs.Err(err))
		}
		return
	}

	if smithyProtocol == "aws.protocols#restXml" || smithyProtocol == "aws.protocols#ec2Query" || smithyProtocol == "aws.protocols#awsQuery" {
		var err error
		if smithyProtocol == "aws.protocols#awsQuery" || smithyProtocol == "aws.protocols#ec2Query" {
			err = protocol.EncodeQueryXMLResponse(w, opName, resp)
		} else if strings.HasPrefix(r.URL.Path, "/2020-05-31/") {
			if isCloudFrontPayloadOperation(opName) {
				payloadRoot := getCloudFrontPayloadRoot(opName)
				extractCloudFrontETag(w, resp, payloadRoot)
				err = protocol.EncodeRestXMLPayloadResponse(w, payloadRoot, resp)
			} else {
				err = protocol.EncodeRestXMLResponseWithNamespace(w, opName, "https://cloudfront.amazonaws.com/doc/2020-05-31/", resp)
			}
		} else if strings.HasPrefix(r.URL.Path, "/2013-04-01/") {
			err = protocol.EncodeRestXMLResponseWithNamespace(w, opName, "https://route53.amazonaws.com/doc/2013-04-01/", resp)
		} else {
			err = protocol.EncodeRestXMLResponse(w, opName, resp)
		}
		if err != nil {
			logs.Error("Failed to encode XML response", logs.Err(err))
		}
		return
	}

	if err := protocol.EncodeJSONResponse(w, resp); err != nil {
		logs.Error("Failed to encode JSON response", logs.Err(err))
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

	logs.Error("Unhandled error", logs.String("type", fmt.Sprintf("%T", err)), logs.Err(err))
	awserrors.WriteAWSError(w, awserrors.ErrInternal, contentType)
}

func (d *Dispatcher) handleErrorForRequest(w http.ResponseWriter, r *http.Request, err error) {
	contentType := d.getErrorContentTypeForRequest(r)

	if customErr, ok := err.(awserrors.CustomJSONMarshaler); ok {
		awserrors.WriteCustomJSONError(w, customErr, contentType)
		return
	}

	awsErr := d.extractAWSError(err)
	if awsErr != nil {
		awserrors.WriteAWSError(w, awsErr, contentType)
		return
	}

	logs.Error("Unhandled error", logs.String("type", fmt.Sprintf("%T", err)), logs.Err(err))
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
		return reqContentType
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
