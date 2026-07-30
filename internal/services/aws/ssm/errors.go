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
	// ErrInvalidAllowedPattern is returned when the AllowedPattern is not a valid regular expression.
	ErrInvalidAllowedPattern = awserrors.NewAWSError("InvalidAllowedPatternException", "The regular expression is invalid", http.StatusBadRequest)
	// ErrInvalidResourceType is returned when a tag operation is given a
	// ResourceType that this implementation does not support.
	ErrInvalidResourceType = awserrors.NewAWSError("InvalidResourceType", "The resource type is not supported", http.StatusBadRequest)
	// ErrInvalidFilterKey is returned when a DescribeParameters/GetParametersByPath
	// filter Key is not in the Smithy ParameterStringFilterKey pattern.
	ErrInvalidFilterKey = awserrors.NewAWSError("InvalidFilterKey", "The filter key is not valid", http.StatusBadRequest)
	// ErrInvalidFilterOption is returned when a filter Option is not one of
	// Equals/BeginsWith/Contains.
	ErrInvalidFilterOption = awserrors.NewAWSError("InvalidFilterOption", "The filter option is not valid", http.StatusBadRequest)
)
