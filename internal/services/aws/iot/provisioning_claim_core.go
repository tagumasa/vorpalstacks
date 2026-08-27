package iot

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"time"

	"github.com/google/uuid"

	iotstore "vorpalstacks/internal/store/aws/iot"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// ---------------------------------------------------------------------------
// Provisioning Claim Core. A provisioning claim is a temporary credential
// bound to an existing provisioning template; AWS uses the claim for
// just-in-time provisioning during device manufacturing. The claim consists
// of a short-lived self-signed X.509 certificate and its private key. The
// per-account registration code (used to register CA certificates) is
// lazily created on first GetRegistrationCode call and persists until
// explicitly deleted.
// ---------------------------------------------------------------------------

// ProvisioningClaimResult is the transport-agnostic result of
// CreateProvisioningClaim.
type ProvisioningClaimResult struct {
	CertificateID  string
	CertificatePem string
	PublicKeyPEM   string
	PrivateKeyPEM  string
	Expiration     int64
}

// createProvisioningClaimCore issues a short-lived provisioning claim
// certificate for an existing provisioning template. An unknown template
// yields ErrTemplateNotFound.
func (s *IoTService) createProvisioningClaimCore(store iotstore.IotStoreInterface, templateName string) (*ProvisioningClaimResult, error) {
	if templateName == "" {
		return nil, iotstore.ErrMissingParam
	}
	tmpl, err := store.GetProvisioningTemplate(templateName)
	if err != nil {
		return nil, err
	}
	if tmpl == nil || tmpl.TemplateName == "" {
		return nil, iotstore.ErrTemplateNotFound
	}

	// Generate an ECDSA P-256 private key for the claim certificate.
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}

	// Serial numbers must be unique per issuer; every other certificate
	// minted here uses the crypto/rand generator, and a clock value would
	// tie uniqueness to nanosecond resolution.
	serial, err := vcrypto.GenerateSerialNumber()
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}

	// Build a self-signed X.509 certificate valid for 1 hour.
	certID := uuid.New().String()
	now := time.Now().UTC()
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   certID,
			Organization: []string{"vorpalstacks-iot-provisioning-claim"},
		},
		NotBefore:             now.Add(-1 * time.Minute),
		NotAfter:              now.Add(1 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}

	certPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}))
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}
	pubKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyBytes}))
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, iotstore.ErrInternalFailure
	}
	privKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privKeyBytes}))

	return &ProvisioningClaimResult{
		CertificateID:  certID,
		CertificatePem: certPEM,
		PublicKeyPEM:   pubKeyPEM,
		PrivateKeyPEM:  privKeyPEM,
		Expiration:     now.Add(time.Hour).UTC().Unix(),
	}, nil
}

// getRegistrationCodeCore returns the per-account registration code,
// lazily creating it on first use.
func (s *IoTService) getRegistrationCodeCore(store iotstore.IotStoreInterface) (map[string]interface{}, error) {
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("config/registrationCode", &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		// Serialise lazy creation with a re-check under the lock: racing
		// first callers would otherwise mint different codes, and every
		// caller but the last would return a code the store no longer
		// holds.
		s.registrationCodeMu.Lock()
		exists, err = store.GetGenericExists("config/registrationCode", &rec)
		if err == nil && !exists {
			rec["registrationCode"] = uuid.New().String()
			err = store.PutGeneric("config/registrationCode", rec)
		}
		s.registrationCodeMu.Unlock()
		if err != nil {
			return nil, err
		}
	}
	return rec, nil
}

// deleteRegistrationCodeCore removes the per-account registration code.
func (s *IoTService) deleteRegistrationCodeCore(store iotstore.IotStoreInterface) error {
	return store.DeleteGeneric("config/registrationCode")
}
