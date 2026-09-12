package cognitoidentityprovider

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"

	"vorpalstacks/pkg/vsjwt"
)

// Core functions for the pool signing-certificate family.

// signingCertificateValidityYears is the certificate lifetime the
// GetSigningCertificate model documents: issued certificates are valid for
// 10 years from the date of issue.
const signingCertificateValidityYears = 10

// getSigningCertificateCore returns the pool's SAML federation signing
// certificate: a self-signed X.509 certificate over the pool's signing key.
// The certificate is derived deterministically from the stored pool record —
// serial 1, subject CN = pool ID, validity running from the pool's creation
// for the documented ten years — so every call returns the identical
// certificate without persisting a separate record.
func (s *CognitoService) getSigningCertificateCore(region, userPoolID string) (string, error) {
	if userPoolID == "" {
		return "", ErrInvalidParameter
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return "", err
	}

	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return "", ErrResourceNotFound
	}

	privateKey, err := vsjwt.DecodePrivateKeyFromPEM(pool.JwtPrivateKey)
	if err != nil {
		return "", ErrInternalError
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: pool.ID},
		NotBefore:    pool.CreationDate,
		NotAfter:     pool.CreationDate.AddDate(signingCertificateValidityYears, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", ErrInternalError
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})), nil
}
