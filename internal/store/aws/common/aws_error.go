package common

import (
	awserrors "vorpalstacks/internal/common/errors"
)

// AWSError represents an AWS-compatible error.
type AWSError = awserrors.AWSError

// NewAWSError creates a new AWS-compatible error.
var NewAWSError = awserrors.NewAWSError
