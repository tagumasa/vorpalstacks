package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
	vcrypto "vorpalstacks/internal/utils/crypto"
)

func (s *IoTService) RegisterCACertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	caCertPem := request.GetParamCaseInsensitive(req.Parameters, "caCertificate")
	if caCertPem == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Deterministic ID from certificate content hash (prevents duplicate
	// registration of the same CA certificate, matching AWS behaviour).
	caCertID := vcrypto.FingerprintPEM(caCertPem)
	verificationCert := request.GetParamCaseInsensitive(req.Parameters, "verificationCertificate")
	rec := map[string]interface{}{
		"certificateId":              caCertID,
		"status":                     "ACTIVE",
		"caCertificatePem":           caCertPem,
		"verificationCertificatePem": verificationCert,
		"creationDate":               time.Now().Unix(),
	}
	if err := store.PutGeneric("caCert/"+caCertID, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"certificateArn": iotstore.BuildCACertificateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), caCertID),
		"certificateId":  caCertID,
	}, nil
}

func (s *IoTService) DescribeCACertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	caCertID := request.GetParamCaseInsensitive(req.Parameters, "certificateId")
	if caCertID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("caCert/"+caCertID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCertificateNotFound
	}
	return map[string]interface{}{
		"certificateDescription": rec,
	}, nil
}

func (s *IoTService) ListCACertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("caCert/")
	if err != nil {
		return nil, err
	}
	certs := make([]map[string]interface{}, 0, len(items))
	for _, c := range items {
		id, _ := c["certificateId"].(string)
		certs = append(certs, map[string]interface{}{
			"certificateArn":  iotstore.BuildCACertificateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), id),
			"certificateId":   c["certificateId"],
			"status":          c["status"],
			"caCertificateId": c["certificateId"],
		})
	}
	return paginatedMaps("certificates", certs, req.Parameters), nil
}

func (s *IoTService) UpdateCACertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	caCertID := request.GetParamCaseInsensitive(req.Parameters, "certificateId")
	if caCertID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	caExists, err := store.GetGenericExists("caCert/"+caCertID, &rec)
	if err != nil {
		return nil, err
	}
	if !caExists {
		return nil, iotstore.ErrCertificateNotFound
	}
	if newStatus := request.GetParamCaseInsensitive(req.Parameters, "newStatus"); newStatus != "" {
		rec["status"] = newStatus
		if err := store.PutGeneric("caCert/"+caCertID, rec); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) DeleteCACertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	caCertID := request.GetParamCaseInsensitive(req.Parameters, "certificateId")
	if caCertID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("caCert/"+caCertID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCertificateNotFound
	}
	// AWS requires the CA certificate to be INACTIVE before deletion.
	if status, _ := rec["status"].(string); status == "ACTIVE" {
		return nil, iotstore.ErrCertHasAttachments
	}
	if err := store.DeleteGeneric("caCert/" + caCertID); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ListCertificatesByCA(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	caCertID := request.GetParamCaseInsensitive(req.Parameters, "caCertificateId")
	if caCertID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	certs, err := store.ListCertificates(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0)
	for _, c := range certs.Items {
		if c.CaCertificateID == caCertID {
			items = append(items, map[string]interface{}{
				"certificateArn":   c.CertificateARN,
				"certificateId":    c.CertificateID,
				"certificateMode":  c.CertificateMode,
				"certificatePem":   c.CertificatePEM,
				"status":           c.Status,
				"creationDate":     c.CreationDate.Unix(),
				"lastModifiedDate": c.LastModifiedDate.Unix(),
			})
		}
	}
	return listResponse("certificates", items, certs.NextMarker), nil
}

func (s *IoTService) RegisterCertificateWithoutCA(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certPem := request.GetParamCaseInsensitive(req.Parameters, "certificatePem")
	if certPem == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	certID := vcrypto.FingerprintPEM(certPem)
	status := request.GetParamCaseInsensitive(req.Parameters, "status")
	if status == "" {
		status = "INACTIVE"
	}
	cert := &iotstore.Certificate{
		CertificateID:    certID,
		CertificatePEM:   certPem,
		Status:           status,
		CertificateMode:  "DEFAULT",
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
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

func (s *IoTService) TransferCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certID := request.GetParamCaseInsensitive(req.Parameters, "certificateId")
	targetAccount := request.GetParamCaseInsensitive(req.Parameters, "targetAwsAccount")
	if certID == "" || targetAccount == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// AWS rejects TransferCertificate on a non-existent certificate with
	// ResourceNotFoundException rather than silently creating a transfer record.
	if _, err := store.GetCertificate(certID); err != nil {
		return nil, err
	}
	// Move the certificate into PENDING_TRANSFER so DescribeCertificate
	// reflects the transfer, matching AWS behaviour.
	if _, err := store.UpdateCertificate(certID, iotstore.CertificateUpdateOpts{NewStatus: "PENDING_TRANSFER"}); err != nil {
		return nil, err
	}
	rec := map[string]interface{}{
		"certificateId":    certID,
		"status":           "PENDING_ACCEPTANCE",
		"targetAwsAccount": targetAccount,
	}
	if err := store.PutGeneric("certTransfer/"+certID, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"transferredCertificateArn": iotstore.BuildCertificateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), certID),
	}, nil
}

func (s *IoTService) AcceptCertificateTransfer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// AWS: default status is INACTIVE; ACTIVE only when setAsActive=true.
	certStatus := "INACTIVE"
	if request.GetBoolParam(req.Parameters, "setAsActive") {
		certStatus = "ACTIVE"
	}
	return s.completeCertTransfer(reqCtx, req, certStatus)
}

func (s *IoTService) CancelCertificateTransfer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// AWS: cancelled transfer leaves certificate INACTIVE.
	return s.completeCertTransfer(reqCtx, req, "INACTIVE")
}

func (s *IoTService) RejectCertificateTransfer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.completeCertTransfer(reqCtx, req, "INACTIVE")
}

func (s *IoTService) ListOutgoingCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("certTransfer/")
	if err != nil {
		return nil, err
	}
	certs := make([]map[string]interface{}, 0, len(items))
	for _, c := range items {
		status, _ := c["status"].(string)
		if status != "PENDING_ACCEPTANCE" {
			continue
		}
		certID, _ := c["certificateId"].(string)
		certs = append(certs, map[string]interface{}{
			"certificateArn": iotstore.BuildCertificateARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), certID),
			"certificateId":  certID,
			"transferStatus": status,
		})
	}
	return paginatedMaps("outgoingCertificates", certs, req.Parameters), nil
}

// completeCertTransfer resolves a pending certificate transfer: it sets the
// certificate to the terminal status (ACTIVE or INACTIVE per the caller) and
// removes the transfer record. A missing transfer record yields
// ResourceNotFoundException, matching AWS.
func (s *IoTService) completeCertTransfer(reqCtx *request.RequestContext, req *request.ParsedRequest, certStatus string) (interface{}, error) {
	certID := request.GetParamCaseInsensitive(req.Parameters, "certificateId")
	if certID == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("certTransfer/"+certID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCertificateNotFound
	}
	// Transition the certificate status and clear the transfer record.
	if _, err := store.UpdateCertificate(certID, iotstore.CertificateUpdateOpts{NewStatus: certStatus}); err != nil {
		return nil, err
	}
	if err := store.DeleteGeneric("certTransfer/" + certID); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
