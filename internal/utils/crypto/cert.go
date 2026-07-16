package crypto

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
)

// EncodeCertificatePEM encodes DER-encoded certificate bytes to a PEM string.
func EncodeCertificatePEM(certDER []byte) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}))
}

// ParseCertificatePEM decodes a PEM-encoded certificate into an x509.Certificate.
func ParseCertificatePEM(pemData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("no PEM data found")
	}
	return x509.ParseCertificate(block.Bytes)
}

// CreateCertificate wraps x509.CreateCertificate, signing template with
// parent's private key. publicKey is the subject's public key.
func CreateCertificate(template, parent *x509.Certificate, publicKey any, parentPrivateKey any) ([]byte, error) {
	return x509.CreateCertificate(rand.Reader, template, parent, publicKey, parentPrivateKey)
}

// VerifyCertificateAgainstRoot verifies that cert was signed by root.
func VerifyCertificateAgainstRoot(cert *x509.Certificate, root *x509.Certificate) error {
	roots := x509.NewCertPool()
	roots.AddCert(root)
	_, err := cert.Verify(x509.VerifyOptions{Roots: roots})
	return err
}

// GenerateSerialNumber generates a random 128-bit serial number suitable
// for X.509 certificates.
func GenerateSerialNumber() (*big.Int, error) {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, limit)
}

// FingerprintPEM computes the SHA-256 fingerprint of a PEM-encoded
// certificate. If PEM decoding fails, the raw input is hashed instead
// (best-effort, maintaining backward compatibility with callers that
// may pass non-PEM data).
func FingerprintPEM(certPEM string) string {
	block, _ := pem.Decode([]byte(certPEM))
	data := []byte(certPEM)
	if block != nil {
		data = block.Bytes
	}
	return SHA256Hash(data)
}

// FingerprintX509 computes the SHA-256 fingerprint of an x509.Certificate
// from its raw DER bytes.
func FingerprintX509(cert *x509.Certificate) string {
	return SHA256Hash(cert.Raw)
}
