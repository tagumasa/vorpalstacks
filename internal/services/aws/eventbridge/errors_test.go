package eventbridge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventsErrors(t *testing.T) {
	t.Run("EventsError Unwrap", func(t *testing.T) {
		err := NewValidationException("test")
		assert.NotNil(t, err.Unwrap())
	})

	t.Run("predefined errors", func(t *testing.T) {
		assert.Equal(t, "ValidationException: Validation error", ErrValidation.Error())
		assert.Equal(t, 400, ErrValidation.GetHTTPStatusCode())

		assert.Equal(t, "ResourceNotFoundException: Resource  not found", ErrResourceNotFound.Error())
		assert.Equal(t, 404, ErrResourceNotFound.GetHTTPStatusCode())

		assert.Equal(t, "ResourceAlreadyExistsException: Resource already exists", ErrResourceAlreadyExists.Error())
		assert.Equal(t, 409, ErrResourceAlreadyExists.GetHTTPStatusCode())

		assert.Equal(t, "InvalidParameterException: Invalid parameter", ErrInvalidParameter.Error())
		assert.Equal(t, 400, ErrInvalidParameter.GetHTTPStatusCode())
	})

	t.Run("NewValidationException", func(t *testing.T) {
		err := NewValidationException("invalid rule name")
		assert.Equal(t, "ValidationException: invalid rule name", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())
	})

	t.Run("NewResourceNotFoundException", func(t *testing.T) {
		err := NewResourceNotFoundException("rule not found")
		assert.Equal(t, "ResourceNotFoundException: rule not found", err.Error())
		assert.Equal(t, 404, err.GetHTTPStatusCode())
	})

	t.Run("NewResourceAlreadyExistsException", func(t *testing.T) {
		err := NewResourceAlreadyExistsException("rule")
		assert.Equal(t, "ResourceAlreadyExistsException: rule already exists", err.Error())
		assert.Equal(t, 409, err.GetHTTPStatusCode())
	})

	t.Run("NewInvalidParameterException", func(t *testing.T) {
		err := NewInvalidParameterException("invalid parameter value")
		assert.Equal(t, "InvalidParameterException: invalid parameter value", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())
	})
}

func TestBuildEventBusARN(t *testing.T) {
	t.Run("builds correct ARN", func(t *testing.T) {
		arn := BuildEventBusARN("123456789012", "us-east-1", "my-event-bus")
		assert.Equal(t, "arn:aws:events:us-east-1:123456789012:event-bus/my-event-bus", arn)
	})

	t.Run("handles special characters in name", func(t *testing.T) {
		arn := BuildEventBusARN("123456789012", "eu-west-1", "my_event_bus")
		assert.Equal(t, "arn:aws:events:eu-west-1:123456789012:event-bus/my_event_bus", arn)
	})
}
