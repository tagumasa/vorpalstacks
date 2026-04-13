package ssm

import (
	"net/http"

	"vorpalstacks/internal/common/errors"
)

// NewSSMError creates a new SSMError with the specified code, message and status code.
func NewSSMError(code, message string, statusCode int) *errors.AWSError {
	return errors.NewAWSError(code, message, statusCode)
}

var (
	// ErrParameterNotFound is returned when the specified parameter does not exist.
	ErrParameterNotFound = NewSSMError("ParameterNotFound", "Parameter not found", http.StatusNotFound)
	// ErrParameterAlreadyExists is returned when attempting to create a parameter that already exists.
	ErrParameterAlreadyExists = NewSSMError("ParameterAlreadyExists", "Parameter already exists", http.StatusConflict)
	// ErrInvalidParameterName is returned when the parameter name is invalid.
	ErrInvalidParameterName = NewSSMError("InvalidParameter", "Invalid parameter name", http.StatusBadRequest)
	// ErrInvalidParameterValue is returned when the parameter value is invalid.
	ErrInvalidParameterValue = NewSSMError("InvalidParameter", "Invalid parameter value", http.StatusBadRequest)
	// ErrInvalidParameterType is returned when the parameter type is invalid.
	ErrInvalidParameterType = NewSSMError("InvalidParameterType", "Invalid parameter type", http.StatusBadRequest)
	// ErrInvalidParameterVersion is returned when the parameter version is invalid.
	ErrInvalidParameterVersion = NewSSMError("InvalidParameterVersion", "Invalid parameter version", http.StatusBadRequest)
	// ErrInvalidParameterLabel is returned when the parameter label is invalid.
	ErrInvalidParameterLabel = NewSSMError("InvalidParameter", "Invalid parameter label", http.StatusBadRequest)
	// ErrParameterVersionNotFound is returned when the specified parameter version does not exist.
	ErrParameterVersionNotFound = NewSSMError("ParameterVersionNotFound", "Parameter version not found", http.StatusNotFound)
)
