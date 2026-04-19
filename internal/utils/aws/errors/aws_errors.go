package errors

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"vorpalstacks/internal/core/logs"
)

// AWSError represents an AWS service error.
type AWSError struct {
	Code           string
	Message        string
	HTTPStatus     int
	RequestID      string
	Fault          string
	QueryErrorCode string
}

// NewAWSError creates a new AWS error with the given code, message, and HTTP status.
func NewAWSError(code, message string, httpStatus int) *AWSError {
	return &AWSError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		RequestID:  uuid.New().String(),
		Fault:      "Server",
	}
}

// SetQueryErrorCode sets the x-amzn-query-error header value and returns the error for chaining.
func (e *AWSError) SetQueryErrorCode(code string) *AWSError {
	e.QueryErrorCode = code
	return e
}

// ToXML converts the error to XML format.
func (e *AWSError) ToXML() string {
	var faultType string
	if e.Fault == "Client" {
		faultType = "Sender"
	} else {
		faultType = "Receiver"
	}

	resp := XMLErrorResponse{}
	resp.Error.Type = faultType
	resp.Error.Code = e.Code
	resp.Error.Message = e.Message
	resp.RequestID = e.RequestID

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(resp); err != nil {
		logs.Error("Failed to encode XML error response", logs.String("code", e.Code), logs.Err(err))
	}
	return buf.String()
}

// ToJSON converts the error to JSON format.
func (e *AWSError) ToJSON() string {
	data := map[string]string{
		"__type":    e.Code,
		"message":   e.Message,
		"requestId": e.RequestID,
	}
	result, _ := json.Marshal(data)
	return string(result)
}

// ToCBOR converts the error to CBOR format.
func (e *AWSError) ToCBOR() []byte {
	data := map[string]interface{}{
		"__type":    e.Code,
		"message":   e.Message,
		"requestId": e.RequestID,
	}
	result, _ := cbor.Marshal(data)
	return result
}

// Error returns the error message.
func (e *AWSError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// GetCode returns the error code.
func (e *AWSError) GetCode() string {
	return e.Code
}

// GetMessage returns the error message.
func (e *AWSError) GetMessage() string {
	return e.Message
}

// GetHTTPStatusCode returns the HTTP status code.
func (e *AWSError) GetHTTPStatusCode() int {
	return e.HTTPStatus
}

// ToJSONWithFormat converts the error to JSON format with the specified format.
func (e *AWSError) ToJSONWithFormat(format string) string {
	switch format {
	case "lambda":
		faultType := "User"
		if e.Fault == "Server" {
			faultType = "Server"
		}
		data := map[string]string{
			"Type":    faultType,
			"Code":    e.Code,
			"Message": e.Message,
		}
		result, _ := json.Marshal(data)
		return string(result)
	case "rest-json":
		data := map[string]string{
			"type":    e.Code,
			"message": e.Message,
		}
		result, _ := json.Marshal(data)
		return string(result)
	default:
		return e.ToJSON()
	}
}

// XMLErrorResponse represents an XML error response.
type XMLErrorResponse struct {
	XMLName xml.Name `xml:"ErrorResponse"`
	Error   struct {
		Type    string `xml:"Type"`
		Code    string `xml:"Code"`
		Message string `xml:"Message"`
	} `xml:"Error"`
	RequestID string `xml:"RequestId"`
}

var (
	// ErrAccessDenied is returned when access is denied.
	ErrAccessDenied = NewAWSError("AccessDenied", "Access denied", http.StatusForbidden)
	// ErrInvalidAction is returned when an invalid action is requested.
	ErrInvalidAction = NewAWSError("InvalidAction", "Invalid action", http.StatusBadRequest)
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = NewAWSError("InvalidParameter", "Invalid parameter", http.StatusBadRequest)
	// ErrMissingParameter is returned when a required parameter is missing.
	ErrMissingParameter = NewAWSError("MissingParameter", "Missing required parameter", http.StatusBadRequest)
	// ErrResourceNotFound is returned when a resource is not found.
	ErrResourceNotFound = NewAWSError("ResourceNotFoundException", "Resource not found", http.StatusNotFound)
	// ErrServiceUnavailable is returned when the service is unavailable.
	ErrServiceUnavailable = NewAWSError("ServiceUnavailable", "Service unavailable", http.StatusServiceUnavailable)
	// ErrThrottling is returned when the request is throttled.
	ErrThrottling = NewAWSError("ThrottlingException", "Request was throttled", http.StatusTooManyRequests)
	// ErrValidation is returned when validation fails.
	ErrValidation = NewAWSError("ValidationException", "Validation error", http.StatusBadRequest)
	// ErrInternal is returned when an internal server error occurs.
	ErrInternal = NewAWSError("InternalFailure", "Internal server error", http.StatusInternalServerError)
	// ErrNotImplemented is returned when an operation is not implemented.
	ErrNotImplemented = NewAWSError("NotImplementedException", "Operation not implemented", http.StatusNotImplemented)
)

// CustomJSONMarshaler is an interface for types that can marshal themselves to JSON.
type CustomJSONMarshaler interface {
	ToJSON() string
	GetHTTPStatusCode() int
	GetCode() string
}

// WriteAWSError writes an AWS error to the response writer.
func WriteAWSError(w http.ResponseWriter, err *AWSError, contentType string) {
	if contentType == "application/xml" || contentType == "text/xml" || contentType == "application/x-www-form-urlencoded" {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(err.HTTPStatus)
		if _, writeErr := w.Write([]byte(err.ToXML())); writeErr != nil {
			logs.Error("Failed to write XML error response", logs.Err(writeErr))
		}
	} else if contentType == "application/cbor" {
		w.Header().Set("Content-Type", "application/cbor")
		w.Header().Set("smithy-protocol", "rpc-v2-cbor")
		w.WriteHeader(err.HTTPStatus)
		if _, writeErr := w.Write(err.ToCBOR()); writeErr != nil {
			logs.Error("Failed to write CBOR error response", logs.Err(writeErr))
		}
	} else {
		ct := contentType
		if ct == "" || ct == "text/plain" {
			ct = "application/json"
		}
		w.Header().Set("Content-Type", ct)
		w.Header().Set("X-Amzn-ErrorType", err.Code)
		if err.QueryErrorCode != "" {
			w.Header().Set("x-amzn-query-error", err.QueryErrorCode)
		}
		w.WriteHeader(err.HTTPStatus)
		if _, writeErr := w.Write([]byte(err.ToJSON())); writeErr != nil {
			logs.Error("Failed to write JSON error response", logs.Err(writeErr))
		}
	}
}

// WriteCustomJSONError writes a custom JSON error to the response writer.
func WriteCustomJSONError(w http.ResponseWriter, err CustomJSONMarshaler, contentType string) {
	if contentType == "application/xml" || contentType == "text/xml" || contentType == "application/x-www-form-urlencoded" {
		if awsErr, ok := err.(*AWSError); ok {
			WriteAWSError(w, awsErr, contentType)
			return
		}
	}
	if contentType == "application/cbor" {
		if awsErr, ok := err.(*AWSError); ok {
			WriteAWSError(w, awsErr, contentType)
			return
		}
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Amzn-ErrorType", err.GetCode())
	w.WriteHeader(err.GetHTTPStatusCode())
	if _, writeErr := w.Write([]byte(err.ToJSON())); writeErr != nil {
		logs.Error("Failed to write JSON error response", logs.Err(writeErr))
	}
}

// WriteAWSErrorWithFormat writes an AWS error with a specific JSON format.
func WriteAWSErrorWithFormat(w http.ResponseWriter, err *AWSError, contentType, jsonFormat string) {
	if contentType == "application/xml" || contentType == "text/xml" || contentType == "application/x-www-form-urlencoded" {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(err.HTTPStatus)
		if _, writeErr := w.Write([]byte(err.ToXML())); writeErr != nil {
			logs.Error("Failed to write XML error response", logs.Err(writeErr))
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.HTTPStatus)
		if _, writeErr := w.Write([]byte(err.ToJSONWithFormat(jsonFormat))); writeErr != nil {
			logs.Error("Failed to write JSON error response", logs.Err(writeErr))
		}
	}
}
