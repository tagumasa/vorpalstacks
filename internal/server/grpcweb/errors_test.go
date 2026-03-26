package grpcweb

import (
	"net/http"
	"testing"
)

func TestServerError_Error(t *testing.T) {
	err := &ServerError{
		Message:    "test error",
		HTTPStatus: http.StatusBadRequest,
	}

	if err.Error() != "test error" {
		t.Errorf("expected 'test error', got '%s'", err.Error())
	}
}

func TestNewServerError(t *testing.T) {
	err := NewServerError("test error", http.StatusBadRequest)

	if err.Message != "test error" {
		t.Errorf("expected Message 'test error', got '%s'", err.Message)
	}
	if err.HTTPStatus != http.StatusBadRequest {
		t.Errorf("expected HTTPStatus %d, got %d", http.StatusBadRequest, err.HTTPStatus)
	}
}

func TestErrInvalidRequest(t *testing.T) {
	err := ErrInvalidRequest("invalid request")

	if err.Message != "invalid request" {
		t.Errorf("expected Message 'invalid request', got '%s'", err.Message)
	}
	if err.HTTPStatus != http.StatusBadRequest {
		t.Errorf("expected HTTPStatus %d, got %d", http.StatusBadRequest, err.HTTPStatus)
	}
}

func TestErrUnauthorized(t *testing.T) {
	err := ErrUnauthorized("unauthorized")

	if err.Message != "unauthorized" {
		t.Errorf("expected Message 'unauthorized', got '%s'", err.Message)
	}
	if err.HTTPStatus != http.StatusUnauthorized {
		t.Errorf("expected HTTPStatus %d, got %d", http.StatusUnauthorized, err.HTTPStatus)
	}
}

func TestErrForbidden(t *testing.T) {
	err := ErrForbidden("forbidden")

	if err.Message != "forbidden" {
		t.Errorf("expected Message 'forbidden', got '%s'", err.Message)
	}
	if err.HTTPStatus != http.StatusForbidden {
		t.Errorf("expected HTTPStatus %d, got %d", http.StatusForbidden, err.HTTPStatus)
	}
}

func TestErrNotFound(t *testing.T) {
	err := ErrNotFound("not found")

	if err.Message != "not found" {
		t.Errorf("expected Message 'not found', got '%s'", err.Message)
	}
	if err.HTTPStatus != http.StatusNotFound {
		t.Errorf("expected HTTPStatus %d, got %d", http.StatusNotFound, err.HTTPStatus)
	}
}

func TestErrInternal(t *testing.T) {
	err := ErrInternal("internal error")

	if err.Message != "internal error" {
		t.Errorf("expected Message 'internal error', got '%s'", err.Message)
	}
	if err.HTTPStatus != http.StatusInternalServerError {
		t.Errorf("expected HTTPStatus %d, got %d", http.StatusInternalServerError, err.HTTPStatus)
	}
}
