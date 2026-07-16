package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/tls"
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
	vcrypto "vorpalstacks/internal/utils/crypto"
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
		cert, certErr := vcrypto.ParseCertificatePEM(certData)
		if certErr == nil {
			keyAny, keyErr := vcrypto.ParsePrivateKeyPEM(keyData)
			if keyErr == nil {
				if ecKey, ok := keyAny.(*ecdsa.PrivateKey); ok {
					ca.rootCA = cert
					ca.rootKey = ecKey
					return nil
				}
			}
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

	key, err := vcrypto.GenerateECDSAKey(elliptic.P256())
	if err != nil {
		return fmt.Errorf("failed to generate CA key: %w", err)
	}

	certBytes, err := vcrypto.CreateCertificate(template, template, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("failed to generate CA cert: %w", err)
	}

	certPEMStr := vcrypto.EncodeCertificatePEM(certBytes)
	if err := ca.bucket.Put([]byte(caCertPEMKey), []byte(certPEMStr)); err != nil {
		return fmt.Errorf("failed to persist CA cert: %w", err)
	}
	keyPEMStr, err := vcrypto.EncodePrivateKeyPEM(key)
	if err != nil {
		return fmt.Errorf("failed to encode CA key: %w", err)
	}
	if err := ca.bucket.Put([]byte(caKeyPEMKey), []byte(keyPEMStr)); err != nil {
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

	serial, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to generate serial: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: serial,
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

	key, err := vcrypto.GenerateECDSAKey(elliptic.P256())
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to generate device key: %w", err)
	}

	certBytes, err := vcrypto.CreateCertificate(template, ca.rootCA, &key.PublicKey, ca.rootKey)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to issue certificate: %w", err)
	}

	certPEMStr := vcrypto.EncodeCertificatePEM(certBytes)
	keyPEMStr, err := vcrypto.EncodePrivateKeyPEM(key)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to encode device key: %w", err)
	}
	pubKeyPEMStr, err := vcrypto.EncodePublicKeyPEM(&key.PublicKey)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to encode device public key: %w", err)
	}

	return certPEMStr, keyPEMStr, pubKeyPEMStr, vcrypto.SHA256Hash(certBytes), nil
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

	serial, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate serial: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      csr.Subject,
		NotBefore:    time.Now().UTC(),
		NotAfter:     time.Now().UTC().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
	}

	certBytes, err := vcrypto.CreateCertificate(template, ca.rootCA, csr.PublicKey, ca.rootKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to issue certificate from CSR: %w", err)
	}

	return vcrypto.EncodeCertificatePEM(certBytes), vcrypto.SHA256Hash(certBytes), nil
}

func (ca *CertificateAuthority) VerifyCertificate(certPEM string) error {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	if ca.rootCA == nil {
		return fmt.Errorf("CA not initialised")
	}

	cert, err := vcrypto.ParseCertificatePEM([]byte(certPEM))
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	return vcrypto.VerifyCertificateAgainstRoot(cert, ca.rootCA)
}

// IssueServerCertificate issues a dedicated TLS server certificate
// signed by the CA root. Not persisted — clients trust the CA root,
// not the server certificate, so regeneration on each broker restart
// is safe.
func (ca *CertificateAuthority) IssueServerCertificate(commonName string) (tls.Certificate, error) {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	if ca.rootCA == nil {
		return tls.Certificate{}, fmt.Errorf("CA not initialised")
	}

	serial, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to generate serial: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore: time.Now().UTC().Add(-1 * time.Hour),
		NotAfter:  time.Now().UTC().Add(365 * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		DNSNames: []string{commonName},
	}

	key, err := vcrypto.GenerateECDSAKey(elliptic.P256())
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to generate server key: %w", err)
	}

	certDER, err := vcrypto.CreateCertificate(template, ca.rootCA, &key.PublicKey, ca.rootKey)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to issue server certificate: %w", err)
	}

	serverCert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to parse server certificate: %w", err)
	}

	return tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  key,
		Leaf:        serverCert,
	}, nil
}
