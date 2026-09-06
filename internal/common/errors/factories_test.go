package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactories_ErrorTypes(t *testing.T) {
	tests := []struct {
		name       string
		err        *AWSError
		wantCode   string
		wantStatus int
		wantMsg    string
		wantFault  string
	}{
		{"ValidationException", NewValidationException("bad input"), "ValidationException", http.StatusBadRequest, "bad input", "Client"},
		{"ResourceNotFoundException", NewResourceNotFoundException("Queue", "q1"), "ResourceNotFoundException", http.StatusNotFound, "Queue q1 not found", "Client"},
		{"NotFoundException", NewNotFoundException("Thing"), "NotFoundException", http.StatusNotFound, "Thing not found", "Client"},
		{"InvalidParameterException", NewInvalidParameterException("param x"), "InvalidParameterException", http.StatusBadRequest, "param x", "Client"},
		{"AccessDeniedException", NewAccessDeniedException("nope"), "AccessDeniedException", http.StatusForbidden, "nope", "Client"},
		{"ThrottlingException", NewThrottlingException("slow down"), "ThrottlingException", http.StatusTooManyRequests, "slow down", "Client"},
		{"ServiceUnavailableException", NewServiceUnavailableException("down"), "ServiceUnavailableException", http.StatusServiceUnavailable, "down", "Server"},
		{"ConflictException", NewConflictException("exists"), "ConflictException", http.StatusConflict, "exists", "Client"},
		{"ResourceAlreadyExistsException", NewResourceAlreadyExistsException("q1"), "ResourceAlreadyExistsException", http.StatusConflict, "q1 already exists", "Client"},
		{"LimitExceededException", NewLimitExceededException("too many"), "LimitExceededException", http.StatusBadRequest, "too many", "Client"},
		{"BadRequestException", NewBadRequestException("bad req"), "BadRequestException", http.StatusBadRequest, "bad req", "Client"},
		{"InternalFailureException", NewInternalFailureException("fail"), "InternalFailure", http.StatusInternalServerError, "fail", "Server"},
		{"NoSuchEntityException", NewNoSuchEntityException("User", "bob"), "NoSuchEntity", http.StatusNotFound, "The User with name bob cannot be found.", "Client"},
		{"EntityAlreadyExistsException", NewEntityAlreadyExistsException("u1"), "EntityAlreadyExists", http.StatusConflict, "u1 already exists", "Client"},
		{"DeleteConflictException", NewDeleteConflictException("in use"), "DeleteConflict", http.StatusConflict, "in use", "Client"},
		{"InvalidInputException", NewInvalidInputException("inv"), "InvalidInput", http.StatusBadRequest, "inv", "Client"},
		{"InvalidParameterValueException", NewInvalidParameterValueException("pv"), "InvalidParameterValue", http.StatusBadRequest, "pv", "Client"},
		{"ResourceInUseException", NewResourceInUseException("busy"), "ResourceInUseException", http.StatusBadRequest, "busy", "Client"},
		{"MissingParameter", NewMissingParameter("req"), "MissingParameter", http.StatusBadRequest, "req", "Client"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.err)
			assert.Equal(t, tt.wantCode, tt.err.Code)
			assert.Equal(t, tt.wantStatus, tt.err.HTTPStatus)
			assert.Contains(t, tt.err.Message, tt.wantMsg)
			assert.NotEmpty(t, tt.err.RequestID)
			assert.Equal(t, tt.wantFault, tt.err.Fault)
		})
	}
}
