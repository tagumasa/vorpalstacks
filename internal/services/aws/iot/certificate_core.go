package iot

import (
	"context"
	"fmt"
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// ---------------------------------------------------------------------------
// Certificate Core — issue, register, and CSR-signing paths
// ---------------------------------------------------------------------------

// CreateKeysAndCertificateInput carries the fields for
// CreateKeysAndCertificate.
type CreateKeysAndCertificateInput struct {
	SetAsActive bool
}

// CreateKeysAndCertificateResult is the transport-agnostic result of
// CreateKeysAndCertificate. The private key is never persisted, so it must
// travel in the result.
type CreateKeysAndCertificateResult struct {
	Certificate   *iotstore.Certificate
	PublicKeyPEM  string
	PrivateKeyPEM string
}

// RegisterCertificateInput carries the fields for RegisterCertificate. The
// status member is the current documented control; setAsActive is the
// deprecated predecessor and applies only when status is absent.
type RegisterCertificateInput struct {
	CertificatePEM   string
	CACertificatePEM string
	SetAsActive      bool
	Status           string
	StatusProvided   bool
}

// RegisterCertificateResult is the transport-agnostic result of
// RegisterCertificate.
type RegisterCertificateResult struct {
	Certificate *iotstore.Certificate
}

// CreateCertificateFromCsrInput carries the fields for
// CreateCertificateFromCsr.
type CreateCertificateFromCsrInput struct {
	CertificateSigningRequest string
	SetAsActive               bool
}

// CreateCertificateFromCsrResult is the transport-agnostic result of
// CreateCertificateFromCsr. CertificatePEM echoes the signed certificate,
// which may originate from a certificate provider rather than the internal CA.
type CreateCertificateFromCsrResult struct {
	Certificate    *iotstore.Certificate
	CertificatePEM string
}

// createKeysAndCertificateCore issues a fresh key pair and certificate from
// the region's certificate authority and persists the certificate with the
// requested initial status.
func (s *IoTService) createKeysAndCertificateCore(store iotstore.IotStoreInterface, in CreateKeysAndCertificateInput) (*CreateKeysAndCertificateResult, error) {
	ca := s.caForRegion(store.GetRegion())
	if ca == nil {
		return nil, fmt.Errorf("iot: certificate authority not available for the request region")
	}
	certPEM, keyPEM, pubKeyPEM, certID, err := ca.IssueCertificate()
	if err != nil {
		return nil, err
	}

	cert := buildCertificateRecord(certPEM, certID, in.SetAsActive)
	created, err := store.CreateCertificate(cert)
	if err != nil {
		return nil, err
	}
	return &CreateKeysAndCertificateResult{
		Certificate:   created,
		PublicKeyPEM:  pubKeyPEM,
		PrivateKeyPEM: keyPEM,
	}, nil
}

// registerCertificateCore registers an existing PEM certificate without CA
// signing, optionally linking it to a CA certificate.
func (s *IoTService) registerCertificateCore(store iotstore.IotStoreInterface, in RegisterCertificateInput) (*RegisterCertificateResult, error) {
	if in.CertificatePEM == "" {
		return nil, iotstore.ErrMissingParam
	}

	certID := vcrypto.FingerprintPEM(in.CertificatePEM)
	cert := buildCertificateRecord(in.CertificatePEM, certID, in.SetAsActive)
	if in.StatusProvided {
		if err := validateRegistrationStatus(in.Status); err != nil {
			return nil, err
		}
		cert.Status = in.Status
	}

	// When caCertificatePem is provided, link this certificate to the CA.
	if in.CACertificatePEM != "" {
		cert.CaCertificateID = vcrypto.FingerprintPEM(in.CACertificatePEM)
		cert.CertificateMode = "SNI_ONLY"
	}

	created, err := store.CreateCertificate(cert)
	if err != nil {
		return nil, err
	}
	return &RegisterCertificateResult{Certificate: created}, nil
}

// createCertificateFromCsrCore signs a CSR — via the registered certificate
// provider when one is active, otherwise the region's certificate authority —
// and persists the resulting certificate.
func (s *IoTService) createCertificateFromCsrCore(ctx context.Context, store iotstore.IotStoreInterface, in CreateCertificateFromCsrInput) (*CreateCertificateFromCsrResult, error) {
	if in.CertificateSigningRequest == "" {
		return nil, iotstore.ErrMissingParam
	}

	// If a CertificateProvider is registered for CreateCertificateFromCsr,
	// invoke its Lambda function to sign the CSR instead of the internal CA.
	// Per AWS spec, the provider fully replaces the default signing flow.
	var certPEM string
	var certID string
	if providerCertPEM, invoked, pErr := s.invokeCertProviderCore(ctx, store, in.CertificateSigningRequest); invoked {
		if pErr != nil {
			return nil, pErr
		}
		certPEM = providerCertPEM
		certID = vcrypto.FingerprintPEM(certPEM)
	} else {
		ca := s.caForRegion(store.GetRegion())
		if ca == nil {
			return nil, fmt.Errorf("iot: certificate authority not available for the request region")
		}
		var err error
		certPEM, certID, err = ca.IssueCertificateFromCSR(in.CertificateSigningRequest)
		if err != nil {
			return nil, err
		}
	}

	cert := buildCertificateRecord(certPEM, certID, in.SetAsActive)
	created, err := store.CreateCertificate(cert)
	if err != nil {
		return nil, err
	}
	return &CreateCertificateFromCsrResult{Certificate: created, CertificatePEM: certPEM}, nil
}

// registrationStatuses is the CertificateStatus enum's value set, the
// documented status domain of the certificate registration operations.
var registrationStatuses = map[string]struct{}{
	"ACTIVE": {}, "INACTIVE": {}, "REVOKED": {}, "PENDING_TRANSFER": {},
	"REGISTER_INACTIVE": {}, "PENDING_ACTIVATION": {},
}

// validateRegistrationStatus enforces the CertificateStatus enum on the
// registration operations.
func validateRegistrationStatus(status string) error {
	if _, ok := registrationStatuses[status]; !ok {
		return iotstore.ErrInvalidRequest
	}
	return nil
}

// buildCertificateRecord assembles the persisted certificate record for a
// freshly issued or registered PEM certificate.
func buildCertificateRecord(certPEM, certID string, setActive bool) *iotstore.Certificate {
	return &iotstore.Certificate{
		CertificateID:    certID,
		CertificatePEM:   certPEM,
		Status:           boolToActiveStatus(setActive),
		CertificateMode:  "DEFAULT",
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}
}
