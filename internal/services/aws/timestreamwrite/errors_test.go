package timestreamwrite

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimestreamWriteErrors(t *testing.T) {
	t.Run("ErrResourceNotFound", func(t *testing.T) {
		assert.Equal(t, "ResourceNotFoundException: The requested resource could not be found.", ErrResourceNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrResourceNotFound.GetHTTPStatusCode())
	})

	t.Run("ErrResourceAlreadyExists", func(t *testing.T) {
		assert.Equal(t, "ResourceAlreadyExistsException: The resource already exists.", ErrResourceAlreadyExists.Error())
		assert.Equal(t, http.StatusConflict, ErrResourceAlreadyExists.GetHTTPStatusCode())
	})

	t.Run("ErrInvalidParameter", func(t *testing.T) {
		assert.Equal(t, "ValidationException: The input fails to satisfy the constraints specified by an AWS service.", ErrInvalidParameter.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidParameter.GetHTTPStatusCode())
	})

	t.Run("ErrValidationException", func(t *testing.T) {
		assert.Equal(t, "ValidationException: The input fails to satisfy the constraints specified by an AWS service.", ErrValidationException.Error())
		assert.Equal(t, http.StatusBadRequest, ErrValidationException.GetHTTPStatusCode())
	})

	t.Run("ErrAccessDenied", func(t *testing.T) {
		assert.Equal(t, "AccessDeniedException: You are not authorised to perform this action.", ErrAccessDenied.Error())
		assert.Equal(t, http.StatusForbidden, ErrAccessDenied.GetHTTPStatusCode())
	})

	t.Run("ErrThrottling", func(t *testing.T) {
		assert.Equal(t, "ThrottlingException: The request was denied due to request throttling.", ErrThrottling.Error())
		assert.Equal(t, http.StatusTooManyRequests, ErrThrottling.GetHTTPStatusCode())
	})

	t.Run("ErrInternalServer", func(t *testing.T) {
		assert.Equal(t, "InternalServerException: Internal server error.", ErrInternalServer.Error())
		assert.Equal(t, http.StatusInternalServerError, ErrInternalServer.GetHTTPStatusCode())
	})

	t.Run("ErrServiceQuotaExceeded", func(t *testing.T) {
		assert.Equal(t, "ServiceQuotaExceededException: The request exceeded the service quota.", ErrServiceQuotaExceeded.Error())
		assert.Equal(t, http.StatusTooManyRequests, ErrServiceQuotaExceeded.GetHTTPStatusCode())
	})
}
