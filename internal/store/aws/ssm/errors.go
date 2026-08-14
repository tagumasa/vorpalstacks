package ssm

import "errors"

var (
	// ErrParameterNotFound is returned when the specified Systems Manager
	// parameter does not exist.
	ErrParameterNotFound = errors.New("parameter not found")

	// ErrParameterAlreadyExists is returned when attempting to create a
	// parameter that already exists.
	ErrParameterAlreadyExists = errors.New("parameter already exists")

	// ErrInvalidParameterName is returned when the parameter name does not
	// meet SSM naming requirements.
	ErrInvalidParameterName = errors.New("invalid parameter name")

	// ErrInvalidParameterValue is returned when the parameter value is not valid.
	ErrInvalidParameterValue = errors.New("invalid parameter value")

	// ErrInvalidParameterType is returned when the parameter type is not valid
	// (e.g., String, StringList, SecureString).
	ErrInvalidParameterType = errors.New("invalid parameter type")

	// ErrInvalidParameterVersion is returned when the parameter version is not
	// valid for the operation.
	ErrInvalidParameterVersion = errors.New("invalid parameter version")

	// ErrParameterVersionNotFound is returned when the specified parameter
	// version does not exist.
	ErrParameterVersionNotFound = errors.New("parameter version not found")

	// ErrParameterLabelNotFound is returned when the specified parameter
	// label does not exist.
	ErrParameterLabelNotFound = errors.New("parameter label not found")

	// ErrReservedParameterName is returned when the parameter name uses a
	// reserved prefix such as "aws" or "ssm".
	ErrReservedParameterName = errors.New("parameter name uses a reserved prefix")

	// ErrInvalidAllowedPattern is returned when an AllowedPattern value is
	// not a valid regular expression.
	ErrInvalidAllowedPattern = errors.New("parameter allowed pattern is not a valid regular expression")

	// ErrParameterPatternMismatch is returned when a value would violate the
	// AllowedPattern constraint already attached to a parameter.
	ErrParameterPatternMismatch = errors.New("parameter value does not match allowed pattern")

	// ErrHierarchyLevelLimitExceeded is returned when a parameter name
	// hierarchy exceeds the maximum depth of fifteen levels.
	ErrHierarchyLevelLimitExceeded = errors.New("parameter name hierarchy exceeds the maximum depth")

	// ErrInvalidNextToken is returned when a pagination marker does not
	// parse as a valid version reference.
	ErrInvalidNextToken = errors.New("invalid next token")
)
