package ssm

import (
	"net/http"
	"testing"
)

func TestSSMError(t *testing.T) {
	tests := []struct {
		name         string
		err          *SSMError
		wantCode     string
		wantMsg      string
		wantHTTPCode int
	}{
		{
			name:         "ErrParameterNotFound",
			err:          ErrParameterNotFound,
			wantCode:     "ParameterNotFound",
			wantMsg:      "Parameter not found",
			wantHTTPCode: http.StatusNotFound,
		},
		{
			name:         "ErrParameterAlreadyExists",
			err:          ErrParameterAlreadyExists,
			wantCode:     "ParameterAlreadyExists",
			wantMsg:      "Parameter already exists",
			wantHTTPCode: http.StatusConflict,
		},
		{
			name:         "ErrInvalidParameterName",
			err:          ErrInvalidParameterName,
			wantCode:     "InvalidParameter",
			wantMsg:      "Invalid parameter name",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrInvalidParameterValue",
			err:          ErrInvalidParameterValue,
			wantCode:     "InvalidParameter",
			wantMsg:      "Invalid parameter value",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrInvalidParameterType",
			err:          ErrInvalidParameterType,
			wantCode:     "InvalidParameterType",
			wantMsg:      "Invalid parameter type",
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

func TestNewSSMError(t *testing.T) {
	err := NewSSMError("TestCode", "Test message", http.StatusBadRequest)

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
