// Package lambda provides AWS Lambda store functionality for vorpalstacks.
package lambda

import "errors"

var (
	// ErrFunctionNotFound is returned when the specified Lambda function
	// does not exist.
	ErrFunctionNotFound = errors.New("function not found")

	// ErrFunctionAlreadyExists is returned when attempting to create a function
	// that already exists.
	ErrFunctionAlreadyExists = errors.New("function already exists")

	// ErrVersionNotFound is returned when the specified function version
	// does not exist.
	ErrVersionNotFound = errors.New("version not found")

	// ErrAliasNotFound is returned when the specified function alias
	// does not exist.
	ErrAliasNotFound = errors.New("alias not found")

	// ErrAliasAlreadyExists is returned when attempting to create an alias
	// that already exists.
	ErrAliasAlreadyExists = errors.New("alias already exists")

	// ErrLayerNotFound is returned when the specified Lambda layer
	// does not exist.
	ErrLayerNotFound = errors.New("layer not found")

	// ErrLayerVersionNotFound is returned when the specified layer version
	// does not exist.
	ErrLayerVersionNotFound = errors.New("layer version not found")

	// ErrEventSourceNotFound is returned when the specified event source mapping
	// does not exist.
	ErrEventSourceNotFound = errors.New("event source mapping not found")

	// ErrEventSourceAlreadyExists is returned when attempting to create an event
	// source mapping that already exists for the given event source and function.
	ErrEventSourceAlreadyExists = errors.New("event source mapping already exists for this event source and function")

	// ErrInvalidParameterValue is returned when a parameter value is not valid.
	ErrInvalidParameterValue = errors.New("invalid parameter value")

	// ErrInvalidRuntime is returned when the specified runtime is not valid
	// or not supported.
	ErrInvalidRuntime = errors.New("invalid runtime")

	// ErrCodeSigningNotFound is returned when the specified code signing
	// configuration does not exist.
	ErrCodeSigningNotFound = errors.New("code signing config not found")

	// ErrProvisionedConcurrencyNotFound is returned when the specified
	// provisioned concurrency configuration does not exist.
	ErrProvisionedConcurrencyNotFound = errors.New("provisioned concurrency config not found")

	// ErrEventInvokeConfigNotFound is returned when the specified event invoke
	// configuration does not exist.
	ErrEventInvokeConfigNotFound = errors.New("event invoke config not found")

	// ErrUrlConfigNotFound is returned when the specified function URL
	// configuration does not exist.
	ErrUrlConfigNotFound = errors.New("function url config not found")

	// ErrPolicyNotFound is returned when the specified resource policy
	// does not exist.
	ErrPolicyNotFound = errors.New("resource policy not found")

	// ErrResourceConflict is returned when the requested operation conflicts
	// with the current state of the resource.
	ErrResourceConflict = errors.New("resource conflict")
)
