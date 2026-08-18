// Package iam provides IAM service operations for vorpalstacks.
//
// This file holds the credential-material parsing shared by the signing
// certificate, server certificate and SSH public key operations: X.509
// parsing, certificate fingerprints and key-pair matching.
package iam

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

// parseCertificate decodes a PEM certificate body and parses it as X.509.
func parseCertificate(pemBody string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemBody))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("no CERTIFICATE PEM block found")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse X.509 certificate: %w", err)
	}
	return cert, nil
}

// certificateFingerprint returns a stable fingerprint of the certificate
// contents, computed over the DER encoding so that PEM whitespace changes
// do not hide a duplicate upload.
func certificateFingerprint(cert *x509.Certificate) string {
	sum := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(sum[:])
}

// parsePrivateKey decodes a PEM private key in PKCS#1, PKCS#8 or SEC 1 form.
func parsePrivateKey(pemBody string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemBody))
	if block == nil {
		return nil, errors.New("no private key PEM block found")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	return nil, errors.New("unsupported private key format")
}

// keyPairMatches reports whether the certificate's public key is the public
// half of the given private key.
func keyPairMatches(cert *x509.Certificate, priv crypto.PrivateKey) bool {
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		key, ok := priv.(*rsa.PrivateKey)
		return ok && pub.Equal(&key.PublicKey)
	case *ecdsa.PublicKey:
		key, ok := priv.(*ecdsa.PrivateKey)
		return ok && pub.Equal(&key.PublicKey)
	}
	return false
}

// parseSSHPublicKey parses an authorized-keys style SSH public key body.
func parseSSHPublicKey(body string) (ssh.PublicKey, error) {
	pub, _, _, _, err := ssh.ParseAuthorizedKey([]byte(body))
	if err != nil {
		return nil, fmt.Errorf("parse SSH public key: %w", err)
	}
	return pub, nil
}

// canonicalSSHPublicKeyBody renders the key in the SSH single-line format,
// normalising whitespace and comments of the uploaded body.
func canonicalSSHPublicKeyBody(pub ssh.PublicKey) string {
	return strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pub)))
}

// sshPublicKeyBodyPEM renders the key as a PEM SubjectPublicKeyInfo block.
func sshPublicKeyBodyPEM(pub ssh.PublicKey) (string, error) {
	cryptoPub, ok := pub.(ssh.CryptoPublicKey)
	if !ok {
		return "", errors.New("key type has no crypto public key")
	}
	der, err := x509.MarshalPKIXPublicKey(cryptoPub.CryptoPublicKey())
	if err != nil {
		return "", fmt.Errorf("marshal public key: %w", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})), nil
}
