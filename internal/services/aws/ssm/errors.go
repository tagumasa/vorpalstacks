package ssm

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

var (
	// ErrParameterNotFound is returned when the specified parameter does not exist.
	ErrParameterNotFound = awserrors.NewAWSError("ParameterNotFound", "Parameter not found", http.StatusNotFound)
	// ErrParameterAlreadyExists is returned when attempting to create a parameter that already exists.
	ErrParameterAlreadyExists = awserrors.NewAWSError("ParameterAlreadyExists", "Parameter already exists", http.StatusBadRequest)
	// ErrInvalidParameterName is returned when the parameter name is invalid.
	ErrInvalidParameterName = awserrors.NewAWSError("InvalidParameter", "Invalid parameter name", http.StatusBadRequest)
	// ErrInvalidParameterValue is returned when the parameter value is invalid.
	ErrInvalidParameterValue = awserrors.NewAWSError("InvalidParameter", "Invalid parameter value", http.StatusBadRequest)
	// ErrInvalidParameterType is returned when the parameter type is invalid.
	ErrInvalidParameterType = awserrors.NewAWSError("InvalidParameterType", "Invalid parameter type", http.StatusBadRequest)
	// ErrInvalidParameterVersion is returned when the parameter version is invalid.
	ErrInvalidParameterVersion = awserrors.NewAWSError("InvalidParameterVersion", "Invalid parameter version", http.StatusBadRequest)
	// ErrInvalidParameterLabel is returned when the parameter label is invalid.
	ErrInvalidParameterLabel = awserrors.NewAWSError("InvalidParameter", "Invalid parameter label", http.StatusBadRequest)
	// ErrParameterVersionNotFound is returned when the specified parameter version does not exist.
	ErrParameterVersionNotFound = awserrors.NewAWSError("ParameterVersionNotFound", "Parameter version not found", http.StatusBadRequest)
	// ErrParameterPatternMismatch is returned when the parameter name uses a reserved prefix.
	ErrParameterPatternMismatch = awserrors.NewAWSError("ParameterPatternMismatchException", "The parameter name is not valid", http.StatusBadRequest)
)
