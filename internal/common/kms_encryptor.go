// Package common provides common AWS utilities for vorpalstacks.
package common

import "context"

// KMSEncryptor defines the interface for KMS encryption operations.
type KMSEncryptor interface {
	EncryptString(ctx context.Context, keyID string, plaintext string) (string, error)
	DecryptString(ctx context.Context, keyID string, ciphertext string) (string, error)
}
