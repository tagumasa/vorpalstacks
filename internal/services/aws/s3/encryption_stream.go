package s3

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"

	s3store "vorpalstacks/internal/store/aws/s3"
	arnutil "vorpalstacks/internal/utils/aws/arn"
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

// EncryptStream encrypts plaintext in chunks and returns the combined
// encrypted bytes. Callers that persist through a store reader should
// prefer NewChunkEncryptReader, which avoids buffering the whole object.
func (m *EncryptionManager) EncryptStream(
	src io.Reader,
	encryptionType EncryptionType,
	bucketEncryption *s3store.EncryptionConfig,
	bucket, key, kmsKeyID string,
	customerKey []byte,
) (*StreamEncryptionResult, error) {
	sseMeta := &s3store.SSEObjectMetadata{}
	reader, err := m.NewChunkEncryptReader(src, encryptionType, bucketEncryption, bucket, key, kmsKeyID, customerKey, sseMeta)
	if err != nil {
		return nil, err
	}
	encrypted, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return &StreamEncryptionResult{
		SSEMetadata:     sseMeta,
		EncryptedData:   encrypted,
		UnencryptedSize: sseMeta.UnencryptedSize,
		UnencryptedMD5:  sseMeta.UnencryptedMD5,
	}, nil
}

// chunkEncryptReader encrypts plaintext chunks on-the-fly as they are
// read — the write-side mirror of chunkDecryptReader. The SSE metadata
// accumulates as chunks pass through: sseMetadata holds the encryption
// type and key material from construction, and the digest, unencrypted
// size and part table are final only once the reader has returned io.EOF.
type chunkEncryptReader struct {
	src              io.Reader
	plainKey         []byte
	encryptedDataKey []byte
	h                hash.Hash
	sseMeta          *s3store.SSEObjectMetadata
	pending          []byte
	pendingOff       int
	done             bool
	buf              []byte
}

// NewChunkEncryptReader returns a reader that encrypts src chunk-by-chunk
// with AES-GCM using the appropriate key for the given encryption type.
// The caller-provided sseMetadata is populated as the stream is consumed;
// reading it before the reader has fully drained yields partial metadata.
func (m *EncryptionManager) NewChunkEncryptReader(
	src io.Reader,
	encryptionType EncryptionType,
	bucketEncryption *s3store.EncryptionConfig,
	bucket, key, kmsKeyID string,
	customerKey []byte,
	sseMetadata *s3store.SSEObjectMetadata,
) (io.Reader, error) {
	genKey, err := m.resolveEncryptionKey(encryptionType, bucketEncryption, bucket, key, kmsKeyID, customerKey)
	if err != nil {
		return nil, err
	}
	sseMetadata.EncryptionType = s3store.SSEType(encryptionType)
	sseMetadata.EncryptedDataKey = genKey.EncryptedDataKey
	sseMetadata.KMSKeyID = genKey.KMSKeyID
	return &chunkEncryptReader{
		src:              src,
		plainKey:         genKey.PlaintextKey,
		encryptedDataKey: genKey.EncryptedDataKey,
		h:                md5.New(),
		sseMeta:          sseMetadata,
		buf:              make([]byte, encryptionChunkSize),
	}, nil
}

func (r *chunkEncryptReader) Read(p []byte) (int, error) {
	if r.pendingOff < len(r.pending) {
		n := copy(p, r.pending[r.pendingOff:])
		r.pendingOff += n
		return n, nil
	}
	if r.done {
		return 0, io.EOF
	}

	n, readErr := io.ReadFull(r.src, r.buf)
	if n > 0 {
		chunk := r.buf[:n]
		r.h.Write(chunk)
		r.sseMeta.UnencryptedSize += int64(n)

		nonce, nonceErr := crypto.RandomNonce()
		if nonceErr != nil {
			return 0, fmt.Errorf("failed to generate nonce: %w", nonceErr)
		}
		encChunk, encErr := crypto.AESGCMEncryptWithNonce(r.plainKey, chunk, nonce)
		if encErr != nil {
			return 0, fmt.Errorf("failed to encrypt chunk: %w", encErr)
		}
		r.sseMeta.PartEncryptionInfos = append(r.sseMeta.PartEncryptionInfos, s3store.PartEncryptionInfo{
			EncryptedSize: int64(len(encChunk)),
			PlainSize:     int64(n),
			ContentNonce:  nonce,
			DataKey:       r.encryptedDataKey,
		})
		r.pending = encChunk
		r.pendingOff = 0
	}
	if readErr == io.EOF || readErr == io.ErrUnexpectedEOF {
		r.done = true
		r.sseMeta.UnencryptedMD5 = base64.StdEncoding.EncodeToString(r.h.Sum(nil))
		if n > 0 {
			c := copy(p, r.pending)
			r.pendingOff = c
			return c, nil
		}
		return 0, io.EOF
	}
	if readErr != nil {
		return 0, fmt.Errorf("failed to read plaintext chunk: %w", readErr)
	}

	c := copy(p, r.pending)
	r.pendingOff = c
	return c, nil
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

// DecryptChunkRange decrypts the plaintext window [offset, offset+length)
// of a chunk-encrypted object without materialising the whole object: the
// per-chunk plain sizes locate the chunks the window overlaps, fetchRange
// retrieves only their encrypted bytes, and the window is sliced from the
// concatenated plaintext of those chunks.
func (m *EncryptionManager) DecryptChunkRange(
	sseMetadata *s3store.SSEObjectMetadata,
	bucket, key string,
	customerKey []byte,
	offset, length int64,
	fetchRange func(encOffset, encLength int64) ([]byte, error),
) ([]byte, error) {
	if sseMetadata == nil || len(sseMetadata.PartEncryptionInfos) == 0 {
		return nil, fmt.Errorf("missing part encryption infos for ranged decryption")
	}

	plainKey, err := m.resolveDecryptionKey(sseMetadata, bucket, key, customerKey)
	if err != nil {
		return nil, err
	}

	type chunkSpan struct {
		part       s3store.PartEncryptionInfo
		encOffset  int64
		plainStart int64
	}
	var spans []chunkSpan
	var encCursor, plainCursor int64
	for _, part := range sseMetadata.PartEncryptionInfos {
		if part.EncryptedSize > 0 && plainCursor+part.PlainSize > offset && plainCursor < offset+length {
			spans = append(spans, chunkSpan{part: part, encOffset: encCursor, plainStart: plainCursor})
		}
		encCursor += part.EncryptedSize
		plainCursor += part.PlainSize
	}
	if len(spans) == 0 {
		return nil, fmt.Errorf("range window [%d,%d) does not overlap any encrypted chunk", offset, offset+length)
	}

	fetchStart := spans[0].encOffset
	fetchEnd := spans[len(spans)-1].encOffset + spans[len(spans)-1].part.EncryptedSize
	encrypted, err := fetchRange(fetchStart, fetchEnd-fetchStart)
	if err != nil {
		return nil, err
	}

	var plain bytes.Buffer
	for _, sp := range spans {
		rel := sp.encOffset - fetchStart
		plainChunk, decErr := crypto.AESGCMDecryptWithNonce(plainKey, encrypted[rel:rel+sp.part.EncryptedSize], sp.part.ContentNonce)
		if decErr != nil {
			return nil, fmt.Errorf("failed to decrypt chunk at encrypted offset %d: %w", sp.encOffset, decErr)
		}
		plain.Write(plainChunk)
	}

	start := offset - spans[0].plainStart
	return plain.Bytes()[start : start+length], nil
}

// chunkDecryptReader decrypts chunked encrypted data on-the-fly as it is
// read from the underlying source, avoiding loading the entire encrypted
// object into memory. Each chunk is independently encrypted with AES-GCM,
// so chunks are decrypted one at a time using per-chunk nonce and the
// shared plaintext key.
type chunkDecryptReader struct {
	source     io.Reader
	parts      []s3store.PartEncryptionInfo
	plainKey   []byte
	chunkIdx   int
	pending    []byte
	pendingOff int
	closed     bool
}

// NewChunkDecryptReader creates a streaming reader that decrypts chunked
// encrypted data on-the-fly. The plaintext key is resolved once and reused
// for all chunks, avoiding redundant KMS or bucket-key operations per chunk.
func (m *EncryptionManager) NewChunkDecryptReader(
	source io.Reader,
	sseMetadata *s3store.SSEObjectMetadata,
	bucket, key string,
	customerKey []byte,
) (io.ReadCloser, error) {
	if sseMetadata == nil || len(sseMetadata.PartEncryptionInfos) == 0 {
		return nil, fmt.Errorf("missing part encryption infos for chunked decryption")
	}

	plainKey, err := m.resolveDecryptionKey(sseMetadata, bucket, key, customerKey)
	if err != nil {
		return nil, err
	}

	return &chunkDecryptReader{
		source:   source,
		parts:    sseMetadata.PartEncryptionInfos,
		plainKey: plainKey,
	}, nil
}

func (r *chunkDecryptReader) Read(p []byte) (int, error) {
	if r.closed {
		return 0, io.ErrClosedPipe
	}

	if r.pendingOff < len(r.pending) {
		n := copy(p, r.pending[r.pendingOff:])
		r.pendingOff += n
		return n, nil
	}

	if r.chunkIdx >= len(r.parts) {
		return 0, io.EOF
	}

	part := r.parts[r.chunkIdx]
	r.chunkIdx++

	if part.EncryptedSize == 0 {
		return r.Read(p)
	}

	encChunk := make([]byte, part.EncryptedSize)
	if _, err := io.ReadFull(r.source, encChunk); err != nil {
		return 0, fmt.Errorf("failed to read encrypted chunk %d: %w", r.chunkIdx-1, err)
	}

	plainChunk, err := crypto.AESGCMDecryptWithNonce(r.plainKey, encChunk, part.ContentNonce)
	if err != nil {
		return 0, fmt.Errorf("failed to decrypt chunk %d: %w", r.chunkIdx-1, err)
	}

	r.pending = plainChunk
	r.pendingOff = 0

	n := copy(p, r.pending)
	r.pendingOff = n
	return n, nil
}

func (r *chunkDecryptReader) Close() error {
	r.closed = true
	if c, ok := r.source.(io.Closer); ok {
		return c.Close()
	}
	return nil
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
		bucketArn := arnutil.NewARNBuilder("", "").S3().Bucket(bucket)
		plainKey, err := m.sseKMSEncryptor.kmsClient.Decrypt(sseMetadata.KMSKeyID, sseMetadata.EncryptedDataKey, encryptionContext, bucketArn)
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
