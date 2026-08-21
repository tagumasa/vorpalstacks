// Package kmsutil holds the cross-service KMS contracts. Services that
// need KMS operations without importing the KMS service package depend
// on these interfaces; the eventbus invokers provide the
// implementations.
package kmsutil

import "context"

// Encryptor defines the interface for KMS encryption operations.
type Encryptor interface {
	EncryptString(ctx context.Context, keyID string, plaintext string) (string, error)
	DecryptString(ctx context.Context, keyID string, ciphertext string) (string, error)
}

// Checker validates a KMS key for use by another service (e.g. SQS
// KmsMasterKeyId). Implementations resolve the key by ID/alias/ARN and
// verify that it exists, is enabled, and has the correct key usage.
type Checker interface {
	CheckKey(ctx context.Context, region, keyID string) error
}
