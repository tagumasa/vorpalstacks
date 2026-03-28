package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateKeyID generates a new KMS key ID.
func GenerateKeyID() (string, error) {
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		return "", err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return hex.EncodeToString(uuid), nil
}

// GenerateGrantID generates a new KMS grant ID.
func GenerateGrantID() (string, error) {
	uuid := make([]byte, 16)
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
	return hex.EncodeToString(token), nil
}

// GenerateGrantToken generates a new KMS grant token.
func GenerateGrantToken() (string, error) {
	token := make([]byte, 48)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
