package iam

import (
	"encoding/json"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const signingCertificateBucketName = "iam_signing_certificates"

// SigningCertificateStore provides storage operations for IAM signing certificates.
type SigningCertificateStore struct {
	*common.BaseStore
	kl common.KeyLocker
}

// NewSigningCertificateStore creates a new SigningCertificateStore instance.
func NewSigningCertificateStore(store storage.BasicStorage) *SigningCertificateStore {
	return &SigningCertificateStore{
		BaseStore: common.NewBaseStore(store.Bucket(signingCertificateBucketName), "iam"),
	}
}

// Get retrieves a signing certificate by its ID.
func (s *SigningCertificateStore) Get(certificateId string) (*SigningCertificate, error) {
	var cert SigningCertificate
	if err := s.BaseStore.Get(certificateId, &cert); err != nil {
		if common.IsNotFound(err) {
			return nil, NewStoreError("get_signing_certificate", ErrSigningCertificateNotFound)
		}
		return nil, NewStoreError("get_signing_certificate", err)
	}
	return &cert, nil
}

// Put stores a signing certificate, keyed by its certificate ID.
func (s *SigningCertificateStore) Put(cert *SigningCertificate) error {
	return s.BaseStore.Put(cert.CertificateId, cert)
}

// Delete removes a signing certificate by its certificate ID.
func (s *SigningCertificateStore) Delete(certificateId string) error {
	return s.BaseStore.Delete(certificateId)
}

// Exists reports whether a signing certificate exists with the given certificate ID.
func (s *SigningCertificateStore) Exists(certificateId string) bool {
	return s.BaseStore.Exists(certificateId)
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
	return s.kl.WithLock(certificateId, func() error {
		cert, err := s.Get(certificateId)
		if err != nil {
			return err
		}
		cert.Status = status
		return s.Put(cert)
	})
}

// ListByUserName returns all signing certificates belonging to the given user.
func (s *SigningCertificateStore) ListByUserName(userName string) ([]*SigningCertificate, error) {
	var certs []*SigningCertificate
	err := s.ForEach(func(k string, v []byte) error {
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
	err := s.ForEach(func(k string, v []byte) error {
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
	return s.BaseStore.Count()
}
