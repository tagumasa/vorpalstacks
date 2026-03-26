package pebbledb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	// ErrEncryptionDisabled is returned when attempting to use encryption on a database that has it disabled.
	ErrEncryptionDisabled = errors.New("encryption is disabled")
	// ErrInvalidKeySize is returned when an encryption key has an invalid length.
	ErrInvalidKeySize = errors.New("key must be 16, 24, or 32 bytes")
)

// Encryptor provides an interface for encrypting and decrypting data.
// Implementations should handle key management and encryption algorithms.
type Encryptor interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

// aesGCMEncryptor provides AES-GCM encryption for database values.
// It uses a 96-bit nonce generated randomly for each encryption operation.
type aesGCMEncryptor struct {
	key   []byte
	block cipher.Block
	gcm   cipher.AEAD
}

// NewAESGCMEncryptor creates a new AES-GCM encryptor with the given key.
// The key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256 respectively.
func NewAESGCMEncryptor(key []byte) (*aesGCMEncryptor, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &aesGCMEncryptor{
		key:   key,
		block: block,
		gcm:   gcm,
	}, nil
}

// Encrypt encrypts the plaintext using AES-GCM.
// The ciphertext includes a random nonce, so the same plaintext produces different outputs.
func (e *aesGCMEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := e.gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts the ciphertext using AES-GCM.
// The nonce is extracted from the beginning of the ciphertext.
func (e *aesGCMEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := e.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, cipherData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := e.gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// noopEncryptor is a no-op implementation of the Encryptor interface.
// It returns the input unchanged, useful for testing or when encryption is not required.
type noopEncryptor struct{}

// NewNoopEncryptor creates a no-op encryptor that performs no encryption.
// This is useful for development or when encryption is not required.
func NewNoopEncryptor() *noopEncryptor {
	return &noopEncryptor{}
}

// Encrypt returns the plaintext unchanged.
func (e *noopEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}

// Decrypt returns the ciphertext unchanged.
func (e *noopEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
}
