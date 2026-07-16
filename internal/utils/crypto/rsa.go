package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

// ErrInvalidKeySize is returned when the key size is less than 2048 bits.
var ErrInvalidKeySize = errors.New("invalid key size: must be at least 2048 bits")

// DefaultRSAKeySize is the default RSA key size in bits.
const DefaultRSAKeySize = 2048

// GenerateRSAKey generates an RSA key pair of the specified bit length.
// Uses crypto/rand for secure random generation.
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	if bits < 2048 {
		return nil, ErrInvalidKeySize
	}
	return rsa.GenerateKey(rand.Reader, bits)
}
