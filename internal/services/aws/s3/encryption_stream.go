package s3

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/crypto"
)

const encryptionChunkSize = 64 * 1024 * 1024

// StreamEncryptionResult holds the result of a streaming chunked encryption.
type StreamEncryptionResult struct {
	SSEMetadata     *s3store.SSEObjectMetadata
	EncryptedData   []byte
	UnencryptedSize int64
	UnencryptedMD5  string
}

// EncryptStream reads plaintext from src in chunks, encrypts each chunk with
// AES-GCM using the appropriate key for the given encryption type, and returns
// the combined encrypted bytes along with SSE metadata that records per-chunk
// information for later streaming decryption.
func (m *EncryptionManager) EncryptStream(
	src io.Reader,
	encryptionType EncryptionType,
	bucketEncryption *s3store.EncryptionConfig,
	bucket, key, kmsKeyID string,
	customerKey []byte,
) (*StreamEncryptionResult, error) {
	genKey, err := m.resolveEncryptionKey(encryptionType, bucketEncryption, bucket, key, kmsKeyID, customerKey)
	if err != nil {
		return nil, err
	}

	var encryptedBuf bytes.Buffer
	var parts []s3store.PartEncryptionInfo
	var totalPlain int64
	h := md5.New()

	buf := make([]byte, encryptionChunkSize)
	for {
		n, readErr := io.ReadFull(src, buf)
		if n > 0 {
			chunk := buf[:n]
			totalPlain += int64(n)
			h.Write(chunk)

			nonce, nonceErr := crypto.RandomNonce()
			if nonceErr != nil {
				return nil, fmt.Errorf("failed to generate nonce: %w", nonceErr)
			}

			encChunk, encErr := crypto.AESGCMEncryptWithNonce(genKey.PlaintextKey, chunk, nonce)
			if encErr != nil {
				return nil, fmt.Errorf("failed to encrypt chunk: %w", encErr)
			}

			parts = append(parts, s3store.PartEncryptionInfo{
				EncryptedSize: int64(len(encChunk)),
				PlainSize:     int64(n),
				ContentNonce:  nonce,
				DataKey:       genKey.EncryptedDataKey,
			})

			if _, writeErr := encryptedBuf.Write(encChunk); writeErr != nil {
				return nil, fmt.Errorf("failed to write encrypted chunk: %w", writeErr)
			}
		}
		if readErr == io.EOF || readErr == io.ErrUnexpectedEOF {
			break
		}
		if readErr != nil {
			return nil, fmt.Errorf("failed to read plaintext chunk: %w", readErr)
		}
	}

	encType := s3store.SSEType(encryptionType)
	unencryptedMD5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	sseMeta := &s3store.SSEObjectMetadata{
		EncryptionType:      encType,
		EncryptedDataKey:    genKey.EncryptedDataKey,
		ContentNonce:        nil,
		KMSKeyID:            genKey.KMSKeyID,
		UnencryptedMD5:      unencryptedMD5,
		UnencryptedSize:     totalPlain,
		PartEncryptionInfos: parts,
	}

	return &StreamEncryptionResult{
		SSEMetadata:     sseMeta,
		EncryptedData:   encryptedBuf.Bytes(),
		UnencryptedSize: totalPlain,
		UnencryptedMD5:  unencryptedMD5,
	}, nil
}

// DecryptChunked decrypts combined encrypted data that was produced by chunked
// encryption (PartEncryptionInfos present). It returns the decrypted bytes.
func (m *EncryptionManager) DecryptChunked(
	encryptedData []byte,
	sseMetadata *s3store.SSEObjectMetadata,
	bucket, key string,
	customerKey []byte,
) ([]byte, error) {
	if sseMetadata == nil || len(sseMetadata.PartEncryptionInfos) == 0 {
		return nil, fmt.Errorf("missing part encryption infos for chunked decryption")
	}

	plainKey, err := m.resolveDecryptionKey(sseMetadata, bucket, key, customerKey)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer
	offset := int64(0)

	for i, part := range sseMetadata.PartEncryptionInfos {
		encSize := part.EncryptedSize
		if encSize == 0 {
			continue
		}
		if offset+encSize > int64(len(encryptedData)) {
			return nil, fmt.Errorf("encrypted data truncated at chunk %d", i)
		}

		encChunk := encryptedData[offset : offset+encSize]
		offset += encSize

		plainChunk, decErr := crypto.AESGCMDecryptWithNonce(plainKey, encChunk, part.ContentNonce)
		if decErr != nil {
			return nil, fmt.Errorf("failed to decrypt chunk %d: %w", i, decErr)
		}
		result.Write(plainChunk)
	}

	return result.Bytes(), nil
}

type resolvedKey struct {
	PlaintextKey     []byte
	EncryptedDataKey []byte
	KMSKeyID         string
}

func (m *EncryptionManager) resolveEncryptionKey(
	encryptionType EncryptionType,
	bucketEncryption *s3store.EncryptionConfig,
	bucket, key, kmsKeyID string,
	customerKey []byte,
) (*resolvedKey, error) {
	switch encryptionType {
	case EncryptionTypeSSE_S3:
		gen, err := m.sseS3Encryptor.GenerateKey(bucket)
		if err != nil {
			return nil, err
		}
		return &resolvedKey{
			PlaintextKey:     gen.PlaintextKey,
			EncryptedDataKey: gen.EncryptedDataKey,
		}, nil
	case EncryptionTypeSSE_KMS, EncryptionTypeSSE_DSSE_KMS:
		if m.sseKMSEncryptor == nil {
			return nil, fmt.Errorf("KMS not configured")
		}
		effectiveKMSKeyID := kmsKeyID
		if effectiveKMSKeyID == "" && bucketEncryption != nil {
			effectiveKMSKeyID = bucketEncryption.KMSMasterKeyID
		}
		if effectiveKMSKeyID == "" {
			return nil, fmt.Errorf("KMS key ID is required")
		}
		gen, err := m.sseKMSEncryptor.GenerateKey(effectiveKMSKeyID, bucket, key)
		if err != nil {
			return nil, err
		}
		return &resolvedKey{
			PlaintextKey:     gen.PlaintextKey,
			EncryptedDataKey: gen.EncryptedDataKey,
			KMSKeyID:         effectiveKMSKeyID,
		}, nil
	case EncryptionTypeSSE_C:
		if customerKey == nil {
			return nil, fmt.Errorf("customer key is required for SSE-C")
		}
		return &resolvedKey{PlaintextKey: customerKey}, nil
	default:
		return nil, fmt.Errorf("unsupported encryption type: %s", encryptionType)
	}
}

func (m *EncryptionManager) resolveDecryptionKey(
	sseMetadata *s3store.SSEObjectMetadata,
	bucket, key string,
	customerKey []byte,
) ([]byte, error) {
	switch sseMetadata.EncryptionType {
	case s3store.SSETypeAES256:
		bucketKey, _, err := m.sseS3Encryptor.getOrCreateBucketKey(bucket)
		if err != nil {
			return nil, err
		}
		keyMetaBytes, err := crypto.AESGCMDecrypt(bucketKey, sseMetadata.EncryptedDataKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt object key: %w", err)
		}
		var keyMeta struct {
			Key []byte `json:"key"`
		}
		if err := json.Unmarshal(keyMetaBytes, &keyMeta); err != nil {
			return nil, fmt.Errorf("failed to unmarshal key metadata: %w", err)
		}
		return keyMeta.Key, nil
	case s3store.SSETypeKMS, s3store.SSETypeDSSEKMS:
		if m.sseKMSEncryptor == nil {
			return nil, fmt.Errorf("KMS not configured")
		}
		if sseMetadata.EncryptedDataKey == nil {
			return nil, fmt.Errorf("missing encrypted data key")
		}
		encryptionContext := m.sseKMSEncryptor.buildEncryptionContext(bucket, key)
		plainKey, err := m.sseKMSEncryptor.kmsClient.Decrypt(sseMetadata.KMSKeyID, sseMetadata.EncryptedDataKey, encryptionContext)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data key: %w", err)
		}
		return plainKey, nil
	case s3store.SSETypeCustomer:
		if customerKey == nil {
			return nil, fmt.Errorf("customer key is required for SSE-C decryption")
		}
		return customerKey, nil
	default:
		return nil, fmt.Errorf("unsupported encryption type: %s", sseMetadata.EncryptionType)
	}
}
