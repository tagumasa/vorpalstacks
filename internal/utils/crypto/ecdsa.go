package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

// GenerateECDSAKey generates a new ECDSA private key on the given curve.
// Common choices: elliptic.P256(), elliptic.P384(), elliptic.P521().
func GenerateECDSAKey(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECDSA key: %w", err)
	}
	return key, nil
}
