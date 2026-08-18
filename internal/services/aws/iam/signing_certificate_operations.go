package iam

import (
	"context"
	"errors"

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
	certificateBody := request.GetStringParam(req.Parameters, "CertificateBody")
	if certificateBody == "" {
		return nil, NewValidationError("CertificateBody")
	}
	if len(certificateBody) > 16384 {
		return nil, NewInvalidInputError("CertificateBody", "must be 1 to 16384 characters")
	}

	cert, err := parseCertificate(certificateBody)
	if err != nil {
		return nil, ErrMalformedCertificate
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	created, err := store.SigningCertificates().UploadWithGuards(userName, certificateBody, certificateFingerprint(cert))
	if err != nil {
		if errors.Is(err, iamstore.ErrDuplicateSigningCertificate) {
			return nil, ErrDuplicateCertificate
		}
		if errors.Is(err, iamstore.ErrSigningCertificateLimitExceeded) {
			return nil, ErrLimitExceededSigningCertificates
		}
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

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	paged := pagination.PaginateSlice(certList, marker, maxItems, func(item interface{}) string {
		if m, ok := item.(map[string]interface{}); ok {
			if id, ok := m["CertificateId"].(string); ok {
				return id
			}
		}
		return ""
	})

	resp := map[string]interface{}{
		"Certificates": paged.Items,
		"IsTruncated":  paged.IsTruncated,
	}
	if paged.NextMarker != "" {
		resp["Marker"] = paged.NextMarker
	}
	return resp, nil
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

	// UserName is optional: when omitted, the caller's own user name is
	// determined implicitly from the access key that signed the request.
	owner, err := resolveUserName(reqCtx, request.GetStringParam(req.Parameters, "UserName"))
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := store.SigningCertificates().Get(certificateId)
	if err != nil {
		return nil, NewNoSuchEntityError("signing certificate", certificateId)
	}
	// Any resolved user name other than the certificate owner fails with
	// NoSuchEntity rather than mutating the certificate.
	if cert.UserName != owner {
		return nil, NewNoSuchEntityError("signing certificate", certificateId)
	}

	if err := store.SigningCertificates().UpdateStatus(certificateId, status); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteSigningCertificate deletes a signing certificate associated with the specified IAM user.
func (s *IAMService) DeleteSigningCertificate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	certificateId := request.GetStringParam(req.Parameters, "CertificateId")
	if certificateId == "" {
		return nil, NewValidationError("CertificateId")
	}
	// UserName is optional: when omitted, the caller's own user name is
	// determined implicitly from the access key that signed the request.
	userName, err := resolveUserName(reqCtx, request.GetStringParam(req.Parameters, "UserName"))
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	cert, err := store.SigningCertificates().Get(certificateId)
	if err != nil {
		return nil, NewNoSuchEntityError("signing certificate", certificateId)
	}
	// Any resolved user name other than the certificate owner fails with
	// NoSuchEntity instead of deleting the certificate.
	if cert.UserName != userName {
		return nil, NewNoSuchEntityError("signing certificate", certificateId)
	}

	if err := store.SigningCertificates().Delete(certificateId); err != nil {
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
