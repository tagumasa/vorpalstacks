// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"fmt"

	s3store "vorpalstacks/internal/store/aws/s3"
)

// EncryptionType represents the type of S3 server-side encryption.
type EncryptionType string

const (
	// EncryptionTypeNone indicates no server-side encryption.
	EncryptionTypeNone EncryptionType = ""
	// EncryptionTypeSSE_S3 indicates Amazon S3-managed keys.
	EncryptionTypeSSE_S3 EncryptionType = "AES256"
	// EncryptionTypeSSE_KMS indicates AWS KMS-managed keys.
	EncryptionTypeSSE_KMS EncryptionType = "aws:kms"
	// EncryptionTypeSSE_C indicates customer-provided keys.
	EncryptionTypeSSE_C EncryptionType = "CUSTOMER"
)

// EncryptionResult holds the result of an encryption operation.
type EncryptionResult struct {
	EncryptedData    []byte
	EncryptionType   EncryptionType
	EncryptedDataKey []byte
	ContentNonce     []byte
	KMSKeyID         string
	UnencryptedMD5   string
	UnencryptedSize  int64
}

// DecryptionResult holds the result of a decryption operation.
type DecryptionResult struct {
	DecryptedData []byte
}

// Encryptor defines the interface for S3 encryption operations.
type Encryptor interface {
	Encrypt(plaintext []byte, bucket, key string) (*EncryptionResult, error)
	Decrypt(ciphertext []byte, sseMetadata *s3store.SSEObjectMetadata, bucket, key string) (*DecryptionResult, error)
	GetEncryptionType() EncryptionType
}

// EncryptionManager manages different S3 encryption methods.
type EncryptionManager struct {
	sseS3Encryptor  *SSES3Encryptor
	sseKMSEncryptor *SSEKMSEncryptor
	sseCEncryptor   *SSECEncryptor
}

// NewEncryptionManager creates a new encryption manager.
func NewEncryptionManager() *EncryptionManager {
	return &EncryptionManager{
		sseS3Encryptor: NewSSES3Encryptor(),
		sseCEncryptor:  NewSSECEncryptor(),
	}
}

// NewEncryptionManagerWithKMS creates a new encryption manager with KMS support.
func NewEncryptionManagerWithKMS(kmsClient KMSClient) *EncryptionManager {
	return &EncryptionManager{
		sseS3Encryptor:  NewSSES3Encryptor(),
		sseKMSEncryptor: NewSSEKMSEncryptor(kmsClient),
		sseCEncryptor:   NewSSECEncryptor(),
	}
}

// Encrypt encrypts plaintext using the specified encryption type.
func (m *EncryptionManager) Encrypt(plaintext []byte, encryptionType EncryptionType, bucketEncryption *s3store.EncryptionConfig, bucket, key, kmsKeyID string) (*EncryptionResult, error) {
	return m.EncryptWithCustomerKey(plaintext, encryptionType, bucketEncryption, bucket, key, kmsKeyID, nil)
}

// EncryptWithCustomerKey encrypts plaintext using a customer-provided key.
func (m *EncryptionManager) EncryptWithCustomerKey(plaintext []byte, encryptionType EncryptionType, bucketEncryption *s3store.EncryptionConfig, bucket, key, kmsKeyID string, customerKey []byte) (*EncryptionResult, error) {
	switch encryptionType {
	case EncryptionTypeSSE_S3:
		return m.sseS3Encryptor.Encrypt(plaintext, bucket, key)
	case EncryptionTypeSSE_KMS:
		if m.sseKMSEncryptor == nil {
			return nil, fmt.Errorf("SSE-KMS encryption requested but KMS is not configured")
		}
		effectiveKMSKeyID := kmsKeyID
		if effectiveKMSKeyID == "" && bucketEncryption != nil {
			effectiveKMSKeyID = bucketEncryption.KMSMasterKeyID
		}
		if effectiveKMSKeyID == "" {
			return nil, fmt.Errorf("KMS key ID is required for SSE-KMS encryption")
		}
		return m.sseKMSEncryptor.Encrypt(plaintext, bucket, key, effectiveKMSKeyID)
	case EncryptionTypeSSE_C:
		if m.sseCEncryptor == nil {
			return nil, fmt.Errorf("SSE-C encryption requested but SSE-C is not configured")
		}
		return m.sseCEncryptor.Encrypt(plaintext, customerKey, bucket, key)
	default:
		return &EncryptionResult{
			EncryptedData:  plaintext,
			EncryptionType: EncryptionTypeNone,
		}, nil
	}
}

// Decrypt decrypts ciphertext using the appropriate encryption type.
func (m *EncryptionManager) Decrypt(ciphertext []byte, sseMetadata *s3store.SSEObjectMetadata, bucket, key string) (*DecryptionResult, error) {
	return m.DecryptWithCustomerKey(ciphertext, sseMetadata, bucket, key, nil)
}

// DecryptWithCustomerKey decrypts ciphertext using a customer-provided key.
func (m *EncryptionManager) DecryptWithCustomerKey(ciphertext []byte, sseMetadata *s3store.SSEObjectMetadata, bucket, key string, customerKey []byte) (*DecryptionResult, error) {
	if sseMetadata == nil {
		return &DecryptionResult{DecryptedData: ciphertext}, nil
	}

	switch sseMetadata.EncryptionType {
	case s3store.SSETypeAES256:
		return m.sseS3Encryptor.Decrypt(ciphertext, sseMetadata, bucket, key)
	case s3store.SSETypeKMS:
		if m.sseKMSEncryptor == nil {
			return &DecryptionResult{DecryptedData: ciphertext}, nil
		}
		return m.sseKMSEncryptor.Decrypt(ciphertext, sseMetadata, bucket, key)
	case s3store.SSEType("CUSTOMER"):
		if m.sseCEncryptor == nil {
			return &DecryptionResult{DecryptedData: ciphertext}, nil
		}
		return m.sseCEncryptor.Decrypt(ciphertext, customerKey, sseMetadata, bucket, key)
	default:
		return &DecryptionResult{DecryptedData: ciphertext}, nil
	}
}

// ShouldEncrypt determines whether encryption should be applied based on the encryption type and bucket configuration.
func (m *EncryptionManager) ShouldEncrypt(encryptionType EncryptionType, bucketEncryption *s3store.EncryptionConfig) bool {
	if encryptionType != EncryptionTypeNone {
		return true
	}
	return bucketEncryption != nil && bucketEncryption.SSEAlgorithm != ""
}

// DetermineEncryptionType determines the effective encryption type based on request and bucket configuration.
func (m *EncryptionManager) DetermineEncryptionType(requestEncryption EncryptionType, bucketEncryption *s3store.EncryptionConfig) EncryptionType {
	if requestEncryption != EncryptionTypeNone {
		return requestEncryption
	}
	if bucketEncryption != nil {
		switch bucketEncryption.SSEAlgorithm {
		case "AES256":
			return EncryptionTypeSSE_S3
		case "aws:kms":
			return EncryptionTypeSSE_KMS
		}
	}
	return EncryptionTypeNone
}

// EncryptWithPlaintextKey encrypts plaintext using a provided plaintext key.
func (m *EncryptionManager) EncryptWithPlaintextKey(plaintext []byte, encryptionType EncryptionType, bucket string, plaintextKey []byte, kmsKeyID string, encryptedDataKey []byte) (*EncryptionResult, error) {
	switch encryptionType {
	case EncryptionTypeSSE_S3:
		return m.sseS3Encryptor.EncryptWithPlaintextKey(plaintext, bucket, plaintextKey)
	case EncryptionTypeSSE_KMS:
		if m.sseKMSEncryptor == nil {
			return nil, fmt.Errorf("KMS client not configured")
		}
		return m.sseKMSEncryptor.EncryptWithPlaintextKey(plaintext, plaintextKey, kmsKeyID, encryptedDataKey)
	default:
		return nil, fmt.Errorf("unsupported encryption type for EncryptWithPlaintextKey: %s", encryptionType)
	}
}

// GeneratedKey holds an encryption key and its encrypted form.
type GeneratedKey struct {
	PlaintextKey     []byte
	EncryptedDataKey []byte
	ContentNonce     []byte
	KMSKeyID         string
}

// GenerateKey generates an encryption key for the specified encryption type.
func (m *EncryptionManager) GenerateKey(encryptionType EncryptionType, bucketEncryption *s3store.EncryptionConfig, bucket, key, kmsKeyID string) (*GeneratedKey, error) {
	switch encryptionType {
	case EncryptionTypeSSE_S3:
		return m.sseS3Encryptor.GenerateKey(bucket)
	case EncryptionTypeSSE_KMS:
		if m.sseKMSEncryptor == nil {
			return nil, fmt.Errorf("KMS client not configured")
		}
		effectiveKMSKeyID := kmsKeyID
		if effectiveKMSKeyID == "" && bucketEncryption != nil {
			effectiveKMSKeyID = bucketEncryption.KMSMasterKeyID
		}
		if effectiveKMSKeyID == "" {
			return nil, fmt.Errorf("KMS key ID is required for SSE-KMS")
		}
		return m.sseKMSEncryptor.GenerateKey(effectiveKMSKeyID, bucket, key)
	default:
		return nil, fmt.Errorf("unsupported encryption type for GenerateKey: %s", encryptionType)
	}
}

// ParseCustomerKey parses and validates a base64-encoded customer key.
func (m *EncryptionManager) ParseCustomerKey(encodedKey, encodedMD5 string) ([]byte, error) {
	if m.sseCEncryptor == nil {
		return nil, fmt.Errorf("SSE-C is not configured")
	}
	return m.sseCEncryptor.ParseCustomerKey(encodedKey, encodedMD5)
}

// DeleteBucketKey removes the cached encryption key for a bucket.
func (m *EncryptionManager) DeleteBucketKey(bucket string) {
	m.sseS3Encryptor.DeleteBucketKey(bucket)
}

// LoadSSE3BucketKey restores a persisted SSE-S3 bucket key.
func (m *EncryptionManager) LoadSSE3BucketKey(bucket string, keyData []byte) error {
	return m.sseS3Encryptor.LoadBucketKey(bucket, keyData)
}

// ForEachSSE3BucketKey iterates over all cached SSE-S3 bucket keys.
func (m *EncryptionManager) ForEachSSE3BucketKey(fn func(bucketName string, keyData []byte)) {
	m.sseS3Encryptor.ForEachBucketKey(fn)
}
