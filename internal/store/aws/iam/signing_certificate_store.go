package iam

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const signingCertificateBucketName = "iam_signing_certificates"

// SigningCertificateStore provides storage operations for IAM signing certificates.
type SigningCertificateStore struct {
	uk userKeyed[SigningCertificate]
}

// NewSigningCertificateStore creates a new SigningCertificateStore instance.
func NewSigningCertificateStore(store storage.BasicStorage) *SigningCertificateStore {
	return &SigningCertificateStore{
		uk: newUserKeyed[SigningCertificate](
			common.NewBaseStore(store.Bucket(signingCertificateBucketName), "iam"),
			func(c *SigningCertificate) string { return c.CertificateId },
			func(c *SigningCertificate) string { return c.UserName },
		),
	}
}

// Get retrieves a signing certificate by its ID.
func (s *SigningCertificateStore) Get(certificateId string) (*SigningCertificate, error) {
	return getByKey[SigningCertificate](s.uk.BaseStore, certificateId, "get_signing_certificate", ErrSigningCertificateNotFound)
}

// Put stores a signing certificate, keyed by its certificate ID.
func (s *SigningCertificateStore) Put(cert *SigningCertificate) error {
	return s.uk.BaseStore.Put(cert.CertificateId, cert)
}

// Delete removes a signing certificate by its certificate ID.
func (s *SigningCertificateStore) Delete(certificateId string) error {
	return s.uk.BaseStore.Delete(certificateId)
}

// Exists reports whether a signing certificate exists with the given certificate ID.
func (s *SigningCertificateStore) Exists(certificateId string) bool {
	return s.uk.BaseStore.Exists(certificateId)
}

// MaxSigningCertificatesPerUser is the AWS-enforced quota of signing
// certificates per IAM user.
const MaxSigningCertificatesPerUser = 2

// UploadWithGuards creates a signing certificate after checking, inside a
// single lock scope, that the same certificate is not already registered for
// the user and that the per-user quota is not exceeded.
func (s *SigningCertificateStore) UploadWithGuards(userName, certificateBody, fingerprint string) (*SigningCertificate, error) {
	var created *SigningCertificate
	err := s.uk.kl.WithLock("signing-cert:"+userName, func() error {
		existing, err := s.ListByUserName(userName)
		if err != nil {
			return err
		}
		for _, cert := range existing {
			if cert.Fingerprint != "" && cert.Fingerprint == fingerprint {
				return NewStoreError("upload_signing_certificate", ErrDuplicateSigningCertificate)
			}
		}
		if len(existing) >= MaxSigningCertificatesPerUser {
			return NewStoreError("upload_signing_certificate", ErrSigningCertificateLimitExceeded)
		}

		id, err := GenerateSigningCertificateID()
		if err != nil {
			return NewStoreError("generate_signing_certificate_id", err)
		}
		created = &SigningCertificate{
			CertificateId:   id,
			UserName:        userName,
			CertificateBody: certificateBody,
			Status:          "Active",
			UploadDate:      time.Now().UTC(),
			Fingerprint:     fingerprint,
		}
		return s.Put(created)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

// UpdateStatus changes the status of a signing certificate (e.g. Active/Inactive).
func (s *SigningCertificateStore) UpdateStatus(certificateId, status string) error {
	return s.uk.updateStatus(certificateId, status, s.Get, func(c *SigningCertificate, status string) { c.Status = status })
}

// ListByUserName returns all signing certificates belonging to the given user.
func (s *SigningCertificateStore) ListByUserName(userName string) ([]*SigningCertificate, error) {
	return s.uk.listByUserName(userName, "list_signing_certificates")
}

// DeleteAllForUser removes all signing certificates belonging to the given user.
func (s *SigningCertificateStore) DeleteAllForUser(userName string) error {
	return s.uk.deleteAllForUser(userName, "delete_user_signing_certificates")
}

// MigrateUser updates the UserName field on all signing certificates from
// oldUserName to newUserName. Called during IAM user rename operations.
func (s *SigningCertificateStore) MigrateUser(oldUserName, newUserName string) error {
	return s.uk.migrateUser(oldUserName, newUserName, "migrate_signing_certificates",
		func(c *SigningCertificate, newName string) { c.UserName = newName })
}

// Count returns the total number of signing certificates.
func (s *SigningCertificateStore) Count() int {
	return s.uk.BaseStore.Count()
}
