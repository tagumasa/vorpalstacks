package iam

import (
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
)

// InvalidParameterValueError is returned when a parameter value is invalid.
type InvalidParameterValueError struct {
	Code    string
	Message string
}

// Error returns the error message.
func (e *InvalidParameterValueError) Error() string {
	return e.Message
}

// GetAWSError returns the underlying AWS error for dispatcher error extraction.
func (e *InvalidParameterValueError) GetAWSError() *awserrors.AWSError {
	code := e.Code
	if code == "" {
		code = "InvalidParameterValue"
	}
	return awserrors.NewAWSError(code, e.Message, 400)
}

// InvalidArnError is returned when an Amazon Resource Name (ARN) is invalid.
type InvalidArnError struct {
	Message string
}

// Error returns the error message.
func (e *InvalidArnError) Error() string {
	return e.Message
}

// GetAWSError returns the underlying AWS error for dispatcher error extraction.
func (e *InvalidArnError) GetAWSError() *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidArn", e.Message, 400)
}

// ValidationException is returned when validation fails.
type ValidationException struct {
	Message string
}

// Error returns the error message.
func (e *ValidationException) Error() string {
	return e.Message
}

// GetAWSError returns the underlying AWS error for dispatcher error extraction.
func (e *ValidationException) GetAWSError() *awserrors.AWSError {
	return awserrors.NewAWSError("ValidationException", e.Message, 400)
}

// InvalidParameterCombinationError is returned when parameters cannot be used together.
type InvalidParameterCombinationError struct {
	Message string
}

// Error returns the error message.
func (e *InvalidParameterCombinationError) Error() string {
	return e.Message
}

// GetAWSError returns the underlying AWS error for dispatcher error extraction.
func (e *InvalidParameterCombinationError) GetAWSError() *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidParameterCombination", e.Message, 400)
}

var (
	// ErrRoleNotFound is returned when the specified IAM role does not exist.
	ErrRoleNotFound = &InvalidParameterValueError{Code: "InvalidParameterValueException", Message: "role not found"}

	// ErrRoleCannotBeAssumed is returned when the IAM role cannot be assumed
	// by the service.
	ErrRoleCannotBeAssumed = &ValidationException{Message: "role cannot be assumed"}

	// ErrInvalidRoleARN is returned when the role ARN is not valid.
	ErrInvalidRoleARN = &InvalidArnError{Message: "invalid role ARN"}

	// ErrCrossAccountRole is returned when cross-account role validation is skipped.
	ErrCrossAccountRole = &ValidationException{Message: "cross-account role validation skipped"}
)

// NewLambdaRoleNotFoundError creates a new error indicating that the Lambda
// execution role does not exist.
func NewLambdaRoleNotFoundError(roleArn string) error {
	return &InvalidParameterValueError{
		Code:    "InvalidParameterValueException",
		Message: fmt.Sprintf("The role defined for the function (%s) does not exist.", roleArn),
	}
}

// NewLambdaRoleCannotBeAssumedError creates a new error indicating that the Lambda
// execution role cannot be assumed by Lambda.
func NewLambdaRoleCannotBeAssumedError(roleArn string) error {
	return &InvalidParameterValueError{
		Code:    "InvalidParameterValueException",
		Message: "The role defined for the function cannot be assumed by Lambda.",
	}
}

// NewInvalidRoleArnError creates a new error indicating that the role ARN is invalid.
func NewInvalidRoleArnError(roleArn string) error {
	return &InvalidArnError{
		Message: fmt.Sprintf("Invalid RoleArn: %s", roleArn),
	}
}

// NewRoleCannotBeAssumedError creates a new error indicating that the role
// cannot be assumed.
func NewRoleCannotBeAssumedError(roleArn string) error {
	return &ValidationException{
		Message: fmt.Sprintf("Role %s is invalid or cannot be assumed.", roleArn),
	}
}

// NewEventsRoleNotFoundError creates a new error indicating that the Events service
// role does not exist.
func NewEventsRoleNotFoundError(roleArn string) error {
	return &InvalidParameterValueError{
		Code:    "InvalidParameterValueException",
		Message: fmt.Sprintf("Parameter RoleArn is not valid. Reason: Role %s does not exist.", roleArn),
	}
}

// NewCloudTrailRoleError creates a new error indicating that the CloudTrail
// role does not exist or cannot be assumed.
func NewCloudTrailRoleError(roleArn string) error {
	return &InvalidParameterCombinationError{
		Message: fmt.Sprintf("Role ARN %s does not exist or cannot be assumed by CloudTrail.", roleArn),
	}
}

// NewTimestreamRoleError creates a new error indicating that the Timestream
// role does not exist or cannot be assumed.
func NewTimestreamRoleError(roleArn string) error {
	return &ValidationException{
		Message: fmt.Sprintf("Role %s is invalid or cannot be assumed by Timestream.", roleArn),
	}
}

// NewCognitoRoleError creates a new error indicating that the Cognito
// role ARN is invalid or cannot be assumed.
func NewCognitoRoleError(roleArn string) error {
	return &InvalidParameterValueError{
		Code:    "InvalidParameterException",
		Message: fmt.Sprintf("Role ARN %s is invalid or cannot be assumed.", roleArn),
	}
}
