package eventbridge

import (
	"testing"

	"github.com/stretchr/testify/assert"

	awserrors "vorpalstacks/internal/common/errors"
)

func TestEventsErrors(t *testing.T) {
	t.Run("EventsError has no Unwrap", func(t *testing.T) {
		err := awserrors.NewValidationException("test")
		assert.Equal(t, "ValidationException: test", err.Error())
	})

	t.Run("predefined error constructors", func(t *testing.T) {
		err := awserrors.NewValidationException("Validation error")
		assert.Equal(t, "ValidationException: Validation error", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())

		notFound := awserrors.NewResourceNotFoundException("Resource", "")
		assert.Equal(t, "ResourceNotFoundException: Resource  not found", notFound.Error())
		assert.Equal(t, 404, notFound.GetHTTPStatusCode())
	})

	t.Run("NewValidationException", func(t *testing.T) {
		err := awserrors.NewValidationException("invalid rule name")
		assert.Equal(t, "ValidationException: invalid rule name", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())
	})

	t.Run("NewResourceNotFoundException", func(t *testing.T) {
		err := NewResourceNotFoundException("rule not found")
		assert.Equal(t, "ResourceNotFoundException: rule not found", err.Error())
		assert.Equal(t, 404, err.GetHTTPStatusCode())
	})

	t.Run("NewResourceAlreadyExistsException", func(t *testing.T) {
		err := awserrors.NewResourceAlreadyExistsException("rule")
		assert.Equal(t, "ResourceAlreadyExistsException: rule already exists", err.Error())
		assert.Equal(t, 409, err.GetHTTPStatusCode())
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
