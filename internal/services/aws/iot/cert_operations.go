package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateKeysAndCertificate generates a new key pair and certificate, storing
// the result with ACTIVE status. Returns the certificate ARN, ID, and key pair.
func (s *IoTService) CreateKeysAndCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createKeysAndCertificateCore(store, CreateKeysAndCertificateInput{
		SetAsActive: request.GetBoolParam(req.Parameters, "setAsActive"),
	})
	if err != nil {
		return nil, err
	}

	created := result.Certificate
	return map[string]interface{}{
		"certificateArn": created.CertificateARN,
		"certificateId":  created.CertificateID,
		"certificatePem": created.CertificatePEM,
		"status":         created.Status,
		"keyPair": map[string]interface{}{
			"PublicKey":  result.PublicKeyPEM,
			"PrivateKey": result.PrivateKeyPEM,
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	status := request.GetParamCaseInsensitive(req.Parameters, "status")
	result, err := s.registerCertificateCore(store, RegisterCertificateInput{
		CertificatePEM:   request.GetParamCaseInsensitive(req.Parameters, "certificatePem"),
		CACertificatePEM: request.GetParamCaseInsensitive(req.Parameters, "caCertificatePem"),
		SetAsActive:      request.GetBoolParam(req.Parameters, "setAsActive"),
		Status:           status,
		StatusProvided:   status != "",
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"certificateArn": result.Certificate.CertificateARN,
		"certificateId":  result.Certificate.CertificateID,
	}, nil
}

// CreateCertificateFromCsr signs a CSR and returns the resulting certificate.
func (s *IoTService) CreateCertificateFromCsr(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createCertificateFromCsrCore(ctx, store, CreateCertificateFromCsrInput{
		CertificateSigningRequest: request.GetParamCaseInsensitive(req.Parameters, "certificateSigningRequest"),
		SetAsActive:               request.GetBoolParam(req.Parameters, "setAsActive"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"certificateArn": result.Certificate.CertificateARN,
		"certificateId":  result.Certificate.CertificateID,
		"certificatePem": result.CertificatePEM,
	}, nil
}
