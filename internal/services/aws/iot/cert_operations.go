package iot

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/pem"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// CreateKeysAndCertificate generates a new key pair and certificate, storing
// the result with ACTIVE status. Returns the certificate ARN, ID, and key pair.
func (s *IoTService) CreateKeysAndCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	certPEM, keyPEM, pubKeyPEM, certID, err := s.ca.IssueCertificate()
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	setActive := request.GetBoolParam(req.Parameters, "setAsActive")
	status := "INACTIVE"
	if setActive {
		status = "ACTIVE"
	}

	cert := &iotstore.Certificate{
		CertificateID:    certID,
		CertificatePEM:   certPEM,
		Status:           status,
		CertificateMode:  "DEFAULT",
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}

	created, err := store.CreateCertificate(cert)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
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

	return map[string]interface{}{
		"certificateDescription": map[string]interface{}{
			"certificateArn":   cert.CertificateARN,
			"certificateId":    cert.CertificateID,
			"certificatePem":   cert.CertificatePEM,
			"status":           cert.Status,
			"certificateMode":  cert.CertificateMode,
			"creationDate":     cert.CreationDate.Unix(),
			"lastModifiedDate": cert.LastModifiedDate.Unix(),
			"caCertificateId":  cert.CaCertificateID,
		},
	}, nil
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

	cert, err := store.GetCertificate(certID)
	if err != nil {
		return nil, iotstore.ErrCertificateNotFound
	}

	cert.Status = newStatus
	cert.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateCertificate(cert); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{}, nil
}

// DeleteCertificate removes a certificate that has no attached entities.
// Returns InvalidRequestException if the certificate still has policies or
// things attached.
func (s *IoTService) DeleteCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	if cert.Status != "INACTIVE" {
		return nil, iotstore.ErrInvalidCertStatus
	}

	// Policies may be attached by either certID or certARN; check both to
	// prevent dangling references in policyAttachBase.
	certARN := cert.CertificateARN
	policiesByID, _ := store.ListPoliciesForPrincipal(certID)
	policiesByARN, _ := store.ListPoliciesForPrincipal(certARN)
	if len(policiesByID) > 0 || len(policiesByARN) > 0 {
		return nil, iotstore.ErrCertHasAttachments
	}

	thingsByID, _ := store.ListThingsForPrincipal(certID)
	thingsByARN, _ := store.ListThingsForPrincipal(certARN)
	if len(thingsByID) > 0 || len(thingsByARN) > 0 {
		return nil, iotstore.ErrCertHasAttachments
	}

	if err := store.DeleteCertificate(certID); err != nil {
		return nil, iotstore.ErrCertificateNotFound
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
		items = append(items, map[string]interface{}{
			"certificateArn":  c.CertificateARN,
			"certificateId":   c.CertificateID,
			"status":          c.Status,
			"certificateMode": c.CertificateMode,
			"creationDate":    c.CreationDate.Unix(),
		})
	}

	resp := map[string]interface{}{
		"certificates": items,
	}
	if certs.NextMarker != "" {
		resp["nextToken"] = certs.NextMarker
	}
	return resp, nil
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

	certID := computeCertID(certPEM)

	status := "INACTIVE"
	if request.GetBoolParam(req.Parameters, "setAsActive") {
		status = "ACTIVE"
	}

	cert := &iotstore.Certificate{
		CertificateID:    certID,
		CertificatePEM:   certPEM,
		Status:           status,
		CertificateMode:  "DEFAULT",
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}

	created, err := store.CreateCertificate(cert)
	if err != nil {
		return nil, iotstore.ErrCertificateAlreadyExists
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

	certPEM, certID, err := s.ca.IssueCertificateFromCSR(csrPEM)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	status := "INACTIVE"
	if request.GetBoolParam(req.Parameters, "setAsActive") {
		status = "ACTIVE"
	}

	cert := &iotstore.Certificate{
		CertificateID:    certID,
		CertificatePEM:   certPEM,
		Status:           status,
		CertificateMode:  "DEFAULT",
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}

	created, err := store.CreateCertificate(cert)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"certificateArn": created.CertificateARN,
		"certificateId":  created.CertificateID,
		"certificatePem": certPEM,
	}, nil
}

func isValidCertStatus(status string) bool {
	switch status {
	case "ACTIVE", "INACTIVE", "REVOKED", "PENDING_TRANSFER", "PENDING_ACTIVATION":
		return true
	}
	return false
}

func computeCertID(certPEM string) string {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		h := sha256.Sum256([]byte(certPEM))
		return hex.EncodeToString(h[:])
	}
	h := sha256.Sum256(block.Bytes)
	return hex.EncodeToString(h[:])
}
