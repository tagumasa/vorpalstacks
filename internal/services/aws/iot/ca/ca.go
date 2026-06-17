package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	pebbledb "vorpalstacks/internal/core/storage/pebbledb"
)

const (
	caBucketName = "iot-ca"
	caKeyPEMKey  = "ca-key"
	caCertPEMKey = "ca-cert"
)

type CertificateAuthority struct {
	bucket  storage.Bucket
	mu      sync.RWMutex
	rootCA  *x509.Certificate
	rootKey *ecdsa.PrivateKey
}

func NewCertificateAuthority(s storage.BasicStorage) (*CertificateAuthority, error) {
	ca := &CertificateAuthority{
		bucket: s.Bucket(caBucketName),
	}
	if err := ca.loadOrCreateRootCA(); err != nil {
		return nil, fmt.Errorf("failed to initialise CA: %w", err)
	}
	return ca, nil
}

func (ca *CertificateAuthority) RootCA() *x509.Certificate {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	return ca.rootCA
}

func (ca *CertificateAuthority) RootKey() *ecdsa.PrivateKey {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	return ca.rootKey
}

func (ca *CertificateAuthority) loadOrCreateRootCA() error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	certData, err := ca.bucket.Get([]byte(caCertPEMKey))
	if err != nil && !errors.Is(err, pebbledb.ErrKeyNotFound) {
		return fmt.Errorf("failed to read CA cert: %w", err)
	}
	keyData, err := ca.bucket.Get([]byte(caKeyPEMKey))
	if err != nil && !errors.Is(err, pebbledb.ErrKeyNotFound) {
		return fmt.Errorf("failed to read CA key: %w", err)
	}

	if certData != nil && keyData != nil {
		cert, certErr := parsePEMCertificate(certData)
		key, keyErr := parsePEMPrivateKey(keyData)
		if certErr == nil && keyErr == nil {
			ca.rootCA = cert
			ca.rootKey = key
			return nil
		}
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "Vorpalstacks IoT CA",
			Organization: []string{"vorpalstacks"},
		},
		NotBefore:             time.Now().UTC().Add(-1 * time.Hour),
		NotAfter:              time.Now().UTC().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate CA key: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("failed to generate CA cert: %w", err)
	}

	certPEMBytes := encodeCertPEM(certBytes)
	if err := ca.bucket.Put([]byte(caCertPEMKey), certPEMBytes); err != nil {
		return fmt.Errorf("failed to persist CA cert: %w", err)
	}
	keyPEM := encodeKeyPEM(key)
	if err := ca.bucket.Put([]byte(caKeyPEMKey), keyPEM); err != nil {
		return fmt.Errorf("failed to persist CA key: %w", err)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return fmt.Errorf("failed to parse generated CA cert: %w", err)
	}
	ca.rootCA = cert
	ca.rootKey = key
	return nil
}

func (ca *CertificateAuthority) IssueCertificate() (certPEM string, keyPEM string, pubKeyPEM string, certificateID string, err error) {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	if ca.rootCA == nil {
		return "", "", "", "", fmt.Errorf("CA not initialised")
	}

	template := &x509.Certificate{
		SerialNumber: generateSerialNumber(),
		Subject: pkix.Name{
			CommonName: "Vorpalstacks Device",
		},
		NotBefore: time.Now().UTC(),
		NotAfter:  time.Now().UTC().Add(365 * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to generate device key: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, ca.rootCA, &key.PublicKey, ca.rootKey)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to issue certificate: %w", err)
	}

	return string(encodeCertPEM(certBytes)), string(encodeKeyPEM(key)), string(encodePublicKeyPEM(key)), CertIDFromDER(certBytes), nil
}

func (ca *CertificateAuthority) IssueCertificateFromCSR(csrPEM string) (certPEM string, certificateID string, err error) {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	if ca.rootCA == nil {
		return "", "", fmt.Errorf("CA not initialised")
	}

	csrBlock, _ := pem.Decode([]byte(csrPEM))
	if csrBlock == nil {
		return "", "", fmt.Errorf("failed to decode CSR PEM")
	}
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse CSR: %w", err)
	}

	if err := csr.CheckSignature(); err != nil {
		return "", "", fmt.Errorf("CSR signature verification failed: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: generateSerialNumber(),
		Subject:      csr.Subject,
		NotBefore:    time.Now().UTC(),
		NotAfter:     time.Now().UTC().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, ca.rootCA, csr.PublicKey, ca.rootKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to issue certificate from CSR: %w", err)
	}

	return string(encodeCertPEM(certBytes)), CertIDFromDER(certBytes), nil
}

func (ca *CertificateAuthority) VerifyCertificate(certPEM string) error {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	if ca.rootCA == nil {
		return fmt.Errorf("CA not initialised")
	}

	certBlock, _ := pem.Decode([]byte(certPEM))
	if certBlock == nil {
		return fmt.Errorf("failed to decode certificate PEM")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(ca.rootCA)
	if _, err = cert.Verify(x509.VerifyOptions{Roots: roots}); err != nil {
		return err
	}
	return nil
}

func generateSerialNumber() *big.Int {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	n, _ := rand.Int(rand.Reader, limit)
	return n
}

func parsePEMCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM data found")
	}
	return x509.ParseCertificate(block.Bytes)
}

func parsePEMPrivateKey(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM data found")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}
	return key, nil
}

func encodeCertPEM(certBytes []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
}

func encodeKeyPEM(key *ecdsa.PrivateKey) []byte {
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil
	}
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
}

func encodePublicKeyPEM(key *ecdsa.PrivateKey) []byte {
	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
}
