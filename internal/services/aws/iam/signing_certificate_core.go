// Transport-agnostic Core functions for IAM user signing certificates:
// validation and store operations shared by the AWS-compatible HTTP API
// handlers and any admin plane paths (the xxxCore pattern).
package iam

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// UploadSigningCertificateInput holds the parameters for uploading a
// signing certificate.
type UploadSigningCertificateInput struct {
	UserName        string
	CertificateBody string
}

// UpdateSigningCertificateInput holds the parameters for updating a signing
// certificate status.  The raw UserName is resolved against the request
// context inside the core because AWS resolves it after the member
// validation.
type UpdateSigningCertificateInput struct {
	CertificateId string
	UserName      string
	Status        string
}

// SigningCertificateListResult holds a paginated signing certificate
// listing.
type SigningCertificateListResult struct {
	Certificates []*iamstore.SigningCertificate
	IsTruncated  bool
	NextMarker   string
}

// uploadSigningCertificateCore validates input and uploads an X.509 signing
// certificate for the specified IAM user.
func (s *IAMService) uploadSigningCertificateCore(store *iamstore.IAMStore, input *UploadSigningCertificateInput) (*iamstore.SigningCertificate, error) {
	certificateBody := input.CertificateBody
	if certificateBody == "" {
		return nil, NewValidationError("CertificateBody")
	}
	// certificateBodyType carries a Latin-1 pattern, so lengths count
	// Unicode characters.
	if utf8.RuneCountInString(certificateBody) > maxCertificateBodyLength {
		return nil, NewInvalidInputError("CertificateBody", fmt.Sprintf("must be 1 to %d characters", maxCertificateBodyLength))
	}

	cert, err := parseCertificate(certificateBody)
	if err != nil {
		return nil, ErrMalformedCertificate
	}

	if !store.Users().Exists(input.UserName) {
		return nil, NewNoSuchUserError(input.UserName)
	}

	created, err := store.SigningCertificates().UploadWithGuards(input.UserName, certificateBody, certificateFingerprint(cert))
	if err != nil {
		if errors.Is(err, iamstore.ErrDuplicateSigningCertificate) {
			return nil, ErrDuplicateCertificate
		}
		if errors.Is(err, iamstore.ErrSigningCertificateLimitExceeded) {
			return nil, ErrLimitExceededSigningCertificates
		}
		return nil, err
	}
	return created, nil
}

// listSigningCertificatesCore returns a paginated list of the signing
// certificates associated with the specified IAM user.
func (s *IAMService) listSigningCertificatesCore(store *iamstore.IAMStore, userName, marker string, maxItems int) (*SigningCertificateListResult, error) {
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	certs, err := store.SigningCertificates().ListByUserName(userName)
	if err != nil {
		return nil, err
	}

	paged := pagination.PaginateSlice(certs, marker, maxItems, func(c *iamstore.SigningCertificate) string {
		return c.CertificateId
	})

	return &SigningCertificateListResult{
		Certificates: paged.Items,
		IsTruncated:  paged.IsTruncated,
		NextMarker:   paged.NextMarker,
	}, nil
}

// updateSigningCertificateCore validates input and changes the status of
// the specified signing certificate to Active or Inactive.  The raw
// UserName is optional: when omitted, the caller's own user name is
// determined implicitly from the access key that signed the request, after
// the member validation.  Any resolved user name other than the certificate
// owner fails with NoSuchEntity rather than mutating the certificate.
func (s *IAMService) updateSigningCertificateCore(reqCtx *request.RequestContext, store *iamstore.IAMStore, input *UpdateSigningCertificateInput) error {
	if input.CertificateId == "" {
		return NewValidationError("CertificateId")
	}
	status := input.Status
	if status == "" {
		return NewValidationError("Status")
	}
	if status != "Active" && status != "Inactive" {
		return NewInvalidInputError("Status", "must be Active or Inactive")
	}

	owner, err := resolveUserName(reqCtx, input.UserName)
	if err != nil {
		return err
	}

	cert, err := store.SigningCertificates().Get(input.CertificateId)
	if err != nil {
		return NewNoSuchEntityError("signing certificate", input.CertificateId)
	}
	if cert.UserName != owner {
		return NewNoSuchEntityError("signing certificate", input.CertificateId)
	}

	return store.SigningCertificates().UpdateStatus(input.CertificateId, status)
}

// deleteSigningCertificateCore validates input and deletes a signing
// certificate associated with the specified IAM user.  The raw UserName is
// optional and resolved against the request context; any resolved user
// name other than the certificate owner fails with NoSuchEntity instead of
// deleting the certificate.
func (s *IAMService) deleteSigningCertificateCore(reqCtx *request.RequestContext, store *iamstore.IAMStore, certificateId, rawUserName string) error {
	if certificateId == "" {
		return NewValidationError("CertificateId")
	}

	userName, err := resolveUserName(reqCtx, rawUserName)
	if err != nil {
		return err
	}

	cert, err := store.SigningCertificates().Get(certificateId)
	if err != nil {
		return NewNoSuchEntityError("signing certificate", certificateId)
	}
	if cert.UserName != userName {
		return NewNoSuchEntityError("signing certificate", certificateId)
	}

	return store.SigningCertificates().Delete(certificateId)
}
