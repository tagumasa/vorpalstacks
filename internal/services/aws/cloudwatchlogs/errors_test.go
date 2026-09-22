package cloudwatchlogs

import (
	"net/http"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
)

func TestLogsError(t *testing.T) {
	tests := []struct {
		name         string
		err          *awserrors.AWSError
		wantCode     string
		wantMsg      string
		wantHTTPCode int
	}{
		{
			name:         "ErrLogGroupNotFound",
			err:          ErrLogGroupNotFound,
			wantCode:     "ResourceNotFoundException",
			wantMsg:      "Log group not found",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrLogGroupAlreadyExists",
			err:          ErrLogGroupAlreadyExists,
			wantCode:     "ResourceAlreadyExistsException",
			wantMsg:      "Log group already exists",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrLogStreamNotFound",
			err:          ErrLogStreamNotFound,
			wantCode:     "ResourceNotFoundException",
			wantMsg:      "Log stream not found",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrLogStreamAlreadyExists",
			err:          ErrLogStreamAlreadyExists,
			wantCode:     "ResourceAlreadyExistsException",
			wantMsg:      "Log stream already exists",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrInvalidParameter",
			err:          ErrInvalidParameter,
			wantCode:     "InvalidParameterException",
			wantMsg:      "Invalid parameter",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrLimitExceeded",
			err:          ErrLimitExceeded,
			wantCode:     "LimitExceededException",
			wantMsg:      "Limit exceeded",
			wantHTTPCode: http.StatusBadRequest,
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

func TestNewAWSError(t *testing.T) {
	err := awserrors.NewAWSError("TestCode", "Test message", http.StatusBadRequest)

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

// The per-op identity helpers: a missing member rejects with the
// operation family's declared validation error, and the scheduled-query
// normaliser maps the shared validators' StartQuery vocabulary
// (MalformedQueryException, InvalidParameterException) onto
// ValidationException while letting every other identity pass through.
func TestErrorIdentityHelpers(t *testing.T) {
	if code := logsErrorCode(errRequiredMember("queryId")); code != "InvalidParameterException" {
		t.Fatalf("errRequiredMember: code=%q", code)
	}
	if code := logsErrorCode(errValidationMember("identifier")); code != "ValidationException" {
		t.Fatalf("errValidationMember: code=%q", code)
	}
	malformed := NewLogsError("MalformedQueryException", "the query is malformed", 400)
	norm := asScheduledQueryValidation(malformed)
	if code := logsErrorCode(norm); code != "ValidationException" {
		t.Fatalf("normalised MalformedQuery: code=%q", code)
	}
	if code := logsErrorCode(asScheduledQueryValidation(ErrLogGroupNotFound)); code != "ResourceNotFoundException" {
		t.Fatalf("pass-through identity changed: code=%q", code)
	}
}
