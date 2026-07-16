package resilience

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockAWSError simulates awserrors.AWSError for testing
type MockAWSError struct {
	Code       string
	Message    string
	HTTPStatus int
}

func (e *MockAWSError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *MockAWSError) GetHTTPStatus() int {
	return e.HTTPStatus
}

func TestAWSErrorPropagation_ThroughWrappers(t *testing.T) {
	awsErr := &MockAWSError{
		Code:       "EntityAlreadyExists",
		Message:    "User already exists",
		HTTPStatus: 409,
	}

	t.Run("single wrap with fmt.Errorf", func(t *testing.T) {
		wrappedErr := fmt.Errorf("wrapper error: %w", awsErr)

		var extracted *MockAWSError
		assert.True(t, errors.As(wrappedErr, &extracted), "errors.As should extract AWSError")
		assert.Equal(t, "EntityAlreadyExists", extracted.Code)
	})

	t.Run("double wrap with fmt.Errorf", func(t *testing.T) {
		wrappedErr := fmt.Errorf("wrapper error: %w", awsErr)
		doubleWrapped := fmt.Errorf("outer: %w", wrappedErr)

		var extracted *MockAWSError
		assert.True(t, errors.As(doubleWrapped, &extracted), "errors.As should extract AWSError from double wrap")
		assert.Equal(t, "EntityAlreadyExists", extracted.Code)
	})

	t.Run("through CircuitBreakerWrapper", func(t *testing.T) {
		cfg := DefaultConfig()
		cb := NewCircuitBreaker(cfg)
		wrapper := NewCircuitBreakerWrapper("test", cb)

		err := wrapper.Execute(context.Background(), func() error {
			return awsErr
		})

		var extracted *MockAWSError
		assert.True(t, errors.As(err, &extracted), "errors.As should extract AWSError from CircuitBreakerWrapper")
		assert.Equal(t, "EntityAlreadyExists", extracted.Code)
	})

	t.Run("through BulkheadWrapper", func(t *testing.T) {
		cfg := DefaultBulkheadConfig()
		bh := NewBulkhead(cfg)
		wrapper := NewBulkheadWrapper("test", bh)

		err := wrapper.Execute(context.Background(), func() error {
			return awsErr
		})

		var extracted *MockAWSError
		assert.True(t, errors.As(err, &extracted), "errors.As should extract AWSError from BulkheadWrapper")
		assert.Equal(t, "EntityAlreadyExists", extracted.Code)
	})

	t.Run("through AdaptiveTimeoutWrapper", func(t *testing.T) {
		cfg := DefaultAdaptiveTimeoutConfig()
		at := NewAdaptiveTimeout(cfg)
		wrapper := NewAdaptiveTimeoutWrapper("test", at)

		err := wrapper.Execute(context.Background(), func(ctx context.Context) error {
			return awsErr
		})

		var extracted *MockAWSError
		assert.True(t, errors.As(err, &extracted), "errors.As should extract AWSError from AdaptiveTimeoutWrapper")
		assert.Equal(t, "EntityAlreadyExists", extracted.Code)
	})

	t.Run("through ResilienceWrapper with all components", func(t *testing.T) {
		cbCfg := DefaultConfig()
		cb := NewCircuitBreaker(cbCfg)
		bhCfg := DefaultBulkheadConfig()
		bh := NewBulkhead(bhCfg)
		atCfg := DefaultAdaptiveTimeoutConfig()
		at := NewAdaptiveTimeout(atCfg)

		wrapper := NewResilienceWrapper("test", cb, bh, at, nil)

		_, err := wrapper.ExecuteWithResult(context.Background(), func() (interface{}, error) {
			return nil, awsErr
		})

		var extracted *MockAWSError
		assert.True(t, errors.As(err, &extracted), "errors.As should extract AWSError from full ResilienceWrapper")
		assert.Equal(t, "EntityAlreadyExists", extracted.Code)
	})

	t.Run("through ResilienceWrapper with Retry", func(t *testing.T) {
		cbCfg := DefaultConfig()
		cb := NewCircuitBreaker(cbCfg)
		bhCfg := DefaultBulkheadConfig()
		bh := NewBulkhead(bhCfg)
		atCfg := DefaultAdaptiveTimeoutConfig()
		at := NewAdaptiveTimeout(atCfg)
		rp := NewRetryPolicy()
		rp.SetMaxAttempts(1)

		wrapper := NewResilienceWrapper("test", cb, bh, at, rp)

		_, err := wrapper.ExecuteWithResult(context.Background(), func() (interface{}, error) {
			return nil, awsErr
		})

		var extracted *MockAWSError
		assert.True(t, errors.As(err, &extracted), "errors.As should extract AWSError from Retry wrapper")
		assert.Equal(t, "EntityAlreadyExists", extracted.Code)
	})
}

func TestAWSErrorPropagation_ThroughRetryWithMaxAttempts(t *testing.T) {
	awsErr := &MockAWSError{
		Code:       "InternalFailure",
		Message:    "Internal server error",
		HTTPStatus: 500,
	}

	rp := NewRetryPolicy()
	rp.SetMaxAttempts(2)
	rp.SetInitialDelay(1 * time.Millisecond)

	callCount := 0
	err := rp.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		return awsErr
	})

	assert.Equal(t, 2, callCount, "Server error (5xx) should be retried")
	assert.Error(t, err)

	var extracted *MockAWSError
	assert.True(t, errors.As(err, &extracted), "errors.As should extract AWSError after retry exhaustion")
	assert.Equal(t, "InternalFailure", extracted.Code)
	t.Logf("Error after retry: %v", err)
}

func TestIsRetryable_WithHTTPStatus(t *testing.T) {
	t.Run("client error 4xx is not retryable", func(t *testing.T) {
		clientErrors := []int{400, 401, 403, 404, 409, 422, 429}
		for _, status := range clientErrors {
			err := &MockAWSError{Code: "ClientError", HTTPStatus: status}
			assert.False(t, IsRetryable(err), "4xx error should not be retryable: %d", status)
		}
	})

	t.Run("server error 5xx is retryable", func(t *testing.T) {
		serverErrors := []int{500, 502, 503, 504}
		for _, status := range serverErrors {
			err := &MockAWSError{Code: "ServerError", HTTPStatus: status}
			assert.True(t, IsRetryable(err), "5xx error should be retryable: %d", status)
		}
	})

	t.Run("circuit breaker open is not retryable", func(t *testing.T) {
		err := NewCircuitBreakerOpen("test")
		assert.False(t, IsRetryable(err))
	})

	t.Run("rate limit is not retryable", func(t *testing.T) {
		err := NewRateLimit("rate limited")
		assert.False(t, IsRetryable(err))
	})

	t.Run("timeout is not retryable", func(t *testing.T) {
		err := NewTimeout("timed out")
		assert.False(t, IsRetryable(err))
	})

	t.Run("wrapped client error is not retryable", func(t *testing.T) {
		innerErr := &MockAWSError{Code: "NotFound", HTTPStatus: 404}
		wrappedErr := fmt.Errorf("wrapper: %w", innerErr)
		assert.False(t, IsRetryable(wrappedErr), "wrapped 4xx error should not be retryable")
	})
}

func TestIsClientError(t *testing.T) {
	t.Run("4xx is client error", func(t *testing.T) {
		err := &MockAWSError{Code: "BadRequest", HTTPStatus: 400}
		assert.True(t, IsClientError(err))
	})

	t.Run("5xx is not client error", func(t *testing.T) {
		err := &MockAWSError{Code: "InternalError", HTTPStatus: 500}
		assert.False(t, IsClientError(err))
	})

	t.Run("wrapped error is detected", func(t *testing.T) {
		innerErr := &MockAWSError{Code: "NotFound", HTTPStatus: 404}
		wrappedErr := fmt.Errorf("wrapper: %w", innerErr)
		assert.True(t, IsClientError(wrappedErr))
	})
}

func TestRetryPolicy_ClientErrorNoRetry(t *testing.T) {
	clientErr := &MockAWSError{
		Code:       "EntityAlreadyExists",
		Message:    "User already exists",
		HTTPStatus: 409,
	}

	rp := NewRetryPolicy()
	rp.SetMaxAttempts(3)
	rp.SetInitialDelay(1 * time.Millisecond)

	callCount := 0
	err := rp.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		return clientErr
	})

	assert.Equal(t, 1, callCount, "Client error should NOT be retried")
	assert.Error(t, err)

	var extracted *MockAWSError
	assert.True(t, errors.As(err, &extracted), "errors.As should extract AWSError")
	assert.Equal(t, "EntityAlreadyExists", extracted.Code)
}

func TestRetryPolicy_ServerErrorWithRetry(t *testing.T) {
	serverErr := &MockAWSError{
		Code:       "InternalFailure",
		Message:    "Internal server error",
		HTTPStatus: 500,
	}

	rp := NewRetryPolicy()
	rp.SetMaxAttempts(3)
	rp.SetInitialDelay(1 * time.Millisecond)

	callCount := 0
	err := rp.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		return serverErr
	})

	assert.Equal(t, 3, callCount, "Server error SHOULD be retried")
	assert.Error(t, err)

	var extracted *MockAWSError
	assert.True(t, errors.As(err, &extracted), "errors.As should extract AWSError after retry")
	assert.Equal(t, "InternalFailure", extracted.Code)
}
