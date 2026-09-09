package s3

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"io"
	"testing"

	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/crypto"
)

func randomBytes(t *testing.T, n int) []byte {
	t.Helper()
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		t.Fatalf("failed to read random bytes: %v", err)
	}
	return buf
}

// TestChunkEncryptReaderRoundTrip pins the streaming write path: the reader
// encrypts on the fly, finalises the SSE metadata at EOF, and the streaming
// decrypt reader recovers the exact plaintext.
func TestChunkEncryptReaderRoundTrip(t *testing.T) {
	plaintext := randomBytes(t, 200*1024)
	manager := NewEncryptionManager()

	sseMeta := &s3store.SSEObjectMetadata{}
	encReader, err := manager.NewChunkEncryptReader(bytes.NewReader(plaintext), EncryptionTypeSSE_S3, nil, "bucket", "key", "", nil, sseMeta)
	if err != nil {
		t.Fatalf("NewChunkEncryptReader failed: %v", err)
	}
	encrypted, err := io.ReadAll(encReader)
	if err != nil {
		t.Fatalf("reading encrypted stream failed: %v", err)
	}

	if sseMeta.UnencryptedSize != int64(len(plaintext)) {
		t.Errorf("UnencryptedSize = %d, want %d", sseMeta.UnencryptedSize, len(plaintext))
	}
	sum := md5.Sum(plaintext)
	wantMD5 := base64.StdEncoding.EncodeToString(sum[:])
	if sseMeta.UnencryptedMD5 != wantMD5 {
		t.Errorf("UnencryptedMD5 = %s, want %s", sseMeta.UnencryptedMD5, wantMD5)
	}
	if len(sseMeta.PartEncryptionInfos) != 1 {
		t.Errorf("part count = %d, want 1 for a single-chunk object", len(sseMeta.PartEncryptionInfos))
	}

	decReader, err := manager.NewChunkDecryptReader(bytes.NewReader(encrypted), sseMeta, "bucket", "key", nil)
	if err != nil {
		t.Fatalf("NewChunkDecryptReader failed: %v", err)
	}
	defer decReader.Close()
	decrypted, err := io.ReadAll(decReader)
	if err != nil {
		t.Fatalf("reading decrypted stream failed: %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("round-trip mismatch: decrypted %d bytes, plaintext %d bytes", len(decrypted), len(plaintext))
	}
}

// TestChunkEncryptReaderEmptyObject pins the zero-byte case: no parts, the
// digest of empty input, and immediate EOF.
func TestChunkEncryptReaderEmptyObject(t *testing.T) {
	manager := NewEncryptionManager()
	sseMeta := &s3store.SSEObjectMetadata{}
	encReader, err := manager.NewChunkEncryptReader(bytes.NewReader(nil), EncryptionTypeSSE_S3, nil, "bucket", "key", "", nil, sseMeta)
	if err != nil {
		t.Fatalf("NewChunkEncryptReader failed: %v", err)
	}
	encrypted, err := io.ReadAll(encReader)
	if err != nil {
		t.Fatalf("reading encrypted stream failed: %v", err)
	}
	if len(encrypted) != 0 {
		t.Errorf("encrypted length = %d, want 0", len(encrypted))
	}
	if len(sseMeta.PartEncryptionInfos) != 0 {
		t.Errorf("part count = %d, want 0", len(sseMeta.PartEncryptionInfos))
	}
	if sseMeta.UnencryptedSize != 0 {
		t.Errorf("UnencryptedSize = %d, want 0", sseMeta.UnencryptedSize)
	}
	sum := md5.Sum(nil)
	if sseMeta.UnencryptedMD5 != base64.StdEncoding.EncodeToString(sum[:]) {
		t.Errorf("UnencryptedMD5 = %s, want digest of empty input", sseMeta.UnencryptedMD5)
	}
}

// TestDecryptChunkRangeWindow pins the ranged-decryption window arithmetic
// on a hand-built three-chunk layout: only the overlapping chunks are
// fetched, and the returned window is the exact plaintext sub-slice.
func TestDecryptChunkRangeWindow(t *testing.T) {
	customerKey := randomBytes(t, 32)
	chunks := [][]byte{
		randomBytes(t, 10),
		randomBytes(t, 20),
		randomBytes(t, 15),
	}
	var plaintext, ciphertext []byte
	infos := make([]s3store.PartEncryptionInfo, 0, len(chunks))
	for _, chunk := range chunks {
		nonce, err := crypto.RandomNonce()
		if err != nil {
			t.Fatalf("RandomNonce failed: %v", err)
		}
		enc, err := crypto.AESGCMEncryptWithNonce(customerKey, chunk, nonce)
		if err != nil {
			t.Fatalf("encrypt chunk failed: %v", err)
		}
		infos = append(infos, s3store.PartEncryptionInfo{
			EncryptedSize: int64(len(enc)),
			PlainSize:     int64(len(chunk)),
			ContentNonce:  nonce,
		})
		plaintext = append(plaintext, chunk...)
		ciphertext = append(ciphertext, enc...)
	}
	sseMeta := &s3store.SSEObjectMetadata{
		EncryptionType:      s3store.SSETypeCustomer,
		PartEncryptionInfos: infos,
	}

	manager := NewEncryptionManager()
	fetches := 0
	fetchRange := func(encOffset, encLength int64) ([]byte, error) {
		fetches++
		return ciphertext[encOffset : encOffset+encLength], nil
	}

	window := func(offset, length int64) []byte {
		t.Helper()
		got, err := manager.DecryptChunkRange(sseMeta, "bucket", "key", customerKey, offset, length, fetchRange)
		if err != nil {
			t.Fatalf("DecryptChunkRange(%d, %d) failed: %v", offset, length, err)
		}
		return got
	}

	if got := window(0, 5); !bytes.Equal(got, plaintext[0:5]) {
		t.Errorf("window [0,5) = %v, want %v", got, plaintext[0:5])
	}
	if got := window(12, 13); !bytes.Equal(got, plaintext[12:25]) {
		t.Errorf("window [12,25) spanning chunks 1-2 mismatch")
	}
	if got := window(40, 5); !bytes.Equal(got, plaintext[40:45]) {
		t.Errorf("tail window [40,45) mismatch")
	}
	if fetches != 3 {
		t.Errorf("fetch count = %d, want one fetch per window", fetches)
	}

	if _, err := manager.DecryptChunkRange(sseMeta, "bucket", "key", customerKey, int64(len(plaintext)), 1, fetchRange); err == nil {
		t.Errorf("window past the last byte unexpectedly succeeded")
	}
}
