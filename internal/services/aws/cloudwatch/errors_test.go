package cloudwatch

import (
	"net/http"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
)

func TestCloudWatchErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          *awserrors.AWSError
		wantCode     string
		wantMsg      string
		wantHTTPCode int
	}{
		{
			name:         "ErrInvalidParameter",
			err:          ErrInvalidParameter,
			wantCode:     "InvalidParameterException",
			wantMsg:      "Invalid parameter",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrResourceNotFound",
			err:          ErrResourceNotFound,
			wantCode:     "ResourceNotFoundException",
			wantMsg:      "Resource  not found",
			wantHTTPCode: http.StatusNotFound,
		},
		{
			name:         "ErrResourceAlreadyExists",
			err:          ErrResourceAlreadyExists,
			wantCode:     "ResourceAlreadyExistsException",
			wantMsg:      "Resource already exists",
			wantHTTPCode: http.StatusConflict,
		},
		{
			name:         "ErrInvalidSequenceToken",
			err:          ErrInvalidSequenceToken,
			wantMsg:      "Invalid sequence token",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrDataAlreadyAccepted",
			err:          ErrDataAlreadyAccepted,
			wantMsg:      "Data already accepted",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrOperationAborted",
			err:          ErrOperationAborted,
			wantMsg:      "Operation aborted",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrLimitExceeded",
			err:          ErrLimitExceeded,
			wantMsg:      "Limit exceeded",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrServiceUnavailable",
			err:          ErrServiceUnavailable,
			wantMsg:      "Service unavailable",
			wantHTTPCode: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantCode != "" && tt.err.Code != tt.wantCode {
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

func TestNewInvalidSequenceTokenError(t *testing.T) {
	expectedToken := "test-token-123"
	err := NewInvalidSequenceTokenError(expectedToken)

	if err.Code != "InvalidSequenceTokenException" {
		t.Errorf("Code = %v, want InvalidSequenceTokenException", err.Code)
	}

	if err.HTTPStatus != http.StatusBadRequest {
		t.Errorf("HTTPStatus = %v, want %v", err.HTTPStatus, http.StatusBadRequest)
	}

	expectedMsg := "The given sequenceToken is invalid. The next expected sequenceToken is: " + expectedToken
	if err.Message != expectedMsg {
		t.Errorf("Message = %v, want %v", err.Message, expectedMsg)
	}
}
