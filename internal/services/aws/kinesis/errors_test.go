package kinesis

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKinesisErrors(t *testing.T) {
	t.Run("ErrResourceNotFound", func(t *testing.T) {
		assert.Equal(t, "ResourceNotFoundException: Requested resource not found.", ErrResourceNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrResourceNotFound.GetHTTPStatusCode())
	})

	t.Run("ErrResourceInUse", func(t *testing.T) {
		assert.Equal(t, "ResourceInUseException: The resource is in use.", ErrResourceInUse.Error())
		assert.Equal(t, http.StatusBadRequest, ErrResourceInUse.GetHTTPStatusCode())
	})

	t.Run("ErrInvalidArgument", func(t *testing.T) {
		assert.Equal(t, "InvalidArgumentException: Invalid argument.", ErrInvalidArgument.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidArgument.GetHTTPStatusCode())
	})

	t.Run("ErrLimitExceeded", func(t *testing.T) {
		assert.Equal(t, "LimitExceededException: Rate exceeded for this stream.", ErrLimitExceeded.Error())
		assert.Equal(t, http.StatusTooManyRequests, ErrLimitExceeded.GetHTTPStatusCode())
	})

	t.Run("ErrProvisionedThroughputExceeded", func(t *testing.T) {
		assert.Equal(t, "ProvisionedThroughputExceededException: Rate exceeded for this shard.", ErrProvisionedThroughputExceeded.Error())
		assert.Equal(t, http.StatusTooManyRequests, ErrProvisionedThroughputExceeded.GetHTTPStatusCode())
	})

	t.Run("ErrExpiredIterator", func(t *testing.T) {
		assert.Equal(t, "ExpiredIteratorException: Iterator expired.", ErrExpiredIterator.Error())
		assert.Equal(t, http.StatusBadRequest, ErrExpiredIterator.GetHTTPStatusCode())
	})

	t.Run("ErrInvalidIterator", func(t *testing.T) {
		assert.Equal(t, "InvalidArgumentException: Invalid iterator.", ErrInvalidIterator.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidIterator.GetHTTPStatusCode())
	})

	t.Run("ErrShardClosed", func(t *testing.T) {
		assert.Equal(t, "InvalidArgumentException: Shard is closed.", ErrShardClosed.Error())
		assert.Equal(t, http.StatusBadRequest, ErrShardClosed.GetHTTPStatusCode())
	})

	t.Run("ErrAccessDenied", func(t *testing.T) {
		assert.Equal(t, "AccessDeniedException: Access denied.", ErrAccessDenied.Error())
		assert.Equal(t, http.StatusForbidden, ErrAccessDenied.GetHTTPStatusCode())
	})

	t.Run("ErrResourceAlreadyExists", func(t *testing.T) {
		assert.Equal(t, "ResourceAlreadyExistsException: Resource already exists.", ErrResourceAlreadyExists.Error())
		assert.Equal(t, http.StatusConflict, ErrResourceAlreadyExists.GetHTTPStatusCode())
	})

	t.Run("ErrKMSAccessDenied", func(t *testing.T) {
		assert.Equal(t, "KMSAccessDeniedException: KMS access denied.", ErrKMSAccessDenied.Error())
		assert.Equal(t, http.StatusForbidden, ErrKMSAccessDenied.GetHTTPStatusCode())
	})

	t.Run("ErrKMSDisabled", func(t *testing.T) {
		assert.Equal(t, "KMSDisabledException: KMS key is disabled.", ErrKMSDisabled.Error())
		assert.Equal(t, http.StatusBadRequest, ErrKMSDisabled.GetHTTPStatusCode())
	})

	t.Run("ErrKMSNotFound", func(t *testing.T) {
		assert.Equal(t, "KMSNotFoundException: KMS key not found.", ErrKMSNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrKMSNotFound.GetHTTPStatusCode())
	})

	t.Run("ErrKMSThrottling", func(t *testing.T) {
		assert.Equal(t, "KMSThrottlingException: KMS request throttled.", ErrKMSThrottling.Error())
		assert.Equal(t, http.StatusTooManyRequests, ErrKMSThrottling.GetHTTPStatusCode())
	})

	t.Run("ErrValidation", func(t *testing.T) {
		assert.Equal(t, "ValidationException: Validation error.", ErrValidation.Error())
		assert.Equal(t, http.StatusBadRequest, ErrValidation.GetHTTPStatusCode())
	})
}
