// Package grpcweb provides gRPC-Web server functionality and error types for vorpalstacks.
package grpcweb

import (
	"errors"
	"net/http"
)

// Error variables for grpcweb operations.
var (
	// ErrServerNotFound is returned when attempting to interact with a server
	// that has not been started.
	ErrServerNotStarted = errors.New("server not started")

	// ErrHandlerExists is returned when attempting to register a handler
	// that has already been registered.
	ErrHandlerExists = errors.New("handler already exists")
)

// ServerError represents an error with an associated HTTP status code.
type ServerError struct {
	Message    string
	HTTPStatus int
}

// Error returns the error message.
func (e *ServerError) Error() string {
	return e.Message
}

// NewServerError creates a new ServerError with the given message and HTTP status code.
func NewServerError(message string, httpStatus int) *ServerError {
	return &ServerError{
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// ErrInvalidRequest returns a bad request error with the given message.
func ErrInvalidRequest(message string) *ServerError {
	return &ServerError{
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// ErrUnauthorized returns an unauthorized error with the given message.
func ErrUnauthorized(message string) *ServerError {
	return &ServerError{
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}

// ErrForbidden returns a forbidden error with the given message.
func ErrForbidden(message string) *ServerError {
	return &ServerError{
		Message:    message,
		HTTPStatus: http.StatusForbidden,
	}
}

// ErrNotFound returns a not found error with the given message.
func ErrNotFound(message string) *ServerError {
	return &ServerError{
		Message:    message,
		HTTPStatus: http.StatusNotFound,
	}
}

// ErrInternal returns an internal server error with the given message.
func ErrInternal(message string) *ServerError {
	return &ServerError{
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
	}
}
