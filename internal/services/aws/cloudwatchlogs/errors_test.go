package cloudwatchlogs

import (
	"net/http"
	"testing"
)

func TestLogsError(t *testing.T) {
	tests := []struct {
		name         string
		err          *LogsError
		wantCode     string
		wantMsg      string
		wantHTTPCode int
	}{
		{
			name:         "ErrLogGroupNotFound",
			err:          ErrLogGroupNotFound,
			wantCode:     "ResourceNotFoundException",
			wantMsg:      "Log group not found",
			wantHTTPCode: http.StatusNotFound,
		},
		{
			name:         "ErrLogGroupAlreadyExists",
			err:          ErrLogGroupAlreadyExists,
			wantCode:     "ResourceAlreadyExistsException",
			wantMsg:      "Log group already exists",
			wantHTTPCode: http.StatusConflict,
		},
		{
			name:         "ErrLogStreamNotFound",
			err:          ErrLogStreamNotFound,
			wantCode:     "ResourceNotFoundException",
			wantMsg:      "Log stream not found",
			wantHTTPCode: http.StatusNotFound,
		},
		{
			name:         "ErrLogStreamAlreadyExists",
			err:          ErrLogStreamAlreadyExists,
			wantCode:     "ResourceAlreadyExistsException",
			wantMsg:      "Log stream already exists",
			wantHTTPCode: http.StatusConflict,
		},
		{
			name:         "ErrInvalidParameter",
			err:          ErrInvalidParameter,
			wantCode:     "InvalidParameterException",
			wantMsg:      "Invalid parameter",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrMissingParameter",
			err:          ErrMissingParameter,
			wantCode:     "MissingParameterException",
			wantMsg:      "Missing required parameter",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrAccessDenied",
			err:          ErrAccessDenied,
			wantCode:     "AccessDeniedException",
			wantMsg:      "Access denied",
			wantHTTPCode: http.StatusForbidden,
		},
		{
			name:         "ErrLimitExceeded",
			err:          ErrLimitExceeded,
			wantCode:     "LimitExceededException",
			wantMsg:      "Limit exceeded",
			wantHTTPCode: http.StatusTooManyRequests,
		},
		{
			name:         "ErrOperationAborted",
			err:          ErrOperationAborted,
			wantCode:     "OperationAbortedException",
			wantMsg:      "Operation aborted",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrServiceUnavailable",
			err:          ErrServiceUnavailable,
			wantCode:     "ServiceUnavailableException",
			wantMsg:      "Service unavailable",
			wantHTTPCode: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.wantCode {
				t.Errorf("Code = %v, want %v", tt.err.Code, tt.wantCode)
			}
			if tt.err.Message != tt.wantMsg {
				t.Errorf("Message = %v, want %v", tt.err.Message, tt.wantMsg)
			}
			if tt.err.HTTPStatus != tt.wantHTTPCode {
				t.Errorf("HTTPStatus = %v, want %v", tt.err.HTTPStatus, tt.wantHTTPCode)
			}
		})
	}
}

func TestNewLogsError(t *testing.T) {
	err := NewLogsError("TestCode", "Test message", http.StatusBadRequest)

	if err.Code != "TestCode" {
		t.Errorf("Code = %v, want TestCode", err.Code)
	}
	if err.Message != "Test message" {
		t.Errorf("Message = %v, want Test message", err.Message)
	}
	if err.HTTPStatus != http.StatusBadRequest {
		t.Errorf("HTTPStatus = %v, want %v", err.HTTPStatus, http.StatusBadRequest)
	}
}

func TestLogsErrorImplementsError(t *testing.T) {
	err := ErrLogGroupNotFound

	if err.Error() != "ResourceNotFoundException: Log group not found" {
		t.Errorf("Error() = %v, want 'ResourceNotFoundException: Log group not found'", err.Error())
	}
}
