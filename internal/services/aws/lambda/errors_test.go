package lambda

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLambdaErrors(t *testing.T) {
	t.Run("NewLambdaError", func(t *testing.T) {
		err := NewLambdaError("TestCode", "Test message", 400)
		assert.Equal(t, "TestCode: Test message", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())
	})

	t.Run("LambdaError ToJSON", func(t *testing.T) {
		err := NewLambdaError("ResourceNotFoundException", "Not found", 404)
		json := err.ToJSON()
		assert.Contains(t, json, "ResourceNotFoundException")
		assert.Contains(t, json, "Not found")
	})

	t.Run("predefined errors", func(t *testing.T) {
		assert.Equal(t, "ResourceNotFoundException: The resource specified in the request does not exist.", ErrResourceNotFound.Error())
		assert.Equal(t, http.StatusNotFound, ErrResourceNotFound.GetHTTPStatusCode())

		assert.Equal(t, "ResourceInUseException: The resource is already in use.", ErrResourceInUse.Error())
		assert.Equal(t, http.StatusBadRequest, ErrResourceInUse.GetHTTPStatusCode())

		assert.Equal(t, "InvalidParameterValueException: The value for the parameter is invalid.", ErrInvalidParameterValue.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidParameterValue.GetHTTPStatusCode())

		assert.Equal(t, "CodeStorageExceededException: The total code size for the account exceeds the maximum allowed limit.", ErrCodeStorageExceeded.Error())
		assert.Equal(t, http.StatusBadRequest, ErrCodeStorageExceeded.GetHTTPStatusCode())

		assert.Equal(t, "TooManyRequestsException: Too many requests have been made. Please retry.", ErrTooManyRequests.Error())
		assert.Equal(t, http.StatusTooManyRequests, ErrTooManyRequests.GetHTTPStatusCode())

		assert.Equal(t, "ServiceException: An internal service error occurred.", ErrServiceException.Error())
		assert.Equal(t, http.StatusInternalServerError, ErrServiceException.GetHTTPStatusCode())
	})

	t.Run("NewResourceNotFound", func(t *testing.T) {
		err := NewResourceNotFound("Function", "my-function")
		assert.Equal(t, "ResourceNotFoundException: Function my-function not found", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewInvalidParameter", func(t *testing.T) {
		err := NewInvalidParameter("Handler", "cannot be empty")
		assert.Equal(t, "InvalidParameterValueException: Invalid parameter 'Handler': cannot be empty", err.Error())
		assert.Equal(t, http.StatusBadRequest, err.GetHTTPStatusCode())
	})

	t.Run("NewResourceConflict", func(t *testing.T) {
		err := NewResourceConflict("function is already being updated")
		assert.Equal(t, "ResourceConflictException: function is already being updated", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})
}
