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
	}{
		{"ValidationException", NewValidationException("bad input"), "ValidationException", http.StatusBadRequest, "bad input"},
		{"ResourceNotFoundException", NewResourceNotFoundException("Queue", "q1"), "ResourceNotFoundException", http.StatusNotFound, "Queue q1 not found"},
		{"NotFoundException", NewNotFoundException("Thing"), "NotFoundException", http.StatusNotFound, "Thing not found"},
		{"InvalidParameterException", NewInvalidParameterException("param x"), "InvalidParameterException", http.StatusBadRequest, "param x"},
		{"AccessDeniedException", NewAccessDeniedException("nope"), "AccessDeniedException", http.StatusForbidden, "nope"},
		{"ThrottlingException", NewThrottlingException("slow down"), "ThrottlingException", http.StatusTooManyRequests, "slow down"},
		{"ServiceUnavailableException", NewServiceUnavailableException("down"), "ServiceUnavailableException", http.StatusServiceUnavailable, "down"},
		{"ConflictException", NewConflictException("exists"), "ConflictException", http.StatusConflict, "exists"},
		{"ResourceAlreadyExistsException", NewResourceAlreadyExistsException("q1"), "ResourceAlreadyExistsException", http.StatusConflict, "q1 already exists"},
		{"LimitExceededException", NewLimitExceededException("too many"), "LimitExceededException", http.StatusBadRequest, "too many"},
		{"BadRequestException", NewBadRequestException("bad req"), "BadRequestException", http.StatusBadRequest, "bad req"},
		{"InternalErrorException", NewInternalErrorException("oops"), "InternalError", http.StatusInternalServerError, "oops"},
		{"InternalFailureException", NewInternalFailureException("fail"), "InternalFailure", http.StatusInternalServerError, "fail"},
		{"NoSuchEntityException", NewNoSuchEntityException("User", "bob"), "NoSuchEntity", http.StatusNotFound, "The User with name bob cannot be found."},
		{"EntityAlreadyExistsException", NewEntityAlreadyExistsException("u1"), "EntityAlreadyExists", http.StatusConflict, "u1 already exists"},
		{"DeleteConflictException", NewDeleteConflictException("in use"), "DeleteConflict", http.StatusConflict, "in use"},
		{"InvalidInputException", NewInvalidInputException("inv"), "InvalidInput", http.StatusBadRequest, "inv"},
		{"InvalidParameterValueException", NewInvalidParameterValueException("pv"), "InvalidParameterValue", http.StatusBadRequest, "pv"},
		{"ResourceInUseException", NewResourceInUseException("busy"), "ResourceInUseException", http.StatusBadRequest, "busy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.err)
			assert.Equal(t, tt.wantCode, tt.err.Code)
			assert.Equal(t, tt.wantStatus, tt.err.HTTPStatus)
			assert.Contains(t, tt.err.Message, tt.wantMsg)
			assert.NotEmpty(t, tt.err.RequestID)
			assert.Equal(t, "Server", tt.err.Fault)
		})
	}
}
