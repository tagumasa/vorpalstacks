package crypto

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// EncodePrivateKeyPEM encodes a private key to PEM format. Supports
// *ecdsa.PrivateKey (SEC1) and *rsa.PrivateKey (PKCS1).
func EncodePrivateKeyPEM(key any) (string, error) {
	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		der, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return "", fmt.Errorf("failed to marshal EC private key: %w", err)
		}
		return string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der})), nil
	case *rsa.PrivateKey:
		der := x509.MarshalPKCS1PrivateKey(k)
		return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der})), nil
	default:
		return "", fmt.Errorf("unsupported private key type: %T", key)
	}
}

// EncodePublicKeyPEM encodes a public key to PEM (PKIX) format. Supports
// *ecdsa.PublicKey and *rsa.PublicKey.
func EncodePublicKeyPEM(key any) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})), nil
}

// ParsePrivateKeyPEM decodes a PEM-encoded private key. Accepts EC (SEC1),
// RSA (PKCS1), and PKCS8 formats.
func ParsePrivateKeyPEM(pemData []byte) (any, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("no PEM data found")
	}

	switch block.Type {
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		return x509.ParsePKCS8PrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
}
