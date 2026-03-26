package crypto

// Package crypto provides cryptographic utilities for vorpalstacks, including
// AES encryption and decryption functions.

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// Error variables for AES operations.
var (
	// ErrInvalidKeyLength is returned when the key length is invalid.
	ErrInvalidKeyLength = errors.New("invalid key length: must be 16, 24, or 32 bytes")
	// ErrInvalidCiphertext is returned when the ciphertext is too short.
	ErrInvalidCiphertext = errors.New("invalid ciphertext: too short")
	// ErrAuthentication is returned when authentication fails.
	ErrAuthentication = errors.New("authentication failed")
)

// AESGCMEncrypt encrypts plaintext using AES-GCM with a random nonce.
func AESGCMEncrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKeyLength
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := RandomNonce()
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// AESGCMDecrypt decrypts ciphertext using AES-GCM.
func AESGCMDecrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKeyLength
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrAuthentication
	}

	return plaintext, nil
}

// AESGCMEncryptWithNonce encrypts plaintext using AES-GCM with a provided nonce.
func AESGCMEncryptWithNonce(key, plaintext, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKeyLength
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

// AESGCMDecryptWithNonce decrypts ciphertext using AES-GCM with a provided nonce.
func AESGCMDecryptWithNonce(key, ciphertext, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKeyLength
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrAuthentication
	}

	return plaintext, nil
}
