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

		assert.Equal(t, "InvalidParameterValueException: The runtime parameter is invalid.", ErrInvalidRuntime.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidRuntime.GetHTTPStatusCode())

		assert.Equal(t, "CodeVerificationFailedException: The code signature failed the signature verification check.", ErrCodeVerificationFailed.Error())
		assert.Equal(t, http.StatusBadRequest, ErrCodeVerificationFailed.GetHTTPStatusCode())

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

func TestIsLambdaError(t *testing.T) {
	t.Run("returns true for LambdaError", func(t *testing.T) {
		err := NewLambdaError("TestCode", "Test message", 400)
		assert.True(t, IsLambdaError(err))
	})

	t.Run("returns false for non-LambdaError", func(t *testing.T) {
		err := &testError{msg: "test"}
		assert.False(t, IsLambdaError(err))
	})
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestGetLambdaError(t *testing.T) {
	t.Run("returns LambdaError when passed LambdaError", func(t *testing.T) {
		err := NewLambdaError("TestCode", "Test message", 400)
		result := GetLambdaError(err)
		assert.Equal(t, "TestCode: Test message", result.Error())
	})

	t.Run("returns ErrServiceException for non-LambdaError", func(t *testing.T) {
		err := &testError{msg: "test"}
		result := GetLambdaError(err)
		assert.Equal(t, ErrServiceException, result)
	})
}

func TestValidateRuntime(t *testing.T) {
	t.Run("valid Node.js runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("nodejs22.x"))
		assert.True(t, ValidateRuntime("nodejs20.x"))
	})

	t.Run("valid Python runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("python3.14"))
		assert.True(t, ValidateRuntime("python3.13"))
		assert.True(t, ValidateRuntime("python3.12"))
		assert.True(t, ValidateRuntime("python3.11"))
		assert.True(t, ValidateRuntime("python3.10"))
	})

	t.Run("valid Java runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("java25"))
		assert.True(t, ValidateRuntime("java21"))
		assert.True(t, ValidateRuntime("java17"))
		assert.True(t, ValidateRuntime("java11"))
		assert.True(t, ValidateRuntime("java8.al2"))
	})

	t.Run("valid .NET runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("dotnet10"))
		assert.True(t, ValidateRuntime("dotnet8"))
	})

	t.Run("valid Ruby runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("ruby3.4"))
		assert.True(t, ValidateRuntime("ruby3.3"))
		assert.True(t, ValidateRuntime("ruby3.2"))
	})

	t.Run("valid custom runtimes", func(t *testing.T) {
		assert.True(t, ValidateRuntime("provided.al2023"))
		assert.True(t, ValidateRuntime("provided.al2"))
	})

	t.Run("case insensitive", func(t *testing.T) {
		assert.True(t, ValidateRuntime("PYTHON3.12"))
		assert.True(t, ValidateRuntime("NodeJS20.x"))
	})

	t.Run("invalid runtime", func(t *testing.T) {
		assert.False(t, ValidateRuntime("python3.7"))
		assert.False(t, ValidateRuntime("nodejs14.x"))
		assert.False(t, ValidateRuntime("invalid"))
		assert.False(t, ValidateRuntime(""))
	})
}

func TestValidateHandler(t *testing.T) {
	t.Run("valid Python handler", func(t *testing.T) {
		err := ValidateHandler("python3.12", "myhandler.handle")
		assert.NoError(t, err)
	})

	t.Run("valid Node.js handler", func(t *testing.T) {
		err := ValidateHandler("nodejs20.x", "index.handler")
		assert.NoError(t, err)
	})

	t.Run("valid Java handler", func(t *testing.T) {
		err := ValidateHandler("java17", "com.example.MyHandler::handleRequest")
		assert.NoError(t, err)
	})

	t.Run("valid Java handler with package only", func(t *testing.T) {
		err := ValidateHandler("java17", "com.example.MyHandler.handleRequest")
		assert.NoError(t, err)
	})

	t.Run("empty handler returns error", func(t *testing.T) {
		err := ValidateHandler("python3.12", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Handler cannot be empty")
	})

	t.Run("Python handler without dot returns error", func(t *testing.T) {
		err := ValidateHandler("python3.12", "myhandler")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Python handler must be in the format module.function")
	})

	t.Run("Node.js handler without dot returns error", func(t *testing.T) {
		err := ValidateHandler("nodejs20.x", "myhandler")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Node.js handler must be in the format file.function")
	})

	t.Run("Java handler without proper format returns error", func(t *testing.T) {
		err := ValidateHandler("java17", "myhandler")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Java handler must be in the format package.Class::method")
	})
}

func TestAWSResponse(t *testing.T) {
	t.Run("creates response with status and body", func(t *testing.T) {
		resp := AWSResponse(200, map[string]string{"message": "success"})
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, map[string]string{"message": "success"}, resp.Body)
	})
}
