package hsm

// Package hsm provides Hardware Security Module (HSM) in-memory implementations
// for vorpalstacks KMS operations.

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"sync"
)

type memoryKey struct {
	keySpec   KeySpec
	symmetric []byte
	rsaKey    *rsa.PrivateKey
	ecdsaKey  *ecdsa.PrivateKey
}

// MemoryBackend is an in-memory implementation of a Hardware Security Module (HSM)
// backend for key storage and cryptographic operations.
type MemoryBackend struct {
	mu          sync.RWMutex
	keys        map[string]*memoryKey
	wrappingKey []byte
}

// NewMemoryBackend creates a new in-memory HSM backend for key storage.
// This is useful for development and testing purposes.
func NewMemoryBackend() (*MemoryBackend, error) {
	wrappingKey := make([]byte, 32)
	if _, err := rand.Read(wrappingKey); err != nil {
		return nil, fmt.Errorf("crypto/rand failed: %w", err)
	}
	return &MemoryBackend{
		keys:        make(map[string]*memoryKey),
		wrappingKey: wrappingKey,
	}, nil
}

// GenerateKey generates a new key with the specified key ID and key specification.
// Supports symmetric, HMAC, RSA, and ECDSA key types.
func (b *MemoryBackend) GenerateKey(keyID string, keySpec KeySpec) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.keys[keyID]; exists {
		return ErrKeyAlreadyExists
	}

	key := &memoryKey{keySpec: keySpec}

	switch keySpec {
	case KeySpecSymmetricDefault:
		key.symmetric = make([]byte, 32)
		if _, err := rand.Read(key.symmetric); err != nil {
			return err
		}
	case KeySpecHMAC224:
		key.symmetric = make([]byte, 28)
		if _, err := rand.Read(key.symmetric); err != nil {
			return err
		}
	case KeySpecHMAC256:
		key.symmetric = make([]byte, 32)
		if _, err := rand.Read(key.symmetric); err != nil {
			return err
		}
	case KeySpecHMAC384:
		key.symmetric = make([]byte, 48)
		if _, err := rand.Read(key.symmetric); err != nil {
			return err
		}
	case KeySpecHMAC512:
		key.symmetric = make([]byte, 64)
		if _, err := rand.Read(key.symmetric); err != nil {
			return err
		}
	case KeySpecRSA2048:
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return err
		}
		key.rsaKey = rsaKey
	case KeySpecRSA3072:
		rsaKey, err := rsa.GenerateKey(rand.Reader, 3072)
		if err != nil {
			return err
		}
		key.rsaKey = rsaKey
	case KeySpecRSA4096:
		rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return err
		}
		key.rsaKey = rsaKey
	case KeySpecECCNISTP256, KeySpecECCSECPP256R1:
		ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
		key.ecdsaKey = ecdsaKey
	case KeySpecECCNISTP384:
		ecdsaKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return err
		}
		key.ecdsaKey = ecdsaKey
	case KeySpecECCNISTP521:
		ecdsaKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return err
		}
		key.ecdsaKey = ecdsaKey
	default:
		return ErrInvalidKeySpec
	}

	b.keys[keyID] = key
	return nil
}

// ImportKey imports an existing key into the in-memory HSM backend.
// Only supports symmetric and HMAC key specifications.
func (b *MemoryBackend) ImportKey(keyID string, keyMaterial []byte, keySpec KeySpec) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.keys[keyID]; exists {
		return ErrKeyAlreadyExists
	}

	key := &memoryKey{keySpec: keySpec}
	switch keySpec {
	case KeySpecSymmetricDefault, KeySpecHMAC256:
		key.symmetric = make([]byte, len(keyMaterial))
		copy(key.symmetric, keyMaterial)
	default:
		return ErrInvalidKeySpec
	}

	b.keys[keyID] = key
	return nil
}

// DeleteKey removes a key from the in-memory HSM backend.
// Returns an error if the key does not exist.
func (b *MemoryBackend) DeleteKey(keyID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.keys[keyID]; !exists {
		return ErrKeyNotFound
	}

	delete(b.keys, keyID)
	return nil
}

// Encrypt encrypts plaintext using the specified symmetric key.
// Returns the ciphertext with key ID prepended and CRC32 checksum.
func (b *MemoryBackend) Encrypt(keyID string, plaintext []byte, context map[string]string) (*EncryptResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, ErrEncryptFailed
	}

	ciphertext, err := aesEncrypt(key.symmetric, plaintext, serializeEncryptionContext(context))
	if err != nil {
		return nil, err
	}

	ciphertextWithKeyID := prependKeyIDToCiphertext(keyID, ciphertext)

	return &EncryptResult{
		Ciphertext:      ciphertextWithKeyID,
		CiphertextCRC32: computeCRC32(ciphertextWithKeyID),
	}, nil
}

// Decrypt decrypts ciphertext using the specified symmetric key.
// Extracts the key ID from the ciphertext if present.
func (b *MemoryBackend) Decrypt(keyID string, ciphertext []byte, context map[string]string) (*DecryptResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	extractedKeyID, actualCiphertext, err := extractKeyIDFromCiphertext(ciphertext)
	if err == nil && extractedKeyID != "" {
		keyID = extractedKeyID
		ciphertext = actualCiphertext
	}

	key, exists := b.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, ErrDecryptFailed
	}

	plaintext, err := aesDecrypt(key.symmetric, ciphertext, serializeEncryptionContext(context))
	if err != nil {
		return nil, err
	}

	return &DecryptResult{
		Plaintext:      plaintext,
		PlaintextCRC32: computeCRC32(plaintext),
	}, nil
}

// DecryptWithoutKeyID decrypts ciphertext and extracts the key ID from the ciphertext.
// Returns the plaintext, key ID used, and any error encountered.
func (b *MemoryBackend) DecryptWithoutKeyID(ciphertext []byte, context map[string]string) (*DecryptResult, string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	keyID, actualCiphertext, err := extractKeyIDFromCiphertext(ciphertext)
	if err != nil {
		return nil, "", err
	}

	key, exists := b.keys[keyID]
	if !exists {
		return nil, "", ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, "", ErrDecryptFailed
	}

	plaintext, err := aesDecrypt(key.symmetric, actualCiphertext, serializeEncryptionContext(context))
	if err != nil {
		return nil, "", err
	}

	return &DecryptResult{
		Plaintext:      plaintext,
		PlaintextCRC32: computeCRC32(plaintext),
	}, keyID, nil
}

func extractKeyIDFromCiphertext(ciphertext []byte) (string, []byte, error) {
	if len(ciphertext) < 2 {
		return "", nil, errors.New("ciphertext too short")
	}
	keyIDLen := int(ciphertext[0])<<8 | int(ciphertext[1])
	if len(ciphertext) < 2+keyIDLen {
		return "", nil, errors.New("ciphertext too short for keyID")
	}
	keyID := string(ciphertext[2 : 2+keyIDLen])
	actualCiphertext := ciphertext[2+keyIDLen:]
	return keyID, actualCiphertext, nil
}

func prependKeyIDToCiphertext(keyID string, ciphertext []byte) []byte {
	keyIDBytes := []byte(keyID)
	result := make([]byte, 2+len(keyIDBytes)+len(ciphertext))
	result[0] = byte(len(keyIDBytes) >> 8)
	result[1] = byte(len(keyIDBytes))
	copy(result[2:], keyIDBytes)
	copy(result[2+len(keyIDBytes):], ciphertext)
	return result
}

// Sign signs a message using the specified key and signing algorithm.
// Supports RSA and ECDSA signing algorithms.
func (b *MemoryBackend) Sign(keyID string, message []byte, algorithm SigningAlgorithm) (*SignResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	var signature []byte
	var err error

	if key.rsaKey != nil {
		hash := hashForAlgorithm(algorithm, message)
		switch algorithm {
		case SigningAlgorithmRSAPKCS1SHA256, SigningAlgorithmRSAPKCS1SHA384, SigningAlgorithmRSAPKCS1SHA512:
			signature, err = rsa.SignPKCS1v15(rand.Reader, key.rsaKey, hashAlgorithm(algorithm), hash)
		case SigningAlgorithmRSAPSSSHA256, SigningAlgorithmRSAPSSSHA384, SigningAlgorithmRSAPSSSHA512:
			signature, err = rsa.SignPSS(rand.Reader, key.rsaKey, hashAlgorithm(algorithm), hash, nil)
		default:
			return nil, ErrInvalidAlgorithm
		}
	} else if key.ecdsaKey != nil {
		hash := hashForAlgorithm(algorithm, message)
		var r, s *big.Int
		r, s, err = ecdsa.Sign(rand.Reader, key.ecdsaKey, hash)
		if err != nil {
			return nil, ErrSignFailed
		}
		signature, err = asn1.Marshal(struct {
			R, S *big.Int
		}{r, s})
	} else {
		return nil, ErrInvalidKeySpec
	}

	if err != nil {
		return nil, ErrSignFailed
	}

	return &SignResult{
		Signature:      signature,
		SignatureCRC32: computeCRC32(signature),
	}, nil
}

// Verify verifies a digital signature for a message using the specified key.
// Returns true if the signature is valid, false otherwise.
func (b *MemoryBackend) Verify(keyID string, message, signature []byte, algorithm SigningAlgorithm) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return false, ErrKeyNotFound
	}

	if key.rsaKey != nil {
		hash := hashForAlgorithm(algorithm, message)
		switch algorithm {
		case SigningAlgorithmRSAPKCS1SHA256, SigningAlgorithmRSAPKCS1SHA384, SigningAlgorithmRSAPKCS1SHA512:
			err := rsa.VerifyPKCS1v15(&key.rsaKey.PublicKey, hashAlgorithm(algorithm), hash, signature)
			return err == nil, nil
		case SigningAlgorithmRSAPSSSHA256, SigningAlgorithmRSAPSSSHA384, SigningAlgorithmRSAPSSSHA512:
			err := rsa.VerifyPSS(&key.rsaKey.PublicKey, hashAlgorithm(algorithm), hash, signature, nil)
			return err == nil, nil
		default:
			return false, ErrInvalidAlgorithm
		}
	} else if key.ecdsaKey != nil {
		hash := hashForAlgorithm(algorithm, message)
		var sig struct {
			R, S *big.Int
		}
		if _, err := asn1.Unmarshal(signature, &sig); err != nil {
			return false, err
		}
		return ecdsa.Verify(&key.ecdsaKey.PublicKey, hash, sig.R, sig.S), nil
	}

	return false, ErrInvalidKeySpec
}

// GenerateMAC generates a Message Authentication Code (MAC) for a message
// using the specified symmetric key and MAC algorithm.
func (b *MemoryBackend) GenerateMAC(keyID string, message []byte, algorithm MACAlgorithm) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, ErrInvalidKeySpec
	}

	var hmacHash []byte
	switch algorithm {
	case MACAlgorithmHMACSHA224:
		h := hmac.New(sha256.New224, key.symmetric)
		h.Write(message)
		hmacHash = h.Sum(nil)
	case MACAlgorithmHMACSHA256:
		h := hmac.New(sha256.New, key.symmetric)
		h.Write(message)
		hmacHash = h.Sum(nil)
	case MACAlgorithmHMACSHA384:
		h := hmac.New(sha512.New384, key.symmetric)
		h.Write(message)
		hmacHash = h.Sum(nil)
	case MACAlgorithmHMACSHA512:
		h := hmac.New(sha512.New, key.symmetric)
		h.Write(message)
		hmacHash = h.Sum(nil)
	default:
		h := hmac.New(sha256.New, key.symmetric)
		h.Write(message)
		hmacHash = h.Sum(nil)
	}

	return hmacHash, nil
}

// VerifyMAC verifies a Message Authentication Code (MAC) for a message
// using the specified symmetric key and MAC algorithm.
func (b *MemoryBackend) VerifyMAC(keyID string, message, macValue []byte, algorithm MACAlgorithm) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return false, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return false, ErrInvalidKeySpec
	}

	var expectedMAC []byte
	switch algorithm {
	case MACAlgorithmHMACSHA224:
		h := hmac.New(sha256.New224, key.symmetric)
		h.Write(message)
		expectedMAC = h.Sum(nil)
	case MACAlgorithmHMACSHA256:
		h := hmac.New(sha256.New, key.symmetric)
		h.Write(message)
		expectedMAC = h.Sum(nil)
	case MACAlgorithmHMACSHA384:
		h := hmac.New(sha512.New384, key.symmetric)
		h.Write(message)
		expectedMAC = h.Sum(nil)
	case MACAlgorithmHMACSHA512:
		h := hmac.New(sha512.New, key.symmetric)
		h.Write(message)
		expectedMAC = h.Sum(nil)
	default:
		h := hmac.New(sha256.New, key.symmetric)
		h.Write(message)
		expectedMAC = h.Sum(nil)
	}

	return hmac.Equal(macValue, expectedMAC), nil
}

// GenerateDataKey generates a new data key for envelope encryption.
// Returns both plaintext and encrypted versions of the data key.
func (b *MemoryBackend) GenerateDataKey(keyID string, keySpec string, numberOfBytes int) (*GenerateDataKeyResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, exists := b.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	var keyBytes int
	if numberOfBytes > 0 {
		keyBytes = numberOfBytes
	} else {
		switch keySpec {
		case "AES_128":
			keyBytes = 16
		case "AES_256":
			keyBytes = 32
		default:
			keyBytes = 32
		}
	}

	plaintext := make([]byte, keyBytes)
	if _, err := rand.Read(plaintext); err != nil {
		return nil, err
	}

	ciphertext, err := aesEncrypt(b.wrappingKey, plaintext, nil)
	if err != nil {
		return nil, err
	}

	ciphertextWithKeyID := prependKeyIDToCiphertext(keyID, ciphertext)

	return &GenerateDataKeyResult{
		Plaintext:      plaintext,
		Ciphertext:     ciphertextWithKeyID,
		PlaintextCRC32: computeCRC32(plaintext),
		KeyID:          keyID,
	}, nil
}

// GetPublicKey retrieves the public key for an RSA or ECDSA key.
// Returns the encoded public key bytes.
func (b *MemoryBackend) GetPublicKey(keyID string) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.rsaKey != nil {
		return encodePublicKey(&key.rsaKey.PublicKey)
	} else if key.ecdsaKey != nil {
		return encodePublicKey(&key.ecdsaKey.PublicKey)
	}

	return nil, ErrInvalidKeySpec
}

// IsKeyAvailable checks whether a key is available for cryptographic operations.
// A key is considered available if it has key material (symmetric, RSA, or ECDSA).
func (b *MemoryBackend) IsKeyAvailable(keyID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return false
	}

	return key.symmetric != nil || key.rsaKey != nil || key.ecdsaKey != nil
}

// KeyExists checks whether a key with the specified ID exists in the backend.
func (b *MemoryBackend) KeyExists(keyID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, exists := b.keys[keyID]
	return exists
}

func hashForAlgorithm(algorithm SigningAlgorithm, message []byte) []byte {
	switch algorithm {
	case SigningAlgorithmRSAPKCS1SHA256, SigningAlgorithmRSAPSSSHA256, SigningAlgorithmECDSASHA256:
		h := sha256.Sum256(message)
		return h[:]
	case SigningAlgorithmRSAPKCS1SHA384, SigningAlgorithmRSAPSSSHA384, SigningAlgorithmECDSASHA384:
		h := sha512.Sum384(message)
		return h[:]
	case SigningAlgorithmRSAPKCS1SHA512, SigningAlgorithmRSAPSSSHA512, SigningAlgorithmECDSASHA512:
		h := sha512.Sum512(message)
		return h[:]
	default:
		h := sha256.Sum256(message)
		return h[:]
	}
}

func hashAlgorithm(algorithm SigningAlgorithm) crypto.Hash {
	switch algorithm {
	case SigningAlgorithmRSAPKCS1SHA256, SigningAlgorithmRSAPSSSHA256:
		return crypto.SHA256
	case SigningAlgorithmRSAPKCS1SHA384, SigningAlgorithmRSAPSSSHA384:
		return crypto.SHA384
	case SigningAlgorithmRSAPKCS1SHA512, SigningAlgorithmRSAPSSSHA512:
		return crypto.SHA512
	default:
		return crypto.SHA256
	}
}
