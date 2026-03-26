package crypto

// Package crypto provides cryptographic utilities for vorpalstacks, including
// random byte generation for secure operations.

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// ErrRandomGeneration is returned when random byte generation fails.
var ErrRandomGeneration = errors.New("failed to generate random bytes")

// RandomBytes generates a slice of random bytes of the specified length.
func RandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		return nil, ErrRandomGeneration
	}
	return bytes, nil
}

// RandomHex generates a random hex string of the specified length.
func RandomHex(length int) (string, error) {
	bytes, err := RandomBytes(length)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// RandomIV generates a random initialization vector (IV) of the specified length.
func RandomIV(length int) ([]byte, error) {
	return RandomBytes(length)
}

// RandomKey generates a random encryption key of the specified length.
func RandomKey(length int) ([]byte, error) {
	return RandomBytes(length)
}

// RandomNonce generates a random nonce (12 bytes, suitable for AES-GCM).
func RandomNonce() ([]byte, error) {
	return RandomBytes(12)
}
