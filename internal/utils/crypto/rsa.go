package crypto

// Package crypto provides cryptographic utilities for vorpalstacks, including
// RSA encryption and decryption functions.

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// ErrInvalidKey is returned when the RSA key is invalid.
var (
	// ErrInvalidKey is returned when the RSA key is invalid.
	ErrInvalidKey = errors.New("invalid RSA key")
	// ErrEncryptionFailed is returned when RSA encryption fails.
	ErrEncryptionFailed = errors.New("encryption failed")
	// ErrDecryptionFailed is returned when RSA decryption fails.
	ErrDecryptionFailed = errors.New("decryption failed")
	// ErrInvalidPEMFormat is returned when the PEM format is invalid.
	ErrInvalidPEMFormat = errors.New("invalid PEM format")
	// ErrInvalidKeySize is returned when the key size is less than 2048 bits.
	ErrInvalidKeySize = errors.New("invalid key size: must be at least 2048 bits")
)

// Default constants for RSA operations.
const (
	// DefaultRSAKeySize is the default RSA key size in bits.
	DefaultRSAKeySize = 2048
)

// GenerateRSAKey generates an RSA key pair of the specified bit length.
// Uses crypto/rand for secure random generation.
//
// Parameters:
//   - bits: The key size in bits (must be at least 2048)
//
// Returns:
//   - *rsa.PrivateKey: The generated private key
//   - error: An error if key generation fails
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	if bits < 2048 {
		return nil, ErrInvalidKeySize
	}
	return rsa.GenerateKey(rand.Reader, bits)
}

// RSAPublicKeyToPEM converts an RSA public key to PEM format.
//
// Parameters:
//   - pubKey: The RSA public key to convert
//
// Returns:
//   - []byte: The PEM-encoded public key
func RSAPublicKeyToPEM(pubKey *rsa.PublicKey) []byte {
	pubKeyBytes := x509.MarshalPKCS1PublicKey(pubKey)
	pemBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	return pem.EncodeToMemory(pemBlock)
}

// RSAPrivateKeyToPEM converts an RSA private key to PEM format.
//
// Parameters:
//   - privKey: The RSA private key to convert
//
// Returns:
//   - []byte: The PEM-encoded private key
func RSAPrivateKeyToPEM(privKey *rsa.PrivateKey) []byte {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	}
	return pem.EncodeToMemory(pemBlock)
}

// PEMToRSAPublicKey parses a PEM-encoded RSA public key.
//
// Parameters:
//   - pemData: The PEM-encoded public key data
//
// Returns:
//   - *rsa.PublicKey: The parsed public key
//   - error: An error if parsing fails
func PEMToRSAPublicKey(pemData []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, ErrInvalidPEMFormat
	}

	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, ErrInvalidKey
	}

	return pubKey, nil
}

// PEMToRSAPrivateKey parses a PEM-encoded RSA private key.
//
// Parameters:
//   - pemData: The PEM-encoded private key data
//
// Returns:
//   - *rsa.PrivateKey: The parsed private key
//   - error: An error if parsing fails
func PEMToRSAPrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, ErrInvalidPEMFormat
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, ErrInvalidKey
	}

	return privKey, nil
}

// RSAEncrypt encrypts plaintext using RSA-OAEP with SHA-256.
//
// Parameters:
//   - pubKey: The RSA public key to use for encryption
//   - plaintext: The data to encrypt
//
// Returns:
//   - []byte: The encrypted ciphertext
//   - error: An error if encryption fails
func RSAEncrypt(pubKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptOAEP(NewSHA256Hash(), rand.Reader, pubKey, plaintext, nil)
	if err != nil {
		return nil, ErrEncryptionFailed
	}
	return ciphertext, nil
}

// RSADecrypt decrypts ciphertext using RSA-OAEP with SHA-256.
//
// Parameters:
//   - privKey: The RSA private key to use for decryption
//   - ciphertext: The encrypted data to decrypt
//
// Returns:
//   - []byte: The decrypted plaintext
//   - error: An error if decryption fails
func RSADecrypt(privKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	plaintext, err := rsa.DecryptOAEP(NewSHA256Hash(), rand.Reader, privKey, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}
	return plaintext, nil
}

// RSASign signs a message using RSA-SHA256.
//
// Parameters:
//   - privKey: The RSA private key to use for signing
//   - message: The message to sign
//
// Returns:
//   - []byte: The signature
//   - error: An error if signing fails
func RSASign(privKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	hashed := SHA256HashRaw(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// RSAVerify verifies an RSA-SHA256 signature.
//
// Parameters:
//   - pubKey: The RSA public key to use for verification
//   - message: The original message
//   - signature: The signature to verify
//
// Returns:
//   - error: nil if verification succeeds, otherwise an error
func RSAVerify(pubKey *rsa.PublicKey, message, signature []byte) error {
	hashed := SHA256HashRaw(message)
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed, signature)
}
