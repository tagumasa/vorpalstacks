package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"crypto/rand"
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
