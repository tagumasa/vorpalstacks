package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

const uuidSize = 16

// GenerateKeyID generates a new KMS key ID.
func GenerateKeyID() (string, error) {
	uuid := make([]byte, uuidSize)
	if _, err := rand.Read(uuid); err != nil {
		return "", err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// GenerateGrantID generates a new KMS grant ID.
func GenerateGrantID() (string, error) {
	uuid := make([]byte, uuidSize)
	if _, err := rand.Read(uuid); err != nil {
		return "", err
	}
	return hex.EncodeToString(uuid), nil
}

// GenerateImportToken generates a new KMS key material import token.
func GenerateImportToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

// GenerateGrantToken generates a new KMS grant token.
func GenerateGrantToken() (string, error) {
	token := make([]byte, 48)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

// GenerateWrappingKeyPair generates an RSA key pair for wrapping imported key material.
// Returns the public key (base64-encoded DER PKIX) and the private key (raw PKCS8 bytes).
func GenerateWrappingKeyPair(wrappingKeySpec string) (string, []byte, error) {
	bits := 2048
	switch wrappingKeySpec {
	case "RSA_3072":
		bits = 3072
	case "RSA_4096":
		bits = 4096
	}
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", nil, fmt.Errorf("generate RSA key: %w", err)
	}
	pubDER, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		return "", nil, fmt.Errorf("marshal public key: %w", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", nil, fmt.Errorf("marshal private key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(pubDER), privBytes, nil
}

// DecryptWrappedKeyMaterial decrypts key material that was encrypted with the wrapping public key
// using RSA-OAEP with SHA-256 (matching the RSASSA_OAEP_SHA_256 wrapping algorithm).
func DecryptWrappedKeyMaterial(privKeyBytes []byte, encryptedKeyMaterial []byte) ([]byte, error) {
	key, err := x509.ParsePKCS8PrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaKey, encryptedKeyMaterial, nil)
}
