package route53

import (
	"net/http"
	"testing"

	"vorpalstacks/internal/services/aws/common/errors"
)

func TestRoute53Errors(t *testing.T) {
	tests := []struct {
		name         string
		err          *errors.AWSError
		wantCode     string
		wantMsg      string
		wantHTTPCode int
	}{
		{
			name:         "ErrNoSuchHostedZone",
			err:          ErrNoSuchHostedZone,
			wantCode:     "NoSuchHostedZone",
			wantMsg:      "The hosted zone does not exist.",
			wantHTTPCode: http.StatusNotFound,
		},
		{
			name:         "ErrHostedZoneAlreadyExists",
			err:          ErrHostedZoneAlreadyExists,
			wantCode:     "HostedZoneAlreadyExists",
			wantMsg:      "The hosted zone already exists.",
			wantHTTPCode: http.StatusConflict,
		},
		{
			name:         "ErrHostedZoneNotEmpty",
			err:          ErrHostedZoneNotEmpty,
			wantCode:     "HostedZoneNotEmpty",
			wantMsg:      "The hosted zone must be empty before it can be deleted.",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name:         "ErrNoSuchHealthCheck",
			err:          ErrNoSuchHealthCheck,
			wantCode:     "NoSuchHealthCheck",
			wantMsg:      "The health check does not exist.",
			wantHTTPCode: http.StatusNotFound,
		},
		{
			name:         "ErrHealthCheckAlreadyExists",
			err:          ErrHealthCheckAlreadyExists,
			wantCode:     "HealthCheckAlreadyExists",
			wantMsg:      "The health check already exists.",
			wantHTTPCode: http.StatusConflict,
		},
		{
			name:         "ErrNoSuchChange",
			err:          ErrNoSuchChange,
			wantCode:     "NoSuchChange",
			wantMsg:      "The change does not exist.",
			wantHTTPCode: http.StatusNotFound,
		},
		{
			name:         "ErrInvalidInput",
			err:          ErrInvalidInput,
			wantCode:     "InvalidInput",
			wantMsg:      "The input is invalid.",
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

func TestNewAPIError(t *testing.T) {
	err := NewAPIError("TestCode", "Test message", http.StatusBadRequest)

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

func TestNewNoSuchHostedZoneError(t *testing.T) {
	zoneID := "Z1234567890ABC"
	err := NewNoSuchHostedZoneError(zoneID)

	if err.Code != "NoSuchHostedZone" {
		t.Errorf("Code = %v, want NoSuchHostedZone", err.Code)
	}

	expectedMsg := "No hosted zone found with id: " + zoneID
	if err.Message != expectedMsg {
		t.Errorf("Message = %v, want %v", err.Message, expectedMsg)
	}
	if err.HTTPStatus != http.StatusNotFound {
		t.Errorf("HTTPStatus = %v, want %v", err.HTTPStatus, http.StatusNotFound)
	}
}

func TestNewHostedZoneAlreadyExistsError(t *testing.T) {
	name := "example.com."
	err := NewHostedZoneAlreadyExistsError(name)

	if err.Code != "HostedZoneAlreadyExists" {
		t.Errorf("Code = %v, want HostedZoneAlreadyExists", err.Code)
	}

	expectedMsg := "Hosted zone already exists: " + name
	if err.Message != expectedMsg {
		t.Errorf("Message = %v, want %v", err.Message, expectedMsg)
	}
	if err.HTTPStatus != http.StatusConflict {
		t.Errorf("HTTPStatus = %v, want %v", err.HTTPStatus, http.StatusConflict)
	}
}

func TestNewNoSuchHealthCheckError(t *testing.T) {
	healthCheckID := "abc123"
	err := NewNoSuchHealthCheckError(healthCheckID)

	if err.Code != "NoSuchHealthCheck" {
		t.Errorf("Code = %v, want NoSuchHealthCheck", err.Code)
	}

	expectedMsg := "No health check found with id: " + healthCheckID
	if err.Message != expectedMsg {
		t.Errorf("Message = %v, want %v", err.Message, expectedMsg)
	}
	if err.HTTPStatus != http.StatusNotFound {
		t.Errorf("HTTPStatus = %v, want %v", err.HTTPStatus, http.StatusNotFound)
	}
}
