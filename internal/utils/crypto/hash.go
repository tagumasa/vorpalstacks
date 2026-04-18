// Package crypto provides cryptographic utilities for vorpalstacks.
package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

// MD5Hash computes the MD5 hash of the given data.
func MD5Hash(data []byte) []byte {
	hash := md5.Sum(data)
	return hash[:]
}

// SHA256Hash computes the SHA-256 hash of the given data and returns it as a hex string.
func SHA256Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// SHA256HashString computes the SHA-256 hash of the given string and returns it as a hex string.
func SHA256HashString(data string) string {
	return SHA256Hash([]byte(data))
}

// SHA256HashRaw computes the SHA-256 hash of the given data and returns the raw bytes.
func SHA256HashRaw(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// NewSHA256Hash creates a new SHA-256 hash writer.
func NewSHA256Hash() hash.Hash {
	return sha256.New()
}

// HMACSHA256 computes the HMAC-SHA256 of the given data using the provided key.
func HMACSHA256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// HMACSHA256String computes the HMAC-SHA256 of the given string using the provided key.
func HMACSHA256String(key []byte, data string) []byte {
	return HMACSHA256(key, []byte(data))
}

// HMACSHA256Hex computes the HMAC-SHA256 of the given data and returns it as a hex string.
func HMACSHA256Hex(key []byte, data []byte) string {
	return hex.EncodeToString(HMACSHA256(key, data))
}

// HMACSHA256HexString computes the HMAC-SHA256 of the given string and returns it as a hex string.
func HMACSHA256HexString(key []byte, data string) string {
	return HMACSHA256Hex(key, []byte(data))
}
