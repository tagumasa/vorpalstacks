package errors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	awserrors "vorpalstacks/internal/utils/aws/errors"
)

func TestAWSError(t *testing.T) {
	t.Run("Error message", func(t *testing.T) {
		err := &AWSError{
			Code:    "TestCode",
			Message: "Test message",
		}
		assert.Equal(t, "TestCode: Test message", err.Error())
	})

	t.Run("NewAWSError with defaults", func(t *testing.T) {
		err := NewAWSError("Code", "Message", http.StatusBadRequest)
		assert.Equal(t, "Code", err.Code)
		assert.Equal(t, "Message", err.Message)
		assert.Equal(t, http.StatusBadRequest, err.HTTPStatus)
		assert.NotEmpty(t, err.RequestID)
		assert.Len(t, err.RequestID, 36)
		assert.Equal(t, "Server", err.Fault)
	})

	t.Run("ToXML with Client fault", func(t *testing.T) {
		err := &AWSError{
			Code:    "AccessDenied",
			Message: "Access denied",
			Fault:   "Client",
		}
		xml := err.ToXML()
		assert.Contains(t, xml, "<Type>Sender</Type>")
		assert.Contains(t, xml, "<Code>AccessDenied</Code>")
		assert.Contains(t, xml, "<Message>Access denied</Message>")
	})

	t.Run("ToXML with Server fault", func(t *testing.T) {
		err := &AWSError{
			Code:    "InternalError",
			Message: "Internal error",
			Fault:   "Server",
		}
		xml := err.ToXML()
		assert.Contains(t, xml, "<Type>Receiver</Type>")
		assert.Contains(t, xml, "<Code>InternalError</Code>")
	})

	t.Run("ToJSON", func(t *testing.T) {
		err := &AWSError{
			Code:      "NotFound",
			Message:   "Not found",
			RequestID: "req-123",
		}
		json := err.ToJSON()
		assert.Contains(t, json, "NotFound")
		assert.Contains(t, json, "Not found")
		assert.Contains(t, json, "req-123")
	})

	t.Run("Predefined errors", func(t *testing.T) {
		assert.Equal(t, "AccessDenied: Access denied", ErrAccessDenied.Error())
		assert.Equal(t, http.StatusForbidden, ErrAccessDenied.HTTPStatus)

		assert.Equal(t, "ResourceNotFoundException: Resource not found", ErrResourceNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrResourceNotFound.HTTPStatus)

		assert.Equal(t, "ThrottlingException: Request was throttled", ErrThrottling.Error())
		assert.Equal(t, http.StatusTooManyRequests, ErrThrottling.HTTPStatus)
	})
}

func TestWriteAWSError(t *testing.T) {
	t.Run("XML response", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := NewAWSError("TestCode", "Test message", http.StatusBadRequest)

		WriteAWSError(w, err, "application/xml")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/xml", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Body.String(), "<Error>")
		assert.Contains(t, w.Body.String(), "<Code>TestCode</Code>")
	})

	t.Run("JSON response", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := NewAWSError("TestCode", "Test message", http.StatusInternalServerError)

		WriteAWSError(w, err, "application/json")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Body.String(), "TestCode")
	})

	t.Run("URLEncoded defaults to XML", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := NewAWSError("TestCode", "Test message", http.StatusOK)

		WriteAWSError(w, err, "application/x-www-form-urlencoded")

		assert.Equal(t, "application/xml", w.Header().Get("Content-Type"))
	})

	t.Run("Unknown content type defaults to JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := NewAWSError("TestCode", "Test message", http.StatusOK)

		WriteAWSError(w, err, "text/plain")

		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})
}

func TestXMLErrorResponse(t *testing.T) {
	t.Run("Structure exists", func(t *testing.T) {
		resp := awserrors.XMLErrorResponse{}
		resp.Error.Type = "Sender"
		resp.Error.Code = "AccessDenied"
		resp.Error.Message = "Access denied"
		resp.RequestID = "req-123"

		assert.Equal(t, "Sender", resp.Error.Type)
		assert.Equal(t, "AccessDenied", resp.Error.Code)
		assert.Equal(t, "Access denied", resp.Error.Message)
		assert.Equal(t, "req-123", resp.RequestID)
	})
}

func TestGetHTTPStatus(t *testing.T) {
	t.Run("returns correct HTTP status", func(t *testing.T) {
		err := &AWSError{
			HTTPStatus: http.StatusNotFound,
		}
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("returns 0 for unset HTTPStatus", func(t *testing.T) {
		err := &AWSError{}
		assert.Equal(t, 0, err.GetHTTPStatusCode())
	})
}

func TestToJSONWithFormat(t *testing.T) {
	t.Run("lambda format with Client fault", func(t *testing.T) {
		err := &AWSError{
			Code:    "AccessDenied",
			Message: "Access denied",
			Fault:   "Client",
		}
		json := err.ToJSONWithFormat("lambda")
		assert.Contains(t, json, `"Type":"User"`)
		assert.Contains(t, json, `"Code":"AccessDenied"`)
		assert.Contains(t, json, `"Message":"Access denied"`)
	})

	t.Run("lambda format with Server fault", func(t *testing.T) {
		err := &AWSError{
			Code:    "InternalError",
			Message: "Internal error",
			Fault:   "Server",
		}
		json := err.ToJSONWithFormat("lambda")
		assert.Contains(t, json, `"Type":"Server"`)
		assert.Contains(t, json, `"Code":"InternalError"`)
	})

	t.Run("rest-json format", func(t *testing.T) {
		err := &AWSError{
			Code:    "ValidationException",
			Message: "Validation failed",
			Fault:   "Client",
		}
		json := err.ToJSONWithFormat("rest-json")
		assert.Contains(t, json, `"type":"ValidationException"`)
		assert.Contains(t, json, `"message":"Validation failed"`)
	})

	t.Run("default format", func(t *testing.T) {
		err := &AWSError{
			Code:      "NotFound",
			Message:   "Not found",
			RequestID: "req-456",
		}
		json := err.ToJSONWithFormat("default")
		assert.Contains(t, json, "NotFound")
		assert.Contains(t, json, "Not found")
		assert.Contains(t, json, "req-456")
	})

	t.Run("unknown format defaults to JSON", func(t *testing.T) {
		err := &AWSError{
			Code:    "Unknown",
			Message: "Unknown error",
		}
		json := err.ToJSONWithFormat("unknown")
		assert.Contains(t, json, "Unknown")
		assert.Contains(t, json, "Unknown error")
	})
}

func TestWriteAWSErrorWithFormat(t *testing.T) {
	t.Run("XML with format", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := NewAWSError("TestCode", "Test message", http.StatusBadRequest)

		WriteAWSErrorWithFormat(w, err, "application/xml", "lambda")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/xml", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Body.String(), "<Error>")
	})

	t.Run("JSON with lambda format", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := NewAWSError("ValidationException", "Invalid input", http.StatusBadRequest)
		err.Fault = "Client"

		WriteAWSErrorWithFormat(w, err, "application/json", "lambda")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Body.String(), `"Type":"User"`)
		assert.Contains(t, w.Body.String(), `"Code":"ValidationException"`)
	})

	t.Run("JSON with rest-json format", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := NewAWSError("Throttled", "Rate exceeded", http.StatusTooManyRequests)

		WriteAWSErrorWithFormat(w, err, "application/json", "rest-json")

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), `"type":"Throttled"`)
		assert.Contains(t, w.Body.String(), `"message":"Rate exceeded"`)
	})
}
