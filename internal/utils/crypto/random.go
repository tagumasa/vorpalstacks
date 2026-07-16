package crypto

import (
	"crypto/rand"
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

// RandomKey generates a random encryption key of the specified length.
func RandomKey(length int) ([]byte, error) {
	return RandomBytes(length)
}

// RandomNonce generates a random nonce (12 bytes, suitable for AES-GCM).
func RandomNonce() ([]byte, error) {
	return RandomBytes(12)
}
