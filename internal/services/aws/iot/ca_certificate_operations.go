package iot

import (
	"context"
	"encoding/json"

	"vorpalstacks/internal/common/request"
)

// registrationConfigJSON serialises the wire registrationConfig structure
// for storage; empty when the member is absent.
func registrationConfigJSON(params map[string]interface{}) string {
	m := request.GetMapParamCaseInsensitive(params, "registrationConfig")
	if m == nil {
		return ""
	}
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func (s *IoTService) RegisterCACertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := RegisterCACertificateInput{
		CACertificatePEM:           request.GetParamCaseInsensitive(req.Parameters, "caCertificate"),
		VerificationCertificatePEM: request.GetParamCaseInsensitive(req.Parameters, "verificationCertificate"),
		RegistrationConfig:         registrationConfigJSON(req.Parameters),
		CertificateMode:            request.GetParamCaseInsensitive(req.Parameters, "certificateMode"),
		Tags:                       tagListParam(req.Parameters),
	}
	if hasParam(req.Parameters, "setAsActive") {
		in.SetAsActive = request.GetBoolParam(req.Parameters, "setAsActive")
	}
	if hasParam(req.Parameters, "allowAutoRegistration") {
		in.AllowAutoRegistration = request.GetBoolParam(req.Parameters, "allowAutoRegistration")
	}
	result, err := s.registerCACertificateCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"certificateArn": result.CertificateARN,
		"certificateId":  result.CertificateID,
	}, nil
}

func (s *IoTService) DescribeCACertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeCACertificateCore(store, request.GetParamCaseInsensitive(req.Parameters, "certificateId"))
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"certificateDescription": result.CertificateDescription,
	}
	if result.RegistrationConfig != nil {
		resp["registrationConfig"] = result.RegistrationConfig
	}
	return resp, nil
}

func (s *IoTService) ListCACertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listCACertificatesCore(store)
	if err != nil {
		return nil, err
	}

	return paginatedMaps("certificates", result.Certificates, req.Parameters)
}

func (s *IoTService) UpdateCACertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := UpdateCACertificateInput{
		CertificateID:             request.GetParamCaseInsensitive(req.Parameters, "certificateId"),
		NewStatus:                 request.GetParamCaseInsensitive(req.Parameters, "newStatus"),
		NewAutoRegistrationStatus: request.GetParamCaseInsensitive(req.Parameters, "newAutoRegistrationStatus"),
		RegistrationConfig:        registrationConfigJSON(req.Parameters),
	}
	if hasParam(req.Parameters, "removeAutoRegistration") {
		in.RemoveAutoRegistration = request.GetBoolParam(req.Parameters, "removeAutoRegistration")
	}
	if err := s.updateCACertificateCore(store, in); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DeleteCACertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteCACertificateCore(store, request.GetParamCaseInsensitive(req.Parameters, "certificateId")); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListCertificatesByCA(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	result, err := s.listCertificatesByCACore(store, ListCertificatesByCACertificateInput{
		CACertificateID: request.GetParamCaseInsensitive(req.Parameters, "caCertificateId"),
		Marker:          opts.Marker,
		MaxItems:        opts.MaxItems,
	})
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0)
	for _, c := range result.Certificates {
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
	return listResponse("certificates", items, result.NextMarker), nil
}

func (s *IoTService) RegisterCertificateWithoutCA(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.registerCertificateWithoutCACore(store, RegisterCertificateWithoutCAInput{
		CertificatePEM: request.GetParamCaseInsensitive(req.Parameters, "certificatePem"),
		Status:         request.GetParamCaseInsensitive(req.Parameters, "status"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"certificateArn": result.Certificate.CertificateARN,
		"certificateId":  result.Certificate.CertificateID,
	}, nil
}

func (s *IoTService) TransferCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.transferCertificateCore(store, TransferCertificateInput{
		CertificateID:    request.GetParamCaseInsensitive(req.Parameters, "certificateId"),
		TargetAWSAccount: request.GetParamCaseInsensitive(req.Parameters, "targetAwsAccount"),
		TransferMessage:  request.GetParamCaseInsensitive(req.Parameters, "transferMessage"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"transferredCertificateArn": result.TransferredCertificateARN,
	}, nil
}

func (s *IoTService) AcceptCertificateTransfer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// AWS: default status is INACTIVE; ACTIVE only when setAsActive=true.
	certStatus := "INACTIVE"
	if request.GetBoolParam(req.Parameters, "setAsActive") {
		certStatus = "ACTIVE"
	}
	return s.completeCertTransfer(reqCtx, request.GetParamCaseInsensitive(req.Parameters, "certificateId"), certStatus)
}

func (s *IoTService) CancelCertificateTransfer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// AWS: cancelled transfer leaves certificate INACTIVE.
	return s.completeCertTransfer(reqCtx, request.GetParamCaseInsensitive(req.Parameters, "certificateId"), "INACTIVE")
}

func (s *IoTService) RejectCertificateTransfer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.completeCertTransfer(reqCtx, request.GetParamCaseInsensitive(req.Parameters, "certificateId"), "INACTIVE")
}

func (s *IoTService) ListOutgoingCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listOutgoingCertificatesCore(store)
	if err != nil {
		return nil, err
	}

	return paginatedMaps("outgoingCertificates", result.OutgoingCertificates, req.Parameters)
}

// completeCertTransfer resolves a pending transfer via the shared Core with
// the terminal status chosen by the caller.
func (s *IoTService) completeCertTransfer(reqCtx *request.RequestContext, certID, certStatus string) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.completeCertTransferCore(store, certID, certStatus); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
