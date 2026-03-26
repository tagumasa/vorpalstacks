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

	// ErrInvalidParameterVersion is returned when the parameter version is
	// not valid for the operation.
	ErrInvalidParameterVersion = errors.New("invalid parameter version")

	// ErrParameterVersionNotFound is returned when the specified parameter
	// version does not exist.
	ErrParameterVersionNotFound = errors.New("parameter version not found")

	// ErrParameterLabelNotFound is returned when the specified parameter
	// label does not exist.
	ErrParameterLabelNotFound = errors.New("parameter label not found")
)
