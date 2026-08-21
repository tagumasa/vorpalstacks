package kmsutil

import "errors"

// KMS validation sentinel errors. The KMS adapter returns these to allow
// the caller (e.g. SQS) to map them to service-specific AWS error codes
// without importing the KMS service package.
var (
	ErrKeyNotFound     = errors.New("kms key not found")
	ErrKeyDisabled     = errors.New("kms key is disabled")
	ErrKeyInvalidState = errors.New("kms key is in an invalid state")
	ErrKeyInvalidUsage = errors.New("kms key has invalid key usage")
)
