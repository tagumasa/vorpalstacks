package errors

import (
	awserrors "vorpalstacks/internal/utils/aws/errors"
)

// AWSError is an alias for the core AWS error type, re-exported for
// convenient access by service and server code.
type AWSError = awserrors.AWSError

// CustomJSONMarshaler is an alias for the custom JSON error response
// interface, used by the dispatcher for type-assertion-based error handling.
type CustomJSONMarshaler = awserrors.CustomJSONMarshaler

// Re-exported error construction helpers and sentinel values from the
// core error package. All service-layer code should import this
// package rather than reaching into internal/utils.
var (
	NewAWSError             = awserrors.NewAWSError
	WriteAWSError           = awserrors.WriteAWSError
	WriteCustomJSONError    = awserrors.WriteCustomJSONError
	WriteAWSErrorWithFormat = awserrors.WriteAWSErrorWithFormat
	ErrAccessDenied         = awserrors.ErrAccessDenied
	ErrInvalidAction        = awserrors.ErrInvalidAction
	ErrInvalidParameter     = awserrors.ErrInvalidParameter
	ErrMissingParameter     = awserrors.ErrMissingParameter
	ErrResourceNotFound     = awserrors.ErrResourceNotFound
	ErrServiceUnavailable   = awserrors.ErrServiceUnavailable
	ErrThrottling           = awserrors.ErrThrottling
	ErrValidation           = awserrors.ErrValidation
	ErrInternal             = awserrors.ErrInternal
	ErrNotImplemented       = awserrors.ErrNotImplemented
)
