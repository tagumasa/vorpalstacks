// Package errors provides AWS error types and helpers.
package errors

import (
	"net/http"

	awserrors "vorpalstacks/internal/utils/aws/errors"
)

// AWSError is an alias for awserrors.AWSError.
type AWSError = awserrors.AWSError

// XMLErrorResponse is an alias for awserrors.XMLErrorResponse.
type XMLErrorResponse = awserrors.XMLErrorResponse

// CustomJSONMarshaler is an alias for awserrors.CustomJSONMarshaler.
type CustomJSONMarshaler = awserrors.CustomJSONMarshaler

// Error variables provide common AWS error instances.
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

func init() {
	var _ *AWSError = (*awserrors.AWSError)(nil)
	var _ http.ResponseWriter = nil
}
