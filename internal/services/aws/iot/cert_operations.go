package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/services/aws/iot/ca"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// CreateKeysAndCertificate generates a new key pair and certificate, storing
// the result with ACTIVE status. Returns the certificate ARN, ID, and key pair.
func (s *IoTService) CreateKeysAndCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	certPEM, keyPEM, pubKeyPEM, certID, err := s.deps.CA.IssueCertificate()
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
	certID := request.GetParamCaseInsensitive(req.Parameters, "certificateId")
	if certID == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cert, err := store.GetCertificate(certID)
	if err != nil {
		return nil, iotstore.ErrCertificateNotFound
	}

	return certificateDetailResponse(cert), nil
}

// UpdateCertificate changes the status of a certificate (ACTIVE, INACTIVE, REVOKED).
func (s *IoTService) UpdateCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certID := request.GetParamCaseInsensitive(req.Parameters, "certificateId")
	newStatus := request.GetParamCaseInsensitive(req.Parameters, "newStatus")
	if certID == "" || newStatus == "" {
		return nil, iotstore.ErrMissingParam
	}

	if !isValidCertStatus(newStatus) {
		return nil, iotstore.ErrInvalidCertStatus
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := iotstore.CertificateUpdateOpts{NewStatus: newStatus}
	_, err = store.UpdateCertificate(certID, opts)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// DeleteCertificate removes a certificate that has no attached entities.
func (s *IoTService) DeleteCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certID := request.GetParamCaseInsensitive(req.Parameters, "certificateId")
	if certID == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteCertificate(certID); err != nil {
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

	certs, err := store.ListCertificates(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(certs.Items))
	for _, c := range certs.Items {
		items = append(items, certificateResponse(c))
	}

	return listResponse("certificates", items, certs.NextMarker), nil
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

	certID := ca.ComputeCertID(certPEM)

	cert := buildCertificateRecord(certPEM, certID, request.GetBoolParam(req.Parameters, "setAsActive"))

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

	certPEM, certID, err := s.deps.CA.IssueCertificateFromCSR(csrPEM)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
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
		"certificatePem": certPEM,
	}, nil
}

func isValidCertStatus(status string) bool {
	switch status {
	case "ACTIVE", "INACTIVE", "REVOKED", "PENDING_TRANSFER", "REGISTER_INACTIVE", "PENDING_ACTIVATION":
		return true
	}
	return false
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
