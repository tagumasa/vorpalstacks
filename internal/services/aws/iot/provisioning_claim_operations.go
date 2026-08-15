package iot

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// CreateProvisioningClaim issues a temporary provisioning claim bound to an
// existing provisioning template. AWS uses the claim for just-in-time
// provisioning during device manufacturing. The claim consists of a
// short-lived self-signed X.509 certificate and its private key.
func (s *IoTService) CreateProvisioningClaim(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	templateName := request.GetParamCaseInsensitive(req.Parameters, "templateName")
	if templateName == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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

	return map[string]interface{}{
		"certificateId":  certID,
		"certificatePem": certPEM,
		"keyPair": map[string]interface{}{
			"PublicKey":  pubKeyPEM,
			"PrivateKey": privKeyPEM,
		},
		"expiration": now.Add(time.Hour).UTC().Unix(),
	}, nil
}

// ---------------------------------------------------------------------------
// Registration code operations.
// AWS IoT generates a per-account registration code used to register CA
// certificates. The code is lazily created on first GetRegistrationCode call
// and persists until explicitly deleted.
// ---------------------------------------------------------------------------

func (s *IoTService) GetRegistrationCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
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

func (s *IoTService) DeleteRegistrationCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteGeneric("config/registrationCode"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
