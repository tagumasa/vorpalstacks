package http

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"vorpalstacks/internal/core/storage"
)

// TLSManager manages TLS certificates, including loading from files
// or generating self-signed certificates.
type TLSManager struct {
	certPEM  []byte
	keyPEM   []byte
	storage  storage.BasicStorage
	certPath string
	keyPath  string
	hostname string
}

// NewTLSManager creates a new TLS manager.
func NewTLSManager(storage storage.BasicStorage, certPath, keyPath, hostname string) (*TLSManager, error) {
	mgr := &TLSManager{
		storage:  storage,
		certPath: certPath,
		keyPath:  keyPath,
		hostname: hostname,
	}

	if certPath == "auto" || keyPath == "auto" {
		_ = storage.CreateBucket("tls")

		if err := mgr.loadFromStorage(); err != nil {
			if err := mgr.generateSelfSignedCert(); err != nil {
				return nil, err
			}
		}

		if len(mgr.certPEM) == 0 || len(mgr.keyPEM) == 0 {
			if err := mgr.generateSelfSignedCert(); err != nil {
				return nil, err
			}
		}
	} else {
		certPEM, err := os.ReadFile(certPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read cert file: %w", err)
		}
		keyPEM, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file: %w", err)
		}
		mgr.certPEM = certPEM
		mgr.keyPEM = keyPEM
	}

	return mgr, nil
}

func (m *TLSManager) loadFromStorage() error {
	bucket := m.storage.Bucket("tls")

	certPEM, err := bucket.Get([]byte("cert"))
	if err != nil {
		return err
	}

	keyPEM, err := bucket.Get([]byte("key"))
	if err != nil {
		return err
	}

	m.certPEM = certPEM
	m.keyPEM = keyPEM
	return nil
}

func (m *TLSManager) generateSelfSignedCert() error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %w", err)
	}

	cn := "localhost"
	dnsNames := []string{"localhost"}
	ips := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}

	if m.hostname != "" && m.hostname != "localhost" {
		cn = m.hostname
		dnsNames = append(dnsNames, m.hostname)
		if ip := net.ParseIP(m.hostname); ip != nil {
			ips = append(ips, ip)
		}
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"VorpalStacks"},
			CommonName:   cn,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: ips,
		DNSNames:    dnsNames,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	m.certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	m.keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	if m.certPath == "auto" || m.keyPath == "auto" {
		if err := m.saveToStorage(); err != nil {
			return err
		}
	}

	return nil
}

func (m *TLSManager) saveToStorage() error {
	_ = m.storage.CreateBucket("tls")

	bucket := m.storage.Bucket("tls")

	if err := bucket.Put([]byte("cert"), m.certPEM); err != nil {
		return fmt.Errorf("failed to save cert: %w", err)
	}

	if err := bucket.Put([]byte("key"), m.keyPEM); err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}

	return nil
}

// GetTLSConfig returns a TLS configuration for use with an HTTP server.
func (m *TLSManager) GetTLSConfig() (*tls.Config, error) {
	cert, err := tls.X509KeyPair(m.certPEM, m.keyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS key pair: %w", err)
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2", "http/1.1"},
	}, nil
}
