package iam

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// UploadSigningCertificate uploads an X.509 signing certificate to the specified IAM user.
func (s *IAMService) UploadSigningCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// UserName is optional: omitting it addresses the certificate to the
	// authenticated caller.
	userName, err := resolveUserName(reqCtx, request.GetStringParam(req.Parameters, "UserName"))
	if err != nil {
		return nil, err
	}
	input := &UploadSigningCertificateInput{
		UserName:        userName,
		CertificateBody: request.GetStringParam(req.Parameters, "CertificateBody"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := s.uploadSigningCertificateCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Certificate": s.signingCertificateToResponse(created),
	}, nil
}

// ListSigningCertificates returns information about the signing certificates associated with the specified IAM user.
func (s *IAMService) ListSigningCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// UserName is optional: omitting it lists the caller's own
	// certificates.
	userName, err := resolveUserName(reqCtx, request.GetStringParam(req.Parameters, "UserName"))
	if err != nil {
		return nil, err
	}
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listSigningCertificatesCore(store, userName, marker, maxItems)
	if err != nil {
		return nil, err
	}

	certList := make([]interface{}, len(result.Certificates))
	for i, cert := range result.Certificates {
		certList[i] = s.signingCertificateToResponse(cert)
	}

	resp := map[string]interface{}{
		"Certificates": certList,
		"IsTruncated":  result.IsTruncated,
	}
	if result.NextMarker != "" {
		resp["Marker"] = result.NextMarker
	}
	return resp, nil
}

// UpdateSigningCertificate changes the status of the specified signing certificate to Active or Inactive.
func (s *IAMService) UpdateSigningCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &UpdateSigningCertificateInput{
		CertificateId: request.GetStringParam(req.Parameters, "CertificateId"),
		UserName:      request.GetStringParam(req.Parameters, "UserName"),
		Status:        request.GetStringParam(req.Parameters, "Status"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.updateSigningCertificateCore(reqCtx, store, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteSigningCertificate deletes a signing certificate associated with the specified IAM user.
func (s *IAMService) DeleteSigningCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certificateId := request.GetStringParam(req.Parameters, "CertificateId")
	rawUserName := request.GetStringParam(req.Parameters, "UserName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteSigningCertificateCore(reqCtx, store, certificateId, rawUserName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *IAMService) signingCertificateToResponse(cert *iamstore.SigningCertificate) map[string]interface{} {
	return map[string]interface{}{
		"CertificateId":   cert.CertificateId,
		"UserName":        cert.UserName,
		"CertificateBody": cert.CertificateBody,
		"Status":          cert.Status,
		"UploadDate":      cert.UploadDate.Format(timeutils.ISO8601SimpleFormat),
	}
}
