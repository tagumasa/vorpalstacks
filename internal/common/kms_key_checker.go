package common

import (
	"context"
	"errors"
)

// KMS validation sentinel errors. The KMS adapter returns these to allow the
// caller (e.g. SQS) to map them to service-specific AWS error codes without
// importing the KMS service package.
var (
	ErrKMSKeyNotFound     = errors.New("kms key not found")
	ErrKMSKeyDisabled     = errors.New("kms key is disabled")
	ErrKMSKeyInvalidState = errors.New("kms key is in an invalid state")
	ErrKMSKeyInvalidUsage = errors.New("kms key has invalid key usage")
)

// KMSKeyChecker validates a KMS key for use by another service (e.g. SQS
// KmsMasterKeyId). Implementations resolve the key by ID/alias/ARN and verify
// that it exists, is enabled, and has the correct key usage.
type KMSKeyChecker interface {
	CheckKey(ctx context.Context, region, keyID string) error
}
