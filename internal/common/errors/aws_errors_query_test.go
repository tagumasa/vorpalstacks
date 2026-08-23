package errors

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAWSErrorQueryCodeEmission(t *testing.T) {
	t.Run("ToXML uses the query code when set", func(t *testing.T) {
		err := NewAWSError("QueueDoesNotExist", "The specified queue does not exist.", 400).
			SetQueryErrorCode("AWS.SimpleQueueService.NonExistentQueue")
		xml := err.ToXML()
		assert.Contains(t, xml, "<Code>AWS.SimpleQueueService.NonExistentQueue</Code>")
		assert.Contains(t, xml, "<Type>Sender</Type>")
		assert.NotContains(t, xml, ";Sender")
	})

	t.Run("ToXML falls back to the shape code when unset", func(t *testing.T) {
		err := NewAWSError("InvalidAttributeName", "The attribute name is invalid.", 400)
		xml := err.ToXML()
		assert.Contains(t, xml, "<Code>InvalidAttributeName</Code>")
	})

	t.Run("ToXML maps server faults to Receiver", func(t *testing.T) {
		err := NewAWSError("InternalFailure", "internal", 500).
			SetQueryErrorCode("InternalFailure")
		assert.Contains(t, err.ToXML(), "<Type>Receiver</Type>")
	})

	t.Run("JSON responses carry the query code with the fault suffix", func(t *testing.T) {
		err := NewAWSError("QueueDoesNotExist", "The specified queue does not exist.", 400).
			SetQueryErrorCode("AWS.SimpleQueueService.NonExistentQueue")
		rec := httptest.NewRecorder()
		WriteAWSError(rec, err, "application/x-amz-json-1.0")
		assert.Equal(t, "AWS.SimpleQueueService.NonExistentQueue;Sender",
			rec.Header().Get("x-amzn-query-error"))
		assert.Equal(t, "QueueDoesNotExist", rec.Header().Get("X-Amzn-ErrorType"))
	})

	t.Run("A 404 client error keeps the Sender fault suffix", func(t *testing.T) {
		err := NewAWSError("ResourceNotFoundException", "missing", 404).
			SetQueryErrorCode("ResourceNotFoundException")
		rec := httptest.NewRecorder()
		WriteAWSError(rec, err, "application/x-amz-json-1.0")
		assert.Equal(t, "ResourceNotFoundException;Sender",
			rec.Header().Get("x-amzn-query-error"))
	})

	t.Run("The header is omitted when no query code is set", func(t *testing.T) {
		err := NewAWSError("InvalidAttributeName", "invalid", 400)
		rec := httptest.NewRecorder()
		WriteAWSError(rec, err, "application/x-amz-json-1.0")
		assert.Empty(t, rec.Header().Get("x-amzn-query-error"))
	})

	t.Run("XML responses carry the query code in the body, not the header", func(t *testing.T) {
		err := NewAWSError("QueueDoesNotExist", "The specified queue does not exist.", 400).
			SetQueryErrorCode("AWS.SimpleQueueService.NonExistentQueue")
		rec := httptest.NewRecorder()
		WriteAWSError(rec, err, "text/xml")
		assert.Empty(t, rec.Header().Get("x-amzn-query-error"))
		body := rec.Body.String()
		assert.True(t, strings.Contains(body, "<Code>AWS.SimpleQueueService.NonExistentQueue</Code>"),
			"query-protocol XML must carry the model's query error code: %s", body)
	})
}
