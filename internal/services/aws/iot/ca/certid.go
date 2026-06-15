package ca

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/pem"
)

// ComputeCertID derives a deterministic certificate identifier from a
// PEM-encoded certificate by SHA-256 hashing the raw DER bytes.
// If PEM decoding fails, the raw input is hashed instead (best-effort).
func ComputeCertID(certPEM string) string {
	block, _ := pem.Decode([]byte(certPEM))
	data := []byte(certPEM)
	if block != nil {
		data = block.Bytes
	}
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// CertIDFromDER derives a certificate identifier from already-decoded DER bytes.
func CertIDFromDER(certBytes []byte) string {
	h := sha256.Sum256(certBytes)
	return hex.EncodeToString(h[:])
}
