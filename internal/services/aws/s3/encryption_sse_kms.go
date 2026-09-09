package s3

import (
	"encoding/base64"
	"fmt"

	arnutil "vorpalstacks/internal/utils/aws/arn"

	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/crypto"
)

// KMSClient defines the interface for KMS operations.
type KMSClient interface {
	GenerateDataKey(keyID string, keySpec string, context map[string]string, sourceArn string) (*GenerateDataKeyResult, error)
	Decrypt(keyID string, ciphertext []byte, context map[string]string, sourceArn string) ([]byte, error)
	KeyExists(keyID string) bool
}

// GenerateDataKeyResult holds the result of a GenerateDataKey operation.
type GenerateDataKeyResult struct {
	Plaintext  []byte
	Ciphertext []byte
}

// SSEKMSEncryptor handles S3 server-side encryption with AWS KMS-managed keys.
type SSEKMSEncryptor struct {
	kmsClient KMSClient
}

// NewSSEKMSEncryptor creates a new KMS-based encryptor.
func NewSSEKMSEncryptor(kmsClient KMSClient) *SSEKMSEncryptor {
	return &SSEKMSEncryptor{kmsClient: kmsClient}
}

func (e *SSEKMSEncryptor) buildEncryptionContext(bucket, key string) map[string]string {
	return map[string]string{
		"aws:s3:arn": arnutil.NewARNBuilder("", "").S3().Object(bucket, key),
	}
}

// Encrypt encrypts data using AWS KMS-managed keys.
func (e *SSEKMSEncryptor) Encrypt(plaintext []byte, bucket, key string, kmsKeyID string) (*EncryptionResult, error) {
	if e.kmsClient == nil {
		return nil, fmt.Errorf("KMS client not configured")
	}

	if kmsKeyID == "" {
		return nil, fmt.Errorf("KMS key ID is required for SSE-KMS")
	}

	encryptionContext := e.buildEncryptionContext(bucket, key)
	bucketArn := arnutil.NewARNBuilder("", "").S3().Bucket(bucket)

	dataKeyResult, err := e.kmsClient.GenerateDataKey(kmsKeyID, "AES_256", encryptionContext, bucketArn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate data key: %w", err)
	}

	nonce, err := crypto.RandomNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	encryptedData, err := crypto.AESGCMEncryptWithNonce(dataKeyResult.Plaintext, plaintext, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	unencryptedMD5 := base64.StdEncoding.EncodeToString(crypto.MD5Hash(plaintext))

	return &EncryptionResult{
		EncryptedData:    encryptedData,
		EncryptionType:   EncryptionTypeSSE_KMS,
		EncryptedDataKey: dataKeyResult.Ciphertext,
		ContentNonce:     nonce,
		KMSKeyID:         kmsKeyID,
		UnencryptedMD5:   unencryptedMD5,
		UnencryptedSize:  int64(len(plaintext)),
	}, nil
}

// EncryptWithPlaintextKey encrypts data using a provided plaintext key and KMS key ID.
func (e *SSEKMSEncryptor) EncryptWithPlaintextKey(plaintext []byte, plaintextKey []byte, kmsKeyID string, encryptedDataKey []byte) (*EncryptionResult, error) {
	nonce, err := crypto.RandomNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	encryptedData, err := crypto.AESGCMEncryptWithNonce(plaintextKey, plaintext, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	unencryptedMD5 := base64.StdEncoding.EncodeToString(crypto.MD5Hash(plaintext))

	return &EncryptionResult{
		EncryptedData:    encryptedData,
		EncryptionType:   EncryptionTypeSSE_KMS,
		EncryptedDataKey: encryptedDataKey,
		ContentNonce:     nonce,
		KMSKeyID:         kmsKeyID,
		UnencryptedMD5:   unencryptedMD5,
		UnencryptedSize:  int64(len(plaintext)),
	}, nil
}

// GenerateKey generates an encryption key using AWS KMS.
func (e *SSEKMSEncryptor) GenerateKey(kmsKeyID string, bucket, key string) (*GeneratedKey, error) {
	if e.kmsClient == nil {
		return nil, fmt.Errorf("KMS client not configured")
	}

	encryptionContext := e.buildEncryptionContext(bucket, key)
	bucketArn := arnutil.NewARNBuilder("", "").S3().Bucket(bucket)

	dataKeyResult, err := e.kmsClient.GenerateDataKey(kmsKeyID, "AES_256", encryptionContext, bucketArn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate data key: %w", err)
	}

	return &GeneratedKey{
		PlaintextKey:     dataKeyResult.Plaintext,
		EncryptedDataKey: dataKeyResult.Ciphertext,
		KMSKeyID:         kmsKeyID,
	}, nil
}

// Decrypt decrypts data using AWS KMS-managed keys.
func (e *SSEKMSEncryptor) Decrypt(ciphertext []byte, sseMetadata *s3store.SSEObjectMetadata, bucket, key string) (*DecryptionResult, error) {
	if e.kmsClient == nil {
		return nil, fmt.Errorf("KMS client not configured")
	}

	if sseMetadata == nil {
		return nil, fmt.Errorf("missing SSE metadata for encrypted object")
	}

	if sseMetadata.EncryptedDataKey == nil {
		return nil, fmt.Errorf("missing encrypted data key")
	}

	encryptionContext := e.buildEncryptionContext(bucket, key)
	bucketArn := arnutil.NewARNBuilder("", "").S3().Bucket(bucket)

	plaintextKey, err := e.kmsClient.Decrypt(sseMetadata.KMSKeyID, sseMetadata.EncryptedDataKey, encryptionContext, bucketArn)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data key: %w", err)
	}

	if sseMetadata.ContentNonce == nil {
		return nil, fmt.Errorf("missing content nonce")
	}

	plaintext, err := crypto.AESGCMDecryptWithNonce(plaintextKey, ciphertext, sseMetadata.ContentNonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return &DecryptionResult{DecryptedData: plaintext}, nil
}
