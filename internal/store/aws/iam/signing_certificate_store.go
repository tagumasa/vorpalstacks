package iam

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
)

const signingCertificateBucketName = "iam_signing_certificates"

// SigningCertificateStore provides storage operations for IAM signing certificates.
type SigningCertificateStore struct {
	bucket storage.Bucket
}

// NewSigningCertificateStore creates a new SigningCertificateStore instance.
func NewSigningCertificateStore(store storage.BasicStorage) *SigningCertificateStore {
	return &SigningCertificateStore{
		bucket: store.Bucket(signingCertificateBucketName),
	}
}

// Get retrieves a signing certificate by its ID.
func (s *SigningCertificateStore) Get(certificateId string) (*SigningCertificate, error) {
	data, err := s.bucket.Get([]byte(certificateId))
	if err != nil {
		return nil, NewStoreError("get_signing_certificate", err)
	}
	if data == nil {
		return nil, NewStoreError("get_signing_certificate", ErrSigningCertificateNotFound)
	}
	var cert SigningCertificate
	if err := json.Unmarshal(data, &cert); err != nil {
		return nil, NewStoreError("get_signing_certificate", err)
	}
	return &cert, nil
}

// Put stores a signing certificate, keyed by its certificate ID.
func (s *SigningCertificateStore) Put(cert *SigningCertificate) error {
	data, err := json.Marshal(cert)
	if err != nil {
		return NewStoreError("put_signing_certificate", err)
	}
	if err := s.bucket.Put([]byte(cert.CertificateId), data); err != nil {
		return NewStoreError("put_signing_certificate", err)
	}
	return nil
}

// Delete removes a signing certificate by its certificate ID.
func (s *SigningCertificateStore) Delete(certificateId string) error {
	if err := s.bucket.Delete([]byte(certificateId)); err != nil {
		return NewStoreError("delete_signing_certificate", err)
	}
	return nil
}

// Exists reports whether a signing certificate exists with the given certificate ID.
func (s *SigningCertificateStore) Exists(certificateId string) bool {
	return s.bucket.Has([]byte(certificateId))
}

// Upload uploads a new signing certificate for the given user.
func (s *SigningCertificateStore) Upload(userName, certificateBody string) (*SigningCertificate, error) {
	id, err := generateSigningCertificateID()
	if err != nil {
		return nil, err
	}
	cert := &SigningCertificate{
		CertificateId:   id,
		UserName:        userName,
		CertificateBody: certificateBody,
		Status:          "Active",
		UploadDate:      time.Now().UTC(),
	}
	if err := s.Put(cert); err != nil {
		return nil, err
	}
	return cert, nil
}

// UpdateStatus changes the status of a signing certificate (e.g. Active/Inactive).
func (s *SigningCertificateStore) UpdateStatus(certificateId, status string) error {
	cert, err := s.Get(certificateId)
	if err != nil {
		return err
	}
	cert.Status = status
	return s.Put(cert)
}

// ListByUserName returns all signing certificates belonging to the given user.
func (s *SigningCertificateStore) ListByUserName(userName string) ([]*SigningCertificate, error) {
	var certs []*SigningCertificate
	err := s.bucket.ForEach(func(k, v []byte) error {
		var cert SigningCertificate
		if err := json.Unmarshal(v, &cert); err != nil {
			return err
		}
		if cert.UserName == userName {
			certs = append(certs, &cert)
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_signing_certificates", err)
	}
	return certs, nil
}

// DeleteAllForUser removes all signing certificates belonging to the given user.
func (s *SigningCertificateStore) DeleteAllForUser(userName string) error {
	var toDelete []string
	err := s.bucket.ForEach(func(k, v []byte) error {
		var cert SigningCertificate
		if err := json.Unmarshal(v, &cert); err != nil {
			return err
		}
		if cert.UserName == userName {
			toDelete = append(toDelete, cert.CertificateId)
		}
		return nil
	})
	if err != nil {
		return NewStoreError("delete_user_signing_certificates", err)
	}
	for _, id := range toDelete {
		if err := s.Delete(id); err != nil {
			return err
		}
	}
	return nil
}

// Count returns the total number of signing certificates.
func (s *SigningCertificateStore) Count() int {
	return s.bucket.Count()
}
