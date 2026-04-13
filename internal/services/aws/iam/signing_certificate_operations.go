package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// UploadSigningCertificate uploads an X.509 signing certificate to the specified IAM user.
func (s *IAMService) UploadSigningCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, ErrNoSuchUser
	}
	certificateBody := request.GetStringParam(req.Parameters, "CertificateBody")
	if certificateBody == "" {
		return nil, NewValidationError("CertificateBody")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	cert, err := store.SigningCertificates().Upload(userName, certificateBody)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Certificate": s.signingCertificateToResponse(cert),
	}, nil
}

// ListSigningCertificates returns information about the signing certificates associated with the specified IAM user.
func (s *IAMService) ListSigningCertificates(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, ErrNoSuchUser
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	certs, err := store.SigningCertificates().ListByUserName(userName)
	if err != nil {
		return nil, err
	}

	certList := make([]interface{}, len(certs))
	for i, cert := range certs {
		certList[i] = s.signingCertificateToResponse(cert)
	}

	return map[string]interface{}{
		"Certificates": certList,
		"IsTruncated":  false,
	}, nil
}

// UpdateSigningCertificate changes the status of the specified signing certificate to Active or Inactive.
func (s *IAMService) UpdateSigningCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certificateId := request.GetStringParam(req.Parameters, "CertificateId")
	if certificateId == "" {
		return nil, NewValidationError("CertificateId")
	}
	status := request.GetStringParam(req.Parameters, "Status")
	if status == "" {
		return nil, NewValidationError("Status")
	}
	if status != "Active" && status != "Inactive" {
		return nil, NewInvalidInputError("Status", "must be Active or Inactive")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.SigningCertificates().Exists(certificateId) {
		return nil, NewNoSuchEntityError("signing certificate", certificateId)
	}

	if err := store.SigningCertificates().UpdateStatus(certificateId, status); err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteSigningCertificate deletes a signing certificate associated with the specified IAM user.
func (s *IAMService) DeleteSigningCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certificateId := request.GetStringParam(req.Parameters, "CertificateId")
	if certificateId == "" {
		return nil, NewValidationError("CertificateId")
	}
	userName := request.GetStringParam(req.Parameters, "UserName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.SigningCertificates().Exists(certificateId) {
		return nil, NewNoSuchEntityError("signing certificate", certificateId)
	}

	if userName != "" {
		if !store.Users().Exists(userName) {
			return nil, NewNoSuchUserError(userName)
		}
	}

	if err := store.SigningCertificates().Delete(certificateId); err != nil {
		return nil, err
	}

	return nil, nil
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
