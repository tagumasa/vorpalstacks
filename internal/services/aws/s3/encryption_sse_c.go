// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"

	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/crypto"
)

// SSECEncryptor handles S3 server-side encryption with customer-provided keys.
type SSECEncryptor struct{}

// NewSSECEncryptor creates a new customer-provided key encryptor.
func NewSSECEncryptor() *SSECEncryptor {
	return &SSECEncryptor{}
}

// GetEncryptionType returns the encryption type.
func (e *SSECEncryptor) GetEncryptionType() EncryptionType {
	return EncryptionTypeSSE_C
}

// ParseCustomerKey parses and validates a base64-encoded customer key.
// Returns the parsed key bytes for per-request use.
func (e *SSECEncryptor) ParseCustomerKey(encodedKey, encodedMD5 string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("invalid customer key encoding: %w", err)
	}

	if encodedMD5 != "" {
		expectedMD5, err := base64.StdEncoding.DecodeString(encodedMD5)
		if err != nil {
			return nil, fmt.Errorf("invalid customer key MD5 encoding: %w", err)
		}
		actualMD5 := md5.Sum(key)
		if !equalBytes(expectedMD5, actualMD5[:]) {
			return nil, fmt.Errorf("customer key MD5 mismatch")
		}
	}

	if len(key) != 32 {
		secondKey, err2 := base64.StdEncoding.DecodeString(string(key))
		if err2 != nil || len(secondKey) != 32 {
			return nil, fmt.Errorf("customer key must be 32 bytes (AES-256)")
		}
		key = secondKey
	}

	return key, nil
}

// Encrypt encrypts data using the provided customer key.
func (e *SSECEncryptor) Encrypt(plaintext []byte, customerKey []byte, bucket, key string) (*EncryptionResult, error) {
	if customerKey == nil {
		return nil, fmt.Errorf("customer key not set")
	}

	nonce, err := crypto.RandomNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	encryptedData, err := crypto.AESGCMEncryptWithNonce(customerKey, plaintext, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	unencryptedMD5 := base64.StdEncoding.EncodeToString(crypto.MD5Hash(plaintext))

	return &EncryptionResult{
		EncryptedData:   encryptedData,
		EncryptionType:  EncryptionTypeSSE_C,
		ContentNonce:    nonce,
		UnencryptedMD5:  unencryptedMD5,
		UnencryptedSize: int64(len(plaintext)),
	}, nil
}

// Decrypt decrypts data using the provided customer key.
func (e *SSECEncryptor) Decrypt(ciphertext []byte, customerKey []byte, sseMetadata *s3store.SSEObjectMetadata, bucket, key string) (*DecryptionResult, error) {
	if customerKey == nil {
		return nil, fmt.Errorf("customer key not set")
	}

	if sseMetadata == nil || sseMetadata.ContentNonce == nil {
		return nil, fmt.Errorf("missing SSE-C metadata for encrypted object")
	}

	plaintext, err := crypto.AESGCMDecryptWithNonce(customerKey, ciphertext, sseMetadata.ContentNonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return &DecryptionResult{DecryptedData: plaintext}, nil
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
