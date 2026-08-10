package iot

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

// CreateKeysAndCertificate generates a new key pair and certificate, storing
// the result with ACTIVE status. Returns the certificate ARN, ID, and key pair.
func (s *IoTService) CreateKeysAndCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ca := s.caForReq(reqCtx)
	if ca == nil {
		return nil, fmt.Errorf("iot: certificate authority not available for the request region")
	}
	certPEM, keyPEM, pubKeyPEM, certID, err := ca.IssueCertificate()
	if err != nil {
		return nil, err
	}

	cert := buildCertificateRecord(certPEM, certID, request.GetBoolParam(req.Parameters, "setAsActive"))

	created, err := store.CreateCertificate(cert)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"certificateArn": created.CertificateARN,
		"certificateId":  created.CertificateID,
		"certificatePem": created.CertificatePEM,
		"status":         created.Status,
		"keyPair": map[string]interface{}{
			"PublicKey":  pubKeyPEM,
			"PrivateKey": keyPEM,
		},
		"certificateMode":  created.CertificateMode,
		"creationDate":     created.CreationDate.Unix(),
		"lastModifiedDate": created.LastModifiedDate.Unix(),
	}, nil
}

// DescribeCertificate retrieves the details of a certificate by its ID.
func (s *IoTService) DescribeCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeCertificateCore(store, request.GetParamCaseInsensitive(req.Parameters, "certificateId"))
	if err != nil {
		return nil, err
	}

	return certificateDetailResponse(result.Certificate), nil
}

// UpdateCertificate changes the status of a certificate (ACTIVE, INACTIVE, REVOKED).
func (s *IoTService) UpdateCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.updateCertificateCore(store, UpdateCertificateInput{
		CertificateID: request.GetParamCaseInsensitive(req.Parameters, "certificateId"),
		NewStatus:     request.GetParamCaseInsensitive(req.Parameters, "newStatus"),
	}); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// DeleteCertificate removes a certificate that has no attached entities.
func (s *IoTService) DeleteCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteCertificateCore(store, request.GetParamCaseInsensitive(req.Parameters, "certificateId"), reqCtx.GetAccountID(), reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListCertificates returns all registered certificates.
func (s *IoTService) ListCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	result, err := s.listCertificatesCore(store, opts.Marker, opts.MaxItems)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(result.Certificates))
	for _, c := range result.Certificates {
		items = append(items, certificateResponse(c))
	}

	return listResponse("certificates", items, result.NextToken), nil
}

// RegisterCertificate registers an existing PEM certificate without CA signing.
func (s *IoTService) RegisterCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certPEM := request.GetParamCaseInsensitive(req.Parameters, "certificatePem")
	if certPEM == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	certID := vcrypto.FingerprintPEM(certPEM)

	cert := buildCertificateRecord(certPEM, certID, request.GetBoolParam(req.Parameters, "setAsActive"))

	// When caCertificatePem is provided, link this certificate to the CA.
	if caCertPem := request.GetParamCaseInsensitive(req.Parameters, "caCertificatePem"); caCertPem != "" {
		cert.CaCertificateID = vcrypto.FingerprintPEM(caCertPem)
		cert.CertificateMode = "SNI_ONLY"
	}

	created, err := store.CreateCertificate(cert)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"certificateArn": created.CertificateARN,
		"certificateId":  created.CertificateID,
	}, nil
}

// CreateCertificateFromCsr signs a CSR and returns the resulting certificate.
func (s *IoTService) CreateCertificateFromCsr(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	csrPEM := request.GetParamCaseInsensitive(req.Parameters, "certificateSigningRequest")
	if csrPEM == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// If a CertificateProvider is registered for CreateCertificateFromCsr,
	// invoke its Lambda function to sign the CSR instead of the internal CA.
	// Per AWS spec, the provider fully replaces the default signing flow.
	var certPEM string
	var certID string
	if providerCertPEM, invoked, pErr := s.tryInvokeCertProvider(ctx, reqCtx, store, csrPEM); invoked {
		if pErr != nil {
			return nil, pErr
		}
		certPEM = providerCertPEM
		certID = vcrypto.FingerprintPEM(certPEM)
	} else {
		ca := s.caForReq(reqCtx)
		if ca == nil {
			return nil, fmt.Errorf("iot: certificate authority not available for the request region")
		}
		certPEM, certID, err = ca.IssueCertificateFromCSR(csrPEM)
		if err != nil {
			return nil, err
		}
	}

	cert := buildCertificateRecord(certPEM, certID, request.GetBoolParam(req.Parameters, "setAsActive"))

	created, err := store.CreateCertificate(cert)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"certificateArn": created.CertificateARN,
		"certificateId":  created.CertificateID,
		"certificatePem": certPEM,
	}, nil
}

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
