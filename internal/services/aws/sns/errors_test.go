package sns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSNSErrors(t *testing.T) {
	t.Run("SNSError Unwrap", func(t *testing.T) {
		err := NewInvalidParameterException("test")
		assert.NotNil(t, err.Unwrap())
	})

	t.Run("predefined errors", func(t *testing.T) {
		assert.Equal(t, "InvalidParameter: Invalid parameter", ErrInvalidParameter.Error())
		assert.Equal(t, 400, ErrInvalidParameter.GetHTTPStatusCode())

		assert.Equal(t, "NotFound: Resource not found", ErrNotFound.Error())
		assert.Equal(t, 404, ErrNotFound.GetHTTPStatusCode())

		assert.Equal(t, "AuthorizationError: Authorization error", ErrAuthorizationError.Error())
		assert.Equal(t, 403, ErrAuthorizationError.GetHTTPStatusCode())

		assert.Equal(t, "InternalError: Internal server error", ErrInternalError.Error())
		assert.Equal(t, 500, ErrInternalError.GetHTTPStatusCode())

		assert.Equal(t, "EndpointDisabled: Endpoint is disabled", ErrEndpointDisabled.Error())
		assert.Equal(t, 400, ErrEndpointDisabled.GetHTTPStatusCode())

		assert.Equal(t, "FilterLimitExceeded: Filter limit exceeded", ErrFilterLimitExceeded.Error())
		assert.Equal(t, 400, ErrFilterLimitExceeded.GetHTTPStatusCode())

		assert.Equal(t, "Throttled: Request was throttled", ErrThrottled.Error())
		assert.Equal(t, 429, ErrThrottled.GetHTTPStatusCode())

		assert.Equal(t, "ValidationException: Validation error", ErrValidation.Error())
		assert.Equal(t, 400, ErrValidation.GetHTTPStatusCode())

		assert.Equal(t, "NotFound: Topic does not exist", ErrTopicNotFound.Error())
		assert.Equal(t, 404, ErrTopicNotFound.GetHTTPStatusCode())

		assert.Equal(t, "NotFound: Subscription does not exist", ErrSubscriptionNotFound.Error())
		assert.Equal(t, 404, ErrSubscriptionNotFound.GetHTTPStatusCode())

		assert.Equal(t, "NotFound: Platform application does not exist", ErrPlatformAppNotFound.Error())
		assert.Equal(t, 404, ErrPlatformAppNotFound.GetHTTPStatusCode())

		assert.Equal(t, "NotFound: Endpoint does not exist", ErrEndpointNotFound.Error())
		assert.Equal(t, 404, ErrEndpointNotFound.GetHTTPStatusCode())

		assert.Equal(t, "TagLimitExceeded: Tag limit exceeded", ErrTagLimitExceeded.Error())
		assert.Equal(t, 400, ErrTagLimitExceeded.GetHTTPStatusCode())

		assert.Equal(t, "TagPolicy: Tag policy violation", ErrTagPolicy.Error())
		assert.Equal(t, 400, ErrTagPolicy.GetHTTPStatusCode())

		assert.Equal(t, "BatchEntryIdsNotDistinct: Two or more batch entries have the same ID", ErrBatchEntryIdsNotDistinct.Error())
		assert.Equal(t, 400, ErrBatchEntryIdsNotDistinct.GetHTTPStatusCode())

		assert.Equal(t, "TooManyEntriesInBatchRequest: Maximum number of entries per request are 10", ErrTooManyEntriesInBatch.Error())
		assert.Equal(t, 400, ErrTooManyEntriesInBatch.GetHTTPStatusCode())
	})

	t.Run("NewInvalidParameterException", func(t *testing.T) {
		err := NewInvalidParameterException("invalid parameter value")
		assert.Equal(t, "InvalidParameterException: invalid parameter value", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())
	})

	t.Run("NewNotFoundException", func(t *testing.T) {
		err := NewNotFoundException("topic")
		assert.Equal(t, "NotFoundException: topic not found", err.Error())
		assert.Equal(t, 404, err.GetHTTPStatusCode())
	})

	t.Run("NewAuthorizationErrorException", func(t *testing.T) {
		err := NewAuthorizationErrorException("not authorized")
		assert.Equal(t, "AccessDeniedException: not authorized", err.Error())
		assert.Equal(t, 403, err.GetHTTPStatusCode())
	})

	t.Run("NewInternalErrorException", func(t *testing.T) {
		err := NewInternalErrorException("service error")
		assert.Equal(t, "InternalError: service error", err.Error())
		assert.Equal(t, 500, err.GetHTTPStatusCode())
	})

	t.Run("NewEndpointDisabledException", func(t *testing.T) {
		err := NewEndpointDisabledException("endpoint disabled for push")
		assert.Equal(t, "EndpointDisabled: endpoint disabled for push", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())
	})

	t.Run("NewFilterLimitExceededException", func(t *testing.T) {
		err := NewFilterLimitExceededException("filter policy limit")
		assert.Equal(t, "FilterLimitExceeded: filter policy limit", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())
	})

	t.Run("NewThrottledException", func(t *testing.T) {
		err := NewThrottledException("too many requests")
		assert.Equal(t, "ThrottlingException: too many requests", err.Error())
		assert.Equal(t, 429, err.GetHTTPStatusCode())
	})

	t.Run("NewValidationException", func(t *testing.T) {
		err := NewValidationException("validation failed")
		assert.Equal(t, "ValidationException: validation failed", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())
	})
}
