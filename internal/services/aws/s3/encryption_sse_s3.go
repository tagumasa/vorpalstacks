// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/crypto"
)

// SSES3Encryptor handles S3 server-side encryption with S3-managed keys.
type SSES3Encryptor struct {
	bucketKeys sync.Map
}

// NewSSES3Encryptor creates a new S3-managed key encryptor.
func NewSSES3Encryptor() *SSES3Encryptor {
	return &SSES3Encryptor{}
}

type sseS3KeyMetadata struct {
	Key       []byte `json:"key"`
	KeyID     string `json:"key_id"`
	CreatedAt int64  `json:"created_at"`
}

// GetEncryptionType returns the encryption type.
func (e *SSES3Encryptor) GetEncryptionType() EncryptionType {
	return EncryptionTypeSSE_S3
}

// DeleteBucketKey removes the cached S3-managed key for a bucket.
func (e *SSES3Encryptor) DeleteBucketKey(bucket string) {
	e.bucketKeys.Delete(bucket)
}

func (e *SSES3Encryptor) getOrCreateBucketKey(bucket string) ([]byte, string, error) {
	if val, ok := e.bucketKeys.Load(bucket); ok {
		meta := val.(*sseS3KeyMetadata)
		return meta.Key, meta.KeyID, nil
	}

	key, err := crypto.RandomKey(32)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate bucket key: %w", err)
	}

	keyID := base64.StdEncoding.EncodeToString(key[:8])
	meta := &sseS3KeyMetadata{
		Key:       key,
		KeyID:     keyID,
		CreatedAt: 0,
	}

	actual, _ := e.bucketKeys.LoadOrStore(bucket, meta)
	return actual.(*sseS3KeyMetadata).Key, actual.(*sseS3KeyMetadata).KeyID, nil
}

// Encrypt encrypts data using S3-managed keys.
func (e *SSES3Encryptor) Encrypt(plaintext []byte, bucket, key string) (*EncryptionResult, error) {
	bucketKey, keyID, err := e.getOrCreateBucketKey(bucket)
	if err != nil {
		return nil, err
	}

	objectKey, err := crypto.RandomKey(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate object key: %w", err)
	}

	encryptedData, err := crypto.AESGCMEncrypt(objectKey, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	keyMeta := &sseS3KeyMetadata{
		Key:   objectKey,
		KeyID: keyID,
	}
	keyMetaBytes, err := json.Marshal(keyMeta)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key metadata: %w", err)
	}

	encryptedKey, err := crypto.AESGCMEncrypt(bucketKey, keyMetaBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt object key: %w", err)
	}

	unencryptedMD5 := base64.StdEncoding.EncodeToString(crypto.MD5Hash(plaintext))

	return &EncryptionResult{
		EncryptedData:    encryptedData,
		EncryptionType:   EncryptionTypeSSE_S3,
		EncryptedDataKey: encryptedKey,
		ContentNonce:     nil,
		UnencryptedMD5:   unencryptedMD5,
		UnencryptedSize:  int64(len(plaintext)),
	}, nil
}

// EncryptWithPlaintextKey encrypts data using a provided plaintext key.
func (e *SSES3Encryptor) EncryptWithPlaintextKey(plaintext []byte, bucket string, plaintextKey []byte) (*EncryptionResult, error) {
	encryptedData, err := crypto.AESGCMEncrypt(plaintextKey, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	bucketKey, _, err := e.getOrCreateBucketKey(bucket)
	if err != nil {
		return nil, err
	}

	keyMeta := &sseS3KeyMetadata{
		Key: plaintextKey,
	}
	keyMetaBytes, err := json.Marshal(keyMeta)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key metadata: %w", err)
	}

	encryptedKey, err := crypto.AESGCMEncrypt(bucketKey, keyMetaBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt object key: %w", err)
	}

	unencryptedMD5 := base64.StdEncoding.EncodeToString(crypto.MD5Hash(plaintext))

	return &EncryptionResult{
		EncryptedData:    encryptedData,
		EncryptionType:   EncryptionTypeSSE_S3,
		EncryptedDataKey: encryptedKey,
		ContentNonce:     nil,
		UnencryptedMD5:   unencryptedMD5,
		UnencryptedSize:  int64(len(plaintext)),
	}, nil
}

// GenerateKey generates an encryption key for S3-managed encryption.
func (e *SSES3Encryptor) GenerateKey(bucket string) (*GeneratedKey, error) {
	objectKey, err := crypto.RandomKey(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate object key: %w", err)
	}

	bucketKey, _, err := e.getOrCreateBucketKey(bucket)
	if err != nil {
		return nil, err
	}

	keyMeta := &sseS3KeyMetadata{
		Key: objectKey,
	}
	keyMetaBytes, err := json.Marshal(keyMeta)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key metadata: %w", err)
	}

	encryptedKey, err := crypto.AESGCMEncrypt(bucketKey, keyMetaBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt object key: %w", err)
	}

	return &GeneratedKey{
		PlaintextKey:     objectKey,
		EncryptedDataKey: encryptedKey,
	}, nil
}

// Decrypt decrypts data using S3-managed keys.
func (e *SSES3Encryptor) Decrypt(ciphertext []byte, sseMetadata *s3store.SSEObjectMetadata, bucket, key string) (*DecryptionResult, error) {
	if sseMetadata == nil {
		return nil, fmt.Errorf("missing SSE metadata for encrypted object")
	}

	bucketKey, _, err := e.getOrCreateBucketKey(bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket key: %w", err)
	}

	keyMetaBytes, err := crypto.AESGCMDecrypt(bucketKey, sseMetadata.EncryptedDataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt object key: %w", err)
	}

	var keyMeta sseS3KeyMetadata
	if err := json.Unmarshal(keyMetaBytes, &keyMeta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal key metadata: %w", err)
	}

	plaintext, err := crypto.AESGCMDecrypt(keyMeta.Key, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return &DecryptionResult{DecryptedData: plaintext}, nil
}

// LoadBucketKey loads a bucket key from serialised data.
func (e *SSES3Encryptor) LoadBucketKey(bucket string, keyData []byte) error {
	var meta sseS3KeyMetadata
	if err := json.Unmarshal(keyData, &meta); err != nil {
		return fmt.Errorf("failed to unmarshal bucket key: %w", err)
	}
	e.bucketKeys.Store(bucket, &meta)
	return nil
}

// ExportBucketKey exports a bucket key to serialised data.
func (e *SSES3Encryptor) ExportBucketKey(bucket string) ([]byte, error) {
	val, ok := e.bucketKeys.Load(bucket)
	if !ok {
		return nil, fmt.Errorf("bucket key not found: %s", bucket)
	}
	meta := val.(*sseS3KeyMetadata)
	return json.Marshal(meta)
}
